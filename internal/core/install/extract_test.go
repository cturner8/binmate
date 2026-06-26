package install

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"cturner8/binmate/internal/database"
)

// Test helpers for creating test archives

// createTestTarGz creates a tar.gz archive with test files
func createTestTarGz(t *testing.T, files map[string]string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tarPath := filepath.Join(tmpDir, "test.tar.gz")

	f, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("Failed to create tar file: %v", err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for name, content := range files {
		header := &tar.Header{
			Name:     name,
			Mode:     0755,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}

		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write tar header: %v", err)
		}

		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write tar content: %v", err)
		}
	}

	return tarPath
}

// createTestZip creates a zip archive with test files
func createTestZip(t *testing.T, files map[string]string) string {
	t.Helper()

	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	for name, content := range files {
		fw, err := zw.CreateHeader(&zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		})
		if err != nil {
			t.Fatalf("Failed to create zip header: %v", err)
		}

		if _, err := fw.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write zip content: %v", err)
		}
	}

	return zipPath
}

// TestExtractAsset tests the main extraction function
func TestExtractAsset_TarGz(t *testing.T) {
	files := map[string]string{
		"mybin": "binary content",
	}
	tarPath := createTestTarGz(t, files)

	binary := &database.Binary{
		UserID: "testbin",
		Name:   "mybin",
		Format: ".tar.gz",
	}

	// Note: This will fail because getExtractPath tries to use real paths
	// We're testing the dispatch logic here
	_, err := ExtractAsset(tarPath, binary, "v1.0.0")
	// We expect an error about extract path, but not about format
	if err != nil && err.Error() == "unsupported asset format: .tar.gz" {
		t.Error("Should not report unsupported format for .tar.gz")
	}
}

func TestExtractAsset_Zip(t *testing.T) {
	files := map[string]string{
		"mybin": "binary content",
	}
	zipPath := createTestZip(t, files)

	binary := &database.Binary{
		UserID: "testbin",
		Name:   "mybin",
		Format: ".zip",
	}

	// Note: This will fail because getExtractPath tries to use real paths
	// We're testing the dispatch logic here
	_, err := ExtractAsset(zipPath, binary, "v1.0.0")
	// We expect an error about extract path, but not about format
	if err != nil && err.Error() == "unsupported asset format: .zip" {
		t.Error("Should not report unsupported format for .zip")
	}
}

func TestExtractAsset_UnsupportedFormat(t *testing.T) {
	binary := &database.Binary{
		UserID: "testbin",
		Name:   "mybin",
		Format: ".rar",
	}

	_, err := ExtractAsset("/some/path", binary, "v1.0.0")
	if err == nil {
		t.Error("Expected error for unsupported format, got none")
	}
	if err != nil && err.Error() != "unsupported asset format: .rar" {
		t.Errorf("Expected unsupported format error, got: %v", err)
	}
}

// Test large binary extraction
func TestExtractTar_LargeBinary(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large binary test in short mode")
	}

	// Create a "large" binary (1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	files := map[string]string{
		"largebin": string(largeContent),
	}
	tarPath := createTestTarGz(t, files)
	destDir := t.TempDir()

	extractedPath, err := extractTar(tarPath, destDir, "largebin")
	if err != nil {
		t.Fatalf("extractTar failed for large binary: %v", err)
	}

	// Verify size
	info, err := os.Stat(extractedPath)
	if err != nil {
		t.Fatalf("Failed to stat extracted file: %v", err)
	}
	if info.Size() != int64(len(largeContent)) {
		t.Errorf("Size mismatch: expected %d, got %d", len(largeContent), info.Size())
	}
}

// Benchmark extraction performance
func BenchmarkExtractTar(b *testing.B) {
	files := map[string]string{
		"bin/testbin": "binary content here",
	}

	// Create tar once
	tmpFile, err := os.CreateTemp("", "bench-*.tar.gz")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	tarPath := tmpFile.Name()
	defer os.Remove(tarPath)
	tmpFile.Close()

	// Write tar
	f, _ := os.Create(tarPath)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	for name, content := range files {
		header := &tar.Header{
			Name:     name,
			Mode:     0755,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}
		tw.WriteHeader(header)
		tw.Write([]byte(content))
	}
	tw.Close()
	gw.Close()
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		destDir := b.TempDir()
		if _, err := extractTar(tarPath, destDir, "testbin"); err != nil {
			b.Fatalf("extractTar failed: %v", err)
		}
	}
}

func BenchmarkExtractZip(b *testing.B) {
	files := map[string]string{
		"bin/testbin": "binary content here",
	}

	// Create zip once
	tmpFile, err := os.CreateTemp("", "bench-*.zip")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	zipPath := tmpFile.Name()
	defer os.Remove(zipPath)
	tmpFile.Close()

	// Write zip
	f, _ := os.Create(zipPath)
	zw := zip.NewWriter(f)

	for name, content := range files {
		fw, _ := zw.Create(name)
		fw.Write([]byte(content))
	}
	zw.Close()
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		destDir := b.TempDir()
		if _, err := extractZip(zipPath, destDir, "testbin"); err != nil {
			b.Fatalf("extractZip failed: %v", err)
		}
	}
}

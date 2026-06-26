package install

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeTarGz writes a tar.gz to path with the given files (name -> content).
// Pass a nil modifier to use defaults; modifier receives the header before it is written.
func makeTarGz(t *testing.T, files map[string]string, modifier func(*tar.Header)) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.tar.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create tar.gz: %v", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	for name, content := range files {
		hdr := &tar.Header{
			Name:     name,
			Mode:     0o755,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}
		if modifier != nil {
			modifier(hdr)
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("write header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("write content: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}
	return path
}

// --- extractTar ---

func TestExtractTar_EmptyDestDir(t *testing.T) {
	p := makeTarGz(t, map[string]string{"bin": "x"}, nil)
	_, err := extractTar(p, "", "bin")
	if err == nil {
		t.Errorf("extractTar() expected error for empty destDir, got none")
	}
	if !strings.Contains(err.Error(), "destination directory is required") {
		t.Errorf("extractTar() expected 'destination directory is required' error, got %v", err)
	}
}

func TestExtractTar_EmptyBinaryName(t *testing.T) {
	p := makeTarGz(t, map[string]string{"bin": "x"}, nil)
	_, err := extractTar(p, t.TempDir(), "")
	if err == nil {
		t.Errorf("extractTar() expected error for empty binary name, got none")
	}
	if !strings.Contains(err.Error(), "binary name is required") {
		t.Errorf("extractTar() expected 'binary name is required' error, got %v", err)
	}
}

func TestExtractTar_NonExistentFile(t *testing.T) {
	_, err := extractTar("/nonexistent/archive.tar.gz", t.TempDir(), "bin")
	if err == nil {
		t.Errorf("extractTar() expected error for nonexistent archive, got none")
	}
	if !strings.Contains(err.Error(), "open tar") {
		t.Errorf("extractTar() expected 'open tar' error, got %v", err)
	}
}

func TestExtractTar_NotGzip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.tar.gz")
	if err := os.WriteFile(p, []byte("not gzip data"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := extractTar(p, t.TempDir(), "bin")
	if err == nil {
		t.Errorf("extractTar() expected error for non-gzip data, got none")
	}
	if !strings.Contains(err.Error(), "create gzip reader") {
		t.Errorf("extractTar() expected 'create gzip reader' error, got %v", err)
	}
}

func TestExtractTar_BinaryNotFound(t *testing.T) {
	p := makeTarGz(t, map[string]string{"other": "data"}, nil)
	_, err := extractTar(p, t.TempDir(), "missing")
	if err == nil {
		t.Errorf("extractTar() expected error for binary not found, got none")
	}
	if !strings.Contains(err.Error(), "not found in archive") {
		t.Errorf("extractTar() expected 'not found in archive' error, got %v", err)
	}
}

func TestExtractTar_Success(t *testing.T) {
	content := "#!/bin/sh\necho hi"
	p := makeTarGz(t, map[string]string{"mybin": content}, nil)
	destDir := t.TempDir()

	got, err := extractTar(p, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractTar() unexpected error: %v", err)
	}
	want := filepath.Join(destDir, "mybin")
	if got != want {
		t.Errorf("path = %s, want %s", got, want)
	}
	data, _ := os.ReadFile(got)
	if string(data) != content {
		t.Errorf("content = %q, want %q", string(data), content)
	}
	info, _ := os.Stat(got)
	if info.Mode().Perm()&0o100 == 0 {
		t.Error("extracted binary should be executable")
	}
}

func TestExtractTar_BinaryInSubdirectory(t *testing.T) {
	p := makeTarGz(t, map[string]string{"bin/mybin": "data"}, nil)
	destDir := t.TempDir()

	got, err := extractTar(p, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractTar() unexpected error: %v", err)
	}
	if got != filepath.Join(destDir, "mybin") {
		t.Errorf("path = %s, want %s", got, filepath.Join(destDir, "mybin"))
	}
}

func TestExtractTar_CreatesDestDir(t *testing.T) {
	p := makeTarGz(t, map[string]string{"bin": "data"}, nil)
	destDir := filepath.Join(t.TempDir(), "a", "b")

	_, err := extractTar(p, destDir, "bin")
	if err != nil {
		t.Fatalf("extractTar() should create destDir: %v", err)
	}
}

func TestExtractTar_MkdirAllFails(t *testing.T) {
	p := makeTarGz(t, map[string]string{"bin": "data"}, nil)
	// Place a regular file where destDir's parent should be, so MkdirAll fails
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	destDir := filepath.Join(blocker, "sub")

	_, err := extractTar(p, destDir, "bin")
	if err == nil {
		t.Errorf("extractTar() expected error when MkdirAll fails, got none")
	}
	if !strings.Contains(err.Error(), "create destination") {
		t.Errorf("extractTar() expected 'create destination' error, got %v", err)
	}
}

func TestExtractTar_MultipleFiles(t *testing.T) {
	p := makeTarGz(t, map[string]string{
		"README":     "readme",
		"bin/target": "correct",
		"bin/other":  "wrong",
	}, nil)
	destDir := t.TempDir()

	got, err := extractTar(p, destDir, "target")
	if err != nil {
		t.Fatalf("extractTar() unexpected error: %v", err)
	}
	data, _ := os.ReadFile(got)
	if string(data) != "correct" {
		t.Errorf("extracted wrong file: %q", string(data))
	}
}

// --- extractTarBinary ---

func TestExtractTarBinary_RejectsSymlink(t *testing.T) {
	hdr := &tar.Header{
		Name:     "link",
		Typeflag: tar.TypeSymlink,
		Linkname: "/etc/passwd",
	}
	err := extractTarBinary(nil, hdr, filepath.Join(t.TempDir(), "out"))
	if err == nil {
		t.Errorf("extractTarBinary() expected error for symlink in archive, got none")
	}
	if !strings.Contains(err.Error(), "symlinks are not supported") {
		t.Errorf("extractTarBinary() expected 'symlinks are not supported' error, got %v", err)
	}
}

func TestExtractTarBinary_RejectsHardlink(t *testing.T) {
	hdr := &tar.Header{
		Name:     "link",
		Typeflag: tar.TypeLink,
		Linkname: "/etc/passwd",
	}
	err := extractTarBinary(nil, hdr, filepath.Join(t.TempDir(), "out"))
	if err == nil {
		t.Errorf("extractTarBinary() expected error for hardlink in archive, got none")
	}
	if !strings.Contains(err.Error(), "symlinks are not supported") {
		t.Errorf("extractTarBinary() expected 'symlinks are not supported' error, got %v", err)
	}
}

func TestExtractTarBinary_ZeroModeDefaultsToExecutable(t *testing.T) {
	content := "binary content"
	p := makeTarGz(t, map[string]string{"bin": content}, func(h *tar.Header) {
		h.Mode = 0 // zero mode — should default to 0755
	})
	destDir := t.TempDir()

	_, err := extractTar(p, destDir, "bin")
	if err != nil {
		t.Fatalf("extractTar() unexpected error: %v", err)
	}
	info, _ := os.Stat(filepath.Join(destDir, "bin"))
	if info.Mode().Perm()&0o100 == 0 {
		t.Error("binary with zero mode should default to executable")
	}
}

func TestExtractTarBinary_CannotCreateDest(t *testing.T) {
	// Build a real tar reader pointing at a valid entry
	content := "data"
	p := makeTarGz(t, map[string]string{"bin": content}, nil)
	f, _ := os.Open(p)
	defer f.Close()
	gr, _ := gzip.NewReader(f)
	defer gr.Close()
	tr := tar.NewReader(gr)
	hdr, _ := tr.Next()

	err := extractTarBinary(tr, hdr, "/nonexistent/dir/out")
	if err == nil {
		t.Errorf("extractTarBinary() expected error when destination cannot be created, got none")
	}
	if !strings.Contains(err.Error(), "create file") {
		t.Errorf("extractTarBinary() expected 'create file' error, got %v", err)
	}
}

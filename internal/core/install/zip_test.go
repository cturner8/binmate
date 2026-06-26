package install

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeZip writes a zip archive to a temp file with the given files (name -> content).
// modifier, if non-nil, is called with each FileHeader before the entry is written.
func makeZip(t *testing.T, files map[string]string, modifier func(*zip.FileHeader)) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.zip")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for name, content := range files {
		hdr := &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}
		hdr.SetMode(0o755)
		if modifier != nil {
			modifier(hdr)
		}
		fw, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatalf("create header: %v", err)
		}
		if _, err := fw.Write([]byte(content)); err != nil {
			t.Fatalf("write content: %v", err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	return path
}

// --- extractZip ---

func TestExtractZip_EmptyDestDir(t *testing.T) {
	p := makeZip(t, map[string]string{"bin": "x"}, nil)
	_, err := extractZip(p, "", "bin")
	if err == nil {
		t.Errorf("extractZip() expected error for empty destDir, got none")
	}
	if !strings.Contains(err.Error(), "destination directory is required") {
		t.Errorf("extractZip() expected 'destination directory is required' error, got %v", err)
	}
}

func TestExtractZip_EmptyBinaryName(t *testing.T) {
	p := makeZip(t, map[string]string{"bin": "x"}, nil)
	_, err := extractZip(p, t.TempDir(), "")
	if err == nil {
		t.Errorf("extractZip() expected error for empty binary name, got none")
	}
	if !strings.Contains(err.Error(), "binary name is required") {
		t.Errorf("extractZip() expected 'binary name is required' error, got %v", err)
	}
}

func TestExtractZip_NonExistentFile(t *testing.T) {
	_, err := extractZip("/nonexistent/archive.zip", t.TempDir(), "bin")
	if err == nil {
		t.Errorf("extractZip() expected error for nonexistent archive, got none")
	}
	if !strings.Contains(err.Error(), "open zip") {
		t.Errorf("extractZip() expected 'open zip' error, got %v", err)
	}
}

func TestExtractZip_NotZip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.zip")
	if err := os.WriteFile(p, []byte("not a zip file"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := extractZip(p, t.TempDir(), "bin")
	if err == nil {
		t.Errorf("extractZip() expected error for invalid zip data, got none")
	}
	if !strings.Contains(err.Error(), "open zip") {
		t.Errorf("extractZip() expected 'open zip' error, got %v", err)
	}
}

func TestExtractZip_BinaryNotFound(t *testing.T) {
	p := makeZip(t, map[string]string{"other": "data"}, nil)
	_, err := extractZip(p, t.TempDir(), "missing")
	if err == nil {
		t.Errorf("extractZip() expected error for binary not found, got none")
	}
	if !strings.Contains(err.Error(), "not found in archive") {
		t.Errorf("extractZip() expected 'not found in archive' error, got %v", err)
	}
}

func TestExtractZip_Success(t *testing.T) {
	content := "#!/bin/sh\necho hi"
	p := makeZip(t, map[string]string{"mybin": content}, nil)
	destDir := t.TempDir()

	got, err := extractZip(p, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractZip() unexpected error: %v", err)
	}
	want := filepath.Join(destDir, "mybin")
	if got != want {
		t.Errorf("path = %s, want %s", got, want)
	}
	data, _ := os.ReadFile(got)
	if string(data) != content {
		t.Errorf("content = %q, want %q", string(data), content)
	}
}

func TestExtractZip_BinaryInSubdirectory(t *testing.T) {
	p := makeZip(t, map[string]string{"bin/mybin": "data"}, nil)
	destDir := t.TempDir()

	got, err := extractZip(p, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractZip() unexpected error: %v", err)
	}
	if got != filepath.Join(destDir, "mybin") {
		t.Errorf("path = %s, want %s", got, filepath.Join(destDir, "mybin"))
	}
}

func TestExtractZip_CreatesDestDir(t *testing.T) {
	p := makeZip(t, map[string]string{"bin": "data"}, nil)
	destDir := filepath.Join(t.TempDir(), "a", "b")

	_, err := extractZip(p, destDir, "bin")
	if err != nil {
		t.Fatalf("extractZip() should create destDir: %v", err)
	}
}

func TestExtractZip_MkdirAllFails(t *testing.T) {
	p := makeZip(t, map[string]string{"bin": "data"}, nil)
	// Place a regular file where destDir's parent should be, so MkdirAll fails
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	destDir := filepath.Join(blocker, "sub")

	_, err := extractZip(p, destDir, "bin")
	if err == nil {
		t.Errorf("extractZip() expected error when MkdirAll fails, got none")
	}
	if !strings.Contains(err.Error(), "create destination") {
		t.Errorf("extractZip() expected 'create destination' error, got %v", err)
	}
}

func TestExtractZip_MultipleFiles(t *testing.T) {
	p := makeZip(t, map[string]string{
		"README":     "readme",
		"bin/target": "correct",
		"bin/other":  "wrong",
	}, nil)
	destDir := t.TempDir()

	got, err := extractZip(p, destDir, "target")
	if err != nil {
		t.Fatalf("extractZip() unexpected error: %v", err)
	}
	data, _ := os.ReadFile(got)
	if string(data) != "correct" {
		t.Errorf("extracted wrong file: %q", string(data))
	}
}

// --- extractZipBinary ---

func TestExtractZipBinary_RejectsSymlink(t *testing.T) {
	p := makeZip(t, map[string]string{"link": "target"}, func(hdr *zip.FileHeader) {
		hdr.SetMode(os.ModeSymlink | 0o777)
	})
	r, err := zip.OpenReader(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	err = extractZipBinary(r.File[0], filepath.Join(t.TempDir(), "out"))
	if err == nil {
		t.Errorf("extractZipBinary() expected error for symlink in archive, got none")
	}
	if !strings.Contains(err.Error(), "symlinks are not supported") {
		t.Errorf("extractZipBinary() expected 'symlinks are not supported' error, got %v", err)
	}
}

func TestExtractZipBinary_ZeroModeDefaultsToExecutable(t *testing.T) {
	content := "binary"
	p := makeZip(t, map[string]string{"bin": content}, func(hdr *zip.FileHeader) {
		hdr.SetMode(0) // zero mode — should default to 0755
	})
	destDir := t.TempDir()

	_, err := extractZip(p, destDir, "bin")
	if err != nil {
		t.Fatalf("extractZip() unexpected error: %v", err)
	}
	info, _ := os.Stat(filepath.Join(destDir, "bin"))
	if info.Mode().Perm()&0o100 == 0 {
		t.Error("binary with zero mode should default to executable")
	}
}

func TestExtractZipBinary_CannotCreateDest(t *testing.T) {
	p := makeZip(t, map[string]string{"bin": "data"}, nil)
	r, err := zip.OpenReader(p)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	err = extractZipBinary(r.File[0], "/nonexistent/dir/out")
	if err == nil {
		t.Errorf("extractZipBinary() expected error when destination cannot be created, got none")
	}
	if !strings.Contains(err.Error(), "create file") {
		t.Errorf("extractZipBinary() expected 'create file' error, got %v", err)
	}
}

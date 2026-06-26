package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractRaw_EmptyDestDir(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(tmpFile, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := extractRaw(tmpFile, "", "bin")
	if err == nil {
		t.Errorf("extractRaw() expected error for empty destDir, got none")
	}
	if !strings.Contains(err.Error(), "destination directory is required") {
		t.Errorf("extractRaw() expected 'destination directory is required' error, got %v", err)
	}
}

func TestExtractRaw_EmptyBinaryName(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(tmpFile, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	destDir := t.TempDir()

	_, err := extractRaw(tmpFile, destDir, "")
	if err == nil {
		t.Errorf("extractRaw() expected error for empty binary name, got none")
	}
	if !strings.Contains(err.Error(), "binary name is required") {
		t.Errorf("extractRaw() expected 'binary name is required' error, got %v", err)
	}
}

func TestExtractRaw_SourceNotFound(t *testing.T) {
	destDir := t.TempDir()

	_, err := extractRaw("/nonexistent/path/bin", destDir, "bin")
	if err == nil {
		t.Errorf("extractRaw() expected error for missing source file, got none")
	}
	if !strings.Contains(err.Error(), "source file") {
		t.Errorf("extractRaw() expected 'source file' error, got %v", err)
	}
}

func TestExtractRaw_SourceIsDirectory(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	_, err := extractRaw(srcDir, destDir, "bin")
	if err == nil {
		t.Errorf("extractRaw() expected error when source path is a directory, got none")
	}
	if !strings.Contains(err.Error(), "is a directory") {
		t.Errorf("extractRaw() expected 'is a directory' error, got %v", err)
	}
}

func TestExtractRaw_Success(t *testing.T) {
	content := []byte("#!/bin/sh\necho hello")
	srcFile := filepath.Join(t.TempDir(), "mybin")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatal(err)
	}
	destDir := t.TempDir()

	got, err := extractRaw(srcFile, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractRaw() unexpected error: %v", err)
	}

	want := filepath.Join(destDir, "mybin")
	if got != want {
		t.Errorf("extractRaw() path = %s, want %s", got, want)
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", string(data), string(content))
	}

	info, err := os.Stat(got)
	if err != nil {
		t.Fatalf("failed to stat extracted file: %v", err)
	}
	if info.Mode().Perm()&0o100 == 0 {
		t.Error("extracted binary should be executable")
	}
}

func TestExtractRaw_CreatesDestDir(t *testing.T) {
	content := []byte("binary")
	srcFile := filepath.Join(t.TempDir(), "mybin")
	if err := os.WriteFile(srcFile, content, 0o644); err != nil {
		t.Fatal(err)
	}
	// Use a nested path that does not exist yet
	destDir := filepath.Join(t.TempDir(), "a", "b", "c")

	_, err := extractRaw(srcFile, destDir, "mybin")
	if err != nil {
		t.Fatalf("extractRaw() should create dest dir, got: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "mybin")); os.IsNotExist(err) {
		t.Error("expected binary to exist in created dest dir")
	}
}

func TestExtractRawBinary_CannotCreateDest(t *testing.T) {
	srcFile := filepath.Join(t.TempDir(), "src")
	if err := os.WriteFile(srcFile, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Target path inside a non-existent directory so OpenFile fails
	targetPath := "/nonexistent/dir/out"

	err := extractRawBinary(srcFile, targetPath)
	if err == nil {
		t.Errorf("extractRawBinary() expected error when destination cannot be created, got none")
	}
	if !strings.Contains(err.Error(), "create file") {
		t.Errorf("extractRawBinary() expected 'create file' error, got %v", err)
	}
}

func TestExtractRawBinary_CannotOpenSource(t *testing.T) {
	destDir := t.TempDir()
	targetPath := filepath.Join(destDir, "out")

	err := extractRawBinary("/nonexistent/source", targetPath)
	if err == nil {
		t.Errorf("extractRawBinary() expected error when source cannot be opened, got none")
	}
	if !strings.Contains(err.Error(), "open source file") {
		t.Errorf("extractRawBinary() expected 'open source file' error, got %v", err)
	}
}

func TestExtractRaw_ExtractBinaryFails(t *testing.T) {
	srcFile := filepath.Join(t.TempDir(), "src")
	if err := os.WriteFile(srcFile, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	destDir := t.TempDir()
	// Place a directory at the target path so os.OpenFile(O_WRONLY) fails with EISDIR
	conflict := filepath.Join(destDir, "out")
	if err := os.Mkdir(conflict, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := extractRaw(srcFile, destDir, "out")
	if err == nil {
		t.Errorf("extractRaw() expected error when binary write fails, got none")
	}
	if !strings.Contains(err.Error(), "create file") {
		t.Errorf("extractRaw() expected 'create file' error, got %v", err)
	}
}

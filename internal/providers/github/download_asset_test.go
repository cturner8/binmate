package github

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveAssetWritesFileToDestination(t *testing.T) {
	dir := t.TempDir()
	content := []byte("binary content")
	assetName := "tool-1.0.0-linux_amd64.tar.gz"

	destPath, err := saveAsset(bytes.NewReader(content), dir, assetName)
	if err != nil {
		t.Fatalf("saveAsset() unexpected error: %v", err)
	}

	if destPath != filepath.Join(dir, assetName) {
		t.Errorf("saveAsset() destPath = %q, want %q", destPath, filepath.Join(dir, assetName))
	}

	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("reading dest file: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("saveAsset() file content = %q, want %q", got, content)
	}
}

func TestSaveAssetCreatesDirIfAbsent(t *testing.T) {
	parent := t.TempDir()
	basePath := filepath.Join(parent, "binmate")
	assetName := "tool-1.0.0-linux_amd64.tar.gz"

	if _, err := os.Stat(basePath); !os.IsNotExist(err) {
		t.Fatalf("expected basePath to not exist before saveAsset")
	}

	if _, err := saveAsset(strings.NewReader("data"), basePath, assetName); err != nil {
		t.Fatalf("saveAsset() unexpected error: %v", err)
	}

	if _, err := os.Stat(basePath); err != nil {
		t.Errorf("saveAsset() did not create basePath: %v", err)
	}
}

func TestSaveAssetTempFileInSameDir(t *testing.T) {
	// Verify no leftover temp files remain after a successful save.
	dir := t.TempDir()
	assetName := "tool-1.0.0-linux_amd64.tar.gz"

	if _, err := saveAsset(strings.NewReader("data"), dir, assetName); err != nil {
		t.Fatalf("saveAsset() unexpected error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected exactly 1 file in dir, got %d: %v", len(entries), names)
	}
	if entries[0].Name() != assetName {
		t.Errorf("expected file %q, got %q", assetName, entries[0].Name())
	}
}

func TestSaveAssetWriteError(t *testing.T) {
	dir := t.TempDir()
	// errReader always returns an error on read.
	_, err := saveAsset(&errReader{}, dir, "asset.tar.gz")
	if err == nil {
		t.Fatal("saveAsset() expected error from write failure, got nil")
	}

	// Temp file must be cleaned up.
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected temp file to be cleaned up, found %d entries", len(entries))
	}
}

type errReader struct{}

func (e *errReader) Read(_ []byte) (int, error) {
	return 0, os.ErrInvalid
}

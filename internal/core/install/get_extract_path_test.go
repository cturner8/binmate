package install

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetLocalDataDir_XDGDataHomeSet(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG_DATA_HOME is not used on Windows")
	}
	expected := "/custom/xdg/data"
	t.Setenv("XDG_DATA_HOME", expected)

	got, err := getLocalDataDir()
	if err != nil {
		t.Fatalf("getLocalDataDir() unexpected error: %v", err)
	}
	if got != expected {
		t.Errorf("getLocalDataDir() = %q, want %q", got, expected)
	}
}

func TestGetLocalDataDir_XDGDataHomeNotSet(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG_DATA_HOME is not used on Windows")
	}
	// Explicitly unset so the UserHomeDir fallback is exercised
	orig, wasSet := os.LookupEnv("XDG_DATA_HOME")
	os.Unsetenv("XDG_DATA_HOME")
	t.Cleanup(func() {
		if wasSet {
			os.Setenv("XDG_DATA_HOME", orig)
		} else {
			os.Unsetenv("XDG_DATA_HOME")
		}
	})

	got, err := getLocalDataDir()
	if err != nil {
		t.Fatalf("getLocalDataDir() unexpected error: %v", err)
	}
	if !strings.HasSuffix(got, filepath.Join(".local", "share")) {
		t.Errorf("getLocalDataDir() = %q, want path ending with .local/share", got)
	}
}

func TestGetExtractPath_ReturnsCorrectStructure(t *testing.T) {
	baseDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", baseDir)

	binaryID := "mybinary"
	version := "v1.2.3"

	got, err := getExtractPath(binaryID, version)
	if err != nil {
		t.Fatalf("getExtractPath() unexpected error: %v", err)
	}
	want := filepath.Join(baseDir, "binmate", "versions", binaryID, version)
	if got != want {
		t.Errorf("getExtractPath() = %q, want %q", got, want)
	}
}

func TestGetExtractPath_EmbedsBinaryIDAndVersion(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	binaryID := "specialbin"
	version := "v0.99.0"

	got, err := getExtractPath(binaryID, version)
	if err != nil {
		t.Fatalf("getExtractPath() unexpected error: %v", err)
	}
	if !strings.Contains(got, binaryID) {
		t.Errorf("getExtractPath() = %q, expected path to contain binary ID %q", got, binaryID)
	}
	if !strings.Contains(got, version) {
		t.Errorf("getExtractPath() = %q, expected path to contain version %q", got, version)
	}
}

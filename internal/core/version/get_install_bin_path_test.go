package version

import (
	"runtime"
	"testing"
)

func TestGetInstallBinPath_NonEmptyInstallPath_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows OS")
	}

	installPath, err := getInstallBinPath("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if installPath != "/usr/local/bin" {
		t.Errorf("Expected install path /usr/local/bin, got %s", installPath)
	}
}

func TestGetInstallBinPath_NonEmptyInstallPath_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	t.Setenv("HOME", "/home/go")

	installPath, err := getInstallBinPath("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if installPath != "/home/go/.local/bin" {
		t.Errorf("Expected install path /home/go/.local/bin, got %s", installPath)
	}
}

func TestGetInstallBinPath_UserProvidedInstallPath(t *testing.T) {
	installPath, err := getInstallBinPath("/usr/local/bin")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if installPath != "/usr/local/bin" {
		t.Errorf("Expected install path /usr/local/bin, got %s", installPath)
	}
}

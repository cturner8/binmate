package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetActiveVersion_MultipleSwitch(t *testing.T) {
	// Create test directories
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	installDir1 := filepath.Join(tmpDir, "v1.0.0")
	installDir2 := filepath.Join(tmpDir, "v2.0.0")

	// Create bin directory
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	// Create test binary files
	binary1 := filepath.Join(installDir1, "testbin")
	binary2 := filepath.Join(installDir2, "testbin")
	if err := os.MkdirAll(installDir1, 0755); err != nil {
		t.Fatalf("Failed to create install dir 1: %v", err)
	}
	if err := os.MkdirAll(installDir2, 0755); err != nil {
		t.Fatalf("Failed to create install dir 2: %v", err)
	}
	if err := os.WriteFile(binary1, []byte("#!/bin/bash\necho v1"), 0755); err != nil {
		t.Fatalf("Failed to create test binary 1: %v", err)
	}
	if err := os.WriteFile(binary2, []byte("#!/bin/bash\necho v2"), 0755); err != nil {
		t.Fatalf("Failed to create test binary 2: %v", err)
	}

	// Test switching from v1 to v2
	symlinkPath1, err := SetActiveVersion(binary1, binDir, "testbin", nil)
	if err != nil {
		t.Fatalf("Failed to set active version to v1: %v", err)
	}

	expectedSymlink := filepath.Join(binDir, "testbin")
	if symlinkPath1 != expectedSymlink {
		t.Errorf("Expected symlink path %s, got %s", expectedSymlink, symlinkPath1)
	}

	// Verify symlink points to v1
	target, err := os.Readlink(symlinkPath1)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if target != binary1 {
		t.Errorf("Expected symlink to point to %s, got %s", binary1, target)
	}

	// Now switch to v2 - this is where the bug would occur
	symlinkPath2, err := SetActiveVersion(binary2, binDir, "testbin", nil)
	if err != nil {
		t.Fatalf("Failed to set active version to v2: %v", err)
	}

	if symlinkPath2 != expectedSymlink {
		t.Errorf("Expected symlink path %s, got %s", expectedSymlink, symlinkPath2)
	}

	// Verify symlink now points to v2
	target, err = os.Readlink(symlinkPath2)
	if err != nil {
		t.Fatalf("Failed to read symlink after switch: %v", err)
	}
	if target != binary2 {
		t.Errorf("Expected symlink to point to %s after switch, got %s", binary2, target)
	}

	// Switch back to v1 to test multiple switches
	symlinkPath3, err := SetActiveVersion(binary1, binDir, "testbin", nil)
	if err != nil {
		t.Fatalf("Failed to set active version back to v1: %v", err)
	}

	// Verify symlink points to v1 again
	target, err = os.Readlink(symlinkPath3)
	if err != nil {
		t.Fatalf("Failed to read symlink after second switch: %v", err)
	}
	if target != binary1 {
		t.Errorf("Expected symlink to point to %s after second switch, got %s", binary1, target)
	}
}

func TestSetActiveVersion_NewInstall(t *testing.T) {
	// Create test directories
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	installDir := filepath.Join(tmpDir, "v1.0.0")

	// Create test binary file
	binary := filepath.Join(installDir, "testbin")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatalf("Failed to create install dir: %v", err)
	}
	if err := os.WriteFile(binary, []byte("#!/bin/bash\necho test"), 0755); err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}

	// Test initial installation (no existing symlink)
	symlinkPath, err := SetActiveVersion(binary, binDir, "testbin", nil)
	if err != nil {
		t.Fatalf("Failed to set active version: %v", err)
	}

	expectedSymlink := filepath.Join(binDir, "testbin")
	if symlinkPath != expectedSymlink {
		t.Errorf("Expected symlink path %s, got %s", expectedSymlink, symlinkPath)
	}

	// Verify symlink was created and points to correct target
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if target != binary {
		t.Errorf("Expected symlink to point to %s, got %s", binary, target)
	}
}

func TestSetActiveVersion_WithAlias(t *testing.T) {
	// Create test directories
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	installDir := filepath.Join(tmpDir, "v1.0.0")

	// Create test binary file
	binary := filepath.Join(installDir, "testbin")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		t.Fatalf("Failed to create install dir: %v", err)
	}
	if err := os.WriteFile(binary, []byte("#!/bin/bash\necho test"), 0755); err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}

	// Test with alias
	alias := "myalias"
	symlinkPath, err := SetActiveVersion(binary, binDir, "testbin", &alias)
	if err != nil {
		t.Fatalf("Failed to set active version with alias: %v", err)
	}

	// Verify symlink uses alias name instead of binary name
	expectedSymlink := filepath.Join(binDir, "testbin")
	if symlinkPath != expectedSymlink {
		t.Errorf("Expected symlink path %s, got %s", expectedSymlink, symlinkPath)
	}

	// Verify symlink was created and points to correct target
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Fatalf("Failed to read symlink: %v", err)
	}
	if target != binary {
		t.Errorf("Expected symlink to point to %s, got %s", binary, target)
	}

	// Verify alias symlink
	aliasSymlinkPath := filepath.Join(binDir, alias)
	aliasTarget, err := os.Readlink(aliasSymlinkPath)
	if err != nil {
		t.Fatalf("Failed to read alias symlink: %v", err)
	}
	if aliasTarget != binary {
		t.Errorf("Expected alias symlink to point to %s, got %s", binary, aliasTarget)
	}
}

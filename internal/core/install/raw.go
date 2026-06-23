package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// extractRaw extracts the specified binary into destDir.
func extractRaw(srcPath string, destDir string, binaryName string) (string, error) {
	if destDir == "" {
		return "", fmt.Errorf("destination directory is required")
	}
	if binaryName == "" {
		return "", fmt.Errorf("binary name is required")
	}

	// Check if source file exists and is readable
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return "", fmt.Errorf("source file %s: %w", srcPath, err)
	}
	// Ensure that the source path is a file, not a directory
	if srcInfo.IsDir() {
		return "", fmt.Errorf("source path %s is a directory, expected a binary file", srcPath)
	}

	destDir = filepath.Clean(destDir)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("create destination: %w", err)
	}

	targetPath := filepath.Join(destDir, binaryName)
	if err := extractRawBinary(srcPath, targetPath); err != nil {
		return "", err
	}
	return targetPath, nil
}

func extractRawBinary(srcPath string, targetPath string) error {
	mode := 0o755 // Binaries should be executable

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("create file %s: %w", targetPath, err)
	}
	defer dst.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source file %s: %w", srcPath, err)
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("write file %s: %w", targetPath, err)
	}

	return nil

}

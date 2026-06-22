package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteShim writes a php.bat shim to ~/.pvm/shims/ that forwards
// all calls to the target php.exe. This file stays on PATH permanently —
// only its contents change when the user runs `pvm use`.
func WriteShim(shimsDir string, targetExePath string) error {
	shimPath := filepath.Join(shimsDir, "php.bat")

	// A .bat shim that forwards all arguments to the real php.exe
	content := fmt.Sprintf("@echo off\n\"%s\" %%*\n", targetExePath)

	if err := os.WriteFile(shimPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("could not write shim: %w", err)
	}

	return nil
}

// ShimPath returns the full path to the php.bat shim file.
func ShimPath(shimsDir string) string {
	return filepath.Join(shimsDir, "php.bat")
}

// ShimExists returns true if the shim file is already in place.
func ShimExists(shimsDir string) bool {
	_, err := os.Stat(ShimPath(shimsDir))
	return err == nil
}

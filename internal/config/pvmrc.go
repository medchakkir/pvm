package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const PvmrcFile = ".pvmrc"

// FindPvmrc walks up from startDir toward the filesystem root, returning the
// path to the first .pvmrc file found. Returns an empty string (not an error)
// when no .pvmrc exists anywhere in the tree.
func FindPvmrc(startDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, PvmrcFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root without finding a .pvmrc
			return "", nil
		}
		dir = parent
	}
}

// ReadPvmrc reads a .pvmrc file and returns its contents trimmed of whitespace.
// Returns an error if the file cannot be read or is empty.
func ReadPvmrc(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %w", path, err)
	}

	version := strings.TrimSpace(string(data))
	if version == "" {
		return "", fmt.Errorf("%s is empty — add a version string like '8.3' or '8.3.7-TS'", path)
	}

	return version, nil
}

// WritePvmrc writes a version string to a .pvmrc file in dir.
// Overwrites any existing .pvmrc in that directory.
func WritePvmrc(dir, version string) error {
	path := filepath.Join(dir, PvmrcFile)
	return os.WriteFile(path, []byte(version+"\n"), 0644)
}

// PvmrcExists reports whether a .pvmrc file exists in dir (not in parents).
func PvmrcExists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, PvmrcFile))
	return err == nil
}

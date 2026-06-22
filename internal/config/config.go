package config

import (
	"os"
	"path/filepath"
)

// PVMHome returns the root ~/.pvm directory path.
func PVMHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".pvm"), nil
}

// VersionsDir returns ~/.pvm/versions
func VersionsDir() (string, error) {
	home, err := PVMHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "versions"), nil
}

// ShimsDir returns ~/.pvm/shims
func ShimsDir() (string, error) {
	home, err := PVMHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "shims"), nil
}

// CurrentFile returns the path to ~/.pvm/current
func CurrentFile() (string, error) {
	home, err := PVMHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "current"), nil
}

// Init creates all required PVM directories if they don't exist.
// This runs on every command before anything else.
func Init() error {
	dirs := []func() (string, error){
		PVMHome,
		VersionsDir,
		ShimsDir,
	}

	for _, dirFn := range dirs {
		path, err := dirFn()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

// GetCurrentVersion reads the active PHP version from ~/.pvm/current.
// Returns an empty string if no version is set.
func GetCurrentVersion() (string, error) {
	currentFile, err := CurrentFile()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(currentFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // no version set yet, not an error
		}
		return "", err
	}

	return string(data), nil
}

// SetCurrentVersion writes the active PHP version to ~/.pvm/current.
func SetCurrentVersion(version string) error {
	currentFile, err := CurrentFile()
	if err != nil {
		return err
	}
	return os.WriteFile(currentFile, []byte(version), 0644)
}

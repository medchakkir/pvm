package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// Shim describes a single launcher batch file that pvm maintains on PATH.
// Name is the base file name (without extension); Target is the absolute path
// to the executable it forwards to.
type Shim struct {
	Name   string
	Target string
}

// DefaultShims returns the set of shims pvm writes for an active version,
// resolved against the given version directory. php is the primary launcher;
// php-cgi backs FastCGI/SAPI usage and lays the groundwork for composer.
func DefaultShims(versionDir string) []Shim {
	return []Shim{
		{Name: "php", Target: filepath.Join(versionDir, "php.exe")},
		{Name: "php-cgi", Target: filepath.Join(versionDir, "php-cgi.exe")},
	}
}

// WriteShim writes a set of .bat shims into ~/.pvm/shims/, each forwarding all
// arguments to its target executable. The shims directory stays on PATH
// permanently — only the shim contents change when the user runs `pvm use`.
//
// A shim whose target executable is missing for the active build is skipped,
// and any stale shim left over from a previous version is removed so callers
// never inherit a dangling launcher.
func WriteShim(shimsDir string, shims []Shim) error {
	for _, s := range shims {
		shimPath := ShimPathFor(shimsDir, s.Name)

		if _, err := os.Stat(s.Target); err != nil {
			if rmErr := os.Remove(shimPath); rmErr != nil && !os.IsNotExist(rmErr) {
				return fmt.Errorf("could not remove stale %s shim: %w", s.Name, rmErr)
			}
			continue
		}

		content := fmt.Sprintf("@echo off\n\"%s\" %%*\n", s.Target)
		if err := os.WriteFile(shimPath, []byte(content), 0755); err != nil {
			return fmt.Errorf("could not write %s shim: %w", s.Name, err)
		}
	}

	return nil
}

// ShimPathFor returns the full path to a named shim's .bat file.
func ShimPathFor(shimsDir, name string) string {
	return filepath.Join(shimsDir, name+".bat")
}

// ShimPath returns the full path to the primary php.bat shim file.
func ShimPath(shimsDir string) string {
	return ShimPathFor(shimsDir, "php")
}

// ShimExists returns true if the primary php shim is already in place.
func ShimExists(shimsDir string) bool {
	_, err := os.Stat(ShimPath(shimsDir))
	return err == nil
}

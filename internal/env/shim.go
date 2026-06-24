package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Shim describes a single launcher batch file that pvm maintains on PATH.
//
//	Name     base file name without extension, e.g. "php" -> php.bat
//	Target   absolute path to the executable the shim invokes
//	Args     fixed arguments inserted before the forwarded args (%*),
//	         e.g. the composer.phar path for the composer launcher
//	Requires files that must all exist for the shim to be written; defaults to
//	         just Target when empty
type Shim struct {
	Name     string
	Target   string
	Args     []string
	Requires []string
}

// DefaultShims returns the set of shims pvm writes for an active version,
// resolved against the given version directory. php is the primary launcher,
// php-cgi backs FastCGI/SAPI usage, and composer (written only when
// composer.phar is present) runs `php.exe composer.phar %*`.
func DefaultShims(versionDir string) []Shim {
	phpExe := filepath.Join(versionDir, "php.exe")
	composerPhar := filepath.Join(versionDir, "composer.phar")

	return []Shim{
		{Name: "php", Target: phpExe},
		{Name: "php-cgi", Target: filepath.Join(versionDir, "php-cgi.exe")},
		{
			Name:     "composer",
			Target:   phpExe,
			Args:     []string{composerPhar},
			Requires: []string{phpExe, composerPhar},
		},
	}
}

// WriteShim writes a set of .bat shims into ~/.pvm/shims/, each forwarding all
// arguments to its target executable. The shims directory stays on PATH
// permanently — only the shim contents change when the user runs `pvm use`.
//
// A shim whose required files are missing for the active version is skipped,
// and any stale shim left over from a previous version is removed so callers
// never inherit a dangling launcher.
func WriteShim(shimsDir string, shims []Shim) error {
	for _, s := range shims {
		shimPath := ShimPathFor(shimsDir, s.Name)

		if !shimRequirementsMet(s) {
			if rmErr := os.Remove(shimPath); rmErr != nil && !os.IsNotExist(rmErr) {
				return fmt.Errorf("could not remove stale %s shim: %w", s.Name, rmErr)
			}
			continue
		}

		if err := os.WriteFile(shimPath, []byte(shimContent(s)), 0755); err != nil {
			return fmt.Errorf("could not write %s shim: %w", s.Name, err)
		}
	}

	return nil
}

// shimRequirementsMet reports whether every file the shim depends on exists.
func shimRequirementsMet(s Shim) bool {
	required := s.Requires
	if len(required) == 0 {
		required = []string{s.Target}
	}
	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return true
}

// shimContent builds the .bat body for a shim, quoting the target and any fixed
// arguments and forwarding the caller's arguments via %*.
func shimContent(s Shim) string {
	// Wrap each path in literal double quotes. We must not use %q here because
	// that escapes backslashes (C:\\path), which is invalid in a .bat file.
	parts := []string{`"` + s.Target + `"`}
	for _, a := range s.Args {
		parts = append(parts, `"`+a+`"`)
	}
	return fmt.Sprintf("@echo off\n%s %%*\n", strings.Join(parts, " "))
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

package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteShimWritesPresentTargetsAndCleansStale(t *testing.T) {
	versionDir := t.TempDir()
	shimsDir := t.TempDir()

	// Only php.exe exists in this build; php-cgi.exe is intentionally absent.
	phpExe := filepath.Join(versionDir, "php.exe")
	if err := os.WriteFile(phpExe, []byte("binary"), 0755); err != nil {
		t.Fatalf("setup php.exe: %v", err)
	}

	// A stale php-cgi.bat from a previous version should be removed.
	staleCgi := ShimPathFor(shimsDir, "php-cgi")
	if err := os.WriteFile(staleCgi, []byte("old"), 0755); err != nil {
		t.Fatalf("setup stale shim: %v", err)
	}

	if err := WriteShim(shimsDir, DefaultShims(versionDir)); err != nil {
		t.Fatalf("WriteShim error: %v", err)
	}

	data, err := os.ReadFile(ShimPath(shimsDir))
	if err != nil {
		t.Fatalf("php.bat not written: %v", err)
	}
	if !strings.Contains(string(data), phpExe) {
		t.Errorf("php.bat should forward to %q, got %q", phpExe, string(data))
	}

	if _, err := os.Stat(staleCgi); !os.IsNotExist(err) {
		t.Errorf("stale php-cgi.bat should have been removed, stat err = %v", err)
	}
}

func TestComposerShimWrittenOnlyWithPhar(t *testing.T) {
	versionDir := t.TempDir()
	shimsDir := t.TempDir()

	phpExe := filepath.Join(versionDir, "php.exe")
	if err := os.WriteFile(phpExe, []byte("binary"), 0755); err != nil {
		t.Fatalf("setup php.exe: %v", err)
	}

	// Without composer.phar present, no composer.bat should be written.
	if err := WriteShim(shimsDir, DefaultShims(versionDir)); err != nil {
		t.Fatalf("WriteShim error: %v", err)
	}
	composerShim := ShimPathFor(shimsDir, "composer")
	if _, err := os.Stat(composerShim); !os.IsNotExist(err) {
		t.Fatalf("composer shim should not exist yet, stat err = %v", err)
	}

	// Add composer.phar and rewrite: composer.bat should appear and invoke php.
	phar := filepath.Join(versionDir, "composer.phar")
	if err := os.WriteFile(phar, []byte("phar"), 0644); err != nil {
		t.Fatalf("setup composer.phar: %v", err)
	}
	if err := WriteShim(shimsDir, DefaultShims(versionDir)); err != nil {
		t.Fatalf("WriteShim error: %v", err)
	}

	data, err := os.ReadFile(composerShim)
	if err != nil {
		t.Fatalf("composer shim not written: %v", err)
	}
	body := string(data)
	if !strings.Contains(body, phpExe) || !strings.Contains(body, phar) {
		t.Errorf("composer shim should run php with the phar, got %q", body)
	}
	// Paths must be raw (not Go-escaped) so the .bat is valid.
	if strings.Contains(body, `\\`) {
		t.Errorf("composer shim must not contain escaped backslashes: %q", body)
	}
}

func TestWriteShimWritesPhpCgiWhenPresent(t *testing.T) {
	versionDir := t.TempDir()
	shimsDir := t.TempDir()

	for _, exe := range []string{"php.exe", "php-cgi.exe"} {
		if err := os.WriteFile(filepath.Join(versionDir, exe), []byte("binary"), 0755); err != nil {
			t.Fatalf("setup %s: %v", exe, err)
		}
	}

	if err := WriteShim(shimsDir, DefaultShims(versionDir)); err != nil {
		t.Fatalf("WriteShim error: %v", err)
	}

	if !ShimExists(shimsDir) {
		t.Error("php shim should exist")
	}
	if _, err := os.Stat(ShimPathFor(shimsDir, "php-cgi")); err != nil {
		t.Errorf("php-cgi shim should exist: %v", err)
	}
}

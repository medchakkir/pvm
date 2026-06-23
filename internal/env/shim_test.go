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

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestHome redirects PVM home to a temp directory for testing.
func setupTestHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)        // Linux/macOS
	t.Setenv("USERPROFILE", tmp) // Windows
	return tmp
}

func TestInit_CreatesDirectories(t *testing.T) {
	tmp := setupTestHome(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	expectedDirs := []string{
		filepath.Join(tmp, ".pvm"),
		filepath.Join(tmp, ".pvm", "versions"),
		filepath.Join(tmp, ".pvm", "shims"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("expected directory to exist: %s", dir)
		}
	}
}

func TestInit_IdempotentOnSecondCall(t *testing.T) {
	setupTestHome(t)

	if err := Init(); err != nil {
		t.Fatalf("first Init() error: %v", err)
	}
	if err := Init(); err != nil {
		t.Fatalf("second Init() should not error: %v", err)
	}
}

func TestSetAndGetCurrentVersion(t *testing.T) {
	setupTestHome(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// No version set yet
	version, err := GetCurrentVersion()
	if err != nil {
		t.Fatalf("GetCurrentVersion() on empty state error: %v", err)
	}
	if version != "" {
		t.Errorf("expected empty version, got %q", version)
	}

	// Set a version
	if err := SetCurrentVersion("8.3.7-TS"); err != nil {
		t.Fatalf("SetCurrentVersion() error: %v", err)
	}

	// Read it back
	version, err = GetCurrentVersion()
	if err != nil {
		t.Fatalf("GetCurrentVersion() error: %v", err)
	}
	if version != "8.3.7-TS" {
		t.Errorf("GetCurrentVersion() = %q, want %q", version, "8.3.7-TS")
	}
}

func TestSetCurrentVersion_Overwrite(t *testing.T) {
	setupTestHome(t)
	Init()

	SetCurrentVersion("8.2.0-TS")
	SetCurrentVersion("8.3.7-NTS")

	version, _ := GetCurrentVersion()
	if version != "8.3.7-NTS" {
		t.Errorf("expected overwritten version %q, got %q", "8.3.7-NTS", version)
	}
}

func TestPVMHome_ReturnsCorrectPath(t *testing.T) {
	tmp := setupTestHome(t)

	home, err := PVMHome()
	if err != nil {
		t.Fatalf("PVMHome() error: %v", err)
	}

	expected := filepath.Join(tmp, ".pvm")
	if home != expected {
		t.Errorf("PVMHome() = %q, want %q", home, expected)
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindPvmrc_InCurrentDir(t *testing.T) {
	tmp := t.TempDir()

	if err := WritePvmrc(tmp, "8.3.7-TS"); err != nil {
		t.Fatalf("WritePvmrc: %v", err)
	}

	got, err := FindPvmrc(tmp)
	if err != nil {
		t.Fatalf("FindPvmrc: %v", err)
	}
	if got != filepath.Join(tmp, PvmrcFile) {
		t.Errorf("FindPvmrc() = %q, want %q", got, filepath.Join(tmp, PvmrcFile))
	}
}

func TestFindPvmrc_InParentDir(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "subdir", "project")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// .pvmrc is in parent, not in child
	if err := WritePvmrc(parent, "8.2.18-NTS"); err != nil {
		t.Fatalf("WritePvmrc: %v", err)
	}

	got, err := FindPvmrc(child)
	if err != nil {
		t.Fatalf("FindPvmrc: %v", err)
	}
	if got != filepath.Join(parent, PvmrcFile) {
		t.Errorf("FindPvmrc() = %q, want %q", got, filepath.Join(parent, PvmrcFile))
	}
}

func TestFindPvmrc_NoneExists(t *testing.T) {
	tmp := t.TempDir()
	child := filepath.Join(tmp, "deep", "nested", "dir")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got, err := FindPvmrc(child)
	if err != nil {
		t.Fatalf("FindPvmrc should not error when no .pvmrc found: %v", err)
	}
	if got != "" {
		t.Errorf("FindPvmrc() = %q, want empty string", got)
	}
}

func TestFindPvmrc_ChildOverridesParent(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "project")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := WritePvmrc(parent, "8.1.0-TS"); err != nil {
		t.Fatalf("WritePvmrc parent: %v", err)
	}
	if err := WritePvmrc(child, "8.3.7-TS"); err != nil {
		t.Fatalf("WritePvmrc child: %v", err)
	}

	// Starting from child — should find child's .pvmrc, not parent's
	got, err := FindPvmrc(child)
	if err != nil {
		t.Fatalf("FindPvmrc: %v", err)
	}
	if got != filepath.Join(child, PvmrcFile) {
		t.Errorf("FindPvmrc() = %q, want child .pvmrc %q", got, filepath.Join(child, PvmrcFile))
	}
}

func TestReadPvmrc_ReturnsVersion(t *testing.T) {
	tmp := t.TempDir()
	if err := WritePvmrc(tmp, "8.3.7-TS"); err != nil {
		t.Fatalf("WritePvmrc: %v", err)
	}

	got, err := ReadPvmrc(filepath.Join(tmp, PvmrcFile))
	if err != nil {
		t.Fatalf("ReadPvmrc: %v", err)
	}
	if got != "8.3.7-TS" {
		t.Errorf("ReadPvmrc() = %q, want %q", got, "8.3.7-TS")
	}
}

func TestReadPvmrc_TrimsWhitespace(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, PvmrcFile)
	if err := os.WriteFile(path, []byte("  8.2\r\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ReadPvmrc(path)
	if err != nil {
		t.Fatalf("ReadPvmrc: %v", err)
	}
	if got != "8.2" {
		t.Errorf("ReadPvmrc() = %q, want %q", got, "8.2")
	}
}

func TestReadPvmrc_EmptyFileErrors(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, PvmrcFile)
	if err := os.WriteFile(path, []byte("   \n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := ReadPvmrc(path)
	if err == nil {
		t.Error("ReadPvmrc should error on empty file")
	}
}

func TestWritePvmrc_Overwrites(t *testing.T) {
	tmp := t.TempDir()

	if err := WritePvmrc(tmp, "8.2.0-TS"); err != nil {
		t.Fatalf("first WritePvmrc: %v", err)
	}
	if err := WritePvmrc(tmp, "8.3.7-NTS"); err != nil {
		t.Fatalf("second WritePvmrc: %v", err)
	}

	got, err := ReadPvmrc(filepath.Join(tmp, PvmrcFile))
	if err != nil {
		t.Fatalf("ReadPvmrc: %v", err)
	}
	if got != "8.3.7-NTS" {
		t.Errorf("WritePvmrc overwrite: got %q, want %q", got, "8.3.7-NTS")
	}
}

func TestPvmrcExists(t *testing.T) {
	tmp := t.TempDir()

	if PvmrcExists(tmp) {
		t.Error("PvmrcExists should be false before writing")
	}

	if err := WritePvmrc(tmp, "8.3"); err != nil {
		t.Fatalf("WritePvmrc: %v", err)
	}

	if !PvmrcExists(tmp) {
		t.Error("PvmrcExists should be true after writing")
	}
}

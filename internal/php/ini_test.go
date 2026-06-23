package php

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeExtName(t *testing.T) {
	cases := map[string]string{
		"curl":                  "curl",
		"php_curl.dll":          "curl",
		"PHP_CURL.DLL":          "curl",
		`"php_curl.dll"`:        "curl",
		`  php_gd.dll  `:        "gd",
		`C:\php\ext\php_pdo.dll`: "pdo",
		"opcache":              "opcache",
		"php_opcache.dll":      "opcache",
	}
	for in, want := range cases {
		if got := NormalizeExtName(in); got != want {
			t.Errorf("NormalizeExtName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEntriesParsing(t *testing.T) {
	content := "" +
		"; This is a plain comment about extensions.\r\n" +
		"extension_dir = \"ext\"\r\n" +
		"extension=curl\r\n" +
		";extension=php_gd.dll\r\n" +
		"; extension = \"php_mbstring.dll\"\r\n" +
		"zend_extension=opcache\r\n"

	ini := writeTempIni(t, content)
	entries := ini.Entries()

	if len(entries) != 4 {
		t.Fatalf("expected 4 extension entries, got %d: %+v", len(entries), entries)
	}

	byName := map[string]ExtEntry{}
	for _, e := range entries {
		byName[e.Name] = e
	}

	if e, ok := byName["curl"]; !ok || !e.Enabled {
		t.Errorf("curl should be present and enabled, got %+v (ok=%v)", e, ok)
	}
	if e, ok := byName["gd"]; !ok || e.Enabled {
		t.Errorf("gd should be present and disabled, got %+v (ok=%v)", e, ok)
	}
	if e, ok := byName["mbstring"]; !ok || e.Enabled {
		t.Errorf("mbstring should be present and disabled, got %+v (ok=%v)", e, ok)
	}
	if e, ok := byName["opcache"]; !ok || e.Directive != "zend_extension" {
		t.Errorf("opcache should be a zend_extension, got %+v (ok=%v)", e, ok)
	}
}

func TestSetEnabledPreservesCRLF(t *testing.T) {
	content := "extension=curl\r\n;extension=php_gd.dll\r\n"
	ini := writeTempIni(t, content)

	for _, e := range ini.Entries() {
		switch e.Name {
		case "curl":
			ini.SetEnabled(e.LineIndex, false)
		case "gd":
			ini.SetEnabled(e.LineIndex, true)
		}
	}

	if err := ini.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	data, err := os.ReadFile(ini.Path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	want := ";extension=curl\r\nextension=php_gd.dll\r\n"
	if string(data) != want {
		t.Errorf("after toggle:\n got  %q\n want %q", string(data), want)
	}
}

func writeTempIni(t *testing.T, content string) *IniFile {
	t.Helper()
	path := filepath.Join(t.TempDir(), "php.ini")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("could not write temp ini: %v", err)
	}
	ini, err := LoadIni(path)
	if err != nil {
		t.Fatalf("LoadIni error: %v", err)
	}
	return ini
}

package php

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// extDirectiveRe matches `extension=` / `zend_extension=` lines, whether active
// or commented out (disabled) with one or more leading `;`. Plain comments that
// merely mention the word "extension" do not match because the directive must be
// immediately followed by `=`.
//
//	group 1: leading whitespace
//	group 2: optional disable prefix (`;` + spacing) — empty when enabled
//	group 3: directive name (`extension` or `zend_extension`)
//	group 4: right-hand value (e.g. `php_curl.dll` or `curl`)
var extDirectiveRe = regexp.MustCompile(`^([ \t]*)(;+[ \t]*)?(zend_extension|extension)[ \t]*=[ \t]*(.+?)[ \t]*$`)

// ExtEntry is a parsed extension directive line from a php.ini file.
type ExtEntry struct {
	Name      string // normalized extension name, e.g. "curl", "opcache"
	Directive string // "extension" or "zend_extension"
	Value     string // original right-hand value, e.g. "php_curl.dll"
	Enabled   bool   // false when the line is commented out
	LineIndex int    // index into the underlying line slice
}

// IniFile is an in-memory, editable view of a php.ini file that preserves the
// original lines and line-ending style so rewrites stay minimal.
type IniFile struct {
	Path  string
	lines []string
	crlf  bool
}

// LoadIni reads a php.ini file from disk into an editable IniFile.
func LoadIni(path string) (*IniFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	crlf := strings.Contains(content, "\r\n")

	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	normalized = strings.TrimSuffix(normalized, "\n")

	var lines []string
	if normalized != "" {
		lines = strings.Split(normalized, "\n")
	}

	return &IniFile{Path: path, lines: lines, crlf: crlf}, nil
}

// Entries returns every extension / zend_extension directive found in the file.
func (f *IniFile) Entries() []ExtEntry {
	var entries []ExtEntry
	for i, line := range f.lines {
		m := extDirectiveRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		entries = append(entries, ExtEntry{
			Name:      NormalizeExtName(m[4]),
			Directive: m[3],
			Value:     m[4],
			Enabled:   m[2] == "",
			LineIndex: i,
		})
	}
	return entries
}

// SetEnabled toggles the leading `;` on the directive at lineIndex. The rest of
// the line's formatting is normalized to the canonical `directive=value` style.
func (f *IniFile) SetEnabled(lineIndex int, enabled bool) {
	if lineIndex < 0 || lineIndex >= len(f.lines) {
		return
	}
	m := extDirectiveRe.FindStringSubmatch(f.lines[lineIndex])
	if m == nil {
		return
	}
	rebuilt := m[1] + m[3] + "=" + m[4]
	if !enabled {
		rebuilt = m[1] + ";" + m[3] + "=" + m[4]
	}
	f.lines[lineIndex] = rebuilt
}

// AddExtension appends a new enabled directive line to the file.
func (f *IniFile) AddExtension(directive, value string) {
	f.lines = append(f.lines, directive+"="+value)
}

// Save writes the file back to disk, preserving its original line-ending style.
func (f *IniFile) Save() error {
	ending := "\n"
	if f.crlf {
		ending = "\r\n"
	}
	content := strings.Join(f.lines, ending)
	if len(f.lines) > 0 {
		content += ending
	}
	return os.WriteFile(f.Path, []byte(content), 0644)
}

// NormalizeExtName reduces any extension reference — a directive value or an
// `ext/` filename — to a bare, comparable name. It strips surrounding quotes and
// directories, lowercases, drops a `.dll` suffix, and drops a `php_` prefix so
// that `"C:\php\ext\php_CURL.dll"`, `php_curl.dll`, and `curl` all collapse to
// `curl`.
func NormalizeExtName(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.Trim(s, `"'`)
	s = strings.TrimSpace(s)
	s = filepath.Base(s)
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, ".dll")
	s = strings.TrimPrefix(s, "php_")
	return s
}

// AvailableExtensions scans an `ext` directory and returns a map of normalized
// extension name to the on-disk DLL filename.
func AvailableExtensions(extDir string) (map[string]string, error) {
	entries, err := os.ReadDir(extDir)
	if err != nil {
		return nil, err
	}

	avail := make(map[string]string)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.EqualFold(filepath.Ext(name), ".dll") {
			continue
		}
		avail[NormalizeExtName(name)] = name
	}
	return avail, nil
}

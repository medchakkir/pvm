package php

import (
	"testing"
)

func TestFilterPHPZips(t *testing.T) {
	links := []string{
		"php-8.3.7-Win32-vs16-x64.zip",        // TS — keep
		"php-8.3.7-nts-Win32-vs16-x64.zip",    // NTS — keep
		"php-8.3.7-Win32-vs16-x86.zip",        // x86 — skip
		"php-8.3.7-nts-Win32-vs16-x86.zip",    // x86 NTS — skip
		"php-8.3.7-debug-pack-vs16-x64.zip",   // debug — skip
		"php-8.2.18-Win32-vs16-x64.zip",       // TS — keep
		"php-8.2.18-nts-Win32-vs16-x64.zip",   // NTS — keep
		"not-a-php-file.zip",                   // irrelevant — skip
		"php-8.3.7-Win32-vs16-x64.zip",        // duplicate — skip
	}

	results := filterPHPZips(links, "https://windows.php.net/downloads/releases/")

	// Expect: 8.3.7-TS, 8.3.7-NTS, 8.2.18-TS, 8.2.18-NTS = 4 results
	if len(results) != 4 {
		t.Errorf("expected 4 results, got %d", len(results))
		for _, r := range results {
			t.Logf("  → %s (TS=%v)", r.Version, r.ThreadSafe)
		}
	}

	// Verify TS and NTS are both present for 8.3.7
	var foundTS, foundNTS bool
	for _, r := range results {
		if r.Version.String() == "8.3.7" {
			if r.ThreadSafe {
				foundTS = true
			} else {
				foundNTS = true
			}
		}
	}
	if !foundTS {
		t.Error("expected 8.3.7 TS build to be included")
	}
	if !foundNTS {
		t.Error("expected 8.3.7 NTS build to be included")
	}
}

func TestFindRemoteVersion_ExactMatch(t *testing.T) {
	available := []RemoteVersion{
		{Version: PHPVersion{8, 3, 7}, ThreadSafe: true, ZipName: "php-8.3.7-Win32-vs16-x64.zip"},
		{Version: PHPVersion{8, 3, 7}, ThreadSafe: false, ZipName: "php-8.3.7-nts-Win32-vs16-x64.zip"},
		{Version: PHPVersion{8, 2, 18}, ThreadSafe: true, ZipName: "php-8.2.18-Win32-vs16-x64.zip"},
	}

	// TS exact match
	result, err := FindRemoteVersion("8.3.7", available, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version.String() != "8.3.7" || !result.ThreadSafe {
		t.Errorf("expected 8.3.7 TS, got %s TS=%v", result.Version, result.ThreadSafe)
	}

	// NTS exact match
	result, err = FindRemoteVersion("8.3.7", available, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version.String() != "8.3.7" || result.ThreadSafe {
		t.Errorf("expected 8.3.7 NTS, got %s TS=%v", result.Version, result.ThreadSafe)
	}
}

func TestFindRemoteVersion_MinorResolvesToLatestPatch(t *testing.T) {
	available := []RemoteVersion{
		{Version: PHPVersion{8, 3, 5}, ThreadSafe: true},
		{Version: PHPVersion{8, 3, 7}, ThreadSafe: true},
		{Version: PHPVersion{8, 3, 6}, ThreadSafe: true},
	}

	result, err := FindRemoteVersion("8.3", available, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version.String() != "8.3.7" {
		t.Errorf("expected latest patch 8.3.7, got %s", result.Version)
	}
}

func TestFindRemoteVersion_NotFound(t *testing.T) {
	available := []RemoteVersion{
		{Version: PHPVersion{8, 3, 7}, ThreadSafe: true},
	}

	_, err := FindRemoteVersion("7.4.33", available, true)
	if err == nil {
		t.Error("expected error for missing version, got nil")
	}
}

func TestRemoteVersion_TypeLabel(t *testing.T) {
	ts := RemoteVersion{ThreadSafe: true}
	nts := RemoteVersion{ThreadSafe: false}

	if ts.TypeLabel() != "TS" {
		t.Errorf("expected TS, got %s", ts.TypeLabel())
	}
	if nts.TypeLabel() != "NTS" {
		t.Errorf("expected NTS, got %s", nts.TypeLabel())
	}
}

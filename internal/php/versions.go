package php

import (
	"fmt"
	"strconv"
	"strings"
)

// PHPVersion represents a parsed semantic PHP version (e.g. 8.3.7)
type PHPVersion struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string into a PHPVersion struct.
// Accepts "8.3", "8.3.7" — patch defaults to 0 if omitted.
func ParseVersion(raw string) (PHPVersion, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return PHPVersion{}, fmt.Errorf("version string is empty")
	}

	parts := strings.Split(raw, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return PHPVersion{}, fmt.Errorf("invalid version format %q — expected Major.Minor or Major.Minor.Patch", raw)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return PHPVersion{}, fmt.Errorf("invalid major version in %q", raw)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return PHPVersion{}, fmt.Errorf("invalid minor version in %q", raw)
	}

	patch := 0
	if len(parts) == 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil || patch < 0 {
			return PHPVersion{}, fmt.Errorf("invalid patch version in %q", raw)
		}
	}

	return PHPVersion{Major: major, Minor: minor, Patch: patch}, nil
}

// String returns the full version string e.g. "8.3.7"
func (v PHPVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// MinorString returns "Major.Minor" e.g. "8.3"
func (v PHPVersion) MinorString() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// IsValid checks that the version looks like a real PHP version.
func (v PHPVersion) IsValid() bool {
	return v.Major >= 5 && v.Major <= 9 && v.Minor >= 0 && v.Patch >= 0
}

// Compare returns:
//
//	-1 if v < other
//	 0 if v == other
//	+1 if v > other
func (v PHPVersion) Compare(other PHPVersion) int {
	if v.Major != other.Major {
		return compareInts(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return compareInts(v.Minor, other.Minor)
	}
	return compareInts(v.Patch, other.Patch)
}

// MatchesMinor returns true if this version shares Major.Minor with the other.
// Useful for resolving "8.3" → "8.3.7"
func (v PHPVersion) MatchesMinor(other PHPVersion) bool {
	return v.Major == other.Major && v.Minor == other.Minor
}

func compareInts(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
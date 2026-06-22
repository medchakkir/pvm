package php

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input     string
		wantMajor int
		wantMinor int
		wantPatch int
		wantErr   bool
	}{
		{"8.3.7", 8, 3, 7, false},
		{"8.3", 8, 3, 0, false},
		{"8.1.29", 8, 1, 29, false},
		{"7.4.33", 7, 4, 33, false},
		// edge cases
		{"", 0, 0, 0, true},
		{"invalid", 0, 0, 0, true},
		{"8", 0, 0, 0, true},
		{"8.3.7.1", 0, 0, 0, true},
		{"abc.def.ghi", 0, 0, 0, true},
		{" 8.3.7 ", 8, 3, 7, false}, // trims whitespace
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseVersion(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseVersion(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseVersion(%q) unexpected error: %v", tt.input, err)
				return
			}

			if got.Major != tt.wantMajor || got.Minor != tt.wantMinor || got.Patch != tt.wantPatch {
				t.Errorf("ParseVersion(%q) = %v, want %d.%d.%d",
					tt.input, got, tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		})
	}
}

func TestString(t *testing.T) {
	v := PHPVersion{Major: 8, Minor: 3, Patch: 7}
	if got := v.String(); got != "8.3.7" {
		t.Errorf("String() = %q, want %q", got, "8.3.7")
	}
}

func TestMinorString(t *testing.T) {
	v := PHPVersion{Major: 8, Minor: 3, Patch: 7}
	if got := v.MinorString(); got != "8.3" {
		t.Errorf("MinorString() = %q, want %q", got, "8.3")
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		version PHPVersion
		want    bool
	}{
		{PHPVersion{8, 3, 7}, true},
		{PHPVersion{7, 4, 33}, true},
		{PHPVersion{5, 6, 0}, true},
		{PHPVersion{4, 0, 0}, false}, // too old
		{PHPVersion{10, 0, 0}, false}, // too new / unexpected
	}

	for _, tt := range tests {
		t.Run(tt.version.String(), func(t *testing.T) {
			if got := tt.version.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		a, b PHPVersion
		want int
	}{
		{PHPVersion{8, 3, 7}, PHPVersion{8, 3, 7}, 0},
		{PHPVersion{8, 3, 7}, PHPVersion{8, 3, 6}, 1},
		{PHPVersion{8, 3, 6}, PHPVersion{8, 3, 7}, -1},
		{PHPVersion{8, 3, 0}, PHPVersion{8, 2, 99}, 1},
		{PHPVersion{7, 4, 0}, PHPVersion{8, 0, 0}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+"_vs_"+tt.b.String(), func(t *testing.T) {
			if got := tt.a.Compare(tt.b); got != tt.want {
				t.Errorf("Compare() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMatchesMinor(t *testing.T) {
	v := PHPVersion{8, 3, 7}
	if !v.MatchesMinor(PHPVersion{8, 3, 0}) {
		t.Error("expected 8.3.7 to match minor 8.3.0")
	}
	if v.MatchesMinor(PHPVersion{8, 2, 7}) {
		t.Error("expected 8.3.7 NOT to match minor 8.2.7")
	}
}
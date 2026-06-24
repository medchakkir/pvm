package php

import "testing"

func TestComposerURLForPHP(t *testing.T) {
	cases := []struct {
		version PHPVersion
		wantLTS bool
	}{
		{PHPVersion{5, 6, 40}, true},
		{PHPVersion{7, 1, 33}, true},
		{PHPVersion{7, 2, 0}, false},
		{PHPVersion{7, 4, 33}, false},
		{PHPVersion{8, 3, 7}, false},
	}

	for _, c := range cases {
		got := ComposerURLForPHP(c.version)
		isLTS := got == composerLTS22URL
		if isLTS != c.wantLTS {
			t.Errorf("ComposerURLForPHP(%s) = %q, wantLTS=%v", c.version, got, c.wantLTS)
		}
	}
}

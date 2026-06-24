package php

import (
	"fmt"
	"os"
	"path/filepath"
)

// ComposerPharName is the file Composer is installed as inside a version dir.
const ComposerPharName = "composer.phar"

const (
	// composerLatestStableURL serves the newest Composer release.
	composerLatestStableURL = "https://getcomposer.org/download/latest-stable/composer.phar"
	// composerLTS22URL serves the Composer 2.2 LTS line, which still supports
	// older PHP runtimes (PHP < 7.2).
	composerLTS22URL = "https://getcomposer.org/download/latest-2.2.x/composer.phar"
)

// ComposerURLForPHP returns the composer.phar download URL appropriate for the
// given PHP version. PHP older than 7.2 gets the 2.2 LTS line; everything else
// gets the latest stable Composer.
func ComposerURLForPHP(v PHPVersion) string {
	if v.Major < 7 || (v.Major == 7 && v.Minor < 2) {
		return composerLTS22URL
	}
	return composerLatestStableURL
}

// ComposerPharPath returns the composer.phar path inside a version directory.
func ComposerPharPath(versionDir string) string {
	return filepath.Join(versionDir, ComposerPharName)
}

// ComposerInstalled reports whether composer.phar already exists for a version.
func ComposerInstalled(versionDir string) bool {
	_, err := os.Stat(ComposerPharPath(versionDir))
	return err == nil
}

// InstallComposer downloads the appropriate composer.phar for the given PHP
// version into versionDir.
func InstallComposer(versionDir string, v PHPVersion) error {
	url := ComposerURLForPHP(v)
	dest := ComposerPharPath(versionDir)
	if err := DownloadFile(url, dest, "Downloading Composer"); err != nil {
		return fmt.Errorf("could not download Composer: %w", err)
	}
	return nil
}

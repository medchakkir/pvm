package php

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

// DownloadZip downloads a ZIP from url into destPath, showing a progress bar.
func DownloadZip(url string, destPath string) error {
	client := &http.Client{Timeout: 5 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: server returned %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(resp.ContentLength, "Downloading")

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("download interrupted: %w", err)
	}

	return nil
}

// ExtractZip extracts all contents of zipPath into destDir.
func ExtractZip(zipPath string, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("could not open zip: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("could not create version directory: %w", err)
	}

	cleanDest := filepath.Clean(destDir)

	for _, f := range r.File {
		destPath := filepath.Join(destDir, f.Name)

		if destPath != cleanDest && !strings.HasPrefix(destPath, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0755)
			continue
		}

		if err := extractFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts a single file from a zip archive.
func extractFile(f *zip.File, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("could not open zip entry %s: %w", f.Name, err)
	}
	defer rc.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", destPath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

// VerifyInstall checks that php.exe exists inside the installed version directory.
func VerifyInstall(versionDir string) error {
	phpExe := filepath.Join(versionDir, "php.exe")
	if _, err := os.Stat(phpExe); os.IsNotExist(err) {
		return fmt.Errorf("php.exe not found in %s — installation may be corrupt", versionDir)
	}
	return nil
}

package php

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	phpWindowsDownloadsURL = "https://windows.php.net/downloads/releases/"
	phpWindowsArchiveURL   = "https://windows.php.net/downloads/releases/archives/"
)

// RemoteVersion represents a PHP version available for download.
type RemoteVersion struct {
	Version     PHPVersion
	DownloadURL string
	ZipName     string
}

// String returns a display-friendly label.
func (r RemoteVersion) String() string {
	return r.Version.String()
}

// FetchRemoteVersions fetches available PHP Windows builds from php.net.
// It returns only stable, thread-safe (TS) x64 ZIP builds.
func FetchRemoteVersions() ([]RemoteVersion, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(phpWindowsDownloadsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach php.net: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("php.net returned status %d", resp.StatusCode)
	}

	links, err := extractLinks(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse php.net page: %w", err)
	}

	versions := filterPHPZips(links, phpWindowsDownloadsURL)

	// Sort descending (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version.Compare(versions[j].Version) > 0
	})

	return versions, nil
}

// extractLinks walks the HTML response and collects all href values.
func extractLinks(resp *http.Response) ([]string, error) {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return links, nil
}

// phpZipPattern matches filenames like: php-8.3.7-Win32-vs16-x64.zip
// We want: TS (no "-nts-"), x64, .zip
var phpZipPattern = regexp.MustCompile(
	`^php-(\d+\.\d+\.\d+)-Win32-vs\d+-x64\.zip$`,
)

// filterPHPZips filters raw href links down to valid TS x64 ZIP builds.
func filterPHPZips(links []string, baseURL string) []RemoteVersion {
	seen := map[string]bool{}
	var versions []RemoteVersion

	for _, link := range links {
		// href may be relative or absolute — extract just the filename
		filename := link
		if idx := strings.LastIndex(link, "/"); idx >= 0 {
			filename = link[idx+1:]
		}

		// Skip non-TS (NTS) builds and debug packs
		if strings.Contains(filename, "-nts-") || strings.Contains(filename, "-debug-") {
			continue
		}

		matches := phpZipPattern.FindStringSubmatch(filename)
		if matches == nil {
			continue
		}

		rawVersion := matches[1]
		if seen[rawVersion] {
			continue
		}
		seen[rawVersion] = true

		version, err := ParseVersion(rawVersion)
		if err != nil || !version.IsValid() {
			continue
		}

		// Build the full download URL
		downloadURL := baseURL + filename
		if strings.HasPrefix(link, "http") {
			downloadURL = link
		}

		versions = append(versions, RemoteVersion{
			Version:     version,
			DownloadURL: downloadURL,
			ZipName:     filename,
		})
	}

	return versions
}

// FindRemoteVersion finds the best matching remote version for a user input.
// "8.3" → returns the latest 8.3.x available.
// "8.3.7" → returns exactly 8.3.7 if available.
func FindRemoteVersion(input string, available []RemoteVersion) (RemoteVersion, error) {
	parsed, err := ParseVersion(input)
	if err != nil {
		return RemoteVersion{}, err
	}

	// Try exact match first
	for _, v := range available {
		if v.Version.Compare(parsed) == 0 {
			return v, nil
		}
	}

	// If patch was omitted (defaults to 0), find latest patch for that Major.Minor
	if parsed.Patch == 0 {
		var best *RemoteVersion
		for i, v := range available {
			if v.Version.MatchesMinor(parsed) {
				if best == nil || v.Version.Compare(best.Version) > 0 {
					best = &available[i]
				}
			}
		}
		if best != nil {
			return *best, nil
		}
	}

	return RemoteVersion{}, fmt.Errorf("PHP %s not found — run `pvm list-remote` to see available versions", input)
}
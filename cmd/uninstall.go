package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Remove an installed PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		requested, err := php.ParseVersion(input)
		if err != nil {
			return fmt.Errorf("✗ invalid version %q: %w", input, err)
		}

		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		// Find all installed directories that match (covers both TS and NTS)
		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("✗ could not read versions directory: %w", err)
		}

		var matches []string
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			// Directory names are like "8.3.7-TS" or "8.3.7-NTS"
			// Strip the type suffix before parsing
			name := entry.Name()
			versionPart := strings.TrimSuffix(strings.TrimSuffix(name, "-TS"), "-NTS")

			v, err := php.ParseVersion(versionPart)
			if err != nil {
				continue
			}
			if v.Compare(requested) == 0 || (requested.Patch == 0 && v.MatchesMinor(requested)) {
				matches = append(matches, name)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("✗ PHP %s is not installed", input)
		}

		// Block uninstall if any match is the active version
		current, err := config.GetCurrentVersion()
		if err != nil {
			return err
		}
		for _, match := range matches {
			if match == current {
				return fmt.Errorf(
					"✗ PHP %s is currently active — run `pvm use <another version>` before uninstalling",
					match,
				)
			}
		}

		// Show what will be removed and confirm
		fmt.Println("The following will be removed:")
		for _, match := range matches {
			fmt.Printf("  ~/.pvm/versions/%s\n", match)
		}

		fmt.Print("\nContinue? [y/N] ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if answer != "y" {
			fmt.Println("Aborted.")
			return nil
		}

		// Delete each matched directory
		for _, match := range matches {
			dir := filepath.Join(versionsDir, match)
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("✗ failed to remove %s: %w", match, err)
			}
			fmt.Printf("✓ Removed PHP %s\n", match)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

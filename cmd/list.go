package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show installed PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		// Read all entries inside ~/.pvm/versions/
		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("could not read versions directory: %w", err)
		}

		// Filter to valid version directories only
		var versions []php.PHPVersion
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			v, err := php.ParseVersion(entry.Name())
			if err != nil || !v.IsValid() {
				continue
			}
			// Confirm php.exe actually exists inside
			phpExe := filepath.Join(versionsDir, entry.Name(), "php.exe")
			if _, err := os.Stat(phpExe); os.IsNotExist(err) {
				continue
			}
			versions = append(versions, v)
		}

		if len(versions) == 0 {
			fmt.Println("No PHP versions installed.")
			fmt.Println("Run `pvm list-remote` to see available versions.")
			return nil
		}

		// Sort descending (newest first)
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].Compare(versions[j]) > 0
		})

		// Get active version for highlighting
		current, err := config.GetCurrentVersion()
		if err != nil {
			current = ""
		}

		fmt.Println("Installed PHP versions:")
		for _, v := range versions {
			if v.String() == current {
				fmt.Printf("  → %s  (active)\n", v.String())
			} else {
				fmt.Printf("    %s\n", v.String())
			}
		}

		fmt.Printf("\n%d version(s) installed.\n", len(versions))
		if current == "" {
			fmt.Println("No active version set. Run `pvm use <version>` to activate one.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

// installedVersion holds a parsed version + its directory name (includes type suffix)
type installedVersion struct {
	version  php.PHPVersion
	dirName  string // e.g. "8.3.7-TS"
	typeLabel string // "TS" or "NTS"
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show installed PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		versionsDir, err := config.VersionsDir()
		if err != nil {
			return fmt.Errorf("could not locate versions directory: %w", err)
		}

		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("could not read versions directory: %w", err)
		}

		var installed []installedVersion
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name() // e.g. "8.3.7-TS" or "8.3.7-NTS"

			// Determine type label and strip it before parsing
			typeLabel := "TS"
			versionPart := name
			if strings.HasSuffix(name, "-NTS") {
				typeLabel = "NTS"
				versionPart = strings.TrimSuffix(name, "-NTS")
			} else if strings.HasSuffix(name, "-TS") {
				versionPart = strings.TrimSuffix(name, "-TS")
			}

			v, err := php.ParseVersion(versionPart)
			if err != nil || !v.IsValid() {
				continue
			}

			// Confirm php.exe is actually present
			phpExe := filepath.Join(versionsDir, name, "php.exe")
			if _, err := os.Stat(phpExe); os.IsNotExist(err) {
				if Verbose() {
					fmt.Fprintf(os.Stderr, "[debug] skipping %s — php.exe not found\n", name)
				}
				continue
			}

			installed = append(installed, installedVersion{
				version:   v,
				dirName:   name,
				typeLabel: typeLabel,
			})
		}

		if len(installed) == 0 {
			ui.Info("No PHP versions installed.")
			ui.Detail("Run `pvm list-remote` to see available versions.")
			return nil
		}

		// Sort: newest first, TS before NTS within same version
		sort.Slice(installed, func(i, j int) bool {
			cmp := installed[i].version.Compare(installed[j].version)
			if cmp != 0 {
				return cmp > 0
			}
			return installed[i].typeLabel == "TS" && installed[j].typeLabel != "TS"
		})

		current, _ := config.GetCurrentVersion()

		ui.Title("Installed PHP versions:")
		for _, iv := range installed {
			marker := "   "
			active := ""
			if iv.dirName == current {
				marker = " →"
				active = "  (active)"
			}
			ui.Info("%s  %-10s %-5s%s", marker, iv.version.String(), iv.typeLabel, active)
		}

		ui.Info("\n%d version(s) installed.", len(installed))
		if current == "" {
			ui.Detail("No active version set. Run `pvm use <version>` to activate one.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

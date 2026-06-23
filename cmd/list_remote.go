package cmd

import (
	"fmt"

	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var (
	limitFlag  int
	ntsFlag    bool
	tsFlag     bool
)

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "Show available PHP versions from php.net",
	Long:  `Fetches and displays PHP versions available for download (Windows x64 builds, both TS and NTS).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Fetching available PHP versions from php.net...")

		versions, err := php.FetchRemoteVersions()
		if err != nil {
			return err
		}

		if len(versions) == 0 {
			ui.Warning("No versions found. Check your internet connection.")
			return nil
		}

		// Apply --ts / --nts filter flags
		filtered := versions
		if tsFlag && !ntsFlag {
			filtered = filterByType(versions, true)
		} else if ntsFlag && !tsFlag {
			filtered = filterByType(versions, false)
		}

		// Apply --limit flag
		display := filtered
		if limitFlag > 0 && limitFlag < len(filtered) {
			display = filtered[:limitFlag]
		}

		ui.Title("\n%-12s %-6s %s", "VERSION", "TYPE", "FILENAME")
		ui.Info("----------------------------------------------------------")

		for _, v := range display {
			ui.Info("%-12s %-6s %s",
				v.Version.String(),
				v.TypeLabel(),
				v.ZipName,
			)
		}

		summary := fmt.Sprintf("\n%d build(s) shown", len(display))
		if len(display) < len(filtered) {
			summary += fmt.Sprintf(" (of %d — use --limit 0 to show all)", len(filtered))
		}
		ui.Info("%s", summary)
		ui.Detail("\nRun `pvm install <version>` to install a TS build.")
		ui.Detail("Run `pvm install --nts <version>` to install an NTS build.")

		return nil
	},
}

func filterByType(versions []php.RemoteVersion, threadSafe bool) []php.RemoteVersion {
	var result []php.RemoteVersion
	for _, v := range versions {
		if v.ThreadSafe == threadSafe {
			result = append(result, v)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(listRemoteCmd)
	listRemoteCmd.Flags().IntVar(&limitFlag, "limit", 20, "Maximum versions to display (0 = all)")
	listRemoteCmd.Flags().BoolVar(&tsFlag, "ts", false, "Show only Thread Safe builds")
	listRemoteCmd.Flags().BoolVar(&ntsFlag, "nts", false, "Show only Non-Thread Safe builds")
}

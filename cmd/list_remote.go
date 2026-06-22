package cmd

import (
	"fmt"
	"os"

	"github.com/medchakkir/pvm/internal/php"
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
		fmt.Println("Fetching available PHP versions from php.net...")

		versions, err := php.FetchRemoteVersions()
		if err != nil {
			return fmt.Errorf("✗ %w", err)
		}

		if len(versions) == 0 {
			fmt.Println("No versions found. Check your internet connection.")
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

		fmt.Printf("\n%-12s %-6s %s\n", "VERSION", "TYPE", "FILENAME")
		fmt.Println("----------------------------------------------------------")

		for _, v := range display {
			fmt.Fprintf(os.Stdout, "%-12s %-6s %s\n",
				v.Version.String(),
				v.TypeLabel(),
				v.ZipName,
			)
		}

		fmt.Printf("\n%d build(s) shown", len(display))
		if len(display) < len(filtered) {
			fmt.Printf(" (of %d — use --limit 0 to show all)", len(filtered))
		}
		fmt.Println()
		fmt.Println("\nRun `pvm install <version>` to install a TS build.")
		fmt.Println("Run `pvm install --nts <version>` to install an NTS build.")

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

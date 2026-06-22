package cmd

import (
	"fmt"
	"os"

	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var limitFlag int

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "Show available PHP versions from php.net",
	Long:  `Fetches and displays PHP versions available for download (Windows x64 TS builds).`,
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

		// Apply --limit flag
		display := versions
		if limitFlag > 0 && limitFlag < len(versions) {
			display = versions[:limitFlag]
		}

		fmt.Printf("\n%-12s %s\n", "VERSION", "DOWNLOAD")
		fmt.Println("-------------------------------------------------------")

		for _, v := range display {
			fmt.Fprintf(os.Stdout, "%-12s %s\n", v.Version.String(), v.ZipName)
		}

		fmt.Printf("\n%d version(s) shown", len(display))
		if len(display) < len(versions) {
			fmt.Printf(" (of %d available — use --limit 0 to show all)", len(versions))
		}
		fmt.Println()
		fmt.Println("\nRun `pvm install <version>` to install one.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listRemoteCmd)
	listRemoteCmd.Flags().IntVar(&limitFlag, "limit", 20, "Maximum versions to display (0 = all)")
}
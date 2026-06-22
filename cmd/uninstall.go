package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove an installed PHP version",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("(uninstall %s — coming in Phase 7)\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

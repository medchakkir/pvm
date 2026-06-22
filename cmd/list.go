package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show installed PHP versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("(list — coming in Phase 5)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

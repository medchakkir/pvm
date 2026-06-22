package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the active PHP version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("(current — coming in Phase 7)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

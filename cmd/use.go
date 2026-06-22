package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Switch to a different PHP version",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("(use %s — coming in Phase 6)\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}

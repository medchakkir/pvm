package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the PVM version and build information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("pvm %s\n", versionInfo)
		fmt.Printf("  commit: %s\n", commitInfo)
		fmt.Printf("  built:  %s\n", dateInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

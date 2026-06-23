package cmd

import (
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the PVM version and build information",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Title("pvm %s", versionInfo)
		ui.Info("  commit: %s", commitInfo)
		ui.Info("  built:  %s", dateInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the active PHP version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := config.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("could not read active version: %w", err)
		}

		if version == "" {
			ui.Info("No active PHP version set.")
			ui.Detail("Run `pvm use <version>` to activate one.")
			return nil
		}

		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		phpExePath := filepath.Join(versionsDir, version, "php.exe")
		ui.Title("PHP %s", version)
		ui.Info("Path: %s", phpExePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

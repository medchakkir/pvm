package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the active PHP version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := config.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("✗ could not read active version: %w", err)
		}

		if version == "" {
			fmt.Println("No active PHP version set.")
			fmt.Println("Run `pvm use <version>` to activate one.")
			return nil
		}

		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		phpExePath := filepath.Join(versionsDir, version, "php.exe")
		fmt.Printf("PHP %s\n", version)
		fmt.Printf("Path: %s\n", phpExePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

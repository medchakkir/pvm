package cmd

import (
	"fmt"
	"os"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pvm",
	Short: "PVM — PHP Version Manager for Windows",
	Long: `PVM lets you install, switch, and manage multiple PHP versions
on Windows with a single command.

  pvm install 8.3       Install PHP 8.3
  pvm use 8.2           Switch to PHP 8.2
  pvm list              Show installed versions
  pvm list-remote       Show available versions`,

	// Runs before every subcommand
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to initialize PVM home: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

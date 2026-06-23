package cmd

import (
	"fmt"
	"os"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var verbose bool

// Version metadata, set by SetVersionInfo from main at startup.
var (
	versionInfo = "dev"
	commitInfo  = "none"
	dateInfo    = "unknown"
)

// SetVersionInfo wires build-time version metadata into the CLI.
func SetVersionInfo(version, commit, date string) {
	versionInfo = version
	commitInfo = commit
	dateInfo = date
	rootCmd.Version = version
}

var rootCmd = &cobra.Command{
	Use:          "pvm",
	Short:        "PVM — PHP Version Manager for Windows",
	SilenceUsage: true, // don't print usage on runtime errors, only on wrong args
	Long: `PVM lets you install, switch, and manage multiple PHP versions
on Windows with a single command.

  pvm install 8.3           Install PHP 8.3 (Thread Safe)
  pvm install --nts 8.3     Install PHP 8.3 (Non-Thread Safe)
  pvm use 8.2               Switch to PHP 8.2
  pvm list                  Show installed versions
  pvm list-remote           Show available versions
  pvm current               Show active version
  pvm uninstall 8.1         Remove PHP 8.1
  pvm extensions list       Manage extensions for the active version
  pvm bin                   Print the PATH entry to add to your system`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to initialize PVM home directory: %w", err)
		}
		if verbose {
			home, _ := config.PVMHome()
			fmt.Fprintf(os.Stderr, "[debug] PVM home: %s\n", home)
		}
		return nil
	},
}

// Verbose returns true if the --verbose flag was set.
// Used by other commands to print extra debug info.
func Verbose() bool {
	return verbose
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error("%s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print debug information")
}

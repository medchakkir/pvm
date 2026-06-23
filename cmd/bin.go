package cmd

import (
	"fmt"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/spf13/cobra"
)

var binCmd = &cobra.Command{
	Use:   "bin",
	Short: "Print the path to add to your PATH environment variable",
	Long: `Prints the PVM shims directory that needs to be on your PATH.
Add this path to your system PATH so that PHP is accessible from any terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shimsDir, err := config.ShimsDir()
		if err != nil {
			return fmt.Errorf("✗ could not locate shims directory: %w", err)
		}

		fmt.Println(shimsDir)
		fmt.Println()
		fmt.Println("Add the above path to your PATH. In PowerShell, run:")
		fmt.Printf(
			"  [Environment]::SetEnvironmentVariable(\"Path\", $env:Path + \";%s\", \"User\")\n",
			shimsDir,
		)
		fmt.Println()
		fmt.Println("Then restart your terminal.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(binCmd)
}

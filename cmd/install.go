package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Download and install a specific PHP version",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("(install %s — coming in Phase 4)\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

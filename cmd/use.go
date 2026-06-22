package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/env"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch the active PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		// 1. Parse the requested version
		requested, err := php.ParseVersion(input)
		if err != nil {
			return fmt.Errorf("✗ invalid version %q: %w", input, err)
		}

		// 2. Find the best match among installed versions
		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("✗ could not read versions directory: %w", err)
		}

		var matched *php.PHPVersion
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			v, err := php.ParseVersion(entry.Name())
			if err != nil {
				continue
			}
			// Exact match, or Major.Minor match when patch was omitted
			if v.Compare(requested) == 0 || (requested.Patch == 0 && v.MatchesMinor(requested)) {
				if matched == nil || v.Compare(*matched) > 0 {
					copy := v
					matched = &copy
				}
			}
		}

		if matched == nil {
			return fmt.Errorf(
				"✗ PHP %s is not installed.\n  Run `pvm install %s` first.",
				input, input,
			)
		}

		// 3. Verify php.exe is intact
		phpExePath := filepath.Join(versionsDir, matched.String(), "php.exe")
		if _, err := os.Stat(phpExePath); os.IsNotExist(err) {
			return fmt.Errorf("✗ php.exe missing for version %s — try reinstalling it.", matched)
		}

		// 4. Write the shim
		shimsDir, err := config.ShimsDir()
		if err != nil {
			return err
		}

		if err := env.WriteShim(shimsDir, phpExePath); err != nil {
			return fmt.Errorf("✗ %w", err)
		}

		// 5. Save the active version
		if err := config.SetCurrentVersion(matched.String()); err != nil {
			return fmt.Errorf("✗ could not save active version: %w", err)
		}

		fmt.Printf("✓ Now using PHP %s\n", matched)

		// 6. Warn if shims dir is not on PATH yet
		if !env.IsOnPath(shimsDir) {
			fmt.Println(env.PathInstructions(shimsDir))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/env"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var useNtsFlag bool

var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch the active PHP version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		requested, err := php.ParseVersion(input)
		if err != nil {
			return fmt.Errorf("✗ invalid version %q: %w", input, err)
		}

		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("✗ could not read versions directory: %w", err)
		}

		// Find best match respecting the --nts flag
		wantType := "TS"
		if useNtsFlag {
			wantType = "NTS"
		}

		var matchedDir string
		var matchedVersion php.PHPVersion

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name() // e.g. "8.3.7-TS"

			// Check type suffix matches what was requested
			if !strings.HasSuffix(name, "-"+wantType) {
				continue
			}

			versionPart := strings.TrimSuffix(name, "-"+wantType)
			v, err := php.ParseVersion(versionPart)
			if err != nil {
				continue
			}

			if v.Compare(requested) == 0 || (requested.Patch == 0 && v.MatchesMinor(requested)) {
				if matchedDir == "" || v.Compare(matchedVersion) > 0 {
					matchedDir = name
					matchedVersion = v
				}
			}
		}

		if matchedDir == "" {
			installHint := fmt.Sprintf("pvm install %s", input)
			if useNtsFlag {
				installHint = fmt.Sprintf("pvm install --nts %s", input)
			}
			return fmt.Errorf(
				"✗ PHP %s (%s) is not installed.\n  Run `%s` first.",
				input, wantType, installHint,
			)
		}

		// Verify php.exe is intact
		phpExePath := filepath.Join(versionsDir, matchedDir, "php.exe")
		if _, err := os.Stat(phpExePath); os.IsNotExist(err) {
			return fmt.Errorf("✗ php.exe missing for %s — try reinstalling it", matchedDir)
		}

		// Write the shim
		shimsDir, err := config.ShimsDir()
		if err != nil {
			return err
		}

		if err := env.WriteShim(shimsDir, phpExePath); err != nil {
			return fmt.Errorf("✗ %w", err)
		}

		// Save the active version (store dirName so list can highlight it)
		if err := config.SetCurrentVersion(matchedDir); err != nil {
			return fmt.Errorf("✗ could not save active version: %w", err)
		}

		fmt.Printf("✓ Now using PHP %s (%s)\n", matchedVersion, wantType)

		if !env.IsOnPath(shimsDir) {
			fmt.Println(env.PathInstructions(shimsDir))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().BoolVar(&useNtsFlag, "nts", false, "Switch to the Non-Thread Safe build")
}

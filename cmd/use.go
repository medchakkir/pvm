package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/env"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var useNtsFlag bool

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch the active PHP version",
	Long: `Switch the active PHP version.

  pvm use 8.3        # switch to PHP 8.3 (resolves to latest installed patch)
  pvm use 8.3.7      # switch to exact version
  pvm use --nts 8.3  # switch to the NTS build specifically
  pvm use            # read version from .pvmrc in current or parent directory`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string

		if len(args) == 1 {
			input = args[0]
		} else {
			// No version argument — look for a .pvmrc file
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("could not determine current directory: %w", err)
			}

			pvmrcPath, err := config.FindPvmrc(cwd)
			if err != nil {
				return err
			}
			if pvmrcPath == "" {
				return fmt.Errorf(
					"no version specified and no .pvmrc found in this directory or any parent.\n" +
						"  Run `pvm init` to pin a version here, or specify one: `pvm use <version>`",
				)
			}

			input, err = config.ReadPvmrc(pvmrcPath)
			if err != nil {
				return err
			}

			ui.Detail("Using version from %s", pvmrcPath)

			// If the .pvmrc contains a full dir name like "8.3.7-TS", extract
			// the type suffix so --nts logic works correctly below.
			if !useNtsFlag {
				if strings.HasSuffix(strings.ToUpper(input), "-NTS") {
					useNtsFlag = true
					input = input[:len(input)-4]
				} else if strings.HasSuffix(strings.ToUpper(input), "-TS") {
					input = input[:len(input)-3]
				}
			}
		}

		requested, err := php.ParseVersion(input)
		if err != nil {
			return fmt.Errorf("invalid version %q: %w", input, err)
		}

		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		entries, err := os.ReadDir(versionsDir)
		if err != nil {
			return fmt.Errorf("could not read versions directory: %w", err)
		}

		// Determine which types to search for.
		// If --nts is explicitly passed, look for NTS only.
		// Otherwise search both TS and NTS, preferring TS when both exist.
		searchTypes := []string{"TS", "NTS"}
		if useNtsFlag {
			searchTypes = []string{"NTS"}
		}

		var matchedDir string
		var matchedVersion php.PHPVersion
		var matchedType string

		for _, wantType := range searchTypes {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}

				name := entry.Name() // e.g. "8.3.7-TS"

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
						matchedType = wantType
					}
				}
			}
			// If we found a match at this preference level, stop searching
			if matchedDir != "" {
				break
			}
		}

		if matchedDir == "" {
			installHint := fmt.Sprintf("pvm install %s", input)
			if useNtsFlag {
				installHint = fmt.Sprintf("pvm install --nts %s", input)
			}
			return fmt.Errorf(
				"PHP %s is not installed.\n  Run `%s` first.",
				input, installHint,
			)
		}

		// Show resolution when the user supplied a minor version (no patch)
		if requested.Patch == 0 {
			ui.Detail("Resolved %s → %s", requested.MinorString(), matchedVersion)
		}

		// Verify php.exe is intact
		matchedVersionDir := filepath.Join(versionsDir, matchedDir)
		phpExePath := filepath.Join(matchedVersionDir, "php.exe")
		if _, err := os.Stat(phpExePath); os.IsNotExist(err) {
			return fmt.Errorf("php.exe missing for %s — try reinstalling it", matchedDir)
		}

		// Write the shims (php, php-cgi, …) for the active version
		shimsDir, err := config.ShimsDir()
		if err != nil {
			return err
		}

		if err := env.WriteShim(shimsDir, env.DefaultShims(matchedVersionDir)); err != nil {
			return err
		}

		// Save the active version (store dirName so list can highlight it)
		if err := config.SetCurrentVersion(matchedDir); err != nil {
			return fmt.Errorf("could not save active version: %w", err)
		}

		ui.Success("Now using PHP %s (%s)", matchedVersion, matchedType)

		if !env.IsOnPath(shimsDir) {
			ui.Info("%s", env.PathInstructions(shimsDir))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().BoolVar(&useNtsFlag, "nts", false, "Switch to the Non-Thread Safe build")
}

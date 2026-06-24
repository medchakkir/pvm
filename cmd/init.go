package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [version]",
	Short: "Pin a PHP version for this directory by creating a .pvmrc file",
	Long: `Creates a .pvmrc file in the current directory that records which PHP version
this project uses. Running 'pvm use' without arguments in this directory (or
any subdirectory) will automatically switch to that version.

  pvm init           # pin the currently active version
  pvm init 8.3       # pin a specific version (must be installed)
  pvm init 8.3.7-TS  # pin an exact version + build type`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine current directory: %w", err)
		}

		var version string

		if len(args) == 1 {
			// Version supplied explicitly — validate it is installed
			input := args[0]
			version, err = resolveInstalledVersion(input)
			if err != nil {
				return err
			}
		} else {
			// No argument — use the currently active version
			version, err = config.GetCurrentVersion()
			if err != nil {
				return fmt.Errorf("could not read active version: %w", err)
			}
			if version == "" {
				return fmt.Errorf(
					"no active PHP version set.\n" +
						"  Run `pvm use <version>` first, or supply a version: `pvm init <version>`",
				)
			}
		}

		// Warn and confirm if a .pvmrc already exists in this directory
		if config.PvmrcExists(cwd) {
			existing, _ := config.ReadPvmrc(fmt.Sprintf("%s/%s", cwd, config.PvmrcFile))
			ui.Warning(".pvmrc already exists in this directory (currently: %s)", existing)
			fmt.Printf("Overwrite with %q? [y/N] ", version)

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
			if answer != "y" {
				ui.Warning("Aborted.")
				return nil
			}
		}

		if err := config.WritePvmrc(cwd, version); err != nil {
			return fmt.Errorf("could not write .pvmrc: %w", err)
		}

		ui.Success("Created .pvmrc → %s", version)
		ui.Detail("  Run `pvm use` (no arguments) in this directory to activate it.")
		return nil
	},
}

// resolveInstalledVersion takes a user-supplied version string (e.g. "8.3",
// "8.3.7", "8.3.7-TS", "8.3.7-NTS") and returns the canonical directory name
// of a matching installed version. Errors if nothing installed matches.
func resolveInstalledVersion(input string) (string, error) {
	// Strip a trailing -TS / -NTS suffix if the user included it, then use it
	// to narrow the search.
	wantTS := true   // default to TS
	hasTypeSuffix := false
	stripped := input

	if strings.HasSuffix(strings.ToUpper(input), "-NTS") {
		wantTS = false
		hasTypeSuffix = true
		stripped = input[:len(input)-4]
	} else if strings.HasSuffix(strings.ToUpper(input), "-TS") {
		wantTS = true
		hasTypeSuffix = true
		stripped = input[:len(input)-3]
	}

	requested, err := php.ParseVersion(stripped)
	if err != nil {
		return "", fmt.Errorf("invalid version %q: %w", input, err)
	}

	versionsDir, err := config.VersionsDir()
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return "", fmt.Errorf("could not read versions directory: %w", err)
	}

	// Collect candidates: exact match on version, respecting type suffix if given.
	// When no type suffix is given, prefer TS but accept NTS.
	searchTypes := []string{"TS", "NTS"}
	if hasTypeSuffix {
		if wantTS {
			searchTypes = []string{"TS"}
		} else {
			searchTypes = []string{"NTS"}
		}
	}

	for _, wantType := range searchTypes {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, "-"+wantType) {
				continue
			}
			versionPart := strings.TrimSuffix(name, "-"+wantType)
			v, err := php.ParseVersion(versionPart)
			if err != nil {
				continue
			}
			if v.Compare(requested) == 0 || (requested.Patch == 0 && v.MatchesMinor(requested)) {
				return name, nil
			}
		}
		if hasTypeSuffix {
			break
		}
	}

	return "", fmt.Errorf(
		"PHP %s is not installed.\n  Run `pvm install %s` first.",
		input, stripped,
	)
}

func init() {
	rootCmd.AddCommand(initCmd)
}

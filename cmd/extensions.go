package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var extensionsCmd = &cobra.Command{
	Use:     "extensions",
	Aliases: []string{"ext", "extension"},
	Short:   "Manage PHP extensions for the active version",
	Long: `List, enable, and disable PHP extensions for the currently active version.

  pvm extensions list                 Show every extension and its status
  pvm extensions enable curl,gd       Enable one or more extensions
  pvm extensions disable xdebug       Disable one or more extensions`,
}

var extListCmd = &cobra.Command{
	Use:   "list",
	Short: "List extensions and their status",
	RunE:  runExtList,
}

var extEnableCmd = &cobra.Command{
	Use:   "enable <ext>[,<ext>...]",
	Short: "Enable one or more extensions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExtToggle(args, true)
	},
}

var extDisableCmd = &cobra.Command{
	Use:   "disable <ext>[,<ext>...]",
	Short: "Disable one or more extensions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runExtToggle(args, false)
	},
}

func init() {
	extensionsCmd.AddCommand(extListCmd, extEnableCmd, extDisableCmd)
	rootCmd.AddCommand(extensionsCmd)
}

// activeVersionPaths resolves the active version's directory plus the paths to
// its php.ini and ext/ directory. It errors when no version is active.
func activeVersionPaths() (versionDir, iniPath, extDir string, err error) {
	current, err := config.GetCurrentVersion()
	if err != nil {
		return "", "", "", fmt.Errorf("could not read active version: %w", err)
	}
	current = strings.TrimSpace(current)
	if current == "" {
		return "", "", "", fmt.Errorf("no active PHP version. Run `pvm use <version>` first.")
	}

	versionsDir, err := config.VersionsDir()
	if err != nil {
		return "", "", "", err
	}

	versionDir = filepath.Join(versionsDir, current)
	iniPath = filepath.Join(versionDir, "php.ini")
	extDir = filepath.Join(versionDir, "ext")
	return versionDir, iniPath, extDir, nil
}

func runExtList(cmd *cobra.Command, args []string) error {
	versionDir, iniPath, extDir, err := activeVersionPaths()
	if err != nil {
		return err
	}

	avail, err := php.AvailableExtensions(extDir)
	if err != nil {
		return fmt.Errorf("could not read extensions directory %s: %w", extDir, err)
	}

	entryByName := make(map[string]php.ExtEntry)
	if _, statErr := os.Stat(iniPath); statErr == nil {
		ini, loadErr := php.LoadIni(iniPath)
		if loadErr != nil {
			return fmt.Errorf("could not read php.ini: %w", loadErr)
		}
		for _, e := range ini.Entries() {
			entryByName[e.Name] = e
		}
	} else {
		ui.Warning("No php.ini found in %s — every extension shows as available.", versionDir)
		ui.Detail("Enabling an extension will create php.ini from php.ini-development.")
	}

	names := make(map[string]struct{})
	for name := range avail {
		names[name] = struct{}{}
	}
	for name := range entryByName {
		names[name] = struct{}{}
	}

	if len(names) == 0 {
		ui.Info("No extensions found for the active version.")
		return nil
	}

	sorted := make([]string, 0, len(names))
	for name := range names {
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)

	var enabled, disabled, available, missing int

	ui.Title("Extensions for the active PHP version:")
	for _, name := range sorted {
		entry, hasEntry := entryByName[name]
		_, hasFile := avail[name]

		var label string
		switch {
		case hasEntry && !hasFile:
			label = "[missing file]"
			missing++
		case hasEntry && entry.Enabled:
			label = "[enabled]"
			enabled++
		case hasEntry && !entry.Enabled:
			label = "[disabled]"
			disabled++
		default:
			label = "[available]"
			available++
		}

		printExtLine(name, label)
	}

	ui.Info("")
	ui.Info("%d enabled, %d disabled, %d available, %d missing file(s).",
		enabled, disabled, available, missing)
	return nil
}

// printExtLine renders one status line, coloring it to match its state. The
// ui helpers prefix their own symbol, so the plain-info case keeps a leading
// indent for visual alignment.
func printExtLine(name, label string) {
	switch label {
	case "[enabled]":
		ui.Success("%-24s %s", name, label)
	case "[missing file]":
		ui.Warning("%-24s %s", name, label)
	default:
		ui.Info("  %-24s %s", name, label)
	}
}

func runExtToggle(args []string, enable bool) error {
	requested := parseExtArgs(args)
	if len(requested) == 0 {
		return fmt.Errorf("no extension names provided")
	}

	versionDir, iniPath, extDir, err := activeVersionPaths()
	if err != nil {
		return err
	}

	if err := ensureIniExists(iniPath, versionDir); err != nil {
		return err
	}

	avail, err := php.AvailableExtensions(extDir)
	if err != nil {
		return fmt.Errorf("could not read extensions directory %s: %w", extDir, err)
	}

	ini, err := php.LoadIni(iniPath)
	if err != nil {
		return fmt.Errorf("could not read php.ini: %w", err)
	}

	entryByName := make(map[string]php.ExtEntry)
	for _, e := range ini.Entries() {
		entryByName[e.Name] = e
	}

	changed := false
	for _, raw := range requested {
		name := php.NormalizeExtName(raw)
		entry, hasEntry := entryByName[name]
		_, hasFile := avail[name]

		if enable {
			switch {
			case hasEntry && entry.Enabled:
				ui.Info("%s is already enabled.", name)
			case hasEntry:
				ini.SetEnabled(entry.LineIndex, true)
				changed = true
				ui.Success("Enabled %s", name)
			case hasFile:
				ini.AddExtension("extension", name)
				changed = true
				ui.Success("Enabled %s", name)
			default:
				ui.Error("%s not found — no php.ini entry and no ext/%s.dll", name, name)
			}
			continue
		}

		switch {
		case !hasEntry:
			ui.Warning("%s is not listed in php.ini — nothing to disable.", name)
		case !entry.Enabled:
			ui.Info("%s is already disabled.", name)
		default:
			ini.SetEnabled(entry.LineIndex, false)
			changed = true
			ui.Success("Disabled %s", name)
		}
	}

	if !changed {
		return nil
	}

	if err := ini.Save(); err != nil {
		return fmt.Errorf("could not write php.ini: %w", err)
	}
	ui.Detail("Updated %s", iniPath)
	return nil
}

// parseExtArgs flattens positional args and comma-separated lists into a clean
// slice of non-empty extension tokens.
func parseExtArgs(args []string) []string {
	var out []string
	for _, arg := range args {
		for _, part := range strings.Split(arg, ",") {
			if p := strings.TrimSpace(part); p != "" {
				out = append(out, p)
			}
		}
	}
	return out
}

// ensureIniExists makes sure php.ini is present, creating it from
// php.ini-development when only the template exists.
func ensureIniExists(iniPath, versionDir string) error {
	if _, err := os.Stat(iniPath); err == nil {
		return nil
	}

	template := filepath.Join(versionDir, "php.ini-development")
	data, err := os.ReadFile(template)
	if err != nil {
		return fmt.Errorf(
			"no php.ini found in %s and php.ini-development is unavailable — cannot edit extensions",
			versionDir,
		)
	}

	if err := os.WriteFile(iniPath, data, 0644); err != nil {
		return fmt.Errorf("could not create php.ini: %w", err)
	}
	ui.Info("Created php.ini from php.ini-development.")
	return nil
}

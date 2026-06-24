package cmd

import (
	"os"
	"path/filepath"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/medchakkir/pvm/internal/ui"
	"github.com/spf13/cobra"
)

var (
	installNtsFlag      bool
	installComposerFlag bool
)

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Download and install a PHP version",
	Long: `Downloads and installs a PHP version from php.net.
Installs the Thread Safe (TS) build by default.
Use --nts to install the Non-Thread Safe build instead.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		threadSafe := !installNtsFlag

		// 1. Resolve the exact version from php.net
		ui.Info("Fetching available versions from php.net...")
		remoteVersions, err := php.FetchRemoteVersions()
		if err != nil {
			return err
		}

		target, err := php.FindRemoteVersion(input, remoteVersions, threadSafe)
		if err != nil {
			return err
		}

		// 2. Check if already installed
		versionsDir, err := config.VersionsDir()
		if err != nil {
			return err
		}

		// Include type in directory name to allow both TS and NTS side by side
		dirName := target.Version.String() + "-" + target.TypeLabel()
		versionDir := filepath.Join(versionsDir, dirName)

		if _, err := os.Stat(versionDir); err == nil {
			ui.Success("PHP %s is already installed.", target)
			// Allow adding Composer to an existing install without reinstalling.
			if installComposerFlag && !php.ComposerInstalled(versionDir) {
				if err := installComposer(versionDir, target.Version); err != nil {
					return err
				}
			}
			ui.Detail("  Run `pvm use %s` to switch to it.", target.Version)
			return nil
		}

		ui.Info("Installing PHP %s...", target)

		// 3. Download ZIP to a temp file
		tmpFile := filepath.Join(os.TempDir(), target.ZipName)
		defer os.Remove(tmpFile)

		if err := php.DownloadZip(target.DownloadURL, tmpFile); err != nil {
			return err
		}

		// 4. Extract into ~/.pvm/versions/<version>-<type>/
		ui.Info("\nExtracting...")
		if err := php.ExtractZip(tmpFile, versionDir); err != nil {
			os.RemoveAll(versionDir)
			return err
		}

		// 5. Verify php.exe exists
		if err := php.VerifyInstall(versionDir); err != nil {
			os.RemoveAll(versionDir)
			return err
		}

		ui.Success("PHP %s installed successfully.", target)

		// 6. Optionally install Composer into the version directory.
		if installComposerFlag {
			if err := installComposer(versionDir, target.Version); err != nil {
				return err
			}
		}

		ui.Detail("  Run `pvm use %s` to activate it.", target.Version)
		return nil
	},
}

// installComposer downloads composer.phar into versionDir for the given PHP
// version and reports progress. The composer.bat shim is written later, on `use`.
func installComposer(versionDir string, version php.PHPVersion) error {
	ui.Info("\nInstalling Composer...")
	if err := php.InstallComposer(versionDir, version); err != nil {
		return err
	}
	ui.Success("Composer installed.")
	ui.Detail("  The `composer` command becomes available after `pvm use`.")
	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&installNtsFlag, "nts", false, "Install the Non-Thread Safe build")
	installCmd.Flags().BoolVar(&installComposerFlag, "with-composer", false, "Also download Composer for this version")
}

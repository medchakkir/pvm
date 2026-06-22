package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/medchakkir/pvm/internal/config"
	"github.com/medchakkir/pvm/internal/php"
	"github.com/spf13/cobra"
)

var installNtsFlag bool

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
		fmt.Println("Fetching available versions from php.net...")
		remoteVersions, err := php.FetchRemoteVersions()
		if err != nil {
			return fmt.Errorf("✗ %w", err)
		}

		target, err := php.FindRemoteVersion(input, remoteVersions, threadSafe)
		if err != nil {
			return fmt.Errorf("✗ %w", err)
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
			fmt.Printf("✓ PHP %s is already installed.\n", target)
			fmt.Printf("  Run `pvm use %s` to switch to it.\n", target.Version)
			return nil
		}

		fmt.Printf("Installing PHP %s...\n", target)

		// 3. Download ZIP to a temp file
		tmpFile := filepath.Join(os.TempDir(), target.ZipName)
		defer os.Remove(tmpFile)

		if err := php.DownloadZip(target.DownloadURL, tmpFile); err != nil {
			return fmt.Errorf("✗ %w", err)
		}

		// 4. Extract into ~/.pvm/versions/<version>-<type>/
		fmt.Println("\nExtracting...")
		if err := php.ExtractZip(tmpFile, versionDir); err != nil {
			os.RemoveAll(versionDir)
			return fmt.Errorf("✗ %w", err)
		}

		// 5. Verify php.exe exists
		if err := php.VerifyInstall(versionDir); err != nil {
			os.RemoveAll(versionDir)
			return fmt.Errorf("✗ %w", err)
		}

		fmt.Printf("✓ PHP %s installed successfully.\n", target)
		fmt.Printf("  Run `pvm use %s` to activate it.\n", target.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&installNtsFlag, "nts", false, "Install the Non-Thread Safe build")
}

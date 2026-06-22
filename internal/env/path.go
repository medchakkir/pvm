package env

import (
	"os"
	"strings"
)

// IsOnPath returns true if the given directory is already present in PATH.
func IsOnPath(dir string) bool {
	pathEnv := os.Getenv("PATH")
	for _, p := range strings.Split(pathEnv, ";") {
		if strings.EqualFold(strings.TrimSpace(p), strings.TrimSpace(dir)) {
			return true
		}
	}
	return false
}

// PathInstructions returns a user-friendly message explaining how to
// add the shims directory to PATH manually on Windows.
func PathInstructions(shimsDir string) string {
	return `
  PVM shims are not on your PATH yet. Add them by running this in PowerShell:

  [Environment]::SetEnvironmentVariable("Path", $env:Path + ";` + shimsDir + `", "User")

  Then restart your terminal for the change to take effect.
`
}

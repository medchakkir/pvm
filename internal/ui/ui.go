// Package ui centralizes PVM's terminal output styling so commands don't
// scatter raw fmt.Printf calls with inline symbols. It is backed by
// github.com/fatih/color, which automatically disables colors when output is
// not a terminal and honors the NO_COLOR environment variable.
package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
	warningColor = color.New(color.FgYellow)
	infoColor    = color.New(color.FgCyan)
	titleColor   = color.New(color.Bold)
)

// Success prints a green success line prefixed with a check mark to stdout.
func Success(format string, a ...any) {
	successColor.Fprintf(color.Output, "\u2713 %s\n", fmt.Sprintf(format, a...))
}

// Error prints a red error line prefixed with a cross to stderr.
func Error(format string, a ...any) {
	errorColor.Fprintf(color.Error, "\u2717 %s\n", fmt.Sprintf(format, a...))
}

// Warning prints a yellow warning line to stdout.
func Warning(format string, a ...any) {
	warningColor.Fprintf(color.Output, "! %s\n", fmt.Sprintf(format, a...))
}

// Info prints a plain, unstyled line to stdout. It exists so callers route all
// output through this package rather than mixing in raw fmt calls.
func Info(format string, a ...any) {
	fmt.Fprintf(color.Output, format+"\n", a...)
}

// Detail prints a dimmed/cyan secondary line to stdout, useful for hints and
// follow-up suggestions under a primary message.
func Detail(format string, a ...any) {
	infoColor.Fprintf(color.Output, format+"\n", a...)
}

// Title prints a bold heading line to stdout.
func Title(format string, a ...any) {
	titleColor.Fprintf(color.Output, "%s\n", fmt.Sprintf(format, a...))
}

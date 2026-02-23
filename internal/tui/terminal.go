package tui

import (
	"os"

	"golang.org/x/term"
)

// IsInteractive returns true when the CLI should use full interactive TUI.
// Retained for backward compatibility — delegates to the output mode system.
func IsInteractive() bool {
	return currentMode == ModeInteractive && term.IsTerminal(int(os.Stdin.Fd()))
}

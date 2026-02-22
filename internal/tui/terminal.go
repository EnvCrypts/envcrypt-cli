package tui

import (
	"os"

	"golang.org/x/term"
)

func IsInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

// OutputMode determines how the CLI renders output.
type OutputMode int

const (
	ModeInteractive OutputMode = iota
	ModePlain
	ModeJSON
)

// TermLevel describes the terminal's capability.
type TermLevel int

const (
	TermFull    TermLevel = iota // Full TTY: colors, cursor control, animations
	TermLimited                  // Limited TTY: colors but no animations (e.g. TERM=dumb)
	TermNone                     // Non-TTY: pipes, CI, scripts
)

// Package-level state, initialised by InitOutput.
var (
	currentMode  OutputMode = ModeInteractive
	currentLevel TermLevel  = TermFull
	isQuiet      bool
	colorEnabled bool = true
)

// DetectTermLevel inspects stdout to determine terminal capabilities.
func DetectTermLevel() TermLevel {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return TermNone
	}
	if t := os.Getenv("TERM"); t == "dumb" {
		return TermLimited
	}
	return TermFull
}

// InitOutput configures the global output mode. Call once from root PersistentPreRun.
func InitOutput(jsonMode, noColor, quiet bool) {
	currentLevel = DetectTermLevel()

	switch {
	case jsonMode:
		currentMode = ModeJSON
		colorEnabled = false
	case currentLevel == TermNone:
		currentMode = ModePlain
		colorEnabled = false
	case currentLevel == TermLimited:
		currentMode = ModePlain
		colorEnabled = !noColor
	default:
		currentMode = ModeInteractive
		colorEnabled = !noColor
	}

	// NO_COLOR convention (https://no-color.org)
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}

	isQuiet = quiet
}

// Mode returns the active output mode.
func Mode() OutputMode { return currentMode }

// IsQuiet returns whether --quiet is active.
func IsQuiet() bool { return isQuiet }

// ColorEnabled returns whether ANSI color output is allowed.
func ColorEnabled() bool { return colorEnabled }

// --- JSON helpers ---

// JSONMessage emits a single structured JSON object to stdout.
func JSONMessage(level string, msg string) {
	obj := map[string]any{
		"level":   level,
		"message": msg,
		"time":    time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(obj)
	fmt.Fprintln(os.Stdout, string(data))
}

// JSONData emits arbitrary structured data as JSON to stdout.
func JSONData(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Fprintln(os.Stdout, string(data))
}

// PlainTimestamp returns a bracketed timestamp for plain mode log lines.
func PlainTimestamp() string {
	return fmt.Sprintf("[%s]", time.Now().Format("15:04:05"))
}

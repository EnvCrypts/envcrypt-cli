package tui

import (
	"errors"
	"fmt"
)

// ErrCancelled is returned when the user aborts an operation.
var ErrCancelled = errors.New("operation cancelled")

// Cancelled prints a standardised cancellation message and returns ErrCancelled.
func Cancelled() error {
	if currentMode == ModeJSON {
		JSONMessage("cancelled", "Operation cancelled")
	} else if !isQuiet {
		fmt.Printf("%s %s\n", IconWarn, StyleMuted.Render("Operation cancelled"))
	}
	return ErrCancelled
}

// ErrorWithHint renders an error with an actionable hint.
func ErrorWithHint(msg string, err error, hint string) error {
	return &UIError{
		Message: msg,
		Err:     err,
		Hint:    hint,
	}
}

// UIError is a user-facing error that carries optional detail and hint text.
type UIError struct {
	Message string
	Err     error
	Hint    string
}

func (e *UIError) Error() string {
	if e == nil {
		return "unknown error"
	}
	return e.Message
}

func (e *UIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

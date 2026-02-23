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
	if currentMode == ModeJSON {
		obj := map[string]any{"level": "error", "message": msg}
		if err != nil {
			obj["detail"] = err.Error()
		}
		if hint != "" {
			obj["hint"] = hint
		}
		JSONData(obj)
		if err != nil {
			return fmt.Errorf("%s: %w", msg, err)
		}
		return fmt.Errorf("%s", msg)
	}

	parts := fmt.Sprintf("%s %s", IconCross, msg)
	if err != nil {
		parts += fmt.Sprintf("\n  %s", StyleMuted.Render(err.Error()))
	}
	if hint != "" {
		parts += fmt.Sprintf("\n  %s %s", IconInfo, StyleInfo.Render(hint))
	}
	return fmt.Errorf("%s", parts)
}

// APIErrorDetail represents the structured error returned by the server.
type APIErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint"`
}

// MapAPIError converts an APIErrorDetail into a user-friendly error with hint.
func MapAPIError(detail *APIErrorDetail) error {
	if detail == nil {
		return fmt.Errorf("unknown server error")
	}
	return ErrorWithHint(detail.Message, nil, detail.Hint)
}

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
)

func mapEnvReadError(path string, err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("env file %q does not exist", path)
	}
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("permission denied reading env file %q", path)
	}
	return fmt.Errorf("failed to read env file %q: %w", path, err)
}

func resolveEnvFile(flagPath string) (string, error) {
	// Explicit flag always wins
	if flagPath != "" {
		if fileExists(flagPath) {
			return flagPath, nil
		}
		return "", fmt.Errorf("env file %q does not exist", flagPath)
	}

	if fileExists(".env") {
		return ".env", nil
	}

	return "", errors.New("no .env file found")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func handlePromptError(err error, message, hint string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, tui.ErrNonInteractive) {
		return tui.Error(message, nil, hint)
	}
	if errors.Is(err, tui.ErrCancelled) {
		return tui.Cancelled()
	}
	return tui.Error("prompt failed", err)
}

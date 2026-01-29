package cmd

import (
	"errors"
	"fmt"
	"os"
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

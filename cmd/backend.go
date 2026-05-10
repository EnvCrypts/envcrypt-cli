package cmd

import (
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
)

func applyBackendURL(backendURL string) error {
	backendURL = strings.TrimSpace(backendURL)
	if backendURL == "" {
		return nil
	}

	if err := tui.ValidateBackendURL(backendURL); err != nil {
		return tui.Error(
			"invalid custom backend URL",
			err,
			"Use a full URL like https://api.example.com or http://localhost:8081",
		)
	}

	if err := config.SaveBackendURL(backendURL); err != nil {
		return tui.Error(
			"failed to save custom backend URL",
			err,
			"Your backend URL could not be written to the local config file",
		)
	}

	if Application != nil {
		Application.SetBaseURL(backendURL)
	}

	return nil
}

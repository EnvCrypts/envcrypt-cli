package cmd

import (
	"strings"
	"testing"

	"github.com/envcrypts/envcrypt-cli/internal/app"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestApplyBackendURLUpdatesPersistedAndActiveURL(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	Application = app.NewApp("http://old.example")
	t.Cleanup(func() {
		Application = nil
	})

	err := applyBackendURL("https://api.example.com")
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com", Application.BaseURL())
	require.Equal(t, "https://api.example.com", viper.GetString("api.base_url"))
}

func TestApplyBackendURLRejectsMalformedURL(t *testing.T) {
	Application = app.NewApp("http://old.example")
	t.Cleanup(func() {
		Application = nil
	})

	err := applyBackendURL("localhost:8081")
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid custom backend URL"))
}

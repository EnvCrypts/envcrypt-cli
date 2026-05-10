package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load() error {

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dir, "envcrypt")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(appDir)

	// Defaults
	viper.SetDefault("api.base_url", "https://api-envcrypt.vijayvenkatj.in")

	// Allow ENV override (ENVCRYPT_USER_EMAIL, etc.)
	viper.SetEnvPrefix("envcrypt")
	viper.AutomaticEnv()

	// Read if exists (don’t fail if missing)
	_ = viper.ReadInConfig()

	return nil
}

func SaveBackendURL(url string) error {
	url = strings.TrimSpace(url)
	if err := validateBackendURL(url); err != nil {
		return err
	}

	viper.Set("api.base_url", url)

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dir, "envcrypt")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(appDir, "config.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return viper.WriteConfigAs(path)
	}

	return viper.WriteConfig()
}

func validateBackendURL(raw string) error {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid backend URL: use a full URL like https://api.example.com or http://localhost:8081")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("invalid backend URL: scheme must be http or https")
	}
	return nil
}

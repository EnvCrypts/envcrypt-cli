package config

import (
	"os"
	"path/filepath"

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
	viper.SetDefault("api.base_url", "https://api.envcrypt.dev")

	// Allow ENV override (ENVCRYPT_USER_EMAIL, etc.)
	viper.SetEnvPrefix("envcrypt")
	viper.AutomaticEnv()

	// Read if exists (don’t fail if missing)
	_ = viper.ReadInConfig()

	return nil
}

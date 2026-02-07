package cryptoutils

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

func SavePrivateKey(user string, secret []byte) error {
	encoded := base64.StdEncoding.EncodeToString(secret)
	err := keyring.Set("envcrypt", user, encoded)
	if err != nil {
		return err
	}

	return nil
}

func LoadPrivateKey(user string) ([]byte, error) {
	secret, err := keyring.Get("envcrypt", user)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func DeletePrivateKey(user string) error {
	_ = keyring.Delete("envcrypt", user)
	return nil
}

func SaveUserEmail(email string) error {
	viper.Set("user.email", email)

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dir, "envcrypt")
	path := filepath.Join(appDir, "config.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return viper.WriteConfigAs(path)
	}

	return viper.WriteConfig()
}

func RemoveUserEmail() error {
	viper.Set("user.email", "")
	return viper.WriteConfig()
}

func SaveUserId(id uuid.UUID) error {
	viper.Set("user.id", id.String())

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(dir, "envcrypt")
	path := filepath.Join(appDir, "config.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return viper.WriteConfigAs(path)
	}

	return viper.WriteConfig()
}

func SaveRefreshToken(refreshToken string) error {
	viper.Set("user.refresh_token", refreshToken)
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	appDir := filepath.Join(dir, "envcrypt")
	path := filepath.Join(appDir, "config.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return viper.WriteConfigAs(path)
	}
	return viper.WriteConfig()
}

func RemoveUserId() error {
	viper.Set("user.id", "")
	return viper.WriteConfig()
}

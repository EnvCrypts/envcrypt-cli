package cryptoutils

import (
	"encoding/base64"

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

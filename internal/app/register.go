package app

import (
	"context"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

func (app *App) Register(ctx context.Context, email, password string) error {
	// Create KeyPair
	keypair, err := cryptoutils.GenerateKeyPair(password)
	if err != nil {
		return err
	}

	requestBody := config.CreateRequestBody{
		Email:                   email,
		Password:                password,
		PublicKey:               keypair.PublicKey,
		EncryptedUserPrivateKey: keypair.EncKey.EncryptedUserPrivateKey,
		PrivateKeySalt:          keypair.EncKey.PrivateKeySalt,
		PrivateKeyNonce:         keypair.EncKey.PrivateKeyNonce,
	}
	var responseBody config.CreateResponseBody

	err = app.HttpClient.Do(ctx, "POST", "/users/create", requestBody, &responseBody)
	if err != nil {
		return err
	}

	err = cryptoutils.SavePrivateKey(email, keypair.PrivateKey)
	if err != nil {
		return err
	}

	return nil
}

package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

func (app *App) Login(ctx context.Context, email, password string) error {

	requestBody := config.LoginRequestBody{
		Email:    email,
		Password: password,
	}

	var responseBody config.LoginResponseBody

	err := app.HttpClient.Do(ctx, "POST", "/users/login", requestBody, &responseBody)
	if err != nil {
		return err
	}

	encryptedKey := &config.EncryptedPrivateKey{
		EncryptedUserPrivateKey: responseBody.User.EncryptedUserPrivateKey,
		PrivateKeySalt:          responseBody.User.PrivateKeySalt,
		PrivateKeyNonce:         responseBody.User.PrivateKeyNonce,
	}

	decryptedPrivateKey, err := cryptoutils.DecryptPrivateKey(encryptedKey, password, &config.DefaultArgon2Params)
	if err != nil {
		return err
	}

	err = cryptoutils.SavePrivateKey(email, decryptedPrivateKey)
	if err != nil {
		return err
	}

	err = cryptoutils.SaveUserEmail(email)
	if err != nil {
		return err
	}

	err = cryptoutils.SaveUserId(responseBody.User.Id)
	if err != nil {
		return err
	}

	return nil
}

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

	err = cryptoutils.SaveUserEmail(email)
	if err != nil {
		return err
	}

	err = cryptoutils.SaveUserId(responseBody.User.Id)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Logout(email string) error {
	var errs []error

	if err := cryptoutils.DeletePrivateKey(email); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveUserEmail(); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveUserId(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

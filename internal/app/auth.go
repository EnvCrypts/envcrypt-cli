package app

import (
	"context"
	"crypto/ecdh"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) Login(ctx context.Context, email, password string) error {

	requestBody := config.LoginRequestBody{
		Email:    email,
		Password: password,
	}

	var responseBody config.LoginResponseBody

	err := app.HttpClient.Do(ctx, "POST", "/users/login", requestBody, &responseBody, false)
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

	err = cryptoutils.SaveAccessToken(responseBody.Session.AccessToken.String())
	if err != nil {
		return err
	}

	err = cryptoutils.SaveRefreshToken(responseBody.Session.RefreshToken.String())
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Register(ctx context.Context, email, password string) (*config.KeyPair, error) {
	// Create KeyPair
	keypair, err := cryptoutils.GenerateKeyPair(password)
	if err != nil {
		return nil, err
	}

	requestBody := config.CreateRequestBody{
		Email:                   email,
		Password:                password,
		PublicKey:               keypair.PublicKey,
		EncryptedUserPrivateKey: keypair.EncKey.EncryptedUserPrivateKey,
		PrivateKeySalt:          keypair.EncKey.PrivateKeySalt,
		PrivateKeyNonce:         keypair.EncKey.PrivateKeyNonce,

		RecoveryPrivateKey: keypair.RecoveryEncKey.EncryptedUserPrivateKey,
		RecoverySalt:       keypair.RecoveryEncKey.PrivateKeySalt,
		RecoveryNonce:      keypair.RecoveryEncKey.PrivateKeyNonce,
	}
	var responseBody config.CreateResponseBody

	err = app.HttpClient.Do(ctx, "POST", "/users/create", requestBody, &responseBody, false)
	if err != nil {
		return nil, err
	}

	err = cryptoutils.SavePrivateKey(email, keypair.PrivateKey)
	if err != nil {
		return nil, err
	}

	err = cryptoutils.SaveUserEmail(email)
	if err != nil {
		return nil, err
	}

	err = cryptoutils.SaveUserId(responseBody.User.Id)
	if err != nil {
		return nil, err
	}

	err = cryptoutils.SaveAccessToken(responseBody.Session.AccessToken.String())
	if err != nil {
		return nil, err
	}

	err = cryptoutils.SaveRefreshToken(responseBody.Session.RefreshToken.String())
	if err != nil {
		return nil, err
	}

	return keypair, nil
}

func (app *App) Logout(ctx context.Context, email string) error {
	var errs []error

	userId := viper.GetString("user.id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	var requestBody = config.LogoutRequestBody{
		UserID: uid,
	}
	var responseBody config.LogoutResponseBody
	err = app.HttpClient.Do(ctx, "POST", "/users/logout", requestBody, &responseBody, false)
	if err != nil {
		return err
	}

	if err := cryptoutils.DeletePrivateKey(email); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveUserEmail(); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveUserId(); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveAccessToken(); err != nil {
		errs = append(errs, err)
	}

	if err := cryptoutils.RemoveRefreshToken(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (app *App) RecoverInit(ctx context.Context, email string) (*config.RecoveryInitResponse, error) {
	requestBody := config.RecoveryInitRequest{Email: email}
	var responseBody config.RecoveryInitResponse

	err := app.HttpClient.Do(ctx, "POST", "/users/recovery/init", requestBody, &responseBody, false)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}

func (app *App) RecoverComplete(ctx context.Context, email, newPassword string, plaintextPrivateKey []byte) error {

	curve := ecdh.X25519()
	privateKey, err := curve.NewPrivateKey(plaintextPrivateKey)
	if err != nil {
		return err
	}

	encryptedNewKey, err := cryptoutils.EncryptPrivateKey(privateKey, newPassword, &config.DefaultArgon2Params)
	if err != nil {
		return err
	}

	completeReq := config.RecoveryCompleteRequest{
		Email:                   email,
		Password:                newPassword,
		EncryptedUserPrivateKey: encryptedNewKey.EncryptedUserPrivateKey,
		PrivateKeySalt:          encryptedNewKey.PrivateKeySalt,
		PrivateKeyNonce:         encryptedNewKey.PrivateKeyNonce,
	}

	var responseBody config.RecoveryCompleteResponse
	err = app.HttpClient.Do(ctx, "POST", "/users/recovery/complete", completeReq, &responseBody, false)
	if err != nil {
		return err
	}

	return nil
}

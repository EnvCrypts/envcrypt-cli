package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) PushEnv(ctx context.Context, projectName, envName string, envMap map[string]string) error {

	// Getting Local Information
	userEmail, userId := viper.GetString("user.email"), viper.GetString("user.id")
	if userEmail == "" || userId == "" {
		return errors.New("missing user email or user id")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	privateKey, err := cryptoutils.LoadPrivateKey(userEmail)
	if err != nil {
		return err
	}

	// Get Project ID
	projectRequest := config.GetMemberProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResponse config.GetMemberProjectResponse
	err = app.HttpClient.Do(ctx, "POST", "/projects/get", projectRequest, &projectResponse)
	if err != nil {
		return err
	}

	// Getting Wrapped Keys
	keyRequest := config.GetUserProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var keyResponse config.GetUserProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/keys", keyRequest, &keyResponse); err != nil {
		return errors.New("could not get project keys")
	}

	wrappedKey := &cryptoutils.WrappedKey{
		WrappedPMK:       keyResponse.WrappedPMK,
		WrapNonce:        keyResponse.WrapNonce,
		WrapEphemeralPub: keyResponse.EphemeralPublicKey,
	}

	data, err := cryptoutils.PrepareEnvForStorage(envMap)
	if err != nil {
		return errors.New("could not prepare environment variables")
	}

	pmk, err := cryptoutils.UnwrapPMK(wrappedKey, privateKey)
	if err != nil {
		return errors.New("could not unwrap private key")
	}

	// encrypt using pmk and store the nonce, ciphertext
	encryptedData, nonce, err := cryptoutils.EncryptENV(pmk, data)
	if err != nil {
		return errors.New("could not encrypt data")
	}

	metadata := config.Metadata{
		Type: "env_created",
	}
	createRequest := config.AddEnvRequest{
		ProjectId:  projectResponse.ProjectId,
		UserId:     uid,
		EnvName:    envName,
		CipherText: encryptedData,
		Nonce:      nonce,
		Metadata:   metadata,
	}

	var createResponse config.AddEnvResponse
	if err := app.HttpClient.Do(ctx, "POST", "/env/create", createRequest, &createResponse); err != nil {
		return err
	}

	return nil
}

func (app *App) PullEnv(ctx context.Context, projectName, envName string) (map[string]string, error) {

	userEmail, userId := viper.GetString("user.email"), viper.GetString("user.id")
	if userEmail == "" || userId == "" {
		return nil, errors.New("missing user email or user id")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	userPriv, err := cryptoutils.LoadPrivateKey(userEmail)
	if err != nil {
		return nil, err
	}

	projectRequest := config.GetMemberProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResponse config.GetMemberProjectResponse
	err = app.HttpClient.Do(ctx, "POST", "/projects/get", projectRequest, &projectResponse)
	if err != nil {
		return nil, err
	}

	envRequest := config.GetEnvRequest{
		ProjectId: projectResponse.ProjectId,
		Email:     userEmail,
		EnvName:   envName,
		Version:   nil,
	}

	var envResponse config.GetEnvResponse
	err = app.HttpClient.Do(ctx, "POST", "/env/search", envRequest, &envResponse)
	if err != nil {
		return nil, err
	}

	// Getting Wrapped Keys
	keyRequest := config.GetUserProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var keyResponse config.GetUserProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/keys", keyRequest, &keyResponse); err != nil {
		return nil, errors.New("could not get project keys")
	}

	wrappedKey := &cryptoutils.WrappedKey{
		WrappedPMK:       keyResponse.WrappedPMK,
		WrapNonce:        keyResponse.WrapNonce,
		WrapEphemeralPub: keyResponse.EphemeralPublicKey,
	}

	pmk, err := cryptoutils.UnwrapPMK(wrappedKey, userPriv)
	if err != nil {
		return nil, errors.New("could not unwrap private key")
	}

	envBytes, err := cryptoutils.DecryptENV(pmk, envResponse.CipherText, envResponse.Nonce)
	if err != nil {
		return nil, errors.New("could not decrypt data")
	}

	envMap, err := cryptoutils.ReadCompressedEnv(envBytes)
	if err != nil {
		return nil, errors.New("could not parse environment variables")
	}

	return envMap, nil
}

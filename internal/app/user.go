package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) AddUserToProject(ctx context.Context, memberEmail, projectName, role string) error {

	adminEmail, adminId := viper.GetString("user.email"), viper.GetString("user.id")
	uid, err := uuid.Parse(adminId)
	if err != nil || uid == uuid.Nil {
		return errors.New("user not authenticated")
	}

	projectReq := config.GetUserProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResp config.GetUserProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/keys", projectReq, &projectResp); err != nil {
		return errors.New("could not get project keys")
	}

	wrappedKey := &cryptoutils.WrappedKey{
		WrappedPMK:       projectResp.WrappedPMK,
		WrapNonce:        projectResp.WrapNonce,
		WrapEphemeralPub: projectResp.EphemeralPublicKey,
	}
	privateKey, err := cryptoutils.LoadPrivateKey(adminEmail)
	if err != nil {
		return errors.New("user not authenticated")
	}

	pmk, err := cryptoutils.UnwrapPMK(wrappedKey, privateKey)
	if err != nil {
		return errors.New("forbidden access")
	}

	// Get Member's publicKey
	pubKeyReq := config.UserKeyRequestBody{
		Email: memberEmail,
	}

	var pubKeyResp config.UserKeyResponseBody
	if err := app.HttpClient.Do(ctx, "POST", "/users/search", pubKeyReq, &pubKeyResp); err != nil {
		return errors.New("user not found")
	}

	// Wrap Key for user
	memberWrappedKey, err := cryptoutils.WrapPMKForUser(pmk, pubKeyResp.PublicKey)
	if err != nil {
		return errors.New("unable to wrap user key")
	}

	addReq := config.AddUserToProjectRequest{
		ProjectName:        projectName,
		UserId:             pubKeyResp.UserId,
		AdminId:            uid,
		Role:               role,
		WrappedPMK:         memberWrappedKey.WrappedPMK,
		WrapNonce:          memberWrappedKey.WrapNonce,
		EphemeralPublicKey: memberWrappedKey.WrapEphemeralPub,
	}

	var addResp config.AddUserToProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/addUser", addReq, &addResp); err != nil {
		return err
	}

	return nil
}

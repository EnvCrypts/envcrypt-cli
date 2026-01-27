package app

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) CreateProject(ctx context.Context, projectName string) error {
	email := viper.GetString("user.email")
	if email == "" {
		return errors.New("no user email found")
	}

	userReq := config.UserKeyRequestBody{
		Email: email,
	}

	var userResp config.UserKeyResponseBody
	if err := app.HttpClient.Do(ctx, "POST", "/users/search", userReq, &userResp); err != nil {
		return err
	}

	pmk := make([]byte, 32)
	if _, err := rand.Read(pmk); err != nil {
		return err
	}

	wrappedKey, err := cryptoutils.WrapPMKForUser(pmk, userResp.PublicKey)
	if err != nil {
		return err
	}

	projectReq := config.ProjectCreateRequest{
		Name:               projectName,
		UserId:             userResp.UserId,
		WrappedPMK:         wrappedKey.WrappedPMK,
		WrapNonce:          wrappedKey.WrapNonce,
		EphemeralPublicKey: wrappedKey.WrapEphemeralPub,
	}

	var projectResp config.ProjectCreateResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/create", projectReq, &projectResp); err != nil {
		return err
	}

	return nil
}

func (app *App) DeleteProject(ctx context.Context, projectName string) error {
	email, userId := viper.GetString("user.email"), viper.GetString("user.id")

	uid, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	if email == "" || uid == uuid.Nil {
		return errors.New("user not authenticated")
	}

	deleteReq := config.ProjectDeleteRequest{
		ProjectName: projectName,
		UserId:      uid,
	}
	var deleteResp config.ProjectDeleteResponse

	if err := app.HttpClient.Do(ctx, "POST", "/projects/delete", deleteReq, &deleteResp); err != nil {
		return err
	}

	return nil
}

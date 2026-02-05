package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) ListServiceRoles(ctx context.Context) ([]config.ServiceRole, error) {

	userID := viper.GetString("user.id")
	if userID == "" {
		return []config.ServiceRole{}, errors.New("user id not found")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return []config.ServiceRole{}, errors.New("user id not valid")
	}

	var requestBody = config.ServiceRoleListRequest{
		CreatedBy: uid,
	}

	var responseBody config.ServiceRoleListResponse
	err = app.HttpClient.Do(ctx, "POST", "/service_role/get/all", requestBody, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.ServiceRoles, nil
}

func (app *App) CreateServiceRole(ctx context.Context, name, repoPrincipal string) (*config.ServiceRoleKeyPair, error) {

	userID := viper.GetString("user.id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("user id not valid")
	}

	keypair, err := cryptoutils.GenerateServiceRoleKeyPair()
	if err != nil {
		return nil, err
	}

	var requestBody = config.ServiceRoleCreateRequest{
		ServiceRoleName:      name,
		RepoPrincipal:        repoPrincipal,
		ServiceRolePublicKey: keypair.PublicKey,
		CreatedBy:            uid,
	}

	var responseBody config.ServiceRoleCreateResponse
	err = app.HttpClient.Do(ctx, "POST", "/service_role/create", requestBody, responseBody)
	if err != nil {
		return nil, err
	}

	return keypair, nil
}

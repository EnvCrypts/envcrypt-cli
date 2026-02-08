package app

import (
	"context"
	"errors"
	"log"

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
	err = app.HttpClient.Do(ctx, "POST", "/service_role/get/all", requestBody, &responseBody, true)
	if err != nil {
		log.Print(err.Error())
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
	err = app.HttpClient.Do(ctx, "POST", "/service_role/create", requestBody, &responseBody, true)
	if err != nil {
		return nil, err
	}

	return keypair, nil
}

func (app *App) GetServiceRole(ctx context.Context, repoPrincipal string) (*config.ServiceRole, error) {
	var requestBody = config.ServiceRoleGetRequest{
		RepoPrincipal: repoPrincipal,
	}

	var responseBody config.ServiceRoleGetResponse
	err := app.HttpClient.Do(ctx, "POST", "/service_role/get", requestBody, &responseBody, true)
	if err != nil {
		return nil, err
	}

	return &responseBody.ServiceRole, nil
}

func (app *App) DeleteServiceRole(ctx context.Context, serviceRoleId uuid.UUID) error {

	userID := viper.GetString("user.id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("user id not valid")
	}

	var requestBody = config.ServiceRoleDeleteRequest{
		ServiceRoleId: serviceRoleId,
		CreatedBy:     uid,
	}
	var responseBody config.ServiceRoleDeleteResponse
	err = app.HttpClient.Do(ctx, "POST", "/service_role/delete", requestBody, &responseBody, true)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetPermissions(ctx context.Context, repoPrincipal string) (*config.ServiceRolePermsResponse, error) {
	var requestBody = config.ServiceRolePermsRequest{
		RepoPrincipal: repoPrincipal,
	}

	var responseBody config.ServiceRolePermsResponse
	err := app.HttpClient.Do(ctx, "POST", "/service_role/perms", requestBody, &responseBody, true)
	if err != nil {
		return nil, err
	}
	return &responseBody, nil
}

func (app *App) DelegateAccess(ctx context.Context, repoPrincipal, projectName, env string) error {
	adminEmail, adminId := viper.GetString("user.email"), viper.GetString("user.id")
	uid, err := uuid.Parse(adminId)
	if err != nil || uid == uuid.Nil {
		return errors.New("user not authenticated")
	}

	// 1. Get Project Keys to get WrappedPMK
	projectReq := config.GetUserProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResp config.GetUserProjectResponse
	if err := app.HttpClient.Do(ctx, "POST", "/projects/keys", projectReq, &projectResp, true); err != nil {
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
		return errors.New("forbidden access: cannot unwrap project key")
	}

	// 2. Get Service Role to get its Public Key and ID
	role, err := app.GetServiceRole(ctx, repoPrincipal)
	if err != nil {
		return err
	}

	// 3. Wrap PMK for Service Role
	serviceRoleWrappedKey, err := cryptoutils.WrapPMKForUser(pmk, role.ServiceRolePublicKey)
	if err != nil {
		return errors.New("unable to wrap key for service role")
	}

	// 4. Delegate Access
	delegateReq := config.ServiceRoleDelegateRequest{
		RepoPrincipal:      repoPrincipal,
		ProjectId:          projectResp.ProjectId,
		EnvName:            env,
		WrappedPMK:         serviceRoleWrappedKey.WrappedPMK,
		WrapNonce:          serviceRoleWrappedKey.WrapNonce,
		EphemeralPublicKey: serviceRoleWrappedKey.WrapEphemeralPub,
		DelegatedBy:        uid,
	}

	var delegateResp config.ServiceRoleDelegateResponse
	if err := app.HttpClient.Do(ctx, "POST", "/service_role/delegate", delegateReq, &delegateResp, true); err != nil {
		return err
	}

	return nil
}

func (app *App) GetServiceRoleProjectKeys(ctx context.Context, projectID, sessionID uuid.UUID, env string) (*config.ServiceRollProjectKeyResponse, error) {

	var requestBody = config.ServiceRollProjectKeyRequest{
		ProjectID: projectID,
		Env:       env,
		SessionID: sessionID,
	}
	var responseBody config.ServiceRollProjectKeyResponse
	err := app.HttpClient.Do(ctx, "POST", "/service_role/project-keys", requestBody, &responseBody, false)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (app *App) GetSessionID(ctx context.Context, oidcToken, projectName, env string) (*uuid.UUID, *uuid.UUID, error) {

	userID := viper.GetString("user.id")
	if userID == "" {
		return nil, nil, errors.New("user not authenticated")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, errors.New("user not authenticated")
	}

	var projectRequest = config.GetMemberProjectRequest{
		ProjectName: projectName,
		UserId:      uid,
	}

	var projectResponse config.GetMemberProjectResponse
	err = app.HttpClient.Do(ctx, "POST", "/projects/get", projectRequest, &projectResponse)
	if err != nil {
		return nil, nil, err
	}

	var githubOIDCRequest = config.GithubOIDCLoginRequest{
		ProjectID: projectResponse.ProjectId,
		Env:       env,
		IDToken:   oidcToken,
	}
	var githubOIDCResponse config.GithubOIDCLoginResponse
	err = app.HttpClient.Do(ctx, "POST", "/oidc/github", githubOIDCRequest, &githubOIDCResponse)
	if err != nil {
		return nil, nil, err
	}

	return &githubOIDCResponse.SessionID, &projectResponse.ProjectId, nil
}

func (app *App) PullEnvForCI(ctx context.Context, projectID uuid.UUID, envName string, pmk []byte) (map[string]string, error) {
	envRequest := config.GetEnvForCIRequest{
		ProjectId: projectID,
		EnvName:   envName,
	}

	var envResponse config.GetEnvResponse
	err := app.HttpClient.Do(ctx, "POST", "/env/ci/search", envRequest, &envResponse)
	if err != nil {
		return nil, err
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

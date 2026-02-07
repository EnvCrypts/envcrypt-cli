package app

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
)

func (app *App) GetSessionID(ctx context.Context, oidcToken string) (*uuid.UUID, *uuid.UUID, error) {

	var githubOIDCRequest = config.GithubOIDCLoginRequest{
		IDToken: oidcToken,
	}
	var githubOIDCResponse config.GithubOIDCLoginResponse
	err := app.HttpClient.Do(ctx, "POST", "/oidc/github", githubOIDCRequest, &githubOIDCResponse, false)
	if err != nil {
		return nil, nil, err
	}

	return &githubOIDCResponse.SessionID, &githubOIDCResponse.ProjectID, nil
}

func (app *App) PullEnvForCI(ctx context.Context, projectID uuid.UUID, envName string, pmk []byte) (map[string]string, error) {
	envRequest := config.GetEnvForCIRequest{
		ProjectId: projectID,
		EnvName:   envName,
	}

	var envResponse config.GetEnvResponse
	err := app.HttpClient.Do(ctx, "POST", "/env/ci/search", envRequest, &envResponse, true)
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

package config

import "github.com/google/uuid"

// ServiceRollProjectKeyRequest POST /service_role/project-key
type ServiceRollProjectKeyRequest struct {
	ProjectID uuid.UUID `json:"project_id"`
	SessionID uuid.UUID `json:"session_id"`
	Env       string    `json:"env"`
}
type ServiceRollProjectKeyResponse struct {
	ProjectId          uuid.UUID `json:"project_id"`
	WrappedPMK         []byte    `json:"wrapped_pmk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

// GithubOIDCLoginRequest POST /oidc/github
type GithubOIDCLoginRequest struct {
	IDToken   string    `json:"id_token"`
	ProjectID uuid.UUID `json:"project_id"`
	Env       string    `json:"env"`
}
type GithubOIDCLoginResponse struct {
	SessionID uuid.UUID `json:"session_id"`
}

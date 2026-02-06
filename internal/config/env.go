package config

import "github.com/google/uuid"

type Metadata struct {
	Type string `json:"type"`
}
type AddEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	UserId    uuid.UUID `json:"user_id"`

	EnvName    string `json:"env_name"`
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`

	Metadata Metadata `json:"metadata"`
}

type AddEnvResponse struct {
	Message string `json:"message"`
}

type GetEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
	Version *int32 `json:"version"`
}

type GetEnvResponse struct {
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
}

type GetEnvVersionsRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
}

type EnvResponse struct {
	CipherText []byte   `json:"cipher_text"`
	Nonce      []byte   `json:"nonce"`
	Version    int32    `json:"version"`
	Metadata   Metadata `json:"metadata"`
}
type GetEnvVersionsResponse struct {
	EnvVersions []EnvResponse `json:"env_versions"`
}

type GetEnvForCIRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	EnvName   string    `json:"env_name"`
}
type GetEnvForCIResponse struct {
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
}


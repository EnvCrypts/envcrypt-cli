package config

import "github.com/google/uuid"

type ProjectCreateRequest struct {
	Name               string    `json:"name"`
	UserId             uuid.UUID `json:"user_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type ProjectCreateResponse struct {
	Message string `json:"message"`
}

type ListProjectRequest struct {
	UserId uuid.UUID `json:"user_id"`
}
type Project struct {
	Id        uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsRevoked bool      `json:"is_revoked"`
}

type ListProjectResponse struct {
	Projects []Project `json:"projects"`
}

type ProjectDeleteRequest struct {
	ProjectName string    `json:"project_name"`
	UserId      uuid.UUID `json:"user_id"`
}

type ProjectDeleteResponse struct {
	Message string `json:"message"`
}

type AddUserToProjectRequest struct {
	ProjectName        string    `json:"project_name"`
	AdminId            uuid.UUID `json:"admin_id"`
	UserId             uuid.UUID `json:"user_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}
type AddUserToProjectResponse struct {
	Message string `json:"message"`
}

type SetAccessRequest struct {
	ProjectName string    `json:"project_name"`
	UserEmail   string    `json:"user_email"`
	AdminId     uuid.UUID `json:"admin_id"`
	IsRevoked   bool      `json:"is_revoked"`
}
type SetAccessResponse struct {
	Message string `json:"message"`
}

type GetUserProjectRequest struct {
	ProjectName string    `json:"project_name"`
	UserId      uuid.UUID `json:"user_id"`
}

type GetUserProjectResponse struct {
	ProjectId          uuid.UUID `json:"project_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type GetMemberProjectRequest struct {
	ProjectName string    `json:"project_name"`
	UserId      uuid.UUID `json:"user_id"`
}

type GetMemberProjectResponse struct {
	ProjectId          uuid.UUID `json:"project_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type GetProjectByRepo struct {
	RepoPrincipal string `json:"repo_principal"`
}
type GetProjectByNameResponse struct {
	ProjectID uuid.UUID `json:"project_id"`
}

type WrappedKey struct {
	UserID             uuid.UUID `json:"user_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type MemberPublicKey struct {
	UserID    uuid.UUID `json:"user_id"`
	PublicKey []byte    `json:"public_key"`
}

type WrappedDEK struct {
	EnvVersionID uuid.UUID `json:"env_version_id"`
	WrappedDEK   []byte    `json:"wrapped_dek"`
	DekNonce     []byte    `json:"dek_nonce"`
}

type RotateInitRequest struct {
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type RotateInitResponse struct {
	WrappedPRKs      []WrappedKey      `json:"wrapped_prks"`
	WrappedDEKs      []WrappedDEK      `json:"wrapped_deks"`
	MemberPublicKeys []MemberPublicKey `json:"member_public_keys"`
	PRKVersion       int32             `json:"prk_version"`
}

type NewWrappedDEK struct {
	EnvVersionID  uuid.UUID `json:"env_version_id"`
	NewWrappedDEK []byte    `json:"new_wrapped_dek"`
	NewDekNonce   []byte    `json:"new_dek_nonce"`
}

type RotateCommitRequest struct {
	ProjectID          uuid.UUID       `json:"project_id"`
	UserID             uuid.UUID       `json:"user_id"`
	ExpectedPRKVersion int32           `json:"expected_prk_version"`
	NewWrappedPRKs     []WrappedKey    `json:"new_wrapped_prks"`
	NewWrappedDEKs     []NewWrappedDEK `json:"new_wrapped_deks"`
}

type RotateCommitResponse struct {
	NewPRKVersion int32 `json:"new_prk_version"`
}

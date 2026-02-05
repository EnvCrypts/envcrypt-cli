package config

import (
	"time"

	"github.com/google/uuid"
)

type ServiceRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	ServiceRolePublicKey []byte `json:"service_role_public_key"`
	RepoPrincipal        string `json:"repo_principal"`

	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
type ServiceRoleListRequest struct {
	CreatedBy uuid.UUID `json:"created_by"`
}
type ServiceRoleListResponse struct {
	ServiceRoles []ServiceRole `json:"services"`
}

// ServiceRoleCreateRequest POST /service_role/create
type ServiceRoleCreateRequest struct {
	ServiceRoleName string `json:"service_role_name"`

	ServiceRolePublicKey []byte `json:"service_role_public_key"`

	RepoPrincipal string    `json:"repo_principal"`
	CreatedBy     uuid.UUID `json:"created_by"`
}
type ServiceRoleCreateResponse struct {
	Message     string      `json:"message"`
	ServiceRole ServiceRole `json:"service_role"`
}

// ServiceRoleGetRequest POST /service_role/get
type ServiceRoleGetRequest struct {
	RepoPrincipal string `json:"repo_principal"`
}
type ServiceRoleGetResponse struct {
	ServiceRole ServiceRole `json:"service_role"`
	Message     string      `json:"message"`
}

// ServiceRoleDeleteRequest POST /service_role/delete
type ServiceRoleDeleteRequest struct {
	ServiceRoleId uuid.UUID `json:"service_role_id"`
	CreatedBy     uuid.UUID `json:"created_by"`
}
type ServiceRoleDeleteResponse struct {
	Message string `json:"message"`
}

// ServiceRolePermsRequest POST /service_role/perms
type ServiceRolePermsRequest struct {
	RepoPrincipal string `json:"repo_principal"`
}
type ServiceRolePermsResponse struct {
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Env         string    `json:"env"`
}

// ServiceRoleDelegateRequest POST /service_role/delegate
type ServiceRoleDelegateRequest struct {
	RepoPrincipal string `json:"repo_principal"`

	ProjectId uuid.UUID `json:"project_id"`
	EnvName   string    `json:"env_name"`

	WrappedPMK         []byte `json:"wrapped_pmk"`
	WrapNonce          []byte `json:"wrap_nonce"`
	EphemeralPublicKey []byte `json:"ephemeral_public_key"`

	DelegatedBy uuid.UUID `json:"delegated_by"`
}
type ServiceRoleDelegateResponse struct {
	Message string `json:"message"`
}

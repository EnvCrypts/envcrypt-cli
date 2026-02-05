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

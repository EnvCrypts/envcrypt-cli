package config

import (
	"time"

	"github.com/google/uuid"
)

type SnapshotProjectMetadata struct {
	Name       string `json:"name"`
	PrkVersion int32  `json:"prk_version"`
}

type SnapshotMember struct {
	UserID             uuid.UUID `json:"user_id"`
	WrappedPRK         []byte    `json:"wrapped_prk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type SnapshotEnvVersion struct {
	EnvVersionID      uuid.UUID `json:"env_version_id"`
	EnvName           string    `json:"env_name"`
	Version           int32     `json:"version"`
	Ciphertext        []byte    `json:"ciphertext"`
	Nonce             []byte    `json:"nonce"`
	WrappedDEK        []byte    `json:"wrapped_dek"`
	DekNonce          []byte    `json:"dek_nonce"`
	EncryptionVersion int32     `json:"encryption_version"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         uuid.UUID `json:"created_by"`
	Metadata          []byte    `json:"metadata"`
}

type Snapshot struct {
	Metadata    SnapshotProjectMetadata `json:"metadata"`
	Members     []SnapshotMember        `json:"members"`
	EnvVersions []SnapshotEnvVersion    `json:"env_versions"`
}

type SnapshotExportRequest struct {
	ProjectName string    `json:"project_name"`
	UserID      uuid.UUID `json:"user_id"`
}

type SnapshotExportResponse struct {
	Snapshot Snapshot `json:"snapshot"`
	Checksum string   `json:"checksum"`
}

type SnapshotImportRequest struct {
	NewProjectName string    `json:"new_project_name"`
	UserID         uuid.UUID `json:"user_id"`
	Snapshot       Snapshot  `json:"snapshot"`
	Checksum       string    `json:"checksum"`
}

type SnapshotImportResponse struct {
	NewProjectID uuid.UUID `json:"new_project_id"`
}

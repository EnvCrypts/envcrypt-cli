package config

import (
	"github.com/google/uuid"
)

type CreateRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	PublicKey               []byte `json:"public_key"`
	EncryptedUserPrivateKey []byte `json:"encrypted_user_private_key"`
	PrivateKeySalt          []byte `json:"private_key_salt"`
	PrivateKeyNonce         []byte `json:"private_key_nonce"`
}

type CreateResponseBody struct {
	Message string `json:"message"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserBody struct {
	Id                      uuid.UUID      `json:"id"`
	Email                   string         `json:"email"`
	PublicKey               []byte         `json:"public_key"`
	EncryptedUserPrivateKey []byte         `json:"encrypted_user_private_key"`
	PrivateKeySalt          []byte         `json:"private_key_salt"`
	PrivateKeyNonce         []byte         `json:"private_key_nonce"`
	ArgonParams             Argon2idParams `json:"argon_params"`
}
type LoginResponseBody struct {
	Message string   `json:"message"`
	User    UserBody `json:"user"`
}

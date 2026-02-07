package config

import "github.com/google/uuid"

type UserBody struct {
	Id                      uuid.UUID      `json:"id"`
	Email                   string         `json:"email"`
	PublicKey               []byte         `json:"public_key"`
	EncryptedUserPrivateKey []byte         `json:"encrypted_user_private_key"`
	PrivateKeySalt          []byte         `json:"private_key_salt"`
	PrivateKeyNonce         []byte         `json:"private_key_nonce"`
	ArgonParams             Argon2idParams `json:"argon_params"`
}
type SessionBody struct {
	AccessToken  uuid.UUID `json:"access_token"`
	RefreshToken uuid.UUID `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
}
type CreateRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	PublicKey               []byte `json:"public_key"`
	EncryptedUserPrivateKey []byte `json:"encrypted_user_private_key"`
	PrivateKeySalt          []byte `json:"private_key_salt"`
	PrivateKeyNonce         []byte `json:"private_key_nonce"`
}
type CreateResponseBody struct {
	Message string      `json:"message"`
	User    UserBody    `json:"user"`
	Session SessionBody `json:"session"`
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponseBody struct {
	Message string      `json:"message"`
	User    UserBody    `json:"user"`
	Session SessionBody `json:"session"`
}

type UserKeyRequestBody struct {
	Email string `json:"email"`
}
type UserKeyResponseBody struct {
	Message   string    `json:"message"`
	UserId    uuid.UUID `json:"user_id"`
	PublicKey []byte    `json:"public_key"`
}

type RefreshRequestBody struct {
	UserID uuid.UUID `json:"user_id"`
}
type RefreshResponseBody struct {
	Message string      `json:"message"`
	Session SessionBody `json:"session"`
}

type LogoutRequestBody struct {
	UserID uuid.UUID `json:"user_id"`
}
type LogoutResponseBody struct {
	Message string `json:"message"`
}

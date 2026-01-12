package config

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

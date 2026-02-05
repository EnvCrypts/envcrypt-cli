package config

type Argon2idParams struct {
	Time        uint32 `json:"time"`
	Memory      uint32 `json:"memory"`
	Parallelism uint8  `json:"parallelism"`
	KeyLength   uint32 `json:"key_length"`
}

var DefaultArgon2Params = Argon2idParams{
	Time:        3,
	Memory:      64 * 1024,
	Parallelism: 1,
	KeyLength:   32,
}

type KeyPair struct {
	PublicKey  []byte              `json:"public_key"`
	PrivateKey []byte              `json:"private_key"`
	EncKey     EncryptedPrivateKey `json:"encrypted_private_key"`
}

type ServiceRoleKeyPair struct {
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}
type EncryptedPrivateKey struct {
	EncryptedUserPrivateKey []byte `json:"encrypted_user_private_key"`
	PrivateKeySalt          []byte `json:"private_key_salt"`
	PrivateKeyNonce         []byte `json:"private_key_nonce"`
}

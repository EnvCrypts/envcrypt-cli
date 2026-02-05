package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"golang.org/x/crypto/argon2"
)

func EncryptPrivateKey(privateKey *ecdh.PrivateKey, password string, argonParams *config.Argon2idParams) (*config.EncryptedPrivateKey, error) {

	// Generating Salt for Argon using crypto/rand
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	// Create EncryptionKey with the user password
	encryptionKey := argon2.IDKey(
		[]byte(password),
		salt,
		argonParams.Time,
		argonParams.Memory,
		argonParams.Parallelism,
		argonParams.KeyLength,
	)

	// Create AES block cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Create GCM instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create AES nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	// AES_GCM encryption for the private key
	encryptedPrivateKey := gcm.Seal(nil, nonce, privateKey.Bytes(), nil)

	return &config.EncryptedPrivateKey{
		EncryptedUserPrivateKey: encryptedPrivateKey,
		PrivateKeySalt:          salt,
		PrivateKeyNonce:         nonce,
	}, nil
}

func DecryptPrivateKey(
	encryptedPrivateKey *config.EncryptedPrivateKey,
	password string,
	argonParams *config.Argon2idParams,
) ([]byte, error) {

	// Derive the same encryption key using Argon2id
	encryptionKey := argon2.IDKey(
		[]byte(password),
		encryptedPrivateKey.PrivateKeySalt,
		argonParams.Time,
		argonParams.Memory,
		argonParams.Parallelism,
		argonParams.KeyLength,
	)

	// Create AES block cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Create GCM instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 4. Decrypt (Open = authenticated decrypt)
	plaintextPrivateKey, err := gcm.Open(
		nil,
		encryptedPrivateKey.PrivateKeyNonce,
		encryptedPrivateKey.EncryptedUserPrivateKey,
		nil,
	)
	if err != nil {
		// This error covers:
		// - wrong password
		// - corrupted ciphertext
		// - wrong nonce
		// - tampering
		return nil, err
	}

	// X25519 private keys are 32 bytes
	if len(plaintextPrivateKey) != 32 {
		return nil, errors.New("invalid private key length")
	}

	return plaintextPrivateKey, nil
}

func GenerateKeyPair(password string) (*config.KeyPair, error) {
	curve := ecdh.X25519()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := EncryptPrivateKey(privateKey, password, &config.DefaultArgon2Params)
	if err != nil {
		return nil, err
	}

	return &config.KeyPair{
		PrivateKey: privateKey.Bytes(),
		PublicKey:  privateKey.PublicKey().Bytes(),
		EncKey:     *encryptedKey,
	}, nil
}

func GenerateServiceRoleKeyPair() (*config.ServiceRoleKeyPair, error) {
	curve := ecdh.X25519()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &config.ServiceRoleKeyPair{
		PrivateKey: privateKey.Bytes(),
		PublicKey:  privateKey.PublicKey().Bytes(),
	}, nil
}

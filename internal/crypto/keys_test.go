package cryptoutils

import (
	"crypto/ecdh"
	"crypto/rand"
	"testing"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyCrypto(t *testing.T) {
	t.Run("PrivateKeyEncryption", func(t *testing.T) {
		curve := ecdh.X25519()
		privateKey, _ := curve.GenerateKey(rand.Reader)
		password := "secure_password_123"
		argonParams := &config.Argon2idParams{Time: 1, Memory: 16 * 1024, Parallelism: 1, KeyLength: 32}

		encryptedKey, err := EncryptPrivateKey(privateKey, password, argonParams)
		require.NoError(t, err)
		assert.NotEmpty(t, encryptedKey.EncryptedUserPrivateKey)

		decryptedBytes, err := DecryptPrivateKey(encryptedKey, password, argonParams)
		require.NoError(t, err)
		assert.Equal(t, privateKey.Bytes(), decryptedBytes)

		_, err = DecryptPrivateKey(encryptedKey, "wrong_password", argonParams)
		require.Error(t, err)
	})

	t.Run("KeyPairs", func(t *testing.T) {
		password := "my_test_password"

		keypair, err := GenerateKeyPair(password)
		require.NoError(t, err)
		assert.Len(t, keypair.PrivateKey, 32)
		assert.Len(t, keypair.PublicKey, 32)

		decryptedBytes, err := DecryptPrivateKey(&keypair.EncKey, password, &config.DefaultArgon2Params)
		require.NoError(t, err)
		assert.Equal(t, keypair.PrivateKey, decryptedBytes)

		assert.NotEmpty(t, keypair.RecoveryKey)
		assert.NotEmpty(t, keypair.RecoveryEncKey.EncryptedUserPrivateKey)
		assert.NotEmpty(t, keypair.RecoveryEncKey.PrivateKeySalt)
		assert.NotEmpty(t, keypair.RecoveryEncKey.PrivateKeyNonce)

		recoveryDecryptedBytes, err := DecryptPrivateKey(&keypair.RecoveryEncKey, keypair.RecoveryKey, &config.DefaultArgon2Params)
		require.NoError(t, err)
		assert.Equal(t, keypair.PrivateKey, recoveryDecryptedBytes)

		roleKeypair, err := GenerateServiceRoleKeyPair()
		require.NoError(t, err)
		assert.Len(t, roleKeypair.PrivateKey, 32)
	})
}

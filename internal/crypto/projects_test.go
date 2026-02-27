package cryptoutils

import (
	"testing"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCrypto(t *testing.T) {
	t.Run("PRKKeyExchange", func(t *testing.T) {
		alice, _ := GenerateEphemeralKeyPair()
		bob, _ := GenerateEphemeralKeyPair()

		secretAlice, err := X25519SharedSecret(alice.PrivateKey, bob.PublicKey)
		require.NoError(t, err)

		secretBob, _ := X25519SharedSecret(bob.PrivateKey, alice.PublicKey)
		assert.Equal(t, secretAlice, secretBob)

		_, err = X25519SharedSecret(alice.PrivateKey[:31], bob.PublicKey)
		require.Error(t, err)

		wrapKey, err := DeriveWrapKey(secretAlice)
		require.NoError(t, err)
		assert.Len(t, wrapKey, 32)

		prk := make([]byte, 32)
		wrapped, err := WrapPRKForUser(prk, alice.PublicKey)
		require.NoError(t, err)
		assert.NotEmpty(t, wrapped.WrappedPRK)

		unwrappedPRK, err := UnwrapPRK(wrapped, alice.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, prk, unwrappedPRK)
	})

	t.Run("EnvEncryption", func(t *testing.T) {
		prk := make([]byte, 32)
		data := []byte("SECRET_ENV_DATA")

		encrypted, nonce, err := EncryptENV(prk, data)
		require.NoError(t, err)
		assert.NotEqual(t, data, encrypted)

		decrypted, err := DecryptENV(prk, encrypted, nonce)
		require.NoError(t, err)
		assert.Equal(t, data, decrypted)
	})

	t.Run("DEKLifecycle", func(t *testing.T) {
		dek, _ := GenerateDEK()
		prk := make([]byte, 32)
		
		wrapped, nonce, err := WrapDEK(prk, dek)
		require.NoError(t, err)

		unwrapped, err := UnwrapDEK(prk, wrapped, nonce)
		require.NoError(t, err)
		assert.Equal(t, dek, unwrapped)

		newPRK := make([]byte, 32)
		newPRK[0] = 1

		wrappedDEKs := []config.WrappedDEK{{EnvVersionID: uuid.New(), WrappedDEK: wrapped, DekNonce: nonce}}
		rewrapped, err := RewrapDEKs(prk, newPRK, wrappedDEKs)
		require.NoError(t, err)

		reUnwrapped, err := UnwrapDEK(newPRK, rewrapped[0].NewWrappedDEK, rewrapped[0].NewDekNonce)
		require.NoError(t, err)
		assert.Equal(t, dek, reUnwrapped)
	})
}

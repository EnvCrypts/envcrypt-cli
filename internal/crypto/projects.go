package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
)

type EphemeralKeyPair struct {
	PrivateKey []byte // 32 bytes (destroy after use)
	PublicKey  []byte // 32 bytes (sent to server)
}

func GenerateEphemeralKeyPair() (*EphemeralKeyPair, error) {
	curve := ecdh.X25519()

	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &EphemeralKeyPair{
		PrivateKey: priv.Bytes(),
		PublicKey:  priv.PublicKey().Bytes(),
	}, nil
}

func X25519SharedSecret(
	privateKeyBytes []byte,
	peerPublicKeyBytes []byte,
) ([]byte, error) {

	if len(privateKeyBytes) != 32 {
		return nil, errors.New("invalid private key length")
	}
	if len(peerPublicKeyBytes) != 32 {
		return nil, errors.New("invalid public key length")
	}

	curve := ecdh.X25519()

	priv, err := curve.NewPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	pub, err := curve.NewPublicKey(peerPublicKeyBytes)
	if err != nil {
		return nil, err
	}

	// 32-byte shared secret
	return priv.ECDH(pub)
}

func DeriveWrapKey(sharedSecret []byte) ([]byte, error) {
	h := hkdf.New(
		sha256.New,
		sharedSecret,
		nil,
		[]byte("envcrypt-pmk-wrap"),
	)

	key := make([]byte, 32)
	if _, err := io.ReadFull(h, key); err != nil {
		return nil, err
	}

	return key, nil
}

type WrappedKey struct {
	WrappedPMK       []byte `json:"wrapped_pmk"`        // AES-GCM ciphertext
	WrapNonce        []byte `json:"wrap_nonce"`         // 12 bytes
	WrapEphemeralPub []byte `json:"wrap_ephemeral_pub"` // 32 bytes
}

func WrapPMKForUser(
	pmk []byte,
	recipientUserPublicKey []byte,
) (*WrappedKey, error) {

	if len(pmk) != 32 {
		return nil, errors.New("invalid PMK length")
	}
	if len(recipientUserPublicKey) != 32 {
		return nil, errors.New("invalid recipient public key length")
	}

	// 1. Generate ephemeral keypair
	ephemeral, err := GenerateEphemeralKeyPair()
	if err != nil {
		return nil, err
	}
	defer zero(ephemeral.PrivateKey)

	// 2. Derive shared secret
	sharedSecret, err := X25519SharedSecret(
		ephemeral.PrivateKey,
		recipientUserPublicKey,
	)
	if err != nil {
		return nil, err
	}

	// 3. Derive symmetric wrap key via HKDF
	wrapKey, err := DeriveWrapKey(sharedSecret)
	if err != nil {
		return nil, err
	}

	// 4. Encrypt PMK using AES-256-GCM
	block, err := aes.NewCipher(wrapKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	wrappedPMK := gcm.Seal(nil, nonce, pmk, nil)

	return &WrappedKey{
		WrappedPMK:       wrappedPMK,
		WrapNonce:        nonce,
		WrapEphemeralPub: ephemeral.PublicKey,
	}, nil
}

func UnwrapPMK(
	wrapped *WrappedKey,
	userPrivateKey []byte,
) ([]byte, error) {

	if len(userPrivateKey) != 32 {
		return nil, errors.New("invalid user private key length")
	}

	//fmt.Printf("EphemeralPub (%d): %x\n", len(wrapped.WrapEphemeralPub), wrapped.WrapEphemeralPub)
	//fmt.Printf("UserPriv (%d): %x\n", len(userPrivateKey), userPrivateKey)
	//
	//sharedSecret, _ := X25519SharedSecret(userPrivateKey, wrapped.WrapEphemeralPub)
	//fmt.Printf("SharedSecret (%d): %x\n", len(sharedSecret), sharedSecret)
	//
	//wrapKey, _ := DeriveWrapKey(sharedSecret)
	//fmt.Printf("WrapKey (%d): %x\n", len(wrapKey), wrapKey)
	//
	//fmt.Printf("Nonce (%d): %x\n", len(wrapped.WrapNonce), wrapped.WrapNonce)
	//fmt.Printf("Ciphertext (%d): %x\n", len(wrapped.WrappedPMK), wrapped.WrappedPMK)

	// 1. Derive shared secret
	sharedSecret, err := X25519SharedSecret(
		userPrivateKey,
		wrapped.WrapEphemeralPub,
	)
	if err != nil {
		return nil, err
	}

	// 2. Derive wrap key
	wrapKey, err := DeriveWrapKey(sharedSecret)
	if err != nil {
		return nil, err
	}

	// 3. Decrypt PMK
	block, err := aes.NewCipher(wrapKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	pmk, err := gcm.Open(
		nil,
		wrapped.WrapNonce,
		wrapped.WrappedPMK,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return pmk, nil
}

func EncryptENV(pmk []byte, data []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(pmk)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	return gcm.Seal(nil, nonce, data, nil), nonce, nil
}

func DecryptENV(pmk []byte, encryptedData []byte, nonce []byte) ([]byte, error) {

	block, err := aes.NewCipher(pmk)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, encryptedData, nil)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

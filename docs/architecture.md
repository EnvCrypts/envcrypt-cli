# Architecture

EnvCrypt is built on a **Zero-Trust, Client-Side Encryption** model. This ensures that the central server—and anyone with access to it—can never see your raw environment variables.

## Cryptographic Hierarchy

EnvCrypt uses a hybrid cryptosystem:

1.  **Project Master Key (PMK)**
    Symmetric Encryption: Environment variables for a project are encrypted locally using a per-project AES-256-GCM key (the PMK).
2.  **User Keypairs (X25519)**
    Key Wrapping: The PMK is encrypted ("wrapped") for each authorized user using their public X25519 key. 
3.  **Local Storage**
    Private keys never leave your device unencrypted. They are stored securely in your OS native keyring.
4.  **Authentication**
    All API requests are signed and authenticated.

## Data Flow

### Push (Encrypting Data)
1. CLI reads the local `.env` file.
2. CLI unwraps the PMK using your locally stored private X25519 key.
3. CLI encrypts the `.env` file contents using AES-256-GCM and the PMK.
4. CLI uploads the **ciphertext** to the server.
5. The server records an immutable version of this ciphertext.

### Pull (Decrypting Data)
1. CLI downloads the requested ciphertext version from the server.
2. CLI unwraps the PMK using your locally stored private X25519 key.
3. CLI decrypts the ciphertext using the PMK.
4. CLI writes the plaintext to your local machine.

### Granting Access
1. Admin unwraps the PMK using their private key.
2. Admin fetches the new user's public X25519 key from the server.
3. Admin wraps the PMK using the new user's public key.
4. Admin uploads the newly wrapped PMK to the server.
The server never sees the unwrapped PMK during this process.

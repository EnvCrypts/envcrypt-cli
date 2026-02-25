# Architecture

EnvCrypt is built on a **Zero-Trust, Client-Side Encryption** model. This ensures that the central server—and anyone with access to it—can never see your raw environment variables.

## Cryptographic Hierarchy

EnvCrypt uses a hybrid cryptosystem with a three-tier key hierarchy:

1.  **Data Encryption Key (DEK)**
    Symmetric Encryption: Each version of your environment variables is encrypted with a unique, randomly generated AES-256-GCM key (the DEK).
2.  **Project Root Key (PRK)**
    Key Wrapping: The DEK is encrypted ("wrapped") using the project's root key. The PRK is a symmetric key that never changes for the lifetime of the project (unless rotated).
3.  **User Keypairs (X25519)**
    Key Wrapping: The PRK is encrypted ("wrapped") for each authorized user using their public X25519 key.
4.  **Local Storage**
    Private keys never leave your device unencrypted. They are stored securely in your OS native keyring.
5.  **Authentication**
    All API requests are signed and authenticated.

## Data Flow

### Push (Encrypting Data)
1. CLI reads the local `.env` file.
2. CLI generates a new random **DEK**.
3. CLI encrypts the `.env` file contents using AES-256-GCM and the DEK.
4. CLI unwraps the **PRK** using your locally stored private X25519 key.
5. CLI encrypts (wraps) the DEK using the PRK.
6. CLI uploads the **ciphertext** (encrypted env vars) and the **wrapped DEK** to the server.
7. The server records an immutable version.

### Pull (Decrypting Data)
1. CLI downloads the requested ciphertext version and wrapped DEK from the server.
2. CLI unwraps the **PRK** using your locally stored private X25519 key.
3. CLI unwraps the **DEK** using the PRK.
4. CLI decrypts the ciphertext using the DEK.
5. CLI writes the plaintext to your local machine.

### Granting Access
1. Admin unwraps the **PRK** using their private key.
2. Admin fetches the new user's public X25519 key from the server.
3. Admin wraps the **PRK** using the new user's public key.
4. Admin uploads the newly wrapped PRK to the server.
The server never sees the unwrapped PRK or DEK during this process.
DEKs do not need to be re-encrypted when adding users, as they are encrypted with the PRK, not user keys.

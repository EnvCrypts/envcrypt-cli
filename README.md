# EnvCrypt CLI

**Secure, end-to-end encrypted environment variable management for modern teams.**

EnvCrypt CLI is the client-side tool for the EnvCrypt platform. It implements a Zero-Trust architecture where all secrets are encrypted locally on your machine before they are ever sent to the server. This ensures that the server—and anyone with access to it—can never see your raw environment variables.

## Features

-   **End-to-End Encryption**: Secrets are encrypted locally using AES-256-GCM. The server only sees ciphertext.
-   **Zero-Trust Model**: Your private key is stored only on your device (in the system keyring).
-   **Immutable Versioning**: Every `push` creates a new, immutable version. Rollback to any previous state instantly.
-   **Granular Access Control**: Manage access for team members and robustly handle user revocation.
-   **Service Roles**: Securely inject secrets into CI/CD pipelines using dedicated machine identities.
-   **Cross-Platform**: Works on Linux, macOS, and Windows.

## Installation

### Automated Install (Recommended)

Run the following command to install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/envcrypts/envcrypt-cli/main/install.sh | bash
```

### Prebuilt Binaries

Download the latest release for your platform from the [Releases](https://github.com/envcrypts/envcrypt-cli/releases) page.

### Building From Source

Requires Go 1.22+:

```bash
go install github.com/envcrypts/envcrypt-cli@latest
```

Ensure your `$GOPATH/bin` is in your system `$PATH`.

## Quick Start

### 1. Account Setup

Create an account. This generates a local X25519 keypair and securely stores the private key in your OS keyring.

```bash
envcrypt register
envcrypt login
```

### 2. Create a Project

Initialize a project. You become the admin and the Project Master Key (PMK) is generated and wrapped for you.

```bash
envcrypt create my-app
```

### 3. Push Secrets

Encrypt and upload your local `.env` file.

```bash
# Push to 'dev' environment
envcrypt push my-app --env dev --env-file .env
```

### 4. Pull Secrets

Decrypt and retrieve secrets on another machine or in production.

```bash
# Pull 'dev' secrets to a local .env file
envcrypt pull my-app --env dev
```

## Advanced Usage

### Team Management

Grant access to other users. The CLI handles the secure re-wrapping of the Project Master Key for the new user.

```bash
envcrypt grant my-app colleague@example.com
```

### Service Roles (CI/CD)

Create restricted machine users for your deployment pipelines.

1.  **Create Role**: `envcrypt service-role create my-ci-role` (Save the output private key!)
2.  **Delegate Access**: `envcrypt service-role grant my-ci-role my-app dev`
3.  **In CI**: Use `envcrypt ci login` with the private key to authenticate.

### Rollbacks

Mistake in production? Revert instantly.

```bash
envcrypt rollback
```

## Security Architecture

EnvCrypt uses a **hybrid cryptosystem**:

1.  **Symmetric Encryption**: Environment variables are encrypted with a per-project AES-256 key (PMK).
2.  **Key Wrapping**: The PMK is encrypted ("wrapped") for each user using their public X25519 key.
3.  **Authentication**: All requests are signed and authenticated.
4.  **Local Storage**: Private keys never leave your device unencrypted.


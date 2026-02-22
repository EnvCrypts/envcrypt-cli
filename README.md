# EnvCrypt CLI 🛡️

**Secure, end-to-end encrypted environment variable management for modern teams.**

EnvCrypt CLI is the client-side tool for the EnvCrypt platform. It implements a Zero-Trust architecture where all secrets are encrypted locally on your machine before they are ever sent to the server. This ensures that the server—and anyone with access to it—can never see your raw environment variables.

## Documentation

- [Getting Started](docs/getting-started.md)
- [CLI Guide](docs/cli-guide.md)
- [CI/CD and Service Roles](docs/ci-cd.md)
- [Architecture](docs/architecture.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)

## Features

-   **End-to-End Encryption**: Secrets are encrypted locally using AES-256-GCM. The server only sees ciphertext.
-   **Zero-Trust Model**: Your private key is stored only on your device (in the system keyring).
-   **Immutable Versioning**: Every `push` creates a new, immutable version. Rollback to any previous state instantly.
-   **Granular Access Control**: Manage access for team members and robustly handle user revocation.
-   **Service Roles**: Securely inject secrets into CI/CD pipelines using dedicated machine identities and OIDC.
-   **Cross-Platform**: Works on Linux, macOS, and Windows.

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/envcrypts/envcrypt-cli/main/install.sh | bash
```
For more installation options, see the [Getting Started Guide](docs/getting-started.md).

## License
MIT License.

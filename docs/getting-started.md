# Getting Started

Welcome to EnvCrypt CLI! This guide will help you install the CLI, create an account, and set up your first encrypted project.

## 1. Installation

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

## 2. Account Setup

To use EnvCrypt, you need to register an account. This generates a local X25519 keypair and securely stores the private key in your OS keyring.

```bash
# Register a new account
envcrypt register

# Log in to your account
envcrypt login
```

## 3. Create a Project

Initialize a new project. You become the admin and the Project Master Key (PMK) is generated and wrapped for you locally.

```bash
envcrypt create my-app
```

## 4. Push Secrets

Encrypt and upload your local `.env` file to the server. The server receives only ciphertext.

```bash
# Push to 'dev' environment
envcrypt push my-app --env dev --env-file .env
```

## 5. Pull Secrets

Decrypt and retrieve secrets on another machine or in production.

```bash
# Pull 'dev' secrets to a local .env file
envcrypt pull my-app --env dev
```

You're now ready to use EnvCrypt! See the [CLI Guide](cli-guide.md) for a full list of commands.

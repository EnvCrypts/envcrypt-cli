# EnvCrypt CLI

**Secure, end-to-end encrypted environment variable management for modern teams.**

EnvCrypt CLI provides a zero-trust architecture for managing secrets. All environment variables are encrypted client-side using robust asymmetric cryptography before ever leaving your machine. The server never sees your raw secrets.

## Features

- **End-to-End Encryption**: Secrets are encrypted locally with per-project, per-user keys. The server acts only as a blind store.
- **Immutable Versioning**: Every change creates a new, immutable version of your environment.
- **Atomic Rollbacks**: Instantaneously revert to any previous version with guaranteed consistency.
- **Granular Access Control**: Grant and revoke access to projects for specific team members.
- **Diffing**: Visualize changes between environment versions before applying them.
- **Developer Experience**: Minimalist, efficient CLI designed for seamless integration into your workflow.

## Installation

### From Source

Requires Go 1.22+:

```bash
go install github.com/envcrypts/envcrypt-cli@latest
```

Ensure your `$GOPATH/bin` is in your system `$PATH`.

## Quick Start

### 1. Account Setup

First, create an account and log in. keys are generated locally during registration.

```bash
envcrypt register
envcrypt login
```

### 2. Create a Project

Create a new project namespace for your application.

```bash
envcrypt create my-app
```

### 3. Push Secrets

Navigate to your project directory containing a `.env` file and push it to the `dev` environment.

```bash
# Push the local .env file to the 'dev' environment
envcrypt push my-app --env dev --env-file .env
```

### 4. Pull Secrets

On another machine (or for deployment), pull the secrets down.

```bash
# Pull 'dev' secrets to a local .env file
envcrypt pull my-app --env dev
```

## Command Reference

### Authentication

- **`register`**: Create a new account and generate local key pairs.
- **`login`**: Authenticate with your credentials.
- **`logout`**: Clear local session.
- **`whoami`**: Display current user information.

### Project Management

- **`create [name]`**: Initialize a new project. You automatically become the owner.
- **`list`**: Show all projects you have access to.
- **`delete [name]`**: Permanently delete a project and all its environments (Ownership required).

### Secrets Management

- **`push [project]`**: Encrypt and upload variables from a local file.
    - Flags: `--env` (default: dev), `--env-file` (default: .env)
- **`pull [project]`**: Download and decrypt variables to a local file.
    - Flags: `--env` (default: dev), `--env-file` (default: .env), `--yes` (skip confirmation)
- **`add [project] [key=value]`**: Add or update a single variable without a file.

### Versioning & History

- **`diff [old] [new]`**: Compare two environment versions.
    - Flags: `--env`, `--project`, `--show-secrets` (reveal values)
    - If versions are omitted, interactive mode is launched.
- **`rollback`**: Revert an environment to a previous version.
    - Interactive prompts guide you through version selection and confirmation.

### Access Control

- **`grant [project] [email]`**: Authorize a user to access a project. They must have an EnvCrypt account.
- **`revoke [project] [email]`**: Remove a user's access to a project.

## Security Architecture

EnvCrypt uses a **Zero-Trust** model:

1.  **Client-Side Key Generation**: When you register, a Public/Private key pair is generated on your machine. The Private key is stored locally, encryption-protected by your password. Only the Public key is sent to the *server*
2.  **Project Keys**: Each project has a unique Project Master Key (PMK).
3.  **Envelope Encryption**: The PMK is encrypted individually for each authorized user using their Public Key.
4.  **Secret Encryption**: Environment variables are encrypted using the PMK and a nonce locally.
5.  **Storage**: The server stores only the encrypted secrets and the encrypted PMKs. It cannot decrypt any data.

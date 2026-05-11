# CLI Guide

This document provides an overview of the most commonly used `envcrypt` commands.

## Basic Commands

### `envcrypt register`
Create a new EnvCrypt user and cryptographic identity. Generates an X25519 keypair and creates your user profile.

### `envcrypt login`
Authenticate and unlock your EnvCrypt session.

### `envcrypt logout`
Lock your EnvCrypt session and remove credentials from active memory.

### `envcrypt whoami`
Show the current authenticated user identity.

## Project Management

### `envcrypt create [project-name]`
Create a new project. Generates a new random Project Root Key (PRK). If left blank, defaults to the git repository name.

### `envcrypt list`
List all projects you have access to.

### `envcrypt delete <project-name>`
Delete a project permanently (admin only).

### `envcrypt audit project <project-name>`
View detailed project-level audit logs for all security and management events.
Example:
```bash
envcrypt audit project my-app
```

## Environment Variable Operations

### `envcrypt push [project-name] --env <env-name> --env-file <path>`
Encrypt and upload environment variables from a local file. If `[project-name]` is left blank, the selector defaults to the git repository name.
Example:
```bash
envcrypt push my-app --env prod --env-file .env.prod
```

### `envcrypt pull [project-name] --env <env-name>`
Download and decrypt environment variables to your local machine. If `[project-name]` is left blank, the selector defaults to the git repository name.

### `envcrypt run [project-name] --env <env-name> -- <command> [args...]`
Run a command with secrets injected at runtime. If `[project-name]` is left blank, the selector defaults to the git repository name.
Example:
```/dev/null/cli-guide.md#L1-2
envcrypt run my-app --env prod -- npm start
```

### `envcrypt diff`
Diff two environment versions or show changes between local and remote.

### `envcrypt rollback`
Rollback to a previous immutable version of an environment.

## Access Control

### `envcrypt grant <project-name> <email>`
Grant a user access to a project. This decrypts the PRK and re-wraps it for the new user's public key.
Example:
```bash
envcrypt grant my-app colleague@example.com
```

### `envcrypt revoke <project-name> <email>`
Revoke a user's access to a project.

## Service Roles (CI/CD)

See the [CI/CD Guide](ci-cd.md) for details.
Commands under `envcrypt service-role`:
- `create`: Create a new service role
- `grant`: Grant CI access to a project/env
- `list`: List service roles
- `delete`: Delete a service role

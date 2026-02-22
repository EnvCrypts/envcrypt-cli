# CI/CD and Service Roles

EnvCrypt uses **Service Roles** to securely inject secrets into CI/CD pipelines (such as GitHub Actions) using dedicated machine identities and OIDC (OpenID Connect).

This allows your CI/CD runner to authenticate without needing a long-lived user password, relying instead on short-lived OIDC tokens and a specific service role private key.

## Concept

A Service Role represents a machine identity (e.g., your GitHub Actions workflow). To allow a pipeline to pull secrets:
1. You create a Service Role bound to a repository and branch (the "Principal").
2. You grant this Service Role access to a specific Project and Environment.
3. In your CI/CD pipeline, the runner authenticates using its GitHub OIDC token and the Service Role's private key.

## Step 1: Create a Service Role

As an admin, create a service role for your repository.

```bash
envcrypt service-role create \
  --repo github:acme/billing-backend \
  --branch main \
  --name sp-billing-backend
```
*Note: Make sure to save the output private key! It is only shown once.*

If you don't provide the flags, the CLI can attempt to auto-detect them from your current Git context.

## Step 2: Grant Access

Grant the new service role access to the specific project and environment it needs to pull.

```bash
envcrypt service-role grant \
  --service-role sp-billing-backend \
  --project billing-service \
  --env main
```

## Step 3: Setup CI/CD (GitHub Actions Example)

Add the Service Role Private Key you saved in Step 1 as a GitHub Repository Secret (e.g., `ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY`).
Your CI pipeline will use the `envcrypt ci login` command to authenticate and pull secrets.

Here is an example workflow incorporating the steps:

```yaml
name: Deploy
on:
  push:
    branches: [ "main" ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write # Required for OIDC token fetching
      contents: read
    
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Install EnvCrypt CLI
        run: |
          curl -fsSL https://raw.githubusercontent.com/envcrypts/envcrypt-cli/main/install.sh | bash

      - name: Get GitHub OIDC Token
        id: oidc
        uses: actions/github-script@v6
        with:
          script: |
            const token = await core.getIDToken('envcrypts/envcrypt')
            core.setOutput('token', token)

      - name: Pull Environment Variables
        env:
          ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY: ${{ secrets.ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY }}
        run: |          
          envcrypt ci login \
            --oidc-token "${{ steps.oidc.outputs.token }}" \
            --env main \
            --output .env
            
      - name: Build and Deploy (Using Secrets)
        run: |
          # The secrets are automatically injected into $GITHUB_ENV and available to subsequent steps!
          # You can also use the generated .env file directly if needed.
          make deploy```

> **Note**: The `envcrypt ci login` command specifically requires the `--oidc-token` and `--env` flags, and expects the `ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY` environment variable to be set in the runner's context.

package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	ciOIDCToken string
	ciProject   string
	ciEnv       string
	ciOutput    string
)

var ciLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and pull secrets in CI environment",
	Long: `Authenticate using GitHub OIDC token and pull secrets for CI/CD.

Example:
  envcrypt ci login \
    --oidc-token $ACTIONS_ID_TOKEN \
    --project my-app \
    --env prod \
    --output .env`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		if ciOIDCToken == "" {
			return Error("--oidc-token is required", nil)
		}
		if ciProject == "" {
			return Error("--project is required", nil)
		}
		if ciEnv == "" {
			return Error("--env is required", nil)
		}

		outputPath := ciOutput
		if outputPath == "" {
			outputPath = ".env"
		}

		Info(fmt.Sprintf("Project: %s", ciProject))
		Info(fmt.Sprintf("Environment: %s", ciEnv))

		sessionID, projectID, err := Application.GetSessionID(cmd.Context(), ciOIDCToken)
		if err != nil {
			return Error("OIDC authentication failed", err)
		}

		Info("OIDC authentication successful")

		keysResp, err := Application.GetServiceRoleProjectKeys(cmd.Context(), *projectID, *sessionID, ciEnv)
		if err != nil {
			return Error("failed to get project keys", err)
		}

		privateKeyB64 := os.Getenv("ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY")
		if privateKeyB64 == "" {
			return Error("ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY environment variable is required", nil)
		}

		privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
		if err != nil {
			return Error("failed to decode service role private key", err)
		}

		wrappedKey := &cryptoutils.WrappedKey{
			WrappedPMK:       keysResp.WrappedPMK,
			WrapNonce:        keysResp.WrapNonce,
			WrapEphemeralPub: keysResp.EphemeralPublicKey,
		}

		pmk, err := cryptoutils.UnwrapPMK(wrappedKey, privateKey)
		if err != nil {
			return Error("failed to unwrap project key", err)
		}

		envMap, err := Application.PullEnvForCI(cmd.Context(), *projectID, ciEnv, pmk)
		if err != nil {
			return Error("failed to pull environment variables", err)
		}

		if len(envMap) == 0 {
			return Error("no environment variables found", fmt.Errorf("environment %q has no variables", ciEnv))
		}

		printEnvSummary(envMap)

		envBytes, err := cryptoutils.EncodeEnv(envMap)
		if err != nil {
			return Error("failed to encode env file", err)
		}

		if err := os.WriteFile(outputPath, envBytes, 0600); err != nil {
			return Error("failed to write env file", fmt.Errorf("could not write to %q: %w", outputPath, err))
		}

		Success(fmt.Sprintf("Pulled %d secrets to %s", len(envMap), outputPath))
		return nil
	},
}

func init() {
	ciLoginCmd.Flags().StringVar(&ciOIDCToken, "oidc-token", "", "GitHub OIDC token (required)")
	ciLoginCmd.Flags().StringVar(&ciProject, "project", "", "Project name (required)")
	ciLoginCmd.Flags().StringVar(&ciEnv, "env", "", "Environment name: dev|stage|prod (required)")
	ciLoginCmd.Flags().StringVarP(&ciOutput, "output", "o", "", "Output path for .env file (default: .env)")
	ciCmd.AddCommand(ciLoginCmd)
}

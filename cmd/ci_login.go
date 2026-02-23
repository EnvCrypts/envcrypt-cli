package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	ciOIDCToken string
	ciEnv       string
	ciOutput    string
)

var ciLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and pull secrets in CI environment",
	Long: `Authenticate using GitHub OIDC token and pull secrets for CI/CD.

Examples:
  envcrypt ci login \
    --oidc-token $ACTIONS_ID_TOKEN \
    --env prod \
    --output .env`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		if ciOIDCToken == "" {
			return tui.Error("--oidc-token is required", nil)
		}
		if ciEnv == "" {
			return tui.Error("--env is required", nil)
		}

		outputPath := ciOutput
		if outputPath == "" {
			outputPath = ".env"
		}

		tui.Info(fmt.Sprintf("Environment: %s", ciEnv))

		sessionID, projectID, err := Application.GetSessionID(cmd.Context(), ciOIDCToken)
		if err != nil {
			return tui.Error("OIDC authentication failed", err)
		}

		tui.Info("OIDC authentication successful")

		keysResp, err := Application.GetServiceRoleProjectKeys(cmd.Context(), *projectID, *sessionID, ciEnv)
		if err != nil {
			return tui.Error("failed to get project keys", err)
		}

		privateKeyB64 := os.Getenv("ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY")
		if privateKeyB64 == "" {
			return tui.Error("ENVCRYPT_SERVICE_ROLE_PRIVATE_KEY environment variable is required", nil)
		}

		privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
		if err != nil {
			return tui.Error("failed to decode service role private key", err)
		}

		wrappedKey := &cryptoutils.WrappedKey{
			WrappedPRK:       keysResp.WrappedPRK,
			WrapNonce:        keysResp.WrapNonce,
			WrapEphemeralPub: keysResp.EphemeralPublicKey,
		}

		prk, err := cryptoutils.UnwrapPRK(wrappedKey, privateKey)
		if err != nil {
			return tui.Error("failed to unwrap project key", err)
		}

		envMap, err := Application.PullEnvForCI(cmd.Context(), *projectID, ciEnv, prk)
		if err != nil {
			return tui.Error("failed to pull environment variables", err)
		}

		if len(envMap) == 0 {
			tui.Info(fmt.Sprintf("No environment variables found for %s. Creating empty .env file.", ciEnv))
		}

		tui.PrintEnvSummary(envMap)

		envBytes, err := cryptoutils.EncodeEnv(envMap)
		if err != nil {
			return tui.Error("failed to encode env file", err)
		}

		if err := os.WriteFile(outputPath, envBytes, 0600); err != nil {
			return tui.Error("failed to write env file", fmt.Errorf("could not write to %q: %w", outputPath, err))
		}

		tui.Success(fmt.Sprintf("Pulled %d secrets to %s", len(envMap), outputPath))

		if githubEnv := os.Getenv("GITHUB_ENV"); githubEnv != "" {
			f, err := os.OpenFile(githubEnv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				defer f.Close()
				for k, v := range envMap {
					if strings.Contains(v, "\n") {
						delimiter := "EOF"
						fmt.Fprintf(f, "%s<<%s\n%s\n%s\n", k, delimiter, v, delimiter)
					} else {
						fmt.Fprintf(f, "%s=%s\n", k, v)
					}
				}
				tui.Success("Injected secrets into GITHUB_ENV")
			}
		}

		return nil
	},
}

func init() {
	ciLoginCmd.Flags().StringVar(&ciOIDCToken, "oidc-token", "", "GitHub OIDC token (required)")
	ciLoginCmd.Flags().StringVar(&ciEnv, "env", "", "Environment name: dev|stage|prod (required)")
	ciLoginCmd.Flags().StringVarP(&ciOutput, "output", "o", "", "Output path for .env file (default: .env)")
	ciCmd.AddCommand(ciLoginCmd)
}

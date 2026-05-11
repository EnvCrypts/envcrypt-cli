package cmd

import (
	"context"
	"errors"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "Recover account access using your recovery key",
	Long: `Recover performs a zero-knowledge recovery flow to restore access
to your account if you lost your password, but saved your recovery key.

Examples:
  envcrypt recover
  envcrypt recover --email user@example.com`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var collectedEmail string

		if email == "" {
			vals, err := tui.RunForm([]tui.FormField{
				{Label: "Email", Required: true, Validate: tui.ValidateEmail},
			}, []string{""})
			if err != nil {
				return handlePromptError(err, "email is required in non-interactive mode", "Provide an email address to recover")
			}
			if len(vals) == 0 {
				return tui.Cancelled()
			}
			collectedEmail = vals[0]
		} else {
			collectedEmail = email
		}

		var initResp *config.RecoveryInitResponse
		err := tui.RunActionWithSpinner("Initializing recovery...", func() error {
			var initErr error
			initResp, initErr = Application.RecoverInit(context.Background(), collectedEmail)
			return initErr
		})
		if err != nil {
			return tui.Error("Failed to initialize recovery", err)
		}

		if initResp == nil || len(initResp.RecoveryPrivateKey) == 0 || len(initResp.RecoverySalt) == 0 || len(initResp.RecoveryNonce) == 0 {
			return tui.Error("Invalid recovery data", errors.New("incomplete recovery data received from server. Are you sure you have a recovery key?"))
		}

		vals, err := tui.RunForm([]tui.FormField{
			{Label: "Recovery Key", Secret: true, Required: true},
		}, []string{""})
		if err != nil {
			return handlePromptError(err, "recovery key is required in non-interactive mode", "Run in an interactive terminal to enter the recovery key")
		}
		if len(vals) == 0 {
			return tui.Cancelled()
		}
		recoveryKey := vals[0]

		encryptedKeyParams := &config.EncryptedPrivateKey{
			EncryptedUserPrivateKey: initResp.RecoveryPrivateKey,
			PrivateKeySalt:          initResp.RecoverySalt,
			PrivateKeyNonce:         initResp.RecoveryNonce,
		}

		privateKeyPlain, err := cryptoutils.DecryptPrivateKey(encryptedKeyParams, recoveryKey, &config.DefaultArgon2Params)
		if err != nil {
			return tui.Error("Invalid recovery key", err)
		}

		tui.Success("Recovery key accepted!")

		vals, err = tui.RunForm([]tui.FormField{
			{Label: "New Password", Secret: true, Required: true},
		}, []string{""})
		if err != nil {
			return handlePromptError(err, "password is required in non-interactive mode", "Run in an interactive terminal to enter a new password")
		}
		if len(vals) == 0 {
			return tui.Cancelled()
		}
		newPassword := vals[0]

		err = tui.RunActionWithSpinner("Finalizing recovery...", func() error {
			return Application.RecoverComplete(context.Background(), collectedEmail, newPassword, privateKeyPlain)
		})
		if err != nil {
			return tui.Error("Failed to complete recovery", err)
		}

		tui.Success("Account recovered and password updated successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(recoverCmd)
	recoverCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
}

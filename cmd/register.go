package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/config"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var email string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new EnvCrypt user and cryptographic identity",
	Long: `Register creates a local encryption key pair and associates it with your EnvCrypt account using end-to-end encryption.

Examples:
  envcrypt register --email user@example.com`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var fields []tui.FormField
		var prefills []string

		if email == "" {
			fields = []tui.FormField{
				{Label: "Email", Required: true, Validate: tui.ValidateEmail},
				{Label: "Password", Secret: true, Required: true},
				{Label: "Custom Backend URL (leave blank for default)", Required: false, Validate: tui.ValidateBackendURL},
			}
			prefills = []string{"", "", ""}
		} else {
			fields = []tui.FormField{
				{Label: "Password", Secret: true, Required: true},
				{Label: "Custom Backend URL (leave blank for default)", Required: false, Validate: tui.ValidateBackendURL},
			}
			prefills = []string{"", ""}
		}

		vals, err := tui.RunForm(fields, prefills)
		if err != nil {
			return handlePromptError(err, "registration requires interactive input", "Run in an interactive terminal to enter account details")
		}

		var collectedEmail, password, backendURL string
		if email == "" {
			collectedEmail = vals[0]
			password = vals[1]
			backendURL = vals[2]
		} else {
			collectedEmail = email
			password = vals[0]
			backendURL = vals[1]
		}

		if err := applyBackendURL(backendURL); err != nil {
			return err
		}

		var keypair *config.KeyPair
		err = tui.RunActionWithSpinner("Registering account...", func() error {
			var regErr error
			keypair, regErr = Application.Register(cmd.Context(), collectedEmail, password)
			return regErr
		})
		if err != nil {
			return tui.Error("Registration failed", err)
		}

		tui.Success("Registration successful!")
		tui.Info("\nIMPORTANT: Store your recovery key safely. It is the only way to recover your account if you forget your password:")
		tui.Success(keypair.RecoveryKey)
		tui.Info("This key will never be shown again.\n")
		tui.Info("Next: run 'envcrypt login' to start using your secrets.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
}

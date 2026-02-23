package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:          "login",
	Short:        "Authenticate and unlock your EnvCrypt session",
	Long: `Login unlocks your local encryption keys and authorizes access to encrypted environment variables without exposing plaintext secrets.

Examples:
  envcrypt login
  envcrypt login --email user@example.com`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var fields []tui.FormField
		var prefills []string

		if email == "" {
			fields = []tui.FormField{
				{Label: "Email", Required: true, Validate: tui.ValidateEmail},
				{Label: "Password", Secret: true, Required: true},
			}
			prefills = []string{"", ""}
		} else {
			fields = []tui.FormField{
				{Label: fmt.Sprintf("Password for %s", email), Secret: true, Required: true},
			}
			prefills = []string{""}
		}

		vals, err := tui.RunForm(fields, prefills)
		if err != nil {
			return tui.Cancelled()
		}

		var collectedEmail, password string
		if email == "" {
			collectedEmail = vals[0]
			password = vals[1]
		} else {
			collectedEmail = email
			password = vals[0]
		}

		err = tui.RunActionWithSpinner("Authenticating...", func() error {
			return Application.Login(context.Background(), collectedEmail, password)
		})
		if err != nil {
			return tui.Error("login failed", err)
		}

		tui.Success("Login successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and unlock your EnvCrypt session",
	Long: `Login unlocks your local encryption keys and authorizes access
to encrypted environment variables without exposing plaintext secrets.`,
	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		if email == "" {
			return Error("email is required", nil)
		}

		// Ensure we are in an interactive terminal
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return Error("login requires an interactive terminal", nil)
		}

		Info("Authenticating…")

		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			return Error("failed to read password", err)
		}

		if len(password) == 0 {
			return Error("password cannot be empty", nil)
		}

		if err := Application.Login(
			cmd.Context(),
			email,
			string(password),
		); err != nil {
			return Error("login failed", err)
		}

		Success("Login successful")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(
		&email,
		"email",
		"e",
		"",
		"Email address",
	)
	loginCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(loginCmd)
}

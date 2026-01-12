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
	Short: "Authenticate and unlock your EnvCrypt session.",
	Long: `Login unlocks your local encryption keys and authorizes access
to encrypted environment variables without exposing plaintext secrets.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			return prettyError("Failed to read password", err)
		}

		err = Application.Login(cmd.Context(), email, string(password))
		if err != nil {
			return prettyError("login failed", err)
		}

		prettySuccess("login successful!")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	loginCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(loginCmd)
}

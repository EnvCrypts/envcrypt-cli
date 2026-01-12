package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	email string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new EnvCrypt user and cryptographic identity",
	Long: `Register creates a local encryption key pair and associates
it with your EnvCrypt account using end-to-end encryption.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()

		if err != nil {
			return prettyError("Failed to read password", err)
		}

		err = Application.Register(cmd.Context(), email, string(password))
		if err != nil {
			return prettyError("Registration failed", err)
		}

		prettySuccess("Registration successful!")
		return nil
	},
}

func init() {
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	registerCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(registerCmd)
}

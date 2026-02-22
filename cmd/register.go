package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var email string

var registerCmd = &cobra.Command{
	Use:          "register",
	Short:        "Create a new EnvCrypt user and cryptographic identity",
	Long:         `Register creates a local encryption key pair and associates it with your EnvCrypt account using end-to-end encryption.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		vals, err := tui.RunForm([]tui.FormField{
			{Label: "Password", Secret: true, Required: true},
		}, []string{""})
		if err != nil {
			return tui.Error("cancelled", nil)
		}

		if err := Application.Register(cmd.Context(), email, vals[0]); err != nil {
			return tui.Error("Registration failed", err)
		}

		tui.Success("Registration successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	registerCmd.MarkFlagRequired("email")
}

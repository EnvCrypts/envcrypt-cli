package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and unlock your EnvCrypt session",
	Long: `Login unlocks your local encryption keys and authorizes access
to encrypted environment variables without exposing plaintext secrets.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var password string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email").
					Value(&email).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("email is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Password").
					Value(&password).
					EchoMode(huh.EchoModePassword).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("password is required")
						}
						return nil
					}),
			),
		)

		if email != "" {
			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title(fmt.Sprintf("Password for %s", email)).
						Value(&password).
						EchoMode(huh.EchoModePassword).
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("password is required")
							}
							return nil
						}),
				),
			)
		}

		if err := form.Run(); err != nil {
			return Error("cancelled", nil)
		}

		if err := Application.Login(cmd.Context(), email, password); err != nil {
			return Error("login failed", err)
		}

		Success("Login successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
}

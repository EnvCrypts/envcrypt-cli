package cmd

import (
	"fmt"
	
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
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
		// Ensure we are in an interactive terminal
		// (huh handles this check internally mostly, but good to keep)

		var password string

		// If email is not provided via flag, ask for it
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

		// If email was set by flag, we might want to skip that input?
		// Huh forms are static. We can build it dynamically.
		if email != "" {
			// Only ask for password
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

		err := form.Run()
		if err != nil {
			return Error("cancelled", nil)
		}

		Info("Authenticating...")
		// TODO: Add Spinner here later

		if err := Application.Login(
			cmd.Context(),
			email,
			password,
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


	rootCmd.AddCommand(loginCmd)
}

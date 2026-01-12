package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Lock your EnvCrypt identity and clear local authentication state.",
	Long: `Logout securely ends your EnvCrypt session by discarding any
in-memory keys and authentication tokens.

This ensures that encrypted environment variables cannot be accessed
again without re-authenticating and unlocking your local identity.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := Application.Logout(email)
		if err != nil {
			return errors.New("user not logged in")
		}
		prettySuccess("Logged out successfully")
		return nil
	},
}

func init() {
	logoutCmd.Flags().StringVarP(&email, "email", "e", "", "Email address")
	logoutCmd.MarkFlagRequired("email")
	rootCmd.AddCommand(logoutCmd)
}

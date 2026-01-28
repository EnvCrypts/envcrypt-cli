package cmd

import (
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Lock your EnvCrypt session",
	Long: `Logout securely ends your EnvCrypt session by discarding any
in-memory keys and clearing local authentication state.

Encrypted environment variables cannot be accessed again without
re-authenticating.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := Application.Logout(email); err != nil {
			return Error("not logged in", err)
		}

		Success("Logged out successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

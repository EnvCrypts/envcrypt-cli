package cmd

import (
	"context"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Lock your EnvCrypt session",
	Long: `Logout securely ends your EnvCrypt session by discarding any
in-memory keys and clearing local authentication state.

Encrypted environment variables cannot be accessed again without
re-authenticating.

Examples:
  envcrypt logout`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.ConfirmDangerousAction("Are you sure you want to logout? You will lose access to decrypt envs until you login again.", "yes") {
			return tui.Cancelled()
		}

		err := tui.RunActionWithSpinner("Logging out...", func() error {
			return Application.Logout(context.Background(), email)
		})
		
		if err != nil {
			return tui.Error("not logged in", err)
		}

		tui.Success("Logged out successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

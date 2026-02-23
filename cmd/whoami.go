package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var whoamiCmd = &cobra.Command{
	Use:          "whoami",
	Short:        "Show the current authenticated user",
	Long: `Display the identity currently logged into EnvCrypt.

Examples:
  envcrypt whoami
  envcrypt whoami --json`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString("user.email")
		userID := viper.GetString("user.id")

		if tui.Mode() == tui.ModeJSON {
			tui.JSONData(map[string]any{
				"authenticated": email != "",
				"email":         email,
				"user_id":       userID,
			})
			return nil
		}

		if email == "" {
			return tui.Error(
				"not logged in",
				nil,
				"Run 'envcrypt login' to authenticate")
		}

		tui.Success("Logged in as " + email)

		if userID != "" {
			tui.Info("User ID: " + userID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:          "whoami",
	Short:        "Show the current authenticated user",
	Long:         "Display the identity currently logged into EnvCrypt.",
	Args:         cobra.NoArgs,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString("user.email")
		userID := viper.GetString("user.id")

		if email == "" {
			return tui.Error(
				"not logged in",
				nil)

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

package cmd

import (
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current auth state and configuration",
	Long: `Display the current authentication state, server URL, and configuration.

Examples:
  envcrypt status
  envcrypt status --json`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString("user.email")
		userID := viper.GetString("user.id")
		serverURL := viper.GetString("api.base_url")

		if tui.Mode() == tui.ModeJSON {
			status := map[string]any{
				"authenticated": email != "",
				"email":         email,
				"user_id":       userID,
				"server_url":    serverURL,
			}
			tui.JSONData(status)
			return nil
		}

		tui.Spacer()

		if email == "" {
			tui.Warn("Authentication: Not logged in")
			tui.Info("Action: Run 'envcrypt login' to authenticate")
		} else {
			header := tui.StyleSuccess.Bold(true).Render(" SESSION ACTIVE ")
			content := fmt.Sprintf(
				"%s %s\n%s %s\n%s %s",
				tui.StyleMuted.Render("User:  "), email,
				tui.StyleMuted.Render("ID:    "), userID,
				tui.StyleMuted.Render("Server:"), serverURL,
			)
			fmt.Println(tui.BoxStylePrimary.Render(fmt.Sprintf("%s\n\n%s", header, content)))
		}

		tui.Spacer()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

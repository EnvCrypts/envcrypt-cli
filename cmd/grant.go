package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	grantProject string
	grantEmail   string
)

var grantCmd = &cobra.Command{
	Use:          "grant [project]",
	Short:        "Grant a user's access to a project",
	Long: `Grant or restore a user's access to a project without re-adding the member.

Examples:
  envcrypt grant my-project --email user@example.com
  envcrypt grant`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := grantProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		// Pick project from list if not provided
		if projectName == "" {
			resp, err := Application.ListProjects(context.Background())
			if err != nil {
				return tui.Error("failed to fetch projects", err)
			}
			names := make([]string, len(resp.Projects))
			for i, p := range resp.Projects {
				names[i] = p.Name
			}
			projectName, err = tui.RunPicker("Select a project to grant access on", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		// Prompt for email if not provided
		if grantEmail == "" {
			vals, err := tui.RunForm([]tui.FormField{
				{Label: fmt.Sprintf("Email to grant on %q", projectName), Required: true, Validate: tui.ValidateEmail},
			}, []string{""})
			if err != nil {
				return tui.Cancelled()
			}
			grantEmail = vals[0]
		}

		if grantEmail == "" {
			return tui.Error("email is required", nil)
		}

		if err := Application.GiveAccess(cmd.Context(), projectName, grantEmail); err != nil {
			return tui.Error("failed to grant access", err)
		}

		tui.Success("Granted access for " + grantEmail + " on project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)

	grantCmd.Flags().StringVar(&grantProject, "project", "", "Project name")
	grantCmd.Flags().StringVar(&grantEmail, "email", "", "Email address of the user")
}

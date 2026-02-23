package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	revokeProject string
	revokeEmail   string
	revokeForce   bool
)

var revokeCmd = &cobra.Command{
	Use:          "revoke [project]",
	Short:        "Revoke a user's access to a project",
	Long: `Revoke a user's access to a project without removing the member.

Use --force to skip the confirmation prompt.

Examples:
  envcrypt revoke my-project --email user@example.com
  envcrypt revoke --force`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := revokeProject
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
			projectName, err = tui.RunPicker("Select a project to revoke access on", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		// Prompt for email if not provided
		if revokeEmail == "" {
			vals, err := tui.RunForm([]tui.FormField{
				{Label: fmt.Sprintf("Email to revoke on %q", projectName), Required: true, Validate: tui.ValidateEmail},
			}, []string{""})
			if err != nil {
				return tui.Cancelled()
			}
			revokeEmail = vals[0]
		}

		if revokeEmail == "" {
			return tui.Error("email is required", nil)
		}

		ok := tui.ConfirmDangerousAction(
			fmt.Sprintf("Revoke access for %s on project %q?", revokeEmail, projectName),
			revokeEmail,
			revokeForce,
		)
		if !ok {
			return tui.Cancelled()
		}

		if err := Application.RevokeAccess(cmd.Context(), projectName, revokeEmail); err != nil {
			return tui.Error("failed to revoke access", err)
		}

		tui.Success("Revoked access for " + revokeEmail + " on project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)

	revokeCmd.Flags().StringVar(&revokeProject, "project", "", "Project name")
	revokeCmd.Flags().StringVar(&revokeEmail, "email", "", "Email address of the user")
	revokeCmd.Flags().BoolVar(&revokeForce, "force", false, "Revoke without confirmation prompt")
}

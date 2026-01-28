package cmd

import (
	"github.com/spf13/cobra"
)

var (
	revokeMemberEmail   string
	revokeMemberProject string
)

// revokeCmd represents the project member revoke command
var revokeCmd = &cobra.Command{
	Use:          "revoke [project]",
	Short:        "Revoke a user's access to a project",
	Long:         "Revoke a user's access to a project without removing the member.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Resolve project name
		projectName := revokeMemberProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}

		if err := Application.RevokeAccess(
			cmd.Context(),
			projectName,
			revokeMemberEmail,
		); err != nil {
			return Error("failed to revoke access", err)
		}

		Success(
			"Revoked access for " + revokeMemberEmail +
				" on project " + projectName,
		)

		return nil
	},
}

func init() {
	memberCmd.AddCommand(revokeCmd)

	revokeCmd.Flags().StringVar(
		&revokeMemberProject,
		"project",
		"",
		"Project name",
	)

	revokeCmd.Flags().StringVar(
		&revokeMemberEmail,
		"member-email",
		"",
		"Email address of the user",
	)

	revokeCmd.MarkFlagRequired("member-email")
}

package cmd

import "github.com/spf13/cobra"

var (
	revokeProject string
	revokeEmail   string
)

var revokeCmd = &cobra.Command{
	Use:          "revoke [project]",
	Short:        "Revoke a user's access to a project",
	Long:         "Revoke a user's access to a project without removing the member.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := revokeProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}

		if err := Application.RevokeAccess(cmd.Context(), projectName, revokeEmail); err != nil {
			return Error("failed to revoke access", err)
		}

		Success("Revoked access for " + revokeEmail + " on project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)

	revokeCmd.Flags().StringVar(&revokeProject, "project", "", "Project name")
	revokeCmd.Flags().StringVar(&revokeEmail, "email", "", "Email address of the user")
	revokeCmd.MarkFlagRequired("email")
}

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	grantMemberEmail   string
	grantMemberProject string
)

// grantCmd represents the project member grant command
var grantCmd = &cobra.Command{
	Use:          "grant [project]",
	Short:        "Grant a user's access to a project",
	Long:         "Grant or restore a user's access to a project without re-adding the member.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Resolve project name (flag > arg)
		projectName := grantMemberProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}

		if err := Application.GiveAccess(
			cmd.Context(),
			projectName,
			grantMemberEmail,
		); err != nil {
			return Error("failed to grant access", err)
		}

		Success(
			"Granted access for " + grantMemberEmail +
				" on project " + projectName,
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)

	grantCmd.Flags().StringVar(
		&grantMemberProject,
		"project",
		"",
		"Project name",
	)

	grantCmd.Flags().StringVar(
		&grantMemberEmail,
		"member-email",
		"",
		"Email address of the user",
	)

	grantCmd.MarkFlagRequired("member-email")
}

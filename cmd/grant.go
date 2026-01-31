package cmd

import "github.com/spf13/cobra"

var (
	grantProject string
	grantEmail   string
)

var grantCmd = &cobra.Command{
	Use:          "grant [project]",
	Short:        "Grant a user's access to a project",
	Long:         "Grant or restore a user's access to a project without re-adding the member.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := grantProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}

		if err := Application.GiveAccess(cmd.Context(), projectName, grantEmail); err != nil {
			return Error("failed to grant access", err)
		}

		Success("Granted access for " + grantEmail + " on project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grantCmd)

	grantCmd.Flags().StringVar(&grantProject, "project", "", "Project name")
	grantCmd.Flags().StringVar(&grantEmail, "email", "", "Email address of the user")
	grantCmd.MarkFlagRequired("email")
}

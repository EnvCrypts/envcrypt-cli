package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	addMemberEmail   string
	addMemberRole    string
	addMemberProject string
)

// addCmd represents the project member add command
var addCmd = &cobra.Command{
	Use:          "add [project]",
	Short:        "Add a user to a project",
	Long:         "Add a user to a project and assign them a role.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		// Resolve project name
		projectName := addMemberProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}

		role := strings.ToLower(addMemberRole)
		if role != "admin" && role != "member" {
			return Error("invalid role (must be admin or member)", nil)
		}

		if err := Application.AddUserToProject(
			cmd.Context(),
			addMemberEmail,
			projectName,
			role,
		); err != nil {
			return Error("failed to add member", err)
		}

		Success(
			"Added " + addMemberEmail +
				" as " + role +
				" to project " + projectName,
		)

		return nil
	},
}

func init() {
	memberCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(
		&addMemberProject,
		"project",
		"",
		"Project name",
	)

	addCmd.Flags().StringVar(
		&addMemberEmail,
		"member-email",
		"",
		"Email address of the user to add",
	)

	addCmd.Flags().StringVar(
		&addMemberRole,
		"member-role",
		"member",
		"Role to assign (admin, member)",
	)

	addCmd.MarkFlagRequired("member-email")
}

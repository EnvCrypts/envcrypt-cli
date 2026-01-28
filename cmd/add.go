package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
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
		// Resolve project name (arg or flag)
		projectName := addMemberProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		// Check if we need to prompts
		// We prompt if critical info is missing
		needsPrompt := false
		if projectName == "" || addMemberEmail == "" {
			needsPrompt = true
		}

		if needsPrompt {
			// Interactive Form
			var role = addMemberRole
			if role == "" {
				role = "member"
			}

			// Build the form
			// We only ask for what's missing, or everything if we are in "interactive mode" logic
			// A simple approach: if entered interactive mode, show fields that are empty, or confirm others?
			// Let's just ask for missing parts.

			// Build the form
			var fields []huh.Field

			if projectName == "" {
				fields = append(fields, huh.NewInput().
					Title("Project Name").
					Value(&projectName).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("project name is required")
						}
						return nil
					}))
			}

			if addMemberEmail == "" {
				fields = append(fields, huh.NewInput().
					Title("Member Email").
					Value(&addMemberEmail).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("email is required")
						}
						return nil
					}))
			}

			// Always offer Role selection in interactive mode
			fields = append(fields, huh.NewSelect[string]().
				Title("Role").
				Options(
					huh.NewOption("Member", "member"),
					huh.NewOption("Admin", "admin"),
				).
				Value(&addMemberRole))

			group := huh.NewGroup(fields...)
			form := huh.NewForm(group)
			err := form.Run()
			if err != nil {
				return Error("cancelled", nil)
			}
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}
		if addMemberEmail == "" {
			return Error("email is required", nil)
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
	rootCmd.AddCommand(addCmd)

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


}

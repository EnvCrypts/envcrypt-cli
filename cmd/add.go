package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	addProject string
	addEmail   string
	addRole    string
)

var addCmd = &cobra.Command{
	Use:          "add [project]",
	Short:        "Add a user to a project",
	Long:         "Add a user to a project and assign them a role.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := addProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		needsPrompt := projectName == "" || addEmail == ""

		if needsPrompt {
			role := addRole
			if role == "" {
				role = "member"
			}

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

			if addEmail == "" {
				fields = append(fields, huh.NewInput().
					Title("Member Email").
					Value(&addEmail).
					Validate(func(str string) error {
						if str == "" {
							return fmt.Errorf("email is required")
						}
						return nil
					}))
			}

			fields = append(fields, huh.NewSelect[string]().
				Title("Role").
				Options(
					huh.NewOption("Member", "member"),
					huh.NewOption("Admin", "admin"),
				).
				Value(&addRole))

			form := huh.NewForm(huh.NewGroup(fields...))
			if err := form.Run(); err != nil {
				return Error("cancelled", nil)
			}
		}

		if projectName == "" {
			return Error("project name is required", nil)
		}
		if addEmail == "" {
			return Error("email is required", nil)
		}

		role := strings.ToLower(addRole)
		if role != "admin" && role != "member" {
			return Error("invalid role (must be admin or member)", nil)
		}

		if err := Application.AddUserToProject(cmd.Context(), addEmail, projectName, role); err != nil {
			return Error("failed to add member", err)
		}

		Success("Added " + addEmail + " as " + role + " to project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&addProject, "project", "", "Project name")
	addCmd.Flags().StringVar(&addEmail, "email", "", "Email address of the user to add")
	addCmd.Flags().StringVar(&addRole, "role", "member", "Role to assign (admin, member)")
}

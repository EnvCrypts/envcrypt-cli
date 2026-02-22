package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	addProject string
	addEmail   string
)

var addCmd = &cobra.Command{
	Use:          "add [project]",
	Short:        "Add a user to a project",
	Long:         "Add a user to a project.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := addProject
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
			projectName, err = tui.RunPicker("Select a project", names)
			if err != nil {
				return tui.Error("cancelled", nil)
			}
		}

		// Prompt for email if not provided
		if addEmail == "" {
			vals, err := tui.RunForm([]tui.FormField{
				{Label: fmt.Sprintf("Member Email for %q", projectName), Required: true},
			}, []string{""})
			if err != nil {
				return tui.Error("cancelled", nil)
			}
			addEmail = vals[0]
		}

		if addEmail == "" {
			return tui.Error("email is required", nil)
		}

		if err := Application.AddUserToProject(cmd.Context(), addEmail, projectName); err != nil {
			return tui.Error("failed to add member", err)
		}

		tui.Success("Added " + addEmail + " to project " + projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&addProject, "project", "", "Project name")
	addCmd.Flags().StringVar(&addEmail, "email", "", "Email address of the user to add")
}

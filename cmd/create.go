package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:          "create [project]",
	Short:        "Create a new project",
	Long:         "Create a new encrypted project.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		if len(args) == 1 {
			projectName = args[0]
		}

		// Prompt for project name if not provided
		if projectName == "" {
			vals, err := tui.RunForm([]tui.FormField{
				{Label: "Project Name", Required: true},
			}, []string{""})
			if err != nil {
				return tui.Error("cancelled", nil)
			}
			projectName = vals[0]
		}

		if projectName == "" {
			return tui.Error("project name is required", nil)
		}

		err := tui.RunActionWithSpinner(fmt.Sprintf("Creating project %q...", projectName), func() error {
			return Application.CreateProject(context.Background(), projectName)
		})
		if err != nil {
			return tui.Error(fmt.Sprintf("failed to create project %q", projectName), err)
		}

		tui.Success(fmt.Sprintf("Project %q created", projectName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

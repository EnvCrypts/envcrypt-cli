package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:           "delete [project]",
	Short:         "Delete a project",
	Long: `Delete a project and all associated encrypted data.

Use --force to skip the confirmation prompt.

Examples:
  envcrypt delete my-project
  envcrypt delete my-project --force
  envcrypt delete`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		if len(args) == 1 {
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
			projectName, err = tui.RunPicker("Select a project to delete", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		if !tui.ConfirmDangerousAction(
			fmt.Sprintf("This will permanently delete project %q and all its data.", projectName),
			projectName,
			deleteForce,
		) {
			return tui.Cancelled()
		}

		err := tui.RunActionWithSpinner(fmt.Sprintf("Deleting project %q...", projectName), func() error {
			return Application.DeleteProject(cmd.Context(), projectName)
		})
		if err != nil {
			return tui.Error(fmt.Sprintf("failed to delete project %q", projectName), err)
		}

		tui.Success(fmt.Sprintf("Project %q deleted", projectName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Delete without confirmation prompt")
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var force bool

var deleteCmd = &cobra.Command{
	Use:           "delete <project>",
	Short:         "Delete a project",
	Long:          "Delete a project and all associated encrypted data.",
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if !force {
			ok := ConfirmDangerousAction(
				fmt.Sprintf(
					"This will permanently delete project %q.",
					projectName,
				),
				projectName,
			)

			if !ok {
				Info("Aborted.")
				return nil
			}
		}

		if err := Application.DeleteProject(cmd.Context(), projectName); err != nil {
			return Error(
				fmt.Sprintf("failed to delete project %q", projectName),
				err,
			)
		}

		Success(fmt.Sprintf("Project %q deleted", projectName))
		return nil
	},
}

func init() {
	projectCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(
		&force,
		"force",
		false,
		"Delete without confirmation",
	)
}

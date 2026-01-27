package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var force bool

// deleteCmd represents the project delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <project>",
	Short:        "Delete a project",
	Long:         "Delete a project and all associated encrypted data.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		if !force {
			fmt.Printf("⚠️  This will permanently delete project %q.\n", projectName)
			fmt.Print("Are you sure? (y/N): ")

			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := Application.DeleteProject(cmd.Context(), projectName); err != nil {
			return fmt.Errorf("failed to delete project %q: %w", projectName, err)
		}

		fmt.Printf("Project %q deleted\n", projectName)
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

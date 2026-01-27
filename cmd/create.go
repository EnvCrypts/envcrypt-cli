package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the project create command
var createCmd = &cobra.Command{
	Use:          "create <project>",
	Short:        "Create a new project",
	Long:         "Create a new encrypted project.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		err := Application.CreateProject(context.Background(), projectName)
		if err != nil {
			return fmt.Errorf("failed to create project %q: %w", projectName, err)
		}

		fmt.Printf("Project %q created\n", projectName)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(createCmd)
}

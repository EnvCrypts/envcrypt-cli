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
			return Error(fmt.Sprintf("failed to create project %q", projectName), err)
		}

		Success(fmt.Sprintf("Project %q created\n", projectName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

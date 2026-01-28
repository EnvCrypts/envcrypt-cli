package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the project list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List projects",
	Long:         "List all projects you have access to.",
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		projectResp, err := Application.ListProjects(cmd.Context())
		if err != nil {
			return err
		}

		PrintProjects(projectResp.Projects)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

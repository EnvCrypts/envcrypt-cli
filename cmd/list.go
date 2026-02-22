package cmd

import "github.com/envcrypts/envcrypt-cli/internal/tui"

import "github.com/spf13/cobra"

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

		tui.RunProjectsTable(projectResp.Projects)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

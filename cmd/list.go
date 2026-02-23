package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long: `List all projects you have access to.

Examples:
  envcrypt list
  envcrypt list --json
  envcrypt list --no-table`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectResp, err := Application.ListProjects(cmd.Context())
		if err != nil {
			return tui.Error("failed to list projects", err)
		}

		return tui.RunProjectsTable(projectResp.Projects)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

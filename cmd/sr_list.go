package cmd

import (
	"context"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var serviceRoleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List service roles",
	Long: `List all service roles you have access to.

Examples:
  envcrypt service-role list
  envcrypt service-role list --json
  envcrypt service-role list --no-table`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		serviceRoles, err := Application.ListServiceRoles(context.Background())
		if err != nil {
			return tui.Error("failed to list service roles", err)
		}

		return tui.RunServiceRolesTable(serviceRoles)
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleListCmd)
}

package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"

	"context"

	"github.com/spf13/cobra"
)

// list
var serviceRoleListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List service roles",
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {

		serviceRoles, err := Application.ListServiceRoles(context.Background())
		if err != nil {
			return err
		}

		tui.RunServiceRolesTable(serviceRoles)

		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleListCmd)
}

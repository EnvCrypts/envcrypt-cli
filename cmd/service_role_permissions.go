package cmd

import (
	"github.com/spf13/cobra"
)

// permissions
var serviceRolePermissionsCmd = &cobra.Command{
	Use:   "permissions <name>",
	Short: "View what a service role can access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPrincipal := args[0]
		
		perm, err := Application.GetPermissions(cmd.Context(), repoPrincipal)
		if err != nil {
			return err
		}

		PrintServiceRolePermissions(perm, repoPrincipal)

		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRolePermissionsCmd)
}

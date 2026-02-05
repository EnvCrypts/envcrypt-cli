package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// revoke
var serviceRoleRevokeCmd = &cobra.Command{
	Use:   "revoke <name>",
	Short: "Revoke a service role (breaks all CI access)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		// Logic to revoke service role would go here
		
		Success(fmt.Sprintf("Service role %q revoked", name))
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleRevokeCmd)
}

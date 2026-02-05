package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// revoke-access
var serviceRoleRevokeAccessCmd = &cobra.Command{
	Use:   "revoke-access",
	Short: "Revoke access",
	Long: `Revoke CI access.

Example:
  envcrypt service-role revoke-access \
    --service-role sp-billing-backend \
    --project billing-service \
    --env prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("service-role")
		project, _ := cmd.Flags().GetString("project")
		env, _ := cmd.Flags().GetString("env")

		// Logic to revoke access would go here

		Success(fmt.Sprintf("Revoked access to %q from service role %q on env %q", project, role, env))
		return nil
	},
}

func init() {
	serviceRoleRevokeAccessCmd.Flags().String("service-role", "", "Service role name (required)")
	serviceRoleRevokeAccessCmd.Flags().String("project", "", "Project name (required)")
	serviceRoleRevokeAccessCmd.Flags().String("env", "", "Environment name (required)")
	serviceRoleRevokeAccessCmd.MarkFlagRequired("service-role")
	serviceRoleRevokeAccessCmd.MarkFlagRequired("project")
	serviceRoleRevokeAccessCmd.MarkFlagRequired("env")
	serviceRoleCmd.AddCommand(serviceRoleRevokeAccessCmd)
}

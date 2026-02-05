package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// grant
var serviceRoleGrantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant CI access to a project/env",
	Long: `Grant CI access to a project/env.

Example:
  envcrypt service-role grant \
    --service-role sp-billing-backend \
    --project billing-service \
    --env prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("service-role")
		project, _ := cmd.Flags().GetString("project")
		env, _ := cmd.Flags().GetString("env")

		// Logic to grant permissions would go here

		Success(fmt.Sprintf("Granted access to %q for service role %q on env %q", project, role, env))
		return nil
	},
}

func init() {
	serviceRoleGrantCmd.Flags().String("service-role", "", "Service role name (required)")
	serviceRoleGrantCmd.Flags().String("project", "", "Project name (required)")
	serviceRoleGrantCmd.Flags().String("env", "", "Environment name (required)")
	serviceRoleGrantCmd.MarkFlagRequired("service-role")
	serviceRoleGrantCmd.MarkFlagRequired("project")
	serviceRoleGrantCmd.MarkFlagRequired("env")
	serviceRoleCmd.AddCommand(serviceRoleGrantCmd)
}

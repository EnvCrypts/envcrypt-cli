package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// create
var serviceRoleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service role",
	Long: `Create a new service role.
  
Example:
  envcrypt service-role create \
    --repo github:acme/billing-backend:ref:refs/heads/main \
    --name sp-billing-backend`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		repo, _ := cmd.Flags().GetString("repo")
		
		// Logic to create service role would go here
		
		Success(fmt.Sprintf("Service role %q created for repo %q", name, repo))
		return nil
	},
}

func init() {
	serviceRoleCreateCmd.Flags().String("repo", "", "Repository identifier (required)")
	serviceRoleCreateCmd.Flags().String("name", "", "Name of the service role (required)")
	serviceRoleCreateCmd.MarkFlagRequired("repo")
	serviceRoleCreateCmd.MarkFlagRequired("name")
	serviceRoleCmd.AddCommand(serviceRoleCreateCmd)
}

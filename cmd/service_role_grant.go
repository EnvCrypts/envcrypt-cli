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
		roleName, _ := cmd.Flags().GetString("service-role")
		project, _ := cmd.Flags().GetString("project")
		env, _ := cmd.Flags().GetString("env")

		if project == "" {
			projectsResp, err := Application.ListProjects(cmd.Context())
			if err != nil {
				return err
			}

			PrintProjects(projectsResp.Projects)
			fmt.Print("Enter project name: ")
			fmt.Scanln(&project)
		}

		if env == "" {
			fmt.Print("Enter environment (e.g., prod, dev): ")
			fmt.Scanln(&env)
		}
		
		if err := Application.DelegateAccess(cmd.Context(), roleName, project, env); err != nil {
			return err
		}

		Success(fmt.Sprintf("Granted access to %q for service role %q on env %q", project, roleName, env))
		return nil
	},
}

func init() {
	serviceRoleGrantCmd.Flags().String("service-role", "", "Service role name (required)")
	serviceRoleGrantCmd.Flags().String("project", "", "Project name (required)")
	serviceRoleGrantCmd.Flags().String("env", "", "Environment name (required)")
	serviceRoleGrantCmd.MarkFlagRequired("service-role")
	serviceRoleCmd.AddCommand(serviceRoleGrantCmd)
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/config"
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
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		roleName, _ := cmd.Flags().GetString("service-role")
		project, _ := cmd.Flags().GetString("project")
		env, _ := cmd.Flags().GetString("env")

		if roleName == "" {
			defPrincipal, _, _, _ := DetectGitContext()
			roleName = PromptWithDefault("Service Role Principal", defPrincipal)
		}

		if project == "" {
			projectsResp, err := Application.ListProjects(cmd.Context())
			if err != nil {
				return err
			}


			var adminProjects []config.Project
			for _, p := range projectsResp.Projects {
				if strings.ToLower(p.Role) == "admin" {
					adminProjects = append(adminProjects, p)
				}
			}

			PrintProjects(adminProjects)
			fmt.Print("Enter project name: ")
			fmt.Scanln(&project)

			found := false
			for _, p := range adminProjects {
				if p.Name == project {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("project %q not found or access denied (admin role required)", project)
			}
		}

		if env == "" {
			fmt.Print("Enter environment (e.g., prod, dev): ")
			fmt.Scanln(&env)
		}

		if roleName == "" {
			return fmt.Errorf("service-role is required (could not auto-detect)")
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
	serviceRoleCmd.AddCommand(serviceRoleGrantCmd)
}

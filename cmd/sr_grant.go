package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

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

		// Prompt for service role if not provided
		if roleName == "" {
			defPrincipal, _, _, _ := DetectGitContext()
			vals, err := tui.RunForm([]tui.FormField{
				{Label: "Service Role Principal", Required: true},
			}, []string{defPrincipal})
			if err != nil {
				return tui.Error("cancelled", nil)
			}
			roleName = vals[0]
		}

		if roleName == "" {
			return tui.Error("service-role is required (could not auto-detect)", nil)
		}

		// Pick project from admin-only list if not provided
		if project == "" {
			projectsResp, err := Application.ListProjects(context.Background())
			if err != nil {
				return tui.Error("failed to fetch projects", err)
			}

			var adminNames []string
			for _, p := range projectsResp.Projects {
				if strings.EqualFold(p.Role, "admin") {
					adminNames = append(adminNames, p.Name)
				}
			}

			if len(adminNames) == 0 {
				return tui.Error("no admin projects found", nil)
			}

			project, err = tui.RunPicker("Select a project (admin only)", adminNames)
			if err != nil {
				return tui.Error("cancelled", nil)
			}
		}

		// Pick env if not provided
		if env == "" {
			picked, err := tui.RunEnvPicker(project)
			if err != nil {
				return tui.Error("cancelled", nil)
			}
			env = picked
		}

		if err := Application.DelegateAccess(cmd.Context(), roleName, project, env); err != nil {
			return tui.Error("failed to grant access", err)
		}

		tui.Success(fmt.Sprintf("Granted %s/%s access to service role %q", project, env, roleName))
		return nil
	},
}

func init() {
	serviceRoleGrantCmd.Flags().String("service-role", "", "Service role name")
	serviceRoleGrantCmd.Flags().String("project", "", "Project name")
	serviceRoleGrantCmd.Flags().String("env", "", "Environment name")
	serviceRoleCmd.AddCommand(serviceRoleGrantCmd)
}

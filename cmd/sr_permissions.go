package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var serviceRolePermissionsCmd = &cobra.Command{
	Use:   "permissions <name>",
	Short: "View what a service role can access",
	Long: `View the project environments a service role has access to.

Examples:
  envcrypt service-role permissions github:acme/backend:ref:refs/heads/main
  envcrypt service-role permissions`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var repoPrincipal string
		if len(args) > 0 {
			repoPrincipal = args[0]
		}

		if repoPrincipal == "" {
			defPrincipal, _, _, _ := DetectGitContext()
			repoPrincipal = tui.PromptWithDefault("Service Role Principal", defPrincipal)
		}

		if repoPrincipal == "" {
			return tui.Error("service role principal is required", nil)
		}

		perm, err := Application.GetPermissions(cmd.Context(), repoPrincipal)
		if err != nil {
			return tui.Error("failed to fetch permissions", err)
		}

		tui.PrintServiceRolePermissions(perm, repoPrincipal)
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRolePermissionsCmd)
}

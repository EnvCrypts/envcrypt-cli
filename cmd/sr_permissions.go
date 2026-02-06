package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// permissions
var serviceRolePermissionsCmd = &cobra.Command{
	Use:          "permissions <name>",
	Short:        "View what a service role can access",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var repoPrincipal string
		if len(args) > 0 {
			repoPrincipal = args[0]
		}

		if repoPrincipal == "" {
			defPrincipal, _, _, _ := DetectGitContext()
			repoPrincipal = PromptWithDefault("Service Role Principal", defPrincipal)
		}

		if repoPrincipal == "" {
			return fmt.Errorf("service role principal is required")
		}

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

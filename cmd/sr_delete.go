package cmd

import (
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var srDeleteForce bool

var serviceRoleDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a service role (rare)",
	Long: `Delete a service role and revoke all of its access.

Examples:
  envcrypt service-role delete github:acme/backend:ref:refs/heads/main
  envcrypt service-role delete --force`,
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

		role, err := Application.GetServiceRole(cmd.Context(), repoPrincipal)
		if err != nil {
			return tui.Error("failed to fetch service role", err)
		}

		if !tui.ConfirmDangerousAction(
			fmt.Sprintf("Are you sure you want to delete service role %q?", role.Name),
			role.Name,
			srDeleteForce,
		) {
			return tui.Cancelled()
		}

		if err := Application.DeleteServiceRole(cmd.Context(), role.ID); err != nil {
			return tui.Error("failed to delete service role", err)
		}

		tui.Success(fmt.Sprintf("Service role %q deleted", role.Name))
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleDeleteCmd)
	serviceRoleDeleteCmd.Flags().BoolVar(&srDeleteForce, "force", false, "Delete without confirmation prompt")
}

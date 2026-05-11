package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var serviceRoleGetCmd = &cobra.Command{
	Use:   "get <repo_identifier>",
	Short: "Show one service role",
	Long: `Show details of a specific service role.

Examples:
  envcrypt service-role get github:acme/backend:ref:refs/heads/main
  envcrypt service-role get`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		var repoPrincipal string
		if len(args) > 0 {
			repoPrincipal = args[0]
		}

		if repoPrincipal == "" {
			if !tui.IsInteractive() {
				return tui.Error("service role principal is required in non-interactive mode", nil, "Use --principal or pass the principal as an argument")
			}
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

		tui.PrintServiceRoleDetail(role)
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleGetCmd)
}

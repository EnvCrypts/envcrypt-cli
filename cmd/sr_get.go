package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"

	"fmt"

	"github.com/spf13/cobra"
)

// get
var serviceRoleGetCmd = &cobra.Command{
	Use:          "get <repo_identifier>",
	Short:        "Show one service role",
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
			return fmt.Errorf("service role principal is required")
		}

		role, err := Application.GetServiceRole(cmd.Context(), repoPrincipal)
		if err != nil {
			return err
		}

		tui.PrintServiceRoleDetail(role)

		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleGetCmd)
}

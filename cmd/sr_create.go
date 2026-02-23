package cmd

import (
	"context"
	"fmt"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var serviceRoleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service role",
	Long: `Create a new service role for CI/CD automation.
  
Examples:
  envcrypt service-role create \
    --repo github:acme/billing-backend:ref:refs/heads/main \
    --name sp-billing-backend
  envcrypt service-role create`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")

		var principal string

		if repo != "" && branch != "" {
			principal = buildRepoPrincipal(repo, branch)
		} else {
			// Try auto-detect for defaults
			_, defRepo, defBranch, _ := DetectGitContext()

			// Prompt if flags weren't provided
			if repo == "" {
				repo = tui.PromptWithDefault("Repository (e.g. acme/backend)", defRepo)
			}
			if branch == "" {
				branch = tui.PromptWithDefault("Branch (e.g. main)", defBranch)
			}

			if repo == "" || branch == "" {
				return tui.Error("repo and branch are required", nil)
			}

			principal = buildRepoPrincipal(repo, branch)

			tui.Info(fmt.Sprintf("Creating service role for principal: %s", principal))

			if !tui.ConfirmDangerousAction(fmt.Sprintf("Create service role %q?", name), "yes") {
				return tui.Cancelled()
			}
		}

		keyPair, err := Application.CreateServiceRole(context.Background(), name, principal)
		if err != nil {
			return tui.Error("failed to create service role", err)
		}

		tui.PrintServiceRoleSecret(keyPair)
		return nil
	},
}

func init() {
	serviceRoleCreateCmd.Flags().String("repo", "", "Repository identifier (e.g. acme/backend)")
	serviceRoleCreateCmd.Flags().String("branch", "", "Branch name (e.g. main)")
	serviceRoleCreateCmd.Flags().String("name", "", "Name of the service role (required)")
	serviceRoleCreateCmd.MarkFlagRequired("name")
	serviceRoleCmd.AddCommand(serviceRoleCreateCmd)
}

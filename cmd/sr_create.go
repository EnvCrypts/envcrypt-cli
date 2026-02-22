package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"

	"context"
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

			// If user didn't provide repo/branch via flags, prompt them
			if repo == "" {
				repo = tui.PromptWithDefault("Repository (e.g. acme/backend)", defRepo)
			}
			if branch == "" {
				branch = tui.PromptWithDefault("Branch (e.g. main)", defBranch)
			}

			if repo == "" || branch == "" {
				return fmt.Errorf("repo and branch are required")
			}

			principal = buildRepoPrincipal(repo, branch)

			// Show what we are about to create
			tui.Info(fmt.Sprintf("Creating service role for principal: %s", principal))

			// Confirm action
			if !tui.ConfirmDangerousAction(fmt.Sprintf("Create service role %q?", name), "yes") {
				return fmt.Errorf("cancelled")
			}
		}

		keyPair, err := Application.CreateServiceRole(context.Background(), name, principal)
		if err != nil {
			return err
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

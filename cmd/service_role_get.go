package cmd

import (
	"github.com/spf13/cobra"
)

// get
var serviceRoleGetCmd = &cobra.Command{
	Use:   "get <repo_identifier>",
	Short: "Show one service role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPrincipal := args[0]
		
		role, err := Application.GetServiceRole(cmd.Context(), repoPrincipal)
		if err != nil {
			return err
		}

		PrintServiceRoleDetail(role)

		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleGetCmd)
}

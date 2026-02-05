package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// delete
var serviceRoleDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a service role (rare)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPrincipal := args[0]
		
		role, err := Application.GetServiceRole(cmd.Context(), repoPrincipal)
		if err != nil {
			return err
		}

		if !ConfirmDangerousAction(fmt.Sprintf("Are you sure you want to delete service role %q?", role.Name), role.Name) {
			return nil
		}

		if err := Application.DeleteServiceRole(cmd.Context(), role.ID); err != nil {
			return err
		}

		Success(fmt.Sprintf("Service role %q deleted", role.Name))
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleDeleteCmd)
}

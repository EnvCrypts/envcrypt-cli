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
		name := args[0]
		// Logic to delete service role would go here

		Success(fmt.Sprintf("Service role %q deleted", name))
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleDeleteCmd)
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// permissions
var serviceRolePermissionsCmd = &cobra.Command{
	Use:   "permissions <name>",
	Short: "View what a service role can access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		// Logic to list permissions would go here
		
		Info(fmt.Sprintf("Permissions for service role %q:", name))
		fmt.Printf(
			"%s  %s\n",
			headerStyle.Render(padRight("PROJECT", 20)),
			headerStyle.Render(padRight("ENV", 10)),
		)
		fmt.Printf(
			"%s  %s\n",
			padRight("billing-service", 20),
			padRight("prod", 10),
		)

		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRolePermissionsCmd)
}

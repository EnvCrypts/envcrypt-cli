package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// list
var serviceRoleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List service roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Logic to list service roles would go here
		
		fmt.Printf(
			"%s  %s\n",
			headerStyle.Render(padRight("NAME", 30)),
			headerStyle.Render(padRight("REPO", 50)),
		)
		// Placeholder data
		fmt.Printf(
			"%s  %s\n",
			padRight("sp-billing-backend", 30),
			padRight("github:acme/billing-backend:ref:refs/heads/main", 50),
		)
		
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleListCmd)
}

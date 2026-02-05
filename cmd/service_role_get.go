package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// get
var serviceRoleGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Show one service role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		// Logic to get service role would go here
		
		Info(fmt.Sprintf("Details for service role %q:", name))
		fmt.Println(mutedStyle.Render("  repo: github:acme/billing-backend:ref:refs/heads/main"))
		fmt.Println(mutedStyle.Render("  created_at: 2023-10-27T10:00:00Z"))
		
		return nil
	},
}

func init() {
	serviceRoleCmd.AddCommand(serviceRoleGetCmd)
}

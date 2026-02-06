package cmd

import (
	"github.com/spf13/cobra"
)

// serviceRoleCmd represents the service-role command
var serviceRoleCmd = &cobra.Command{
	Use:   "service-role",
	Short: "Manage service roles (admin only)",
	Long:  "Manage service roles for automated access to projects.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(serviceRoleCmd)
}

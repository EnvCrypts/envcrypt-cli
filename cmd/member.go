package cmd

import "github.com/spf13/cobra"

// memberCmd represents the project member command
var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage project members and permissions",
	Long: `Manage users who have access to a project.

Use subcommands to add or remove members, assign roles,
and view current project members.`,
}

func init() {
	projectCmd.AddCommand(memberCmd)
}

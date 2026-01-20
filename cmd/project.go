package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage encrypted projects",
	Long:  "Manage isolated projects used for securely storing and versioning encrypted configuration data.",
}

func init() {
	rootCmd.AddCommand(projectCmd)
}

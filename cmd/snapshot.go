package cmd

import "github.com/spf13/cobra"

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Export and import project snapshots",
	Long:  "Commands to securely export a snapshot of a project and import it into a new project.",
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

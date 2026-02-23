package cmd

import (
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long: `Print the current version of the EnvCrypt CLI.

Examples:
  envcrypt version
  envcrypt version --json`,
	Run: func(cmd *cobra.Command, args []string) {
		if tui.Mode() == tui.ModeJSON {
			tui.JSONData(map[string]string{"version": Version})
			return
		}
		tui.Info("envcrypt version " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"context"

	"github.com/envcrypts/envcrypt-cli/internal/client"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var snapshotImportFilename string

var snapshotImportCmd = &cobra.Command{
	Use:   "import [new_project]",
	Short: "Import a project snapshot",
	Long: `Import a previously exported project snapshot as a new project.

Examples:
  envcrypt snapshot import new-project-name
  envcrypt snapshot import new-project-name --file backup.json
  envcrypt snapshot import`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newProjectName := ""
		if len(args) == 1 {
			newProjectName = args[0]
		}

		if newProjectName == "" {
			vals, err := tui.RunForm([]tui.FormField{{Label: "New Project Name", Required: true, Validate: tui.ValidateProjectName}}, nil)
			if err != nil || len(vals) == 0 || vals[0] == "" {
				return tui.Cancelled()
			}
			newProjectName = vals[0]
		}

		filename := snapshotImportFilename
		if filename == "" {
			vals, err := tui.RunForm([]tui.FormField{{Label: "Filename to import from", Required: true, Validate: tui.ValidateFileExists}}, []string{newProjectName + ".json"})
			if err != nil || len(vals) == 0 || vals[0] == "" {
				return tui.Cancelled()
			}
			filename = vals[0]
		}

		var newID string
		err := tui.RunActionWithSpinner("Importing snapshot...", func() error {
			var importErr error
			newID, importErr = Application.ImportSnapshot(context.Background(), newProjectName, filename)
			return importErr
		})

		if err != nil {
			if httpErr, ok := err.(*client.HTTPError); ok {
				return tui.MapAPIError(&tui.APIErrorDetail{
					Code:    httpErr.Code,
					Message: httpErr.Message,
					Hint:    httpErr.Hint,
				})
			}
			return tui.Error("Failed to import snapshot", err)
		}

		tui.Success("Snapshot restored successfully. New project ID: " + newID)
		return nil
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotImportCmd)
	snapshotImportCmd.Flags().StringVarP(&snapshotImportFilename, "file", "f", "", "Filename to import the snapshot from")
}

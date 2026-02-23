package cmd

import (
	"context"
	"net/http"

	"github.com/envcrypts/envcrypt-cli/internal/client"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var snapshotImportFilename string

var snapshotImportCmd = &cobra.Command{
	Use:   "import [new_project_name]",
	Short: "Import a project snapshot",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newProjectName := ""
		if len(args) == 1 {
			newProjectName = args[0]
		}

		if newProjectName == "" {
			vals, err := tui.RunForm([]tui.FormField{{Label: "New Project Name", Required: true}}, nil)
			if err != nil || len(vals) == 0 || vals[0] == "" {
				return tui.Error("cancelled", nil)
			}
			newProjectName = vals[0]
		}

		filename := snapshotImportFilename
		if filename == "" {
			vals, err := tui.RunForm([]tui.FormField{{Label: "Filename to import from", Required: true}}, []string{newProjectName + ".json"})
			if err != nil || len(vals) == 0 || vals[0] == "" {
				return tui.Error("cancelled", nil)
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
				switch httpErr.Status {
				case http.StatusBadRequest:
					return tui.Error("Snapshot validation failed (checksum mismatch or malformed file).", nil)
				case http.StatusForbidden:
					return tui.Error("Permission denied.", nil)
				case http.StatusConflict:
					return tui.Error("Project conflict detected.", nil)
				}
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

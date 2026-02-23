package cmd

import (
	"context"
	"net/http"

	"github.com/envcrypts/envcrypt-cli/internal/client"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var snapshotExportFilename string

var snapshotExportCmd = &cobra.Command{
	Use:   "export [project_name]",
	Short: "Export a project snapshot",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		if len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			projectsResp, err := Application.ListProjects(context.Background())
			if err != nil {
				return tui.Error("failed to fetch projects", err)
			}
			
			var adminProjects []string
			for _, p := range projectsResp.Projects {
				if p.Role == "admin" {
					adminProjects = append(adminProjects, p.Name)
				}
			}
			if len(adminProjects) == 0 {
				return tui.Error("no projects found where you have 'admin' access", nil)
			}

			projectName, err = tui.RunPicker("Select a project to export", adminProjects)
			if err != nil {
				return tui.Error("cancelled", nil)
			}
		}

		filename := snapshotExportFilename
		if filename == "" {
			vals, err := tui.RunForm([]tui.FormField{{Label: "Filename to export to", Required: true}}, []string{projectName + ".json"})
			if err != nil || len(vals) == 0 || vals[0] == "" {
				return tui.Error("cancelled", nil)
			}
			filename = vals[0]
		}

		
		var absPath string
		err := tui.RunActionWithSpinner("Exporting snapshot...", func() error {
			var exportErr error
			absPath, exportErr = Application.ExportSnapshot(context.Background(), projectName, filename)
			return exportErr
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
			return tui.Error("Failed to export snapshot", err)
		}

		tui.Success("Snapshot exported successfully to " + absPath)
		return nil
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotExportCmd)
	snapshotExportCmd.Flags().StringVarP(&snapshotExportFilename, "file", "f", "", "Filename to export the snapshot to")
}

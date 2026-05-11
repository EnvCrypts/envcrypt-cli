package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/envcrypts/envcrypt-cli/internal/client"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	rotateForce bool
)

var rotateCmd = &cobra.Command{
	Use:   "rotate [project]",
	Short: "Rotate a project's Root Key (PRK) and rewrap Data Encryption Keys (DEKs)",
	Long: `Performs a client-side rotation of the Project Root Key (PRK) without exposing plaintext keys to the server.

Use --force to skip the confirmation prompt.

Examples:
  envcrypt rotate my-project
  envcrypt rotate --force
  envcrypt rotate`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := ""
		if len(args) == 1 {
			projectName = args[0]
		}

		projectsResp, err := Application.ListProjects(context.Background())
		if err != nil {
			return tui.Error("failed to fetch projects", err)
		}

		// Filter list to admin projects only
		var adminProjects []string
		var projectID string

		for _, p := range projectsResp.Projects {
			if p.Role == "admin" {
				adminProjects = append(adminProjects, p.Name)
				if p.Name == projectName {
					projectID = p.Id.String()
				}
			}
		}

		// Pick project from list if not provided
		if projectName == "" {
			if len(adminProjects) == 0 {
				return tui.Error("no projects found where you have 'admin' access", nil)
			}
			projectName, err = tui.RunPicker("Select a project to rotate its PRK", adminProjects)
			if err != nil {
				return handlePromptError(err, "project is required in non-interactive mode", "Use --project to select a project")
			}
			for _, p := range projectsResp.Projects {
				if p.Name == projectName {
					projectID = p.Id.String()
				}
			}
		}

		if projectID == "" {
			return tui.Error(fmt.Sprintf("you do not have admin access to project '%s' or it does not exist", projectName), nil)
		}

		if !tui.ConfirmDangerousAction(fmt.Sprintf("Rotate PRK for project %q?", projectName), "yes", rotateForce) {
			return tui.Cancelled()
		}

		var newVersion int32

		err = tui.RunActionWithSpinner(fmt.Sprintf("Rotating PRK for project %q...", projectName), func() error {
			ver, rotateErr := Application.RotatePRK(context.Background(), projectID)
			if rotateErr != nil {
				if httpErr, ok := rotateErr.(*client.HTTPError); ok && httpErr.Status == http.StatusConflict {
					tui.Warn("PRK version changed during rotation. Retrying automatically...")
					ver, rotateErr = Application.RotatePRK(context.Background(), projectID)
					if rotateErr != nil {
						return rotateErr
					}
					newVersion = ver
					return nil
				}
				return rotateErr
			}
			newVersion = ver
			return nil
		})

		if err != nil {
			return tui.Error(fmt.Sprintf("failed to rotate PRK for project %q", projectName), err)
		}

		tui.Success(fmt.Sprintf("Successfully rotated PRK for project %q (New Version: %d)", projectName, newVersion))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rotateCmd)
	rotateCmd.Flags().BoolVar(&rotateForce, "force", false, "Skip confirmation prompt")
}

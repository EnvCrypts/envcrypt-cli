package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	pullProject string
	pullEnvName string
	pullEnvFile string
	pullYes     bool
)

var pullCmd = &cobra.Command{
	Use:          "pull [project]",
	Short:        "Download and decrypt environment variables",
	Long: `Download environment variables from a project and write them to a .env file.
If no project name is provided, the selector defaults to the current git repository name.

Examples:
  envcrypt pull my-project --env prod
  envcrypt pull --project my-project --env dev --env-file .env.local
  envcrypt pull my-project --env prod -y`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := pullProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}

		// Pick project from list if not provided
		if projectName == "" {
			resp, err := Application.ListProjects(context.Background())
			if err != nil {
				return tui.Error("failed to fetch projects", err)
			}
			names := make([]string, len(resp.Projects))
			for i, p := range resp.Projects {
				names[i] = p.Name
			}
			defaultProject := ""
			repo, errGit := getRepoFromGit()
			if errGit == nil && repo != "" {
				defaultProject = path.Base(repo)
			}

			projectName, err = tui.RunPickerWithDefault("Select a project to pull from", names, defaultProject)
			if err != nil {
				return tui.Cancelled()
			}
		}

		// Pick env if not provided via flag
		envName := pullEnvName
		if envName == "" {
			picked, err := tui.RunEnvPicker(projectName)
			if err != nil {
				return tui.Cancelled()
			}
			envName = picked
		}

		envPath := pullEnvFile
		if envPath == "" {
			envPath = ".env"
		}

		tui.Info("Project: " + projectName)
		tui.Info("Environment: " + envName)

		if fileExists(envPath) && !pullYes {
			if !tui.ConfirmOverwrite(envPath) {
				return tui.Cancelled()
			}
		}

		var envMap map[string]string
		err := tui.RunActionWithSpinner(
			fmt.Sprintf("Pulling %s/%s...", projectName, envName),
			func() error {
				var e error
				envMap, e = Application.PullEnv(cmd.Context(), projectName, envName)
				return e
			},
		)
		if err != nil {
			return tui.Error("failed to pull environment variables", err)
		}

		if len(envMap) == 0 {
			tui.Info(fmt.Sprintf("No environment variables found for %s. Creating empty .env file.", envName))
		}

		tui.PrintEnvSummary(envMap)

		envBytes, err := cryptoutils.EncodeEnv(envMap)
		if err != nil {
			return tui.Error("failed to encode env file", err)
		}

		if err := os.WriteFile(envPath, envBytes, 0600); err != nil {
			return tui.Error(
				"failed to write env file",
				fmt.Errorf("could not write to %q: %w", envPath, err))
		}

		tui.Success(fmt.Sprintf("Pulled %s/%s → %s", projectName, envName, envPath))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVar(&pullProject, "project", "", "Project name")
	pullCmd.Flags().StringVar(&pullEnvName, "env", "", "Environment name (dev, staging, prod)")
	pullCmd.Flags().StringVarP(&pullEnvFile, "env-file", "e", "", "Path to write .env file (default: ./.env)")
	pullCmd.Flags().BoolVarP(&pullYes, "yes", "y", false, "Skip confirmation when overwriting .env file")
}

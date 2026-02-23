package cmd

import (
	"context"
	"fmt"
	"os"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	pushProject string
	pushEnvName string
	pushEnvFile string
)

var pushCmd = &cobra.Command{
	Use:          "push [project]",
	Short:        "Encrypt and upload environment variables",
	Long: `Encrypt variables from a .env file and upload them to a project environment.

Examples:
  envcrypt push my-project --env prod
  envcrypt push --project my-project --env dev --env-file .env.local
  envcrypt push`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := pushProject
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
			projectName, err = tui.RunPicker("Select a project to push to", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		// Pick env if not provided via flag
		envName := pushEnvName
		if envName == "" {
			picked, err := tui.RunEnvPicker(projectName)
			if err != nil {
				return tui.Cancelled()
			}
			envName = picked
		}

		envPath, err := resolveEnvFile(pushEnvFile)
		if err != nil {
			return tui.Error("failed to load env file", err)
		}

		tui.Info("Loaded " + envPath)
		tui.Info("Environment: " + envName)

		fileData, err := os.ReadFile(envPath)
		if err != nil {
			return tui.Error("failed to read env file", mapEnvReadError(envPath, err))
		}

		envMap, err := cryptoutils.ParseEnv(fileData)
		if err != nil {
			return tui.Error("failed to parse env file", mapEnvReadError(envPath, err))
		}
		if len(envMap) == 0 {
			return tui.Error(
				"no environment variables found",
				fmt.Errorf("env file %q is empty or contains only comments", envPath))
		}

		tui.PrintEnvSummary(envMap)

		err = tui.RunActionWithSpinner(
			fmt.Sprintf("Uploading to %s/%s...", projectName, envName),
			func() error {
				return Application.PushEnv(cmd.Context(), projectName, envName, envMap)
			},
		)
		if err != nil {
			return tui.Error("failed to upload environment variables", err)
		}

		tui.Success(fmt.Sprintf("Uploaded environment variables to %s/%s", projectName, envName))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&pushProject, "project", "", "Project name")
	pushCmd.Flags().StringVar(&pushEnvName, "env", "", "Environment name (dev, staging, prod)")
	pushCmd.Flags().StringVarP(&pushEnvFile, "env-file", "e", "", "Path to .env file (default: ./.env)")
}

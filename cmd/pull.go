package cmd

import (
	"fmt"
	"os"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
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
	Long:         "Download environment variables from a project and write them to a .env file.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := pullProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}
		if projectName == "" {
			return Error("project name is required", nil)
		}

		envName := pullEnvName
		if envName == "" {
			envName = "dev"
		}

		envPath := pullEnvFile
		if envPath == "" {
			envPath = ".env"
		}

		Info("Project: " + projectName)
		Info("Environment: " + envName)

		if fileExists(envPath) && !pullYes {
			if !ConfirmOverwrite(envPath) {
				return nil
			}
		}

		envMap, err := Application.PullEnv(
			cmd.Context(),
			projectName,
			envName,
		)
		if err != nil {
			return Error("failed to pull environment variables", err)
		}

		if len(envMap) == 0 {
			Info(fmt.Sprintf("No environment variables found for %s. Creating empty .env file.", envName))
		}

		printEnvSummary(envMap)

		envBytes, err := cryptoutils.EncodeEnv(envMap)
		if err != nil {
			return Error("failed to encode env file", err)
		}

		if err := os.WriteFile(envPath, envBytes, 0600); err != nil {
			return Error(
				"failed to write env file",
				fmt.Errorf("could not write to %q: %w", envPath, err),
			)
		}

		Success(
			fmt.Sprintf(
				"Pulled environment variables to %s/%s (%s)",
				projectName,
				envName,
				envPath,
			),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVar(&pullProject, "project", "", "Project name")
	pullCmd.Flags().StringVar(&pullEnvName, "env", "dev", "Environment name (dev, staging, prod)")
	pullCmd.Flags().StringVarP(&pullEnvFile, "env-file", "e", "", "Path to write .env file (default: ./.env)")
	pullCmd.Flags().BoolVarP(&pullYes, "yes", "y", false, "Skip confirmation when overwriting .env file")
}

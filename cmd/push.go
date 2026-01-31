package cmd

import (
	"fmt"
	"os"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
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
	Long:         "Encrypt variables from a .env file and upload them to a project environment.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := pushProject
		if projectName == "" && len(args) == 1 {
			projectName = args[0]
		}
		if projectName == "" {
			return Error("project name is required", nil)
		}

		envName := pushEnvName
		if envName == "" {
			envName = "dev"
		}

		envPath, err := resolveEnvFile(pushEnvFile)
		if err != nil {
			return Error("failed to load env file", err)
		}

		Info("Loaded " + envPath)
		Info("Environment: " + envName)

		fileData, err := os.ReadFile(envPath)
		if err != nil {
			return Error("failed to read env file", mapEnvReadError(envPath, err))
		}

		envMap, err := cryptoutils.ParseEnv(fileData)
		if err != nil {
			return Error("failed to parse env file", mapEnvReadError(envPath, err))
		}
		if len(envMap) == 0 {
			return Error(
				"no environment variables found",
				fmt.Errorf("env file %q is empty or contains only comments", envPath),
			)
		}
		printEnvSummary(envMap)

		if err := Application.PushEnv(
			cmd.Context(),
			projectName,
			envName,
			envMap,
		); err != nil {
			return Error("failed to upload environment variables", err)
		}

		Success(
			fmt.Sprintf(
				"Uploaded environment variables to %s/%s",
				projectName,
				envName,
			),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&pushProject, "project", "", "Project name")
	pushCmd.Flags().StringVar(&pushEnvName, "env", "dev", "Environment name (dev, staging, prod)")
	pushCmd.Flags().StringVarP(&pushEnvFile, "env-file", "e", "", "Path to .env file (default: ./.env)")
}

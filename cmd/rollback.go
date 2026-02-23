package cmd

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	rollbackProject string
	rollbackEnv     string
	rollbackVer     int
	rollbackForce   bool
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [version]",
	Short: "Rollback to a previous version of an environment",
	Long: `Rollback an environment to a specific version.

This command will create a new version that is an exact copy of the specified previous version.
You will see a diff of the changes before confirming the rollback unless --force is used.

Examples:
  envcrypt rollback 3 --project my-project --env prod
  envcrypt rollback --project my-project --env dev --force
  envcrypt rollback`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := rollbackProject
		envName := rollbackEnv

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
			projectName, err = tui.RunPicker("Select a project", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		// Pick env if not provided via flag
		if envName == "" {
			picked, err := tui.RunEnvPicker(projectName)
			if err != nil {
				return tui.Cancelled()
			}
			envName = picked
		}

		// Fetch all versions
		versions, err := Application.PullAllEnv(cmd.Context(), projectName, envName)
		if err != nil {
			return tui.Error("failed to fetch environment versions", err)
		}
		if len(versions) == 0 {
			return tui.Error("no versions found for this environment", nil)
		}

		sort.Slice(versions, func(i, j int) bool {
			return versions[i].Version > versions[j].Version
		})

		currentVersion := versions[0]
		var targetVersion *int32

		// Resolve target version from arg, flag, or picker
		if len(args) > 0 {
			v, err := strconv.Atoi(args[0])
			if err != nil {
				return tui.Error("invalid version number", err)
			}
			v32 := int32(v)
			targetVersion = &v32
		} else if rollbackVer != 0 {
			v32 := int32(rollbackVer)
			targetVersion = &v32
		} else {
			// Build version labels — exclude current
			labels := make([]string, 0, len(versions)-1)
			verMap := map[string]int32{}
			for i, v := range versions {
				if i == 0 {
					continue // skip current
				}
				label := fmt.Sprintf("v%d", v.Version)
				labels = append(labels, label)
				verMap[label] = v.Version
			}

			if len(labels) == 0 {
				return tui.Error("no previous versions to rollback to", nil)
			}

			picked, err := tui.RunPicker(
				fmt.Sprintf("Rollback %s/%s  (current: v%d)", projectName, envName, currentVersion.Version),
				labels,
			)
			if err != nil {
				return tui.Cancelled()
			}
			v32 := verMap[picked]
			targetVersion = &v32
		}

		if *targetVersion == currentVersion.Version {
			tui.Warn(fmt.Sprintf("Environment is already at version %d.", *targetVersion))
			return nil
		}

		// Find target env map
		var targetMap map[string]string
		foundTarget := false
		for _, v := range versions {
			if v.Version == *targetVersion {
				targetMap = v.Env
				foundTarget = true
				break
			}
		}
		if !foundTarget {
			return tui.Error(fmt.Sprintf("version %d not found", *targetVersion), nil)
		}

		// Show diff
		diff := cryptoutils.DiffEnvVersions(currentVersion.Env, targetMap)
		tui.Spacer()
		tui.Info(fmt.Sprintf("Rolling back %s/%s from v%d → v%d", projectName, envName, currentVersion.Version, *targetVersion))
		tui.Spacer()
		tui.RenderDiff(diff, currentVersion.Env, targetMap, showSecrets)

		// Confirm
		if !tui.ConfirmDangerousAction(
			fmt.Sprintf("Rollback %s/%s to v%d?", projectName, envName, *targetVersion),
			"rollback",
			rollbackForce,
		) {
			return tui.Cancelled()
		}

		// Execute
		err = tui.RunActionWithSpinner(
			fmt.Sprintf("Rolling back to v%d...", *targetVersion),
			func() error {
				return Application.RollbackEnv(cmd.Context(), projectName, envName, targetVersion)
			},
		)
		if err != nil {
			return tui.Error("rollback failed", err)
		}

		tui.Success(fmt.Sprintf("Rolled back %s/%s to v%d", projectName, envName, *targetVersion))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)

	rollbackCmd.Flags().StringVarP(&rollbackProject, "project", "p", "", "Project name")
	rollbackCmd.Flags().StringVarP(&rollbackEnv, "env", "e", "", "Environment name")
	rollbackCmd.Flags().IntVarP(&rollbackVer, "version", "v", 0, "Version to rollback to")
	rollbackCmd.Flags().BoolVar(&rollbackForce, "force", false, "Skip confirmation prompt")
	rollbackCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show actual secret values in diff output")
}

/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/huh"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	rollbackProject string
	rollbackEnv     string
	rollbackVer     int
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback [version]",
	Short: "Rollback to a previous version of an environment",
	Long: `Rollback an environment to a specific version.

This command will create a new version that is an exact copy of the specified previous version.
You will see a diff of the changes before confirming the rollback.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Resolve Project and Env
		projectName := rollbackProject
		envName := rollbackEnv

		// Prompt for Project if missing
		if projectName == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Project Name").
						Value(&projectName).
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("project name is required")
							}
							return nil
						}),
				),
			)
			if err := form.Run(); err != nil {
				return Error("cancelled", nil)
			}
		}

		// Prompt for Env if missing
		if envName == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Environment Name").
						Value(&envName).
						Validate(func(str string) error {
							if str == "" {
								return fmt.Errorf("environment name is required")
							}
							return nil
						}),
				),
			)
			if err := form.Run(); err != nil {
				return Error("cancelled", nil)
			}
		}

		// 2. Fetch all versions
		versions, err := Application.PullAllEnv(cmd.Context(), projectName, envName)
		if err != nil {
			return Error("failed to fetch environment versions", err)
		}

		if len(versions) == 0 {
			return Error("no versions found for this environment", nil)
		}

		// Sort versions by version number (descending)
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].Version > versions[j].Version
		})

		currentVersion := versions[0]
		var targetVersion *int32

		// 3. Determine target version
		if len(args) > 0 {
			v, err := strconv.Atoi(args[0])
			if err != nil {
				return Error("invalid version number", err)
			}
			v32 := int32(v)
			targetVersion = &v32
		} else if rollbackVer != 0 {
			v32 := int32(rollbackVer)
			targetVersion = &v32
		} else {
			// Interactive selection
			options := make([]huh.Option[int], 0, len(versions))
			for i, v := range versions {
				label := fmt.Sprintf("v%d", v.Version)
				if i == 0 {
					label += " (Current)"
				}
				// Don't allowing rolling back to the current version?
				// Actually, it might be useful if the current version is somehow messed up in a way that is not purely data but metadata?
				// But generally rollback implies going back.
				// For now let's show all.
				options = append(options, huh.NewOption(label, int(v.Version)))
			}

			var selectedVer int
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[int]().
						Title("Select Version to Rollback to").
						Options(options...).
						Value(&selectedVer),
				),
			)

			if err := form.Run(); err != nil {
				return Error("cancelled", nil)
			}
			v32 := int32(selectedVer)
			targetVersion = &v32
		}

		if *targetVersion == currentVersion.Version {
			Warn(fmt.Sprintf("Environment is already at version %d.", *targetVersion))
			return nil
		}

		// 4. Find the actual maps
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
			return Error(fmt.Sprintf("version %d not found", *targetVersion), nil)
		}

		// 5. Compute Diff
		// We are going FROM current TO target.
		diff := cryptoutils.DiffEnvVersions(currentVersion.Env, targetMap)

		// 6. Preview and Confirm
		Spacer()
		fmt.Printf("Rolling back %s/%s from v%d to v%d\n", projectName, envName, currentVersion.Version, *targetVersion)
		Spacer()

		renderDiff(diff, currentVersion.Env, targetMap, showSecrets)

		if !ConfirmDangerousAction(fmt.Sprintf("Are you sure you want to rollback to v%d?", *targetVersion), "rollback") {
			return nil
		}

		// 7. Execute Rollback
		err = Application.RollbackEnv(cmd.Context(), projectName, envName, targetVersion)
		if err != nil {
			return Error("rollback failed", err)
		}

		Success(fmt.Sprintf("Successfully rolled back to version %d", *targetVersion))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)

	rollbackCmd.Flags().StringVarP(&rollbackProject, "project", "p", "", "Project name")
	rollbackCmd.Flags().StringVarP(&rollbackEnv, "env", "e", "", "Environment name")
	rollbackCmd.Flags().IntVarP(&rollbackVer, "version", "v", 0, "Version to rollback to")
	rollbackCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show actual secret values in diff output")
}

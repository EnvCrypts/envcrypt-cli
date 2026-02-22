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
	diffProject string
	diffEnv     string
	showSecrets bool
)

var diffCmd = &cobra.Command{
	Use:   "diff [old_version] [new_version]",
	Short: "Diff two environment versions",
	Long: `Compare two versions of an environment configuration.

If version numbers are not provided, an interactive prompt will allow you to select the versions to compare.

Use --show-secrets to reveal the actual values in the diff output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := diffProject
		envName := diffEnv

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
				return tui.Error("cancelled", nil)
			}
		}

		// Pick env if not provided via flag
		if envName == "" {
			picked, err := tui.RunEnvPicker(projectName)
			if err != nil {
				return tui.Error("cancelled", nil)
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

		var oldVer, newVer int

		if len(args) == 2 {
			v1, err := strconv.Atoi(args[0])
			if err != nil {
				return tui.Error("invalid old version number", err)
			}
			v2, err := strconv.Atoi(args[1])
			if err != nil {
				return tui.Error("invalid new version number", err)
			}
			oldVer, newVer = v1, v2
		} else {
			// Build version labels for picker
			labels := make([]string, len(versions))
			for i, v := range versions {
				label := fmt.Sprintf("v%d", v.Version)
				if i == 0 {
					label += "  (Current)"
				}
				labels[i] = label
			}

			oldLabel, err := tui.RunPicker("Base version (old)", labels)
			if err != nil {
				return tui.Error("cancelled", nil)
			}
			newLabel, err := tui.RunPicker("Target version (new)", labels)
			if err != nil {
				return tui.Error("cancelled", nil)
			}

			// Map label back to version number
			for i, v := range versions {
				if labels[i] == oldLabel {
					oldVer = int(v.Version)
				}
				if labels[i] == newLabel {
					newVer = int(v.Version)
				}
			}
		}

		var oldMap, newMap map[string]string
		foundOld, foundNew := false, false
		for _, v := range versions {
			if int(v.Version) == oldVer {
				oldMap = v.Env
				foundOld = true
			}
			if int(v.Version) == newVer {
				newMap = v.Env
				foundNew = true
			}
		}

		if !foundOld {
			return tui.Error(fmt.Sprintf("version %d not found", oldVer), nil)
		}
		if !foundNew {
			return tui.Error(fmt.Sprintf("version %d not found", newVer), nil)
		}

		diff := cryptoutils.DiffEnvVersions(oldMap, newMap)
		tui.RenderDiff(diff, oldMap, newMap, showSecrets)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&diffProject, "project", "p", "", "Project name")
	diffCmd.Flags().StringVarP(&diffEnv, "env", "e", "", "Environment name")
	diffCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show actual secret values in diff output")
}

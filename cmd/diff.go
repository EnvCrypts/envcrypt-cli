package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	diffProject string
	diffEnv     string
	showSecrets bool
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [old_version] [new_version]",
	Short: "Diff two environment versions",
	Long: `Compare two versions of an environment configuration.

If version numbers are not provided, an interactive prompt will allow you to select the versions to compare.

Use --show-secrets to reveal the actual values in the diff output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Resolve Project and Env
		projectName := diffProject
		envName := diffEnv

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

		var oldVer, newVer int

		// 3. Determine versions to compare
		if len(args) == 2 {
			v1, err := strconv.Atoi(args[0])
			if err != nil {
				return Error("invalid old version number", err)
			}
			v2, err := strconv.Atoi(args[1])
			if err != nil {
				return Error("invalid new version number", err)
			}
			oldVer = v1
			newVer = v2
		} else {
			// Interactive selection
			options := make([]huh.Option[int], len(versions))
			for i, v := range versions {
				// Show "Current" for the latest version
				label := fmt.Sprintf("v%d", v.Version)
				if i == 0 {
					label += " (Current)"
				}
				options[i] = huh.NewOption(label, int(v.Version))
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[int]().
						Title("Base Version (Old)").
						Options(options...).
						Value(&oldVer),
					huh.NewSelect[int]().
						Title("Target Version (New)").
						Options(options...).
						Value(&newVer),
				),
			)

			if err := form.Run(); err != nil {
				return Error("cancelled", nil)
			}
		}

		// 4. Find the actual maps
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
			return Error(fmt.Sprintf("version %d not found", oldVer), nil)
		}
		if !foundNew {
			return Error(fmt.Sprintf("version %d not found", newVer), nil)
		}

		// 5. Compute Diff
		diff := cryptoutils.DiffEnvVersions(oldMap, newMap)

		// 6. Render Output
		renderDiff(diff, oldMap, newMap, showSecrets)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&diffProject, "project", "p", "", "Project name")
	diffCmd.Flags().StringVarP(&diffEnv, "env", "e", "", "Environment name")
	diffCmd.Flags().BoolVar(&showSecrets, "show-secrets", false, "Show actual secret values in diff output")
}

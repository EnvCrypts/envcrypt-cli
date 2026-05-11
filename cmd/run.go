package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	runProject string
	runEnvName string
)

var runCmd = &cobra.Command{
	Use:   "run [project] -- <command> [args...]",
	Short: "Run a command with injected environment variables",
	Long: `Run a command with EnvCrypt secrets injected into the runtime environment.

Examples:
  envcrypt run my-project --env prod -- npm start
  envcrypt run --project my-project --env dev -- python app.py
  envcrypt run --env staging -- npm run build`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dashIndex := cmd.ArgsLenAtDash()
		if dashIndex == -1 {
			return tui.Error("missing command separator", nil, "Use -- before the command to run")
		}

		argsBeforeDash := args[:dashIndex]
		commandArgs := args[dashIndex:]
		if len(commandArgs) == 0 {
			return tui.Error("missing command to run", nil, "Provide a command after --")
		}
		if runProject != "" && len(argsBeforeDash) > 0 {
			return tui.Error("project specified twice", nil)
		}
		if runProject == "" {
			if len(argsBeforeDash) > 1 {
				return tui.Error("too many arguments", nil, "Provide at most one project before --")
			}
			if len(argsBeforeDash) == 1 {
				runProject = argsBeforeDash[0]
			}
		}

		projectName := runProject
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

			projectName, err = tui.RunPickerWithDefault("Select a project to run", names, defaultProject)
			if err != nil {
				return tui.Cancelled()
			}
		}

		envName := runEnvName
		if envName == "" {
			picked, err := tui.RunEnvPicker(projectName)
			if err != nil {
				return tui.Cancelled()
			}
			envName = picked
		}

		tui.Info("Project: " + projectName)
		tui.Info("Environment: " + envName)

		var envMap map[string]string
		err := tui.RunActionWithSpinner(
			fmt.Sprintf("Injecting %s/%s...", projectName, envName),
			func() error {
				var e error
				envMap, e = Application.PullEnv(cmd.Context(), projectName, envName)
				return e
			},
		)
		if err != nil {
			return tui.Error("failed to pull environment variables", err)
		}

		command := exec.CommandContext(cmd.Context(), commandArgs[0], commandArgs[1:]...)
		command.Env = mergeEnviron(envMap)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if exitErr.ProcessState != nil {
					os.Exit(exitErr.ProcessState.ExitCode())
				}
			}
			return tui.Error("command failed", err)
		}

		return nil
	},
}

func mergeEnviron(envMap map[string]string) []string {
	merged := make(map[string]string, len(os.Environ())+len(envMap))
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			merged[parts[0]] = parts[1]
		}
	}
	for k, v := range envMap {
		merged[k] = v
	}

	out := make([]string, 0, len(merged))
	for k, v := range merged {
		out = append(out, k+"="+v)
	}
	return out
}

func init() {
	runCmd.Flags().StringVar(&runProject, "project", "", "Project name")
	runCmd.Flags().StringVar(&runEnvName, "env", "", "Environment name (dev, staging, prod)")
	rootCmd.AddCommand(runCmd)
}

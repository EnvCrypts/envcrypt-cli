package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var (
	runProject string
	runEnvName string
	runEnvFile string
	runPrint   bool
)

var runCmd = &cobra.Command{
	Use:   "run [project] [env] [--] <command> [args...]",
	Short: "Run a command with injected environment variables",
	Long: `Run a command with EnvCrypt secrets injected into the runtime environment.

Examples:
  envcrypt run my-project prod -- npm start
  envcrypt run my-project prod npm start
  envcrypt run --project my-project --env dev -- python app.py
  envcrypt run --env staging -- npm run build`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dashIndex := cmd.ArgsLenAtDash()
		commandArgs := []string{}

		if dashIndex != -1 {
			argsBeforeDash := args[:dashIndex]
			commandArgs = args[dashIndex:]

			if len(argsBeforeDash) > 0 {
				switch {
				case runProject != "" && runEnvName != "":
					return tui.Error("project and environment specified twice", nil)
				case runProject != "" && runEnvName == "":
					if len(argsBeforeDash) > 1 {
						return tui.Error("too many arguments", nil, "Provide at most one environment before --")
					}
					if len(argsBeforeDash) == 1 {
						runEnvName = argsBeforeDash[0]
					}
				case runProject == "" && runEnvName != "":
					if len(argsBeforeDash) > 1 {
						return tui.Error("too many arguments", nil, "Provide at most one project before --")
					}
					if len(argsBeforeDash) == 1 {
						runProject = argsBeforeDash[0]
					}
				default:
					if len(argsBeforeDash) > 2 {
						return tui.Error("too many arguments", nil, "Provide at most [project] [env] before --")
					}
					if len(argsBeforeDash) >= 1 {
						runProject = argsBeforeDash[0]
					}
					if len(argsBeforeDash) == 2 {
						runEnvName = argsBeforeDash[1]
					}
				}
			}
		} else {
			switch {
			case runProject == "" && runEnvName == "" && len(args) >= 3:
				runProject = args[0]
				runEnvName = args[1]
				commandArgs = args[2:]
			case runProject != "" && runEnvName == "" && len(args) >= 2:
				runEnvName = args[0]
				commandArgs = args[1:]
			case runProject == "" && runEnvName != "" && len(args) >= 2:
				runProject = args[0]
				commandArgs = args[1:]
			default:
				commandArgs = args
			}
		}

		if len(commandArgs) == 0 {
			return tui.Error("missing command to run", nil, "Provide a command after -- or after [project] [env]")
		}

		projectName := runProject
		if projectName == "" && !tui.IsInteractive() {
			return tui.Error("project is required in non-interactive mode", nil, "Use --project or provide [project]")
		}
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
		if envName == "" && !tui.IsInteractive() {
			return tui.Error("environment is required in non-interactive mode", nil, "Use --env or provide [env]")
		}
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

		if runEnvFile != "" {
			fileBytes, err := os.ReadFile(runEnvFile)
			if err != nil {
				return tui.Error("failed to read env file", mapEnvReadError(runEnvFile, err))
			}
			fileEnv, err := cryptoutils.ParseEnv(fileBytes)
			if err != nil {
				return tui.Error("failed to parse env file", err)
			}
			envMap = mergeEnvMaps(fileEnv, envMap)
		}

		if runPrint {
			tui.PrintEnvSummary(envMap)
			return nil
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

func mergeEnvMaps(base, override map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range override {
		merged[k] = v
	}
	return merged
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
	runCmd.Flags().StringVarP(&runEnvFile, "env-file", "e", "", "Path to local .env file to merge before execution")
	runCmd.Flags().BoolVar(&runPrint, "print", false, "Print injected keys and exit")
	rootCmd.AddCommand(runCmd)
}

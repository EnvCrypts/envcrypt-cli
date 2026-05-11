package cmd

import (
	"errors"
	"os"

	"github.com/envcrypts/envcrypt-cli/internal/app"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var Version = "dev"

// Global flags
var (
	globalJSON    bool
	globalNoColor bool
	globalQuiet   bool
	globalNoTable bool
)

var rootCmd = &cobra.Command{
	Version:       Version,
	Use:           "envcrypt",
	Short:         "Zero-trust, end-to-end encrypted environment variable management.",
	Long:          tui.RenderBanner(),
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		tui.InitOutput(globalJSON, globalNoColor, globalQuiet)
		tui.SetNoTable(globalNoTable)
	},
}

var Application *app.App

func Execute(a *app.App) {
	Application = a
	assignGroups()
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			os.Exit(130)
		}
		tui.RenderError(err)
		os.Exit(1)
	}
}

func init() {
	// Command groups for better organization
	workflowGroup := &cobra.Group{ID: "workflow", Title: "WORKFLOW COMMANDS"}
	mgmtGroup := &cobra.Group{ID: "mgmt", Title: "MANAGEMENT COMMANDS"}
	accessGroup := &cobra.Group{ID: "access", Title: "ACCESS & SECURITY"}
	sessionGroup := &cobra.Group{ID: "session", Title: "SESSION COMMANDS"}
	systemGroup := &cobra.Group{ID: "system", Title: "SYSTEM COMMANDS"}

	rootCmd.AddGroup(workflowGroup, mgmtGroup, accessGroup, sessionGroup, systemGroup)

	rootCmd.PersistentFlags().BoolVar(&globalJSON, "json", false, "Output results as structured JSON")
	rootCmd.PersistentFlags().BoolVar(&globalNoColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&globalQuiet, "quiet", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&globalNoTable, "no-table", false, "Use plain tabular output instead of interactive tables")

	// Custom Usage Template for a more premium feel
	rootCmd.SetUsageTemplate(usageTemplate())
}

func assignGroups() {
	for _, c := range rootCmd.Commands() {
		switch c.Name() {
		case "pull", "push", "diff", "rollback", "rotate", "run":
			c.GroupID = "workflow"
		case "create", "delete", "list", "audit":
			c.GroupID = "mgmt"
		case "add", "grant", "revoke", "service-role":
			c.GroupID = "access"
		case "login", "logout", "whoami", "status", "recover":
			c.GroupID = "session"
		case "version", "completion", "snapshot":
			c.GroupID = "system"
		}
	}
}

func usageTemplate() string {
	return `USAGE:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

ALIASES:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

EXAMPLES:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

AVAILABLE COMMANDS:{{range $cmds}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}:{{range $cmds}}{{if and (eq .GroupID $group.ID) (.IsAvailableCommand)}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
{{if not .AllChildCommandsHaveGroup}}OTHER COMMANDS:{{range $cmds}}{{if and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

FLAGS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

GLOBAL FLAGS:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

ADDITIONAL HELP TOPICS:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

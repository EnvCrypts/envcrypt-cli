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
	Version: Version,
	Use:     "envcrypt",
	Short:   "Zero-trust, end-to-end encrypted environment variable management.",
	Long:    tui.RenderBanner(),
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
		case "pull", "push", "diff", "rollback", "rotate":
			c.GroupID = "workflow"
		case "create", "delete", "list", "audit":
			c.GroupID = "mgmt"
		case "add", "grant", "revoke", "service-role":
			c.GroupID = "access"
		case "login", "logout", "whoami", "status":
			c.GroupID = "session"
		case "version", "completion", "snapshot":
			c.GroupID = "system"
		}
	}
}

func usageTemplate() string {
	return `USAGE:
  {{.CommandPath}} [command]

{{if .HasAvailableSubCommands -}}
{{range .Groups -}}
{{$group := . -}}
{{.Title}}:
{{range $c := $.Commands -}}
{{if eq $c.GroupID $group.ID}}  {{rpad $c.Name $c.NamePadding}} {{$c.Short}}
{{end}}{{end}}
{{end}}{{if not .HasParent -}}
OTHER COMMANDS:
{{range .Commands}}{{if not .GroupID}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}
{{end}}{{end}}
FLAGS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`
}

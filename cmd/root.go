package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envcrypt",
	Short: "Zero-trust, end-to-end encrypted environment variable management.",
	Long: lipgloss.NewStyle().
		MarginLeft(1).
		MarginRight(1).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true).Render("EnvCrypt CLI 🛡️"),
				"",
				"Zero-trust, end-to-end encrypted environment variable management.",
				"All secrets are encrypted client-side with immutable versioning.",
			),
		),
}

var Application *app.App

func Execute(a *app.App) {
	Application = a
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.envcrypt-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/app"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Version: Version,
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
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

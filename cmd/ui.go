package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/config"
)

var (
	// Styles
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).PaddingLeft(1) // Green
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true).PaddingLeft(1) // Red
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true).PaddingLeft(1) // Yellow
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).PaddingLeft(1)             // Blue
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))            // Grey
	revokedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Strikethrough(true) // Red Strikethrough

	// Icons (No extra padding here as it's on the message)
	iconCheck = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✔")
	iconCross = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render("✖")
	iconWarn  = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("⚠")
	iconInfo  = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("ℹ")
)

func Success(msg string) {
	fmt.Printf("%s %s\n", iconCheck, successStyle.Render(msg))
}

func Error(msg string, err error) error {
	errMsg := errorStyle.Render(msg)
	if err != nil {
		return fmt.Errorf("%s %s: %w", iconCross, errMsg, err)
	}
	return fmt.Errorf("%s %s", iconCross, errMsg)
}

func Warn(msg string) {
	fmt.Printf("%s %s\n", iconWarn, warnStyle.Render(msg))
}

func Info(msg string) {
	fmt.Printf("%s %s\n", iconInfo, infoStyle.Render(msg))
}

func Spacer() {
	fmt.Println()
}

func ConfirmDangerousAction(prompt, expected string) bool {
	Warn(prompt)
	fmt.Printf("Type %q to confirm: ", expected)

	var input string
	fmt.Scanln(&input)

	return input == expected
}

func PrintProjects(projects []config.Project) {
	if len(projects) == 0 {
		fmt.Println(mutedStyle.Render("No projects found."))
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	// Headers
	fmt.Fprintln(w, lipgloss.NewStyle().Bold(true).Underline(true).Render("PROJECT NAME\tROLE"))

	for _, p := range projects {
		name := p.Name
		role := p.Role
        
		switch p.Role {
		case "admin":
			role = successStyle.Render("admin")
		case "member":
			role = infoStyle.Render("member")
		default:
			role = mutedStyle.Render(p.Role)
		}

		if p.IsRevoked {
			name = revokedStyle.Render(name)
			role = revokedStyle.Render(p.Role + " (revoked)")
		}

		fmt.Fprintf(w, "%s\t%s\n", name, role)
	}

	w.Flush()
}

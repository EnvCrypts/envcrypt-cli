package tui

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/config"
)

// noTable is set by the --no-table global flag.
var noTable bool

// SetNoTable configures table output to use plain tabular format.
func SetNoTable(v bool) { noTable = v }

var tableBaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(ColorPrimary)

type tableModel struct {
	table table.Model
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return "\n" + tableBaseStyle.Render(m.table.View()) +
		"\n" + StyleMuted.Render("  "+TableNavHint) + "\n"
}

func styledTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(min(len(rows)+1, 15)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorPrimary).
		BorderBottom(true).
		Bold(true).
		Foreground(ColorPrimary)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("0")).
		Background(ColorPrimary).
		Bold(false)
	t.SetStyles(s)
	return t
}

func runTable(t table.Model) error {
	p := tea.NewProgram(tableModel{table: t}, tea.WithOutput(os.Stdout))
	_, err := p.Run()
	return err
}

func RunProjectsTable(projects []config.Project) error {
	if len(projects) == 0 {
		if currentMode == ModeJSON {
			JSONData([]any{})
			return nil
		}
		fmt.Println(StyleMuted.Render("No projects found."))
		return nil
	}

	// JSON mode: structured output
	if currentMode == ModeJSON {
		rows := make([]map[string]string, len(projects))
		for i, p := range projects {
			status := "active"
			if p.IsRevoked {
				status = "revoked"
			}
			rows[i] = map[string]string{"project": p.Name, "role": p.Role, "status": status}
		}
		data, _ := json.MarshalIndent(rows, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// Plain mode or --no-table: tabular text
	if !IsInteractive() || noTable {
		fmt.Printf("%-24s  %-12s  %-10s\n", "PROJECT", "ROLE", "STATUS")
		for _, p := range projects {
			status := "active"
			if p.IsRevoked {
				status = "revoked"
			}
			fmt.Printf("%-24s  %-12s  %-10s\n", p.Name, p.Role, status)
		}
		return nil
	}

	rows := make([]table.Row, len(projects))
	for i, p := range projects {
		status := "active"
		if p.IsRevoked {
			status = "revoked"
		}
		rows[i] = table.Row{p.Name, p.Role, status}
	}
	return runTable(styledTable([]table.Column{
		{Title: "Project", Width: 24},
		{Title: "Role", Width: 12},
		{Title: "Status", Width: 10},
	}, rows))
}


func RunServiceRolesTable(roles []config.ServiceRole) error {
	if len(roles) == 0 {
		if currentMode == ModeJSON {
			JSONData([]any{})
			return nil
		}
		fmt.Println(StyleMuted.Render("No service roles found."))
		return nil
	}

	// JSON mode
	if currentMode == ModeJSON {
		rows := make([]map[string]string, len(roles))
		for i, r := range roles {
			rows[i] = map[string]string{"name": r.Name, "repo_principal": r.RepoPrincipal}
		}
		data, _ := json.MarshalIndent(rows, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// Plain mode or --no-table
	if !IsInteractive() || noTable {
		fmt.Printf("%-24s  %-40s\n", "NAME", "REPO PRINCIPAL")
		for _, r := range roles {
			fmt.Printf("%-24s  %-40s\n", r.Name, r.RepoPrincipal)
		}
		return nil
	}

	rows := make([]table.Row, len(roles))
	for i, r := range roles {
		rows[i] = table.Row{r.Name, r.RepoPrincipal}
	}
	return runTable(styledTable([]table.Column{
		{Title: "Name", Width: 24},
		{Title: "Repo Principal", Width: 40},
	}, rows))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

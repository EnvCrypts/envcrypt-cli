package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/config"
)

// noTable is set by the --no-table global flag.
var noTable bool

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

	if currentMode == ModeJSON {
		rows := make([]map[string]string, len(roles))
		for i, r := range roles {
			rows[i] = map[string]string{"name": r.Name, "repo_principal": r.RepoPrincipal}
		}
		data, _ := json.MarshalIndent(rows, "", "  ")
		fmt.Println(string(data))
		return nil
	}

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

func RunAuditTable(logs []config.AuditEntry, total int) error {
	if currentMode == ModeJSON {
		if len(logs) == 0 {
			JSONData([]any{})
		} else {
			data, _ := json.MarshalIndent(logs, "", "  ")
			fmt.Println(string(data))
		}
		return nil
	}

	if !IsInteractive() || noTable {
		if len(logs) == 0 {
			fmt.Println("No audit logs found.")
			return nil
		}
		fmt.Printf("%-15s  %-20s  %-16s  %-8s  %-10s\n", "TIME", "ACTOR", "ACTION", "STATUS", "ENV")
		for _, row := range logs {
			timestamp := row.CreatedAt.Local().Format("Jan 02 15:04")
			env := "-"
			if row.Environment != nil && *row.Environment != "" {
				env = *row.Environment
			}

			fmt.Printf("%-15s  %-20s  %-16s  %-8s  %-10s\n", timestamp, Truncate(row.ActorEmail, 20), row.Action, row.Status, env)

			if strings.ToLower(row.Status) != "success" && row.ErrorMessage != "" {
				fmt.Printf("   └── Error: %s\n", row.ErrorMessage)
			}
		}

		fmt.Printf("\nShowing %d of %d total logs\n", len(logs), total)
		return nil
	}

	if len(logs) == 0 {
		emptyMsg := StyleMuted.Render("No audit logs found for this project.")
		fmt.Printf("\n  %s\n\n", emptyMsg)
		return nil
	}

	rows := make([]table.Row, len(logs))
	for i, row := range logs {
		timestamp := row.CreatedAt.Local().Format("Jan 02 15:04")
		env := "-"
		if row.Environment != nil && *row.Environment != "" {
			env = *row.Environment
		}

		actor := row.ActorEmail
		if actor == "" {
			actor = "System"
		}

		status := row.Status
		if strings.ToLower(row.Status) != "success" {
			if row.ErrorMessage != "" {
				errMsg := row.ErrorMessage
				if len(errMsg) > 15 {
					errMsg = errMsg[:12] + "..."
				}
				status = fmt.Sprintf("failed (%s)", errMsg)
			}
		}

		rows[i] = table.Row{
			timestamp,
			Truncate(actor, 18),
			row.Action,
			env,
			Truncate(status, 22),
		}
	}

	err := runTable(styledTable([]table.Column{
		{Title: "Time", Width: 14},
		{Title: "Actor", Width: 18},
		{Title: "Action", Width: 18},
		{Title: "Env", Width: 8},
		{Title: "Status", Width: 22},
	}, rows))

	if err == nil {
		fmt.Printf("\n  %s\n", StyleMuted.Render(fmt.Sprintf("Showing %d of %d total logs", len(logs), total)))
	}

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

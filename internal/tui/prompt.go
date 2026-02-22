package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	prompt   string
	expected string
	input    string
	done     bool
	abort    bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.abort = true
			return m, tea.Quit
		case tea.KeyEnter:
			if m.input == m.expected {
				m.done = true
			} else {
				m.abort = true
			}
			return m, tea.Quit
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.done || m.abort {
		return ""
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n  %s %s\n\n", IconWarn, StyleWarn.Render(m.prompt)))
	b.WriteString(fmt.Sprintf("  %s\n", StyleMuted.Render(fmt.Sprintf(`Type "%s" to confirm:`, m.expected))))
	b.WriteString(fmt.Sprintf("  %s%s\n", StylePrimary.Render("> "), m.input))
	b.WriteString(fmt.Sprintf("\n  %s\n", StyleMuted.Render("enter to confirm  •  esc to cancel")))
	return b.String()
}

// ConfirmDangerousAction presents a bubbletea inline confirmation prompt.
// Returns true only if the user types `expected` exactly and presses Enter.
func ConfirmDangerousAction(prompt, expected string) bool {
	p := tea.NewProgram(confirmModel{prompt: prompt, expected: expected})
	result, err := p.Run()
	if err != nil {
		return false
	}
	final, ok := result.(confirmModel)
	return ok && final.done
}

// ConfirmOverwrite asks the user to type "yes" to confirm overwriting a file.
func ConfirmOverwrite(path string) bool {
	return ConfirmDangerousAction(fmt.Sprintf("Overwrite %q?", path), "yes")
}

// PromptWithDefault shows a single-field form pre-filled with defaultVal.
// Returns defaultVal unchanged if the user submits an empty value.
func PromptWithDefault(label, defaultVal string) string {
	vals, err := RunForm([]FormField{{Label: label}}, []string{defaultVal})
	if err != nil || len(vals) == 0 || vals[0] == "" {
		return defaultVal
	}
	return vals[0]
}

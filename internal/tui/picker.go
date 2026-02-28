package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	pickerSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(ColorPrimary).
				PaddingLeft(1).PaddingRight(1)
	pickerItemStyle = lipgloss.NewStyle().PaddingLeft(3)
)

type pickerModel struct {
	title  string
	items  []string
	cursor int
	chosen string
	abort  bool
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.abort = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = m.items[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	if m.chosen != "" || m.abort {
		return ""
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n  %s\n\n", activeStyle.Render(m.title)))
	for i, item := range m.items {
		if i == m.cursor {
			b.WriteString("  " + pickerSelectedStyle.Render(item) + "\n")
		} else {
			b.WriteString(pickerItemStyle.Render(item) + "\n")
		}
	}
	b.WriteString(fmt.Sprintf("\n  %s\n", StyleMuted.Render(PickerNavHint)))
	return b.String()
}


func RunPicker(title string, items []string) (string, error) {
	return RunPickerWithDefault(title, items, "")
}

func RunPickerWithDefault(title string, items []string, defaultItem string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to pick from")
	}
	if !IsInteractive() {
		return "", fmt.Errorf("not a terminal: provide required values via flags")
	}

	cursor := 0
	if defaultItem != "" {
		for i, item := range items {
			if item == defaultItem {
				cursor = i
				break
			}
		}
	}

	p := tea.NewProgram(pickerModel{title: title, items: items, cursor: cursor})
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	final, ok := result.(pickerModel)
	if !ok || final.abort {
		return "", ErrCancelled
	}
	return final.chosen, nil
}


func RunEnvPicker(projectName string) (string, error) {
	const other = "Other..."
	picked, err := RunPicker(
		fmt.Sprintf("Select environment for %q", projectName),
		[]string{"dev", "staging", "prod", other},
	)
	if err != nil {
		return "", err
	}
	if picked != other {
		return picked, nil
	}
	vals, err := RunForm([]FormField{{Label: "Environment name", Required: true}}, []string{""})
	if err != nil {
		return "", err
	}
	return vals[0], nil
}

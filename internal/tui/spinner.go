package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type actionMsg struct{ err error }

type spinnerModel struct {
	spinner  spinner.Model
	title    string
	action   func() error
	quitting bool
	err      error
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			return actionMsg{err: m.action()}
		},
	)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			m.err = fmt.Errorf("cancelled")
			return m, tea.Quit
		}
	case actionMsg:
		m.quitting = true
		m.err = msg.err
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s %s\n\n", m.spinner.View(), lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(m.title))
}

// RunActionWithSpinner runs action in the background while showing a spinner.
func RunActionWithSpinner(title string, action func() error) error {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = StylePrimary

	p := tea.NewProgram(spinnerModel{spinner: s, title: title, action: action})
	result, err := p.Run()
	if err != nil {
		return err
	}
	if final, ok := result.(spinnerModel); ok {
		return final.err
	}
	return nil
}

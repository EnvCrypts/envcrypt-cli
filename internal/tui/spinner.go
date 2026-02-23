package tui

import (
	"fmt"
	"time"

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
	started  time.Time
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
			m.err = ErrCancelled
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
		elapsed := time.Since(m.started).Round(time.Millisecond)
		if m.err != nil && m.err != ErrCancelled {
			return fmt.Sprintf(" %s %s (%s)\n", IconCross, m.title, elapsed)
		}
		return fmt.Sprintf(" %s %s (%s)\n", IconCheck, m.title, elapsed)
	}
	elapsed := time.Since(m.started).Round(time.Millisecond)
	return fmt.Sprintf("\n %s %s %s\n\n",
		m.spinner.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(m.title),
		StyleMuted.Render(fmt.Sprintf("(%s)", elapsed)),
	)
}

// RunActionWithSpinner executes an action with a spinner in interactive mode,
// a static progress line in plain mode, or silently in JSON mode.
func RunActionWithSpinner(title string, action func() error) error {
	switch currentMode {
	case ModeJSON:
		// JSON mode: run silently, report nothing (caller handles output)
		return action()

	case ModePlain:
		// Plain mode: static progress line with elapsed time
		start := time.Now()
		fmt.Printf("%s ⏳ %s...\n", PlainTimestamp(), title)
		err := action()
		elapsed := time.Since(start).Round(time.Millisecond)
		if err != nil {
			fmt.Printf("%s %s %s (%s)\n", PlainTimestamp(), PlainIconCross, title, elapsed)
		} else {
			fmt.Printf("%s %s %s (%s)\n", PlainTimestamp(), PlainIconCheck, title, elapsed)
		}
		return err

	default:
		// Interactive mode: animated spinner with elapsed time
		if !IsInteractive() {
			// Fallback to plain if stdin is not a TTY
			start := time.Now()
			fmt.Printf("⏳ %s...\n", title)
			err := action()
			elapsed := time.Since(start).Round(time.Millisecond)
			if err != nil {
				fmt.Printf("%s %s (%s)\n", PlainIconCross, title, elapsed)
			} else {
				fmt.Printf("%s %s (%s)\n", PlainIconCheck, title, elapsed)
			}
			return err
		}

		s := spinner.New()
		s.Spinner = spinner.Dot
		s.Style = StylePrimary

		m := spinnerModel{spinner: s, title: title, action: action, started: time.Now()}
		p := tea.NewProgram(m)
		result, err := p.Run()
		if err != nil {
			return err
		}
		if final, ok := result.(spinnerModel); ok {
			return final.err
		}
		return nil
	}
}

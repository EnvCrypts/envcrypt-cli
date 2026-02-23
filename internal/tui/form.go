package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField defines a single input field in a form.
type FormField struct {
	Label    string
	Secret   bool
	Required bool
}

type formModel struct {
	fields  []FormField
	inputs  []textinput.Model
	focused int
	done    bool
	err     error
}

var (
	labelStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	activeStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	cursorStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
)

func newFormModel(fields []FormField, prefills []string) formModel {
	inputs := make([]textinput.Model, len(fields))
	for i, f := range fields {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.Cursor.SetMode(cursor.CursorBlink)
		t.CharLimit = 256
		if f.Secret {
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		if prefills != nil && i < len(prefills) && prefills[i] != "" {
			t.SetValue(prefills[i])
		}
		inputs[i] = t
	}
	inputs[0].Focus()
	return formModel{fields: fields, inputs: inputs}
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = fmt.Errorf("cancelled")
			return m, tea.Quit
		case tea.KeyEnter:
			if m.fields[m.focused].Required && strings.TrimSpace(m.inputs[m.focused].Value()) == "" {
				return m, nil
			}
			if m.focused == len(m.inputs)-1 {
				m.done = true
				return m, tea.Quit
			}
			m.inputs[m.focused].Blur()
			m.focused++
			m.inputs[m.focused].Focus()
			return m, textinput.Blink
		case tea.KeyTab:
			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				m.inputs[m.focused].Focus()
			}
			return m, textinput.Blink
		case tea.KeyShiftTab:
			if m.focused > 0 {
				m.inputs[m.focused].Blur()
				m.focused--
				m.inputs[m.focused].Focus()
			}
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m formModel) View() string {
	if m.done || m.err != nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n")
	for i, f := range m.fields {
		label := labelStyle.Render(f.Label)
		if i == m.focused {
			label = activeStyle.Render(f.Label)
		}
		b.WriteString(fmt.Sprintf("  %s\n  %s\n", label, m.inputs[i].View()))
		if i < len(m.fields)-1 {
			b.WriteString("\n")
		}
	}
	b.WriteString(fmt.Sprintf("\n  %s\n", StyleMuted.Render("enter to confirm  •  esc to cancel")))
	return b.String()
}

func (m formModel) values() []string {
	vals := make([]string, len(m.inputs))
	for i, inp := range m.inputs {
		vals[i] = inp.Value()
	}
	return vals
}


func RunForm(fields []FormField, prefills []string) ([]string, error) {
	if !IsInteractive() {
		return nil, fmt.Errorf("not a terminal: provide required values via flags")
	}
	m := newFormModel(fields, prefills)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}
	final, ok := result.(formModel)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}
	if final.err != nil {
		return nil, final.err
	}
	return final.values(), nil
}

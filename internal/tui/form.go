package tui

import (
	"fmt"
	"net/mail"
	"os"
	"regexp"
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
	Validate func(string) error // optional validation function
}

type formModel struct {
	fields    []FormField
	inputs    []textinput.Model
	focused   int
	done      bool
	err       error
	validErrs []string // per-field validation error messages
}

var (
	labelStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	activeStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	cursorStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	errStyle    = lipgloss.NewStyle().Foreground(ColorError)
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
	return formModel{fields: fields, inputs: inputs, validErrs: make([]string, len(fields))}
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.err = ErrCancelled
			return m, tea.Quit
		case tea.KeyEnter:
			val := strings.TrimSpace(m.inputs[m.focused].Value())

			// Required check
			if m.fields[m.focused].Required && val == "" {
				m.validErrs[m.focused] = "This field is required"
				return m, nil
			}

			// Custom validation
			if m.fields[m.focused].Validate != nil && val != "" {
				if err := m.fields[m.focused].Validate(val); err != nil {
					m.validErrs[m.focused] = err.Error()
					return m, nil
				}
			}

			m.validErrs[m.focused] = "" // clear any previous error

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

	// Clear validation error on typing
	m.validErrs[m.focused] = ""

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
		if m.validErrs[i] != "" {
			b.WriteString(fmt.Sprintf("  %s\n", errStyle.Render("⚠ "+m.validErrs[i])))
		}
		if i < len(m.fields)-1 {
			b.WriteString("\n")
		}
	}
	b.WriteString(fmt.Sprintf("\n  %s\n", StyleMuted.Render(FormNavHint)))
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


// --- Built-in validators ---

var projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}[a-zA-Z0-9]$`)

// ValidateEmail checks for a valid email format.
func ValidateEmail(s string) error {
	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateFileExists checks that a file exists at the given path.
func ValidateFileExists(s string) error {
	info, err := os.Stat(s)
	if err != nil {
		return fmt.Errorf("file %q does not exist", s)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", s)
	}
	return nil
}

// ValidateProjectName checks that a project name follows naming rules.
func ValidateProjectName(s string) error {
	if !projectNameRegex.MatchString(s) {
		return fmt.Errorf("project name must be 3-64 chars, start/end with alphanumeric, and contain only letters, digits, hyphens, or underscores")
	}
	return nil
}

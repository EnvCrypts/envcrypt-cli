package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorSuccess = lipgloss.Color("2")
	ColorError   = lipgloss.Color("1")
	ColorWarn    = lipgloss.Color("3")
	ColorInfo    = lipgloss.Color("4")
	ColorMuted   = lipgloss.Color("8")
	ColorPrimary = lipgloss.Color("63") // Matches standard EnvCrypt color

	// Base Styles
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess)
	StyleError   = lipgloss.NewStyle().Foreground(ColorError)
	StyleWarn    = lipgloss.NewStyle().Foreground(ColorWarn)
	StyleInfo    = lipgloss.NewStyle().Foreground(ColorInfo)
	StyleMuted   = lipgloss.NewStyle().Foreground(ColorMuted)
	StylePrimary = lipgloss.NewStyle().Foreground(ColorPrimary)

	// Icons
	IconCheck = StyleSuccess.Render("[✓]")
	IconCross = StyleError.Render("[x]")
	IconWarn  = StyleWarn.Render("[!]")
	IconInfo  = StyleInfo.Render(">")

	// Structural Styles
	HeaderStyle = lipgloss.NewStyle().Foreground(ColorMuted).Bold(true)
	
	BoxStyleWarn = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorWarn).
		Padding(1).
		MarginTop(1)
		
	BoxStylePrimary = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1).
		MarginTop(1)
)

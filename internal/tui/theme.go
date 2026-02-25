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

	// Plain-mode icons (no ANSI)
	PlainIconCheck = "[✓]"
	PlainIconCross = "[x]"
	PlainIconWarn  = "[!]"
	PlainIconInfo  = ">"

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

	// Centralised navigation hint strings
	TableNavHint  = "↑/↓ to navigate  •  q to quit"
	PickerNavHint = "↑/↓ to navigate  •  enter to select  •  esc to cancel"
	FormNavHint   = "enter to confirm  •  esc to cancel"
)

// RenderBanner returns a high-fidelity ASCII/styled banner for the CLI.
func RenderBanner() string {
	banner := `
  ______               _____                  _   
 |  ____|             / ____|                | |  
 | |__   _ ____   __ | |     _ __ _   _ _ __ | |_ 
 |  __| | '_ \ \ / / | |    | '__| | | | '_ \| __|
 | |____| | | \ V /  | |____| |  | |_| | |_) | |_ 
 |______|_| |_|\_/    \_____|_|   \__, | .__/ \__|
                                   __/ | |        
                                  |___/|_|        `

	logo := lipgloss.NewStyle().Foreground(ColorPrimary).Render(banner)
	subtitle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true).
		Render("    Zero-trust, end-to-end encrypted environment management")

	return lipgloss.JoinVertical(lipgloss.Left, logo, subtitle)
}

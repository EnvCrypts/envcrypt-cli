package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

//
// ─── COLUMN WIDTHS (TUNE ONCE) ─────────────────────────────────────────────
//

const (
	projectColWidth = 24
	roleColWidth    = 8
	statusColWidth  = 8
)

//
// ─── STYLES ────────────────────────────────────────────────────────────────
//

// Semantic colors (calm, not loud)
var (
	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("160"))

	warnStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220"))

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	mutedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	revokedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("160")).
		Strikethrough(true)

	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Underline(true)
)

//
// ─── ICONS ─────────────────────────────────────────────────────────────────
//

var (
	iconCheck = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")).
		Render("✔")

	iconCross = lipgloss.NewStyle().
		Foreground(lipgloss.Color("160")).
		Render("✖")

	iconWarn = lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Render("⚠")

	iconInfo = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Render("ℹ")
)

//
// ─── BASIC OUTPUT ──────────────────────────────────────────────────────────
//

func Spacer() {
	fmt.Println()
}

func Success(msg string) {
	Spacer()
	fmt.Printf("%s %s\n", iconCheck, successStyle.Render(msg))
}

func Info(msg string) {
	fmt.Printf("%s %s\n", iconInfo, infoStyle.Render(msg))
}

func Warn(msg string) {
	fmt.Printf("%s %s\n", iconWarn, warnStyle.Render(msg))
}

//
// ─── ERROR HANDLING ─────────────────────────────────────────────────────────
//

// Error formats an error for CLI UX.
// Caller should return this error, not print it.
func Error(msg string, err error) error {
	if err != nil {
		return fmt.Errorf(
			"%s %s\n  %s",
			iconCross,
			errorStyle.Render(msg),
			mutedStyle.Render("↳ "+err.Error()),
		)
	}

	return fmt.Errorf(
		"%s %s",
		iconCross,
		errorStyle.Render(msg),
	)
}

//
// ─── CONFIRMATION ──────────────────────────────────────────────────────────
//

func ConfirmDangerousAction(prompt, expected string) bool {
	Spacer()
	Warn(prompt)

	fmt.Printf(
		"%s Type %q to confirm: ",
		mutedStyle.Render("→"),
		expected,
	)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != expected {
		Info("Aborted.")
		return false
	}

	return true
}

// ConfirmOverwrite prompts for a simple y/n confirmation
func ConfirmOverwrite(path string) bool {
	Warn(fmt.Sprintf("File %q already exists. Overwrite?", path))
	fmt.Printf("%s [y/N]: ", mutedStyle.Render("→"))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "y" || input == "yes"
}

//
// ─── TABLE OUTPUT ──────────────────────────────────────────────────────────
//

func PrintProjects(projects []config.Project) {
	if len(projects) == 0 {
		fmt.Println(mutedStyle.Render("No projects found."))
		return
	}

	// Header
	fmt.Printf(
		"%s  %s  %s\n",
		headerStyle.Render(padRight("PROJECT", projectColWidth)),
		headerStyle.Render(padRight("ROLE", roleColWidth)),
		headerStyle.Render(padRight("STATUS", statusColWidth)),
	)

	// Rows
	for _, p := range projects {
		name := truncate(p.Name, projectColWidth)
		role := p.Role
		status := "active"

		switch p.Role {
		case "admin":
			role = successStyle.Render("admin")
		case "member":
			role = infoStyle.Render("member")
		default:
			role = mutedStyle.Render(p.Role)
		}

		if p.IsRevoked {
			name = revokedStyle.Render(name)
			status = errorStyle.Render("revoked")
		} else {
			status = successStyle.Render("active")
		}

		fmt.Printf(
			"%s  %s  %s\n",
			padRight(name, projectColWidth),
			padRight(role, roleColWidth),
			padRight(status, statusColWidth),
		)
	}
}

//
// ─── STRING HELPERS ────────────────────────────────────────────────────────
//

func truncate(s string, max int) string {
	if visibleLen(s) <= max {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}

func padRight(s string, width int) string {
	l := visibleLen(s)
	if l >= width {
		return s
	}
	return s + strings.Repeat(" ", width-l)
}

// visibleLen counts characters ignoring ANSI escape sequences
func visibleLen(s string) int {
	return len([]rune(stripANSI(s)))
}

// minimal ANSI stripper (safe for lipgloss)
func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

func printEnvSummary(env map[string]string) {
	if len(env) == 0 {
		Warn("No environment variables found")
		return
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	Info(fmt.Sprintf("Variables (%d):", len(keys)))

	// Don’t spam the terminal
	const maxShown = 10
	for i, k := range keys {
		if i == maxShown {
			fmt.Printf("  %s\n", mutedStyle.Render("…and more"))
			break
		}
		fmt.Printf("  %s %s\n", mutedStyle.Render("-"), k)
	}
}

func renderDiff(diff cryptoutils.DiffingResult, oldMap, newMap map[string]string, showSecrets bool) {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Modified) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	// Helper to mask values
	mask := func(val string) string {
		if showSecrets {
			return val
		}
		return "********"
	}

	// Styles
	addedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))     // Green
	removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("160"))  // Red
	modifiedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange/Yellow

	// Print Added
	for _, key := range diff.Added {
		fmt.Println(addedStyle.Render(fmt.Sprintf("+ %s=%s", key, mask(newMap[key]))))
	}

	// Print Removed
	for _, key := range diff.Removed {
		fmt.Println(removedStyle.Render(fmt.Sprintf("- %s=%s", key, mask(oldMap[key]))))
	}

	// Print Modified
	for _, key := range diff.Modified {
		fmt.Println(modifiedStyle.Render(fmt.Sprintf("~ %s: %s -> %s", key, mask(oldMap[key]), mask(newMap[key]))))
	}
}

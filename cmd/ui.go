package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"

	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

const (
	projectColWidth = 24
	roleColWidth    = 10
	statusColWidth  = 10
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	
	iconCheck = successStyle.Render("✓")
	iconCross = errorStyle.Render("×")
	iconWarn  = warnStyle.Render("!")
	iconInfo  = infoStyle.Render("•")
)

func Spacer() {
	fmt.Println()
}

func Success(msg string) {
	fmt.Printf("%s %s\n", iconCheck, msg)
}

func Info(msg string) {
	fmt.Printf("%s %s\n", iconInfo, msg)
}

func Warn(msg string) {
	fmt.Printf("%s %s\n", iconWarn, msg)
}

func Error(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s %s\n  %s", iconCross, msg, mutedStyle.Render(err.Error()))
	}
	return fmt.Errorf("%s %s", iconCross, msg)
}

func ConfirmDangerousAction(prompt, expected string) bool {
	Spacer()
	Warn(prompt)
	fmt.Printf("%s To confirm, type \"%s\": ", mutedStyle.Render("?"), expected)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != expected {
		fmt.Println(mutedStyle.Render("Aborted."))
		return false
	}
	return true
}

func ConfirmOverwrite(path string) bool {
	Warn(fmt.Sprintf("Overwrite %q?", path))
	fmt.Printf("%s [y/N]: ", mutedStyle.Render("?"))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	return input == "y" || input == "yes"
}

func PrintProjects(projects []config.Project) {
	if len(projects) == 0 {
		fmt.Println(mutedStyle.Render("No projects found."))
		return
	}

	fmt.Printf(
		"%s  %s  %s\n",
		headerStyle.Render(padRight("PROJECT", projectColWidth)),
		headerStyle.Render(padRight("ROLE", roleColWidth)),
		headerStyle.Render(padRight("STATUS", statusColWidth)),
	)

	for _, p := range projects {
		name := truncate(p.Name, projectColWidth)
		role := p.Role
		status := "active"

		if p.IsRevoked {
			status = errorStyle.Render("revoked")
			name = mutedStyle.Render(name)
		}

		fmt.Printf(
			"%s  %s  %s\n",
			padRight(name, projectColWidth),
			padRight(role, roleColWidth),
			padRight(status, statusColWidth),
		)
	}
}

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

func visibleLen(s string) int {
	return len([]rune(stripANSI(s)))
}

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

	fmt.Printf("%s %s\n", iconInfo, mutedStyle.Render(fmt.Sprintf("%d Environment Variables", len(keys))))

	const maxShown = 10
	for i, k := range keys {
		if i == maxShown {
			fmt.Printf("  %s\n", mutedStyle.Render(fmt.Sprintf("+%d more", len(keys)-maxShown)))
			break
		}
		fmt.Printf("  %s\n", k)
	}
}

func renderDiff(diff cryptoutils.DiffingResult, oldMap, newMap map[string]string, showSecrets bool) {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Modified) == 0 {
		fmt.Println(mutedStyle.Render("No changes."))
		return
	}

	mask := func(val string) string {
		if showSecrets {
			return val
		}
		return "********"
	}

	addedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	removedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	modifiedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	for _, key := range diff.Added {
		fmt.Println(addedStyle.Render(fmt.Sprintf("+ %s=%s", key, mask(newMap[key]))))
	}

	for _, key := range diff.Removed {
		fmt.Println(removedStyle.Render(fmt.Sprintf("- %s=%s", key, mask(oldMap[key]))))
	}

	for _, key := range diff.Modified {
		fmt.Println(modifiedStyle.Render(fmt.Sprintf("~ %s: %s → %s", key, mask(oldMap[key]), mask(newMap[key]))))
	}
}

func PrintServiceRoles(roles []config.ServiceRole) {
	if len(roles) == 0 {
		fmt.Println(mutedStyle.Render("No service roles found."))
		return
	}

	fmt.Printf(
		"%s  %s\n",
		headerStyle.Render(padRight("NAME", 30)),
		headerStyle.Render(padRight("REPO PRINCIPAL", 50)),
	)

	for _, r := range roles {
		fmt.Printf(
			"%s  %s\n",
			padRight(truncate(r.Name, 30), 30),
			padRight(truncate(r.RepoPrincipal, 50), 50),
		)
	}
}

func PrintServiceRoleSecret(keyPair *config.ServiceRoleKeyPair) {
	Spacer()
	Warn("This is a one-time view. Save these credentials securely!")
	fmt.Println(mutedStyle.Render("These keys allow read/write access to project secrets."))

	pub := base64.StdEncoding.EncodeToString(keyPair.PublicKey)
	priv := base64.StdEncoding.EncodeToString(keyPair.PrivateKey)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")).
		Padding(1).
		MarginTop(1)

	content := fmt.Sprintf(
		"ENVCRYPT_SERVICE_ROLE_PUBLIC_KEY=%s\nENVCRYPT_SERVICE_ROLE_PRIVATE_KEY=%s",
		pub, priv,
	)

	fmt.Println(boxStyle.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(content)))
}

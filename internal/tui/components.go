package tui

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

// Spacer prints a blank line (suppressed in JSON/quiet mode).
func Spacer() {
	if currentMode == ModeJSON || isQuiet {
		return
	}
	fmt.Println()
}

// Success prints a success message, dispatched by output mode.
func Success(msg string) {
	if isQuiet {
		return
	}
	switch currentMode {
	case ModeJSON:
		JSONMessage("success", msg)
	case ModePlain:
		fmt.Printf("%s %s %s\n", PlainTimestamp(), PlainIconCheck, msg)
	default:
		fmt.Printf("%s %s\n", IconCheck, msg)
	}
}

// Info prints an informational message, dispatched by output mode.
func Info(msg string) {
	if isQuiet {
		return
	}
	switch currentMode {
	case ModeJSON:
		JSONMessage("info", msg)
	case ModePlain:
		fmt.Printf("%s %s %s\n", PlainTimestamp(), PlainIconInfo, msg)
	default:
		fmt.Printf("%s %s\n", IconInfo, msg)
	}
}

// Warn prints a warning message, dispatched by output mode.
func Warn(msg string) {
	if isQuiet {
		return
	}
	switch currentMode {
	case ModeJSON:
		JSONMessage("warn", msg)
	case ModePlain:
		fmt.Printf("%s %s %s\n", PlainTimestamp(), PlainIconWarn, msg)
	default:
		fmt.Printf("%s %s\n", IconWarn, msg)
	}
}

// Error formats an error for display. Supports an optional hint string.
func Error(msg string, err error, hints ...string) error {
	hint := ""
	if len(hints) > 0 {
		hint = hints[0]
	}
	return ErrorWithHint(msg, err, hint)
}

func Truncate(s string, max int) string {
	if VisibleLen(s) <= max {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}

func PadRight(s string, width int) string {
	l := VisibleLen(s)
	if l >= width {
		return s
	}
	return s + strings.Repeat(" ", width-l)
}

func VisibleLen(s string) int {
	return len([]rune(StripANSI(s)))
}

func StripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}

func RenderDiff(diff cryptoutils.DiffingResult, oldMap, newMap map[string]string, showSecrets bool) {
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Modified) == 0 {
		if currentMode == ModeJSON {
			JSONData(map[string]any{"changes": []string{}})
			return
		}
		fmt.Println(StyleMuted.Render("No changes."))
		return
	}

	if currentMode == ModeJSON {
		jsonDiff := map[string]any{
			"added":    diff.Added,
			"removed":  diff.Removed,
			"modified": diff.Modified,
		}
		JSONData(jsonDiff)
		return
	}

	mask := func(val string) string {
		if showSecrets {
			return val
		}
		return "********"
	}

	// Colorized diff header
	fmt.Println(StyleError.Render("--- old"))
	fmt.Println(StyleSuccess.Render("+++ new"))
	fmt.Println()

	for _, key := range diff.Added {
		fmt.Println(StyleSuccess.Render(fmt.Sprintf("+ %s=%s", key, mask(newMap[key]))))
	}

	for _, key := range diff.Removed {
		fmt.Println(StyleError.Render(fmt.Sprintf("- %s=%s", key, mask(oldMap[key]))))
	}

	for _, key := range diff.Modified {
		fmt.Println(StyleWarn.Render(fmt.Sprintf("~ %s: %s → %s", key, mask(oldMap[key]), mask(newMap[key]))))
	}
}

func PrintEnvSummary(env map[string]string) {
	if len(env) == 0 {
		Warn("No environment variables found")
		return
	}

	if currentMode == ModeJSON {
		JSONData(map[string]any{"count": len(env), "keys": sortedKeys(env)})
		return
	}

	keys := sortedKeys(env)

	fmt.Printf("%s %s\n", IconInfo, StyleMuted.Render(fmt.Sprintf("%d Environment Variables", len(keys))))

	const maxShown = 10
	for i, k := range keys {
		if i == maxShown {
			fmt.Printf("  %s\n", StyleMuted.Render(fmt.Sprintf("+%d more", len(keys)-maxShown)))
			break
		}
		fmt.Printf("  %s\n", k)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func PrintServiceRoleSecret(keyPair *config.ServiceRoleKeyPair) {
	Spacer()
	Warn("This is a one-time view. Save these credentials securely!")
	fmt.Println(StyleMuted.Render("These keys allow read/write access to project secrets."))

	pub := base64.StdEncoding.EncodeToString(keyPair.PublicKey)
	priv := base64.StdEncoding.EncodeToString(keyPair.PrivateKey)

	if currentMode == ModeJSON {
		JSONData(map[string]string{
			"public_key":  pub,
			"private_key": priv,
		})
		return
	}

	content := fmt.Sprintf(
		"ENVCRYPT_SERVICE_ROLE_PUBLIC_KEY=%s\nENVCRYPT_SERVICE_ROLE_PRIVATE_KEY=%s",
		pub, priv,
	)

	fmt.Println(BoxStyleWarn.Render(StyleSuccess.Render(content)))
}

func PrintServiceRoleDetail(role *config.ServiceRole) {
	if currentMode == ModeJSON {
		JSONData(role)
		return
	}

	Info(fmt.Sprintf("Service Role: %s", role.Name))

	fmt.Printf("  %s %s\n", HeaderStyle.Render("ID:"), role.ID)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Repo:"), role.RepoPrincipal)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Created By:"), role.CreatedBy)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Created At:"), role.CreatedAt.Format("2006-01-02 15:04:05"))
}

func PrintServiceRolePermissions(perm *config.ServiceRolePermsResponse, repoPrincipal string) {
	if currentMode == ModeJSON {
		JSONData(map[string]string{
			"repo_principal": repoPrincipal,
			"project":        perm.ProjectName,
			"env":            perm.Env,
		})
		return
	}

	Info(fmt.Sprintf("Permissions for %s", repoPrincipal))

	fmt.Printf("  %s %s\n", HeaderStyle.Render("Project:"), perm.ProjectName)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Env:"), perm.Env)
}

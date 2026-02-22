package tui

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"regexp"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
)

// Replaces standard output functions
func Spacer() {
	fmt.Println()
}

func Success(msg string) {
	fmt.Printf("%s %s\n", IconCheck, msg)
}

func Info(msg string) {
	fmt.Printf("%s %s\n", IconInfo, msg)
}

func Warn(msg string) {
	fmt.Printf("%s %s\n", IconWarn, msg)
}

func Error(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s %s\n  %s", IconCross, msg, StyleMuted.Render(err.Error()))
	}
	return fmt.Errorf("%s %s", IconCross, msg)
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
		fmt.Println(StyleMuted.Render("No changes."))
		return
	}

	mask := func(val string) string {
		if showSecrets {
			return val
		}
		return "********"
	}

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

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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

func PrintServiceRoleSecret(keyPair *config.ServiceRoleKeyPair) {
	Spacer()
	Warn("This is a one-time view. Save these credentials securely!")
	fmt.Println(StyleMuted.Render("These keys allow read/write access to project secrets."))

	pub := base64.StdEncoding.EncodeToString(keyPair.PublicKey)
	priv := base64.StdEncoding.EncodeToString(keyPair.PrivateKey)

	content := fmt.Sprintf(
		"ENVCRYPT_SERVICE_ROLE_PUBLIC_KEY=%s\nENVCRYPT_SERVICE_ROLE_PRIVATE_KEY=%s",
		pub, priv,
	)

	fmt.Println(BoxStyleWarn.Render(StyleSuccess.Render(content)))
}

func PrintServiceRoleDetail(role *config.ServiceRole) {
	Info(fmt.Sprintf("Service Role: %s", role.Name))
	
	fmt.Printf("  %s %s\n", HeaderStyle.Render("ID:"), role.ID)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Repo:"), role.RepoPrincipal)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Created By:"), role.CreatedBy)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Created At:"), role.CreatedAt.Format("2006-01-02 15:04:05"))
}

func PrintServiceRolePermissions(perm *config.ServiceRolePermsResponse, repoPrincipal string) {
	Info(fmt.Sprintf("Permissions for %s", repoPrincipal))

	fmt.Printf("  %s %s\n", HeaderStyle.Render("Project:"), perm.ProjectName)
	fmt.Printf("  %s %s\n", HeaderStyle.Render("Env:"), perm.Env)
}

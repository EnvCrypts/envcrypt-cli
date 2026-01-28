package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/envcrypts/envcrypt-cli/internal/config"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgCyan)
	mutedColor   = color.New(color.FgHiBlack)

	revokedColor = color.New(color.FgRed, color.CrossedOut)
)

func init() {
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}
}

func Success(msg string) {
	successColor.Printf("✔ %s\n", msg)
}

func Error(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s %s: %w", errorColor.Sprint("✖"), msg, err)
	}
	return fmt.Errorf("%s %s", errorColor.Sprint("✖"), msg)
}

func Warn(msg string) {
	warnColor.Printf("⚠ %s\n", msg)
}

func Info(msg string) {
	infoColor.Printf("ℹ %s\n", msg)
}

func ConfirmDangerousAction(prompt, expected string) bool {
	Warn(prompt)
	fmt.Printf("Type %q to confirm: ", expected)

	var input string
	fmt.Scanln(&input)

	return input == expected
}

func PrintProjects(projects []config.Project) {
	if len(projects) == 0 {
		mutedColor.Println("No projects found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "PROJECT NAME\tROLE")
	fmt.Fprintln(w, "────────────\t────")

	for _, p := range projects {
		name := p.Name
		role := p.Role

		switch p.Role {
		case "admin":
			role = successColor.Sprint("admin")
		case "member":
			role = infoColor.Sprint("member")
		default:
			role = mutedColor.Sprint(p.Role)
		}

		if p.IsRevoked {
			if color.NoColor {
				name = name + " [REVOKED]"
				role = p.Role + " [REVOKED]"
			} else {
				name = revokedColor.Sprint(name)
				role = revokedColor.Sprint(p.Role + " (revoked)")
			}
		}

		fmt.Fprintf(w, "%s\t%s\n", name, role)
	}

	w.Flush()
}

package cmd

import "fmt"

func prettySuccess(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

func prettyError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("❌ %s: %w", msg, err)
	}
	return fmt.Errorf("❌ %s", msg)
}

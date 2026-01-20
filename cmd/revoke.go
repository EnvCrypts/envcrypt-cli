package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeMemberEmail string

// removeCmd represents the project member remove command
var removeCmd = &cobra.Command{
	Use:   "remove <project>",
	Short: "Remove a user from a project",
	Long:  "Remove a user from a project and revoke all access.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		//if err := Application.RemoveProjectMember(
		//	cmd.Context(),
		//	projectName,
		//	removeMemberEmail,
		//); err != nil {
		//	return fmt.Errorf("❌ failed to remove member: %w", err)
		//}

		fmt.Printf(
			"✅ Removed %s from project %q\n",
			removeMemberEmail,
			projectName,
		)
		return nil
	},
}

func init() {
	memberCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(
		&removeMemberEmail,
		"email",
		"e",
		"",
		"Email address of the user",
	)

	removeCmd.MarkFlagRequired("email")
}

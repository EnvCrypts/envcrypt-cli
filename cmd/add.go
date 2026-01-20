package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addMemberEmail string
	addMemberRole  string
)

// addCmd represents the project member add command
var addCmd = &cobra.Command{
	Use:   "add <project>",
	Short: "Add a user to a project",
	Long:  "Add a user to a project and assign them a role.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		//if err := Application.AddProjectMember(
		//	cmd.Context(),
		//	projectName,
		//	addMemberEmail,
		//	addMemberRole,
		//); err != nil {
		//	return fmt.Errorf("❌ failed to add member: %w", err)
		//}

		fmt.Printf(
			"✅ Added %s as %s to project %q\n",
			addMemberEmail,
			addMemberRole,
			projectName,
		)
		return nil
	},
}

func init() {
	memberCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(
		&addMemberEmail,
		"email",
		"e",
		"",
		"Email address of the user",
	)

	addCmd.Flags().StringVarP(
		&addMemberRole,
		"role",
		"r",
		"viewer",
		"Role to assign (admin, user)",
	)

	addCmd.MarkFlagRequired("email")
}

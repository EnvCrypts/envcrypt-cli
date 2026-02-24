package cmd

import (
	"fmt"
	"time"

	"github.com/envcrypts/envcrypt-cli/internal/app"
	"github.com/envcrypts/envcrypt-cli/internal/tui"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit operations",
	Long:  "View audit logs.",
}

var (
	auditLimit      int
	auditOffset     int
	auditActorEmail string
	auditAction     string
	auditStatus     string
	auditFrom       string
	auditTo         string
)

var auditProjectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Audit logs for a project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectName string
		if len(args) == 1 {
			projectName = args[0]
		}

		if projectName == "" {
			resp, err := Application.ListProjects(cmd.Context())
			if err != nil {
				return tui.Error("failed to fetch projects", err)
			}
			if len(resp.Projects) == 0 {
				return tui.Error("no projects found", nil)
			}
			names := make([]string, len(resp.Projects))
			for i, p := range resp.Projects {
				names[i] = p.Name
			}
			projectName, err = tui.RunPicker("Select a project to view audit logs", names)
			if err != nil {
				return tui.Cancelled()
			}
		}

		if auditLimit > 200 {
			auditLimit = 200
		}

		if auditFrom != "" {
			_, err := time.Parse(time.RFC3339, auditFrom)
			if err != nil {
				return fmt.Errorf("invalid from time format: %w", err)
			}
		}

		if auditTo != "" {
			_, err := time.Parse(time.RFC3339, auditTo)
			if err != nil {
				return fmt.Errorf("invalid to time format: %w", err)
			}
		}

		opts := app.AuditOptions{
			Limit:      auditLimit,
			Offset:     auditOffset,
			ActorEmail: auditActorEmail,
			Action:     auditAction,
			Status:     auditStatus,
			From:       auditFrom,
			To:         auditTo,
			JSON:       globalJSON,
		}

		resp, err := Application.GetProjectAuditLogs(cmd.Context(), projectName, opts)
		if err != nil {
			return err
		}

		return tui.RunAuditTable(resp.Logs, resp.Pagination.Total)
	},
}

func init() {
	auditProjectCmd.Flags().IntVar(&auditLimit, "limit", 50, "Limit the number of results (max 200)")
	auditProjectCmd.Flags().IntVar(&auditOffset, "offset", 0, "Pagination offset")
	auditProjectCmd.Flags().StringVar(&auditActorEmail, "actor", "", "Filter by actor email")
	auditProjectCmd.Flags().StringVar(&auditAction, "action", "", "Filter by action type (e.g. ENV_CREATE)")
	auditProjectCmd.Flags().StringVar(&auditStatus, "status", "", "Filter by status (success|failed)")
	auditProjectCmd.Flags().StringVar(&auditFrom, "from", "", "From timestamp (RFC3339)")
	auditProjectCmd.Flags().StringVar(&auditTo, "to", "", "To timestamp (RFC3339)")

	auditCmd.AddCommand(auditProjectCmd)
	rootCmd.AddCommand(auditCmd)
}

package config

import "time"

type AuditEntry struct {
	CreatedAt    time.Time `json:"timestamp"`
	ActorEmail   string    `json:"actor_email"`
	Action       string    `json:"action"`
	Status       string    `json:"status"`
	Environment  *string   `json:"environment,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

type ProjectAuditRequest struct {
	ProjectID  string `json:"project_id"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	ActorEmail string `json:"actor_email,omitempty"`
	Action     string `json:"action,omitempty"`
	Status     string `json:"status,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

type ProjectAuditPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type ProjectAuditResponse struct {
	Logs       []AuditEntry           `json:"logs"`
	Pagination ProjectAuditPagination `json:"pagination"`
}

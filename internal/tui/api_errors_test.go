package tui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapAPIErrorIncludesFriendlyTitleAndFields(t *testing.T) {
	err := MapAPIError(&APIErrorDetail{
		Status:  400,
		Code:    "VALIDATION_FAILED",
		Message: "Validation failed",
		Hint:    "Check the missing values",
		Fields: map[string]string{
			"checksum":         "checksum is required",
			"new_project_name": "new_project_name is required",
		},
	})

	msg := err.Error()
	require.Contains(t, msg, "Fix these fields")
	require.Contains(t, msg, "Validation failed")
	require.Contains(t, msg, "Fields:")
	require.Contains(t, msg, "checksum: checksum is required")
	require.Contains(t, msg, "Hint: Check the missing values")
}

func TestMapAPIErrorUsesStatusBasedTitle(t *testing.T) {
	err := MapAPIError(&APIErrorDetail{
		Status:  403,
		Code:    "FORBIDDEN",
		Message: "You don't have permission to access this environment",
	})

	require.Contains(t, err.Error(), "Insufficient permissions")
}

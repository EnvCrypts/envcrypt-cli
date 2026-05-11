package tui

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/client"
)

// APIErrorDetail represents the structured error returned by the server.
type APIErrorDetail struct {
	Status  int               `json:"status,omitempty"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Hint    string            `json:"hint,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// MapAPIError converts an API error into a user-friendly error string.
func MapAPIError(detail *APIErrorDetail) error {
	if detail == nil {
		return fmt.Errorf("unknown server error")
	}
	return errors.New(renderAPIErrorText(detail.Status, detail.Code, detail.Message, detail.Hint, detail.Fields))
}

// RenderError prints an error in the active output mode.
func RenderError(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		renderUIError(&UIError{
			Message: "Request timed out",
			Hint:    "Check your network connection and try again",
			Action:  "Retry the command",
		})
		return
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		renderUIError(&UIError{
			Message: "Request timed out",
			Hint:    "Check your network connection and try again",
			Action:  "Retry the command",
		})
		return
	}

	var uiErr *UIError
	if errors.As(err, &uiErr) {
		renderUIError(uiErr)
		return
	}

	var apiErr *client.HTTPError
	if errors.As(err, &apiErr) {
		renderHTTPError(apiErr)
		return
	}

	if currentMode == ModeJSON {
		JSONData(map[string]any{
			"level":   "error",
			"message": err.Error(),
			"hint":    "",
			"action":  "",
		})
		return
	}

	fmt.Fprintln(os.Stderr, err.Error())
}

func renderUIError(err *UIError) {
	if err == nil {
		return
	}

	if currentMode == ModeJSON {
		obj := map[string]any{
			"level":   "error",
			"message": err.Message,
			"hint":    err.Hint,
			"action":  err.Action,
		}
		if err.Err != nil {
			obj["detail"] = err.Err.Error()
		}
		JSONData(obj)
		return
	}

	if currentMode == ModePlain {
		fmt.Fprintf(os.Stderr, "%s %s %s\n", PlainTimestamp(), PlainIconCross, err.Message)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", IconCross, err.Message)
	}
	if err.Hint != "" {
		fmt.Fprintf(os.Stderr, "  %s %s\n", IconInfo, StyleInfo.Render(err.Hint))
	}
	if err.Action != "" {
		fmt.Fprintf(os.Stderr, "  %s %s\n", IconInfo, StyleInfo.Render(err.Action))
	}
}

func renderHTTPError(err *client.HTTPError) {
	if err == nil {
		return
	}

	if currentMode == ModeJSON {
		action := apiErrorAction(err.Status, err.Code)
		obj := map[string]any{
			"level":   "error",
			"status":  err.Status,
			"code":    err.Code,
			"message": renderAPIErrorTitle(err.Status, err.Code, err.Message),
			"hint":    err.Hint,
			"action":  action,
		}
		if err.Message != "" {
			obj["detail"] = err.Message
		}
		if len(err.Fields) > 0 {
			obj["fields"] = err.Fields
		}
		JSONData(obj)
		return
	}

	title := renderAPIErrorTitle(err.Status, err.Code, err.Message)
	if currentMode == ModePlain {
		fmt.Fprintf(os.Stderr, "%s %s %s\n", PlainTimestamp(), PlainIconCross, title)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", IconCross, title)
	}

	if err.Message != "" && err.Message != title {
		fmt.Fprintf(os.Stderr, "  %s\n", StyleMuted.Render(err.Message))
	}
	if len(err.Fields) > 0 {
		fmt.Fprintln(os.Stderr, "  Fields:")
		keys := sortedErrorFieldKeys(err.Fields)
		for _, key := range keys {
			fmt.Fprintf(os.Stderr, "    %s: %s\n", key, err.Fields[key])
		}
	}
	action := apiErrorAction(err.Status, err.Code)
	if err.Hint != "" {
		fmt.Fprintf(os.Stderr, "  %s %s\n", IconInfo, StyleInfo.Render(err.Hint))
	}
	if action != "" {
		fmt.Fprintf(os.Stderr, "  %s %s\n", IconInfo, StyleInfo.Render(action))
	}
}

func renderAPIErrorText(status int, code, message, hint string, fields map[string]string) string {
	var b strings.Builder
	b.WriteString(renderAPIErrorTitle(status, code, message))
	if message != "" && message != b.String() {
		b.WriteString("\n")
		b.WriteString(message)
	}
	if len(fields) > 0 {
		b.WriteString("\n")
		b.WriteString("Fields:")
		keys := sortedErrorFieldKeys(fields)
		for _, key := range keys {
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(key)
			b.WriteString(": ")
			b.WriteString(fields[key])
		}
	}
	if hint != "" {
		b.WriteString("\n")
		b.WriteString("Hint: ")
		b.WriteString(hint)
	}
	return b.String()
}

func renderAPIErrorTitle(status int, code, message string) string {
	normalizedCode := strings.ToUpper(code)

	switch {
	case normalizedCode == "VALIDATION_FAILED":
		return "Fix these fields"
	case status == http.StatusBadRequest || normalizedCode == "BAD_REQUEST":
		return "Malformed request"
	case status == http.StatusUnauthorized ||
		normalizedCode == "SESSION_MISSING" ||
		normalizedCode == "SESSION_INVALID" ||
		normalizedCode == "SESSION_EXPIRED" ||
		normalizedCode == "INVALID_CREDENTIALS" ||
		normalizedCode == "INVALID_OIDC_TOKEN" ||
		normalizedCode == "MISSING_IDENTITY" ||
		normalizedCode == "USER_NOT_FOUND":
		return "Re-authenticate"
	case status == http.StatusForbidden || normalizedCode == "FORBIDDEN":
		return "Insufficient permissions"
	case status == http.StatusNotFound || strings.HasSuffix(normalizedCode, "_NOT_FOUND"):
		return "Target missing"
	case status == http.StatusConflict || normalizedCode == "CONFLICT":
		return "Retry with a different name or state"
	case status >= http.StatusInternalServerError || normalizedCode == "INTERNAL_ERROR":
		return "Server failure"
	}

	if message != "" {
		return message
	}
	if code != "" {
		return code
	}
	return "Request failed"
}

func apiErrorAction(status int, code string) string {
	normalizedCode := strings.ToUpper(code)

	switch {
	case status == http.StatusUnauthorized ||
		normalizedCode == "SESSION_MISSING" ||
		normalizedCode == "SESSION_INVALID" ||
		normalizedCode == "SESSION_EXPIRED" ||
		normalizedCode == "INVALID_CREDENTIALS" ||
		normalizedCode == "INVALID_OIDC_TOKEN" ||
		normalizedCode == "MISSING_IDENTITY" ||
		normalizedCode == "USER_NOT_FOUND":
		return "Run 'envcrypt login' to re-authenticate"
	case status == http.StatusForbidden || normalizedCode == "FORBIDDEN":
		return "Ask a project admin for access"
	case status == http.StatusNotFound || strings.HasSuffix(normalizedCode, "_NOT_FOUND"):
		return "Verify the name and try again"
	case status == http.StatusConflict || normalizedCode == "CONFLICT":
		return "Retry with a different name or state"
	case normalizedCode == "VALIDATION_FAILED":
		return "Fix the invalid fields and retry"
	}
	return ""
}

func sortedErrorFieldKeys(fields map[string]string) []string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

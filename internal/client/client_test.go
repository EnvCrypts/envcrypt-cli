package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeErrorResponseParsesStructuredBody(t *testing.T) {
	body := []byte(`{"error":{"code":"VALIDATION_FAILED","message":"Validation failed","hint":"Fix the fields","fields":{"name":"name is required"}}}`)

	resp := decodeErrorResponse(body)

	require.Equal(t, "VALIDATION_FAILED", resp.Error.Code)
	require.Equal(t, "Validation failed", resp.Error.Message)
	require.Equal(t, "Fix the fields", resp.Error.Hint)
	require.Equal(t, map[string]string{"name": "name is required"}, resp.Error.Fields)
}

func TestDecodeErrorResponseFallsBackToRawBody(t *testing.T) {
	body := []byte("request failed")

	resp := decodeErrorResponse(body)

	require.Equal(t, "request failed", resp.Error.Message)
}

func TestHTTPErrorStringPrefersCodeAndMessage(t *testing.T) {
	err := &HTTPError{Status: 400, Code: "BAD_REQUEST", Message: "Invalid request body"}

	require.Equal(t, "BAD_REQUEST: Invalid request body", err.Error())
}

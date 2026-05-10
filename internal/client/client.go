package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Client struct {
	baseUrl string
	http    *http.Client
}

func NewClient(baseUrl string, client *http.Client) *Client {
	c := &Client{
		baseUrl: baseUrl,
		http:    client,
	}

	return c
}

func (c *Client) SetBaseURL(url string) {
	c.baseUrl = url
}

func (c *Client) BaseURL() string {
	return c.baseUrl
}

// ErrorDetail represents the structured error body from the server.
type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Hint    string            `json:"hint,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// ErrorResponse wraps the server error envelope: {"error": {...}}.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func (c *Client) Do(
	ctx context.Context,
	method string,
	path string,
	body any,
	out any,
	protected bool,
) error {

	err := c.doOnce(ctx, method, path, body, out, protected)
	if err == nil {
		return nil
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.Status != http.StatusUnauthorized {
		return err
	}

	if protected {
		if err := c.Refresh(ctx); err != nil {
			return fmt.Errorf("refresh failed: %w", err)
		}
		return c.doOnce(ctx, method, path, body, out, protected)
	}

	return err
}

func (c *Client) doOnce(
	ctx context.Context,
	method string,
	path string,
	body any,
	out any,
	protected bool,
) error {

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseUrl+path,
		&buf,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if protected {
		token := viper.GetString("user.access_token")
		req.Header.Set("X-Session-ID", token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		errResp := decodeErrorResponse(body)

		return &HTTPError{
			Status:  resp.StatusCode,
			Code:    errResp.Error.Code,
			Message: errResp.Error.Message,
			Hint:    errResp.Error.Hint,
			Fields:  errResp.Error.Fields,
		}
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

func (c *Client) Refresh(ctx context.Context) error {

	userID := viper.GetString("user.id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	req := config.RefreshRequestBody{
		UserID: uid,
	}

	var resp config.RefreshResponseBody

	err = c.doOnce(ctx, "POST", "/users/refresh", req, &resp, false)
	if err != nil {
		return err
	}

	err = cryptoutils.SaveRefreshToken(resp.Session.RefreshToken.String())
	if err != nil {
		return err
	}

	err = cryptoutils.SaveAccessToken(resp.Session.AccessToken.String())
	if err != nil {
		return err
	}

	return nil
}

type HTTPError struct {
	Status  int
	Code    string
	Message string
	Hint    string
	Fields  map[string]string
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "unknown error"
	}
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return fmt.Sprintf("HTTP %d", e.Status)
}

func decodeErrorResponse(body []byte) ErrorResponse {
	var errResp ErrorResponse
	if len(body) == 0 {
		return errResp
	}

	if err := json.Unmarshal(body, &errResp); err == nil && (errResp.Error.Code != "" || errResp.Error.Message != "" || errResp.Error.Hint != "" || len(errResp.Error.Fields) > 0) {
		return errResp
	}

	msg := strings.TrimSpace(string(body))
	if msg == "" {
		msg = http.StatusText(http.StatusInternalServerError)
	}
	errResp.Error.Message = msg
	return errResp
}

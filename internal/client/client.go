package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/envcrypts/envcrypt-cli/internal/config"
	cryptoutils "github.com/envcrypts/envcrypt-cli/internal/crypto"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Client struct {
	baseUrl     string
	http        *http.Client
	accessToken uuid.UUID
}

func NewClient(baseUrl string, client *http.Client) *Client {
	return &Client{
		baseUrl: baseUrl,
		http:    client,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
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
		req.Header.Set("X-Session-ID", c.accessToken.String())
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResp)

		return &HTTPError{
			Status: resp.StatusCode,
			Body:   errResp.Error,
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

	c.accessToken = resp.Session.AccessToken
	err = cryptoutils.SaveRefreshToken(resp.Session.RefreshToken.String())
	if err != nil {
		return err
	}

	return nil
}

type HTTPError struct {
	Status int
	Body   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s", e.Body)
}

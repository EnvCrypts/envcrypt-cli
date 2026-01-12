package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseUrl string
	http    *http.Client
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

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var err ErrorResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&err)
		if decodeErr != nil {
			return decodeErr
		}
		return fmt.Errorf(err.Error)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}

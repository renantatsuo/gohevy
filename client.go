package hevy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Version string

const (
	V1 Version = "v1"
)

const (
	baseURL        = "https://api.hevyapp.com"
	defaultVersion = V1
	defaultTimeout = 30 * time.Second
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client represents a Hevy API client
type Client struct {
	apiKey     string
	baseURL    string
	version    Version
	httpClient Doer
}

// NewClient creates a new Hevy client with the given API key and optional configuration options.
func NewClient(apiKey string, opts ...func(*Client)) *Client {
	c := &Client{
		apiKey:  apiKey,
		version: defaultVersion,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.baseURL = fmt.Sprintf("%s/%s", baseURL, c.version)
	return c
}

// WithVersion sets the version of the API to use.
func WithVersion(version Version) func(*Client) {
	return func(c *Client) {
		c.version = version
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client Doer) func(*Client) {
	return func(c *Client) {
		c.httpClient = client
	}
}

// request executes an HTTP request with JSON marshaling/unmarshaling
func (c *Client) request(ctx context.Context, method, path string, body, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(bodyBytes),
			Method:     method,
			URL:        req.URL.String(),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

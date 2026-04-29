// Package hevy provides a Go client for the Hevy fitness tracking API (https://api.hevyapp.com/docs).
//
// Hevy is a workout tracking application. This package wraps the Hevy REST API v1,
// which requires a Hevy PRO subscription and an API key obtained from
// https://hevy.com/settings?developer.
//
// # Creating a client
//
//	client := hevy.NewClient(os.Getenv("HEVY_API_KEY"))
//
// All methods accept a context.Context as the first argument for cancellation and timeout control.
//
// # Pagination
//
// List endpoints accept PaginationParams and return a paginated response with Page (1-based)
// and PageCount fields. Iterate by incrementing Page until Page == PageCount.
//
// # Error handling
//
// On non-2xx responses, methods return *APIError. Use errors.As to extract it and
// IsClientError / IsServerError to classify the failure.
//
//	var apiErr *hevy.APIError
//	if errors.As(err, &apiErr) && apiErr.IsClientError() { ... }
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

// Version represents a Hevy REST API version string (e.g. "v1").
type Version string

const (
	// V1 is the current stable Hevy API version and the default used by NewClient.
	V1 Version = "v1"
)

const (
	baseURL        = "https://api.hevyapp.com"
	defaultVersion = V1
	defaultTimeout = 30 * time.Second
)

// Doer is the HTTP transport interface used by Client. It is satisfied by *http.Client.
// Implement this interface to inject test doubles, add middleware, or use a custom transport.
// Pass a Doer implementation to NewClient via WithHTTPClient.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the Hevy API client. Create one with NewClient.
// All methods on Client are safe for concurrent use.
type Client struct {
	apiKey     string
	baseURL    string
	version    Version
	httpClient Doer
}

// NewClient creates a new Hevy API client authenticated with the given API key.
// By default the client targets API version V1 with a 30-second HTTP timeout.
// Use WithVersion or WithHTTPClient to override these defaults.
//
//	client := hevy.NewClient(os.Getenv("HEVY_API_KEY"))
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

// WithVersion overrides the Hevy API version. The default is V1.
//
//	client := hevy.NewClient(apiKey, hevy.WithVersion(hevy.V1))
func WithVersion(version Version) func(*Client) {
	return func(c *Client) {
		c.version = version
	}
}

// WithHTTPClient replaces the default *http.Client (30s timeout) with the given Doer.
// Use this to set a custom timeout, inject a test double, or add transport-level middleware.
//
//	client := hevy.NewClient(apiKey, hevy.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}))
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

package hevy

import "fmt"

type version string

const (
	v1 version = "v1"
)
const (
	baseURL        = "https://api.hevyapp.com"
	defaultVersion = v1
)

// Client represents a Hevy API client
type Client struct {
	apiKey  string
	baseURL string
	version version
}

// NewClient creates a new Hevy client with the given API key and optional configuration options.
func NewClient(apiKey string, opts ...func(*Client)) *Client {
	c := &Client{apiKey: apiKey, version: defaultVersion}
	for _, opt := range opts {
		opt(c)
	}
	c.baseURL = fmt.Sprintf("%s/%s", baseURL, c.version)
	return c
}

// WithVersion sets the version of the API to use.
func WithVersion(version version) func(*Client) {
	return func(c *Client) {
		c.version = version
	}
}

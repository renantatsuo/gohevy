package hevy

import "fmt"

type APIError struct {
	StatusCode int
	Status     string
	Body       string
	Method     string
	URL        string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("%s %s failed: %s (status: %d) - %s", e.Method, e.URL, e.Status, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("%s %s failed: %s (status: %d)", e.Method, e.URL, e.Status, e.StatusCode)
}

// IsClientError returns true if the error is a 4xx client error
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

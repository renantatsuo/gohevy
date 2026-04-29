package hevy

import "fmt"

// APIError is returned when the Hevy API responds with a non-2xx HTTP status code.
// Use errors.As to extract it from a method's error return value.
//
//	var apiErr *hevy.APIError
//	if errors.As(err, &apiErr) {
//	    fmt.Println(apiErr.StatusCode) // e.g. 401, 404, 429
//	}
type APIError struct {
	StatusCode int    // HTTP status code (e.g. 401, 404, 422, 429, 500)
	Status     string // HTTP status text (e.g. "401 Unauthorized")
	Body       string // Raw response body; may contain a JSON error message from the API
	Method     string // HTTP method of the failed request (e.g. "GET", "POST")
	URL        string // Full URL of the failed request
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

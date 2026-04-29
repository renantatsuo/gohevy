package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetRoutines retrieves a paginated list of routines.
// Use PaginationParams to specify the page (1-based) and page size.
// Iterate until the returned PageCount is reached to retrieve all routines.
// Returns *APIError on failure (e.g. StatusCode 401 for invalid API key, 429 for rate limiting).
func (c *Client) GetRoutines(ctx context.Context, params PaginationParams) (res *PaginatedRoutinesResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/routines?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetRoutine retrieves a single routine by its unique ID.
// Returns *APIError with StatusCode 404 if no routine with the given ID exists.
func (c *Client) GetRoutine(ctx context.Context, routineID string) (res *Routine, err error) {
	path := fmt.Sprintf("/routines/%s", routineID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateRoutine creates a new routine and returns the server-assigned record (including ID and timestamps).
// Populate Exercises and their RoutineSets with target values. The ID field of the input is ignored.
func (c *Client) CreateRoutine(ctx context.Context, routine Routine) (res *Routine, err error) {
	err = c.request(ctx, http.MethodPost, "/routines", routine, &res)
	return
}

// UpdateRoutine replaces the routine identified by routineID with the provided data.
// This is a full replacement (PUT), not a partial update — all exercises and sets must be included.
// Returns *APIError with StatusCode 404 if no routine with the given ID exists.
func (c *Client) UpdateRoutine(ctx context.Context, routineID string, routine Routine) (res *Routine, err error) {
	path := fmt.Sprintf("/routines/%s", routineID)

	err = c.request(ctx, http.MethodPut, path, routine, &res)
	return
}

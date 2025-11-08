package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetRoutines retrieves a paginated list of routines
func (c *Client) GetRoutines(ctx context.Context, params PaginationParams) (res *PaginatedRoutinesResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/routines?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetRoutine retrieves a single routine by ID
func (c *Client) GetRoutine(ctx context.Context, routineID string) (res *Routine, err error) {
	path := fmt.Sprintf("/routines/%s", routineID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateRoutine creates a new routine
func (c *Client) CreateRoutine(ctx context.Context, routine Routine) (res *Routine, err error) {
	err = c.request(ctx, http.MethodPost, "/routines", routine, &res)
	return
}

// UpdateRoutine updates an existing routine
func (c *Client) UpdateRoutine(ctx context.Context, routineID string, routine Routine) (res *Routine, err error) {
	path := fmt.Sprintf("/routines/%s", routineID)

	err = c.request(ctx, http.MethodPut, path, routine, &res)
	return
}

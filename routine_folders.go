package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetRoutineFolders retrieves a paginated list of routine folders.
// Use PaginationParams to specify the page (1-based) and page size.
// Iterate until the returned PageCount is reached to retrieve all folders.
// Returns *APIError on failure (e.g. StatusCode 401 for invalid API key, 429 for rate limiting).
func (c *Client) GetRoutineFolders(ctx context.Context, params PaginationParams) (res *PaginatedRoutineFoldersResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/routine_folders?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetRoutineFolder retrieves a single routine folder by its unique ID.
// Returns *APIError with StatusCode 404 if no folder with the given ID exists.
func (c *Client) GetRoutineFolder(ctx context.Context, folderID int) (res *RoutineFolder, err error) {
	path := fmt.Sprintf("/routine_folders/%d", folderID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateRoutineFolder creates a new routine folder and returns the server-assigned record.
// Use folders to organize related routines. The ID field of the input is ignored; the server assigns a new ID.
func (c *Client) CreateRoutineFolder(ctx context.Context, folder RoutineFolder) (res *RoutineFolder, err error) {
	err = c.request(ctx, http.MethodPost, "/routine_folders", folder, &res)
	return
}

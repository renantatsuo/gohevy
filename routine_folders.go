package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetRoutineFolders retrieves a paginated list of routine folders
func (c *Client) GetRoutineFolders(ctx context.Context, params PaginationParams) (res *PaginatedRoutineFoldersResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/routine_folders?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetRoutineFolder retrieves a single routine folder by ID
func (c *Client) GetRoutineFolder(ctx context.Context, folderID int) (res *RoutineFolder, err error) {
	path := fmt.Sprintf("/routine_folders/%d", folderID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateRoutineFolder creates a new routine folder
func (c *Client) CreateRoutineFolder(ctx context.Context, folder RoutineFolder) (res *RoutineFolder, err error) {
	err = c.request(ctx, http.MethodPost, "/routine_folders", folder, &res)
	return
}

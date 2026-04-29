package hevy

import (
	"context"
	"net/http"
)

// GetUserInfo returns profile metadata for the authenticated API key (GET /v1/user/info).
func (c *Client) GetUserInfo(ctx context.Context) (*UserInfoResponse, error) {
	var res UserInfoResponse
	err := c.request(ctx, http.MethodGet, "/user/info", nil, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

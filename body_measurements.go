package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetBodyMeasurements lists body measurements with pagination (GET /v1/body_measurements).
func (c *Client) GetBodyMeasurements(ctx context.Context, params PaginationParams) (res *PaginatedBodyMeasurements, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))
	path := fmt.Sprintf("/body_measurements?%s", urlParams.Encode())
	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateBodyMeasurement creates a measurement row for the date in body.Date (POST /v1/body_measurements).
// Returns *APIError with StatusCode 409 if an entry already exists for that date.
func (c *Client) CreateBodyMeasurement(ctx context.Context, body BodyMeasurement) error {
	return c.request(ctx, http.MethodPost, "/body_measurements", body, nil)
}

// GetBodyMeasurement returns the measurement for date (YYYY-MM-DD) (GET /v1/body_measurements/{date}).
func (c *Client) GetBodyMeasurement(ctx context.Context, date string) (res *BodyMeasurement, err error) {
	path := fmt.Sprintf("/body_measurements/%s", url.PathEscape(date))
	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// UpdateBodyMeasurement replaces all fields for date (PUT /v1/body_measurements/{date}). Omitted fields become null per API.
func (c *Client) UpdateBodyMeasurement(ctx context.Context, date string, body PutBodyMeasurement) error {
	path := fmt.Sprintf("/body_measurements/%s", url.PathEscape(date))
	return c.request(ctx, http.MethodPut, path, body, nil)
}

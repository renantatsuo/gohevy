package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetExerciseTemplates retrieves a paginated list of exercise templates, including both Hevy built-ins and user-created custom templates.
// Use PaginationParams to specify the page (1-based) and page size.
// Iterate until the returned PageCount is reached to retrieve all templates.
// Returns *APIError on failure (e.g. StatusCode 401 for invalid API key, 429 for rate limiting).
func (c *Client) GetExerciseTemplates(ctx context.Context, params PaginationParams) (res *PaginatedExerciseTemplatesResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/exercise_templates?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetExerciseTemplate retrieves a single exercise template by its unique ID.
// Returns *APIError with StatusCode 404 if no template with the given ID exists.
func (c *Client) GetExerciseTemplate(ctx context.Context, templateID string) (res *ExerciseTemplate, err error) {
	path := fmt.Sprintf("/exercise_templates/%s", templateID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateExerciseTemplate creates a new custom exercise template and returns the server-assigned record.
// Use this to define exercises not available in Hevy's built-in template library.
// The ID field of the input is ignored; the server assigns a new ID.
func (c *Client) CreateExerciseTemplate(ctx context.Context, template ExerciseTemplate) (res *ExerciseTemplate, err error) {
	err = c.request(ctx, http.MethodPost, "/exercise_templates", template, &res)
	return
}

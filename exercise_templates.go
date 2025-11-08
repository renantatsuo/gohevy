package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetExerciseTemplates retrieves a paginated list of exercise templates
func (c *Client) GetExerciseTemplates(ctx context.Context, params PaginationParams) (res *PaginatedExerciseTemplatesResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/exercise_templates?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetExerciseTemplate retrieves a single exercise template by ID
func (c *Client) GetExerciseTemplate(ctx context.Context, templateID string) (res *ExerciseTemplate, err error) {
	path := fmt.Sprintf("/exercise_templates/%s", templateID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateExerciseTemplate creates a new custom exercise template
func (c *Client) CreateExerciseTemplate(ctx context.Context, template ExerciseTemplate) (res *ExerciseTemplate, err error) {
	err = c.request(ctx, http.MethodPost, "/exercise_templates", template, &res)
	return
}

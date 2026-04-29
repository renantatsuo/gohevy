package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetExerciseHistory retrieves exercise history rows for exerciseTemplateID (OpenAPI exercise_history list).
// Pass nil params if no date filters are needed.
func (c *Client) GetExerciseHistory(ctx context.Context, exerciseTemplateID string, params *ExerciseHistoryParams) (res []ExerciseHistoryEntry, err error) {
	path := fmt.Sprintf("/exercise_history/%s", exerciseTemplateID)
	if params != nil && (params.StartDate != nil || params.EndDate != nil) {
		q := url.Values{}
		if params.StartDate != nil {
			q.Set("start_date", params.StartDate.UTC().Format(time.RFC3339))
		}
		if params.EndDate != nil {
			q.Set("end_date", params.EndDate.UTC().Format(time.RFC3339))
		}
		path += "?" + q.Encode()
	}

	var wrapped exerciseHistoryAPIResponse
	err = c.request(ctx, http.MethodGet, path, nil, &wrapped)
	if err != nil {
		return nil, err
	}
	return wrapped.ExerciseHistory, nil
}

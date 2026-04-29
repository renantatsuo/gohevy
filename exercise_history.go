package hevy

import (
	"context"
	"fmt"
	"net/http"
)

// GetExerciseHistory retrieves all recorded sets for the exercise identified by exerciseTemplateID,
// grouped by workout. Returns a flat slice of ExerciseHistoryEntry (not a paginated response).
// Use this to review an athlete's progression over time for a given exercise.
// Returns *APIError with StatusCode 404 if no exercise template with the given ID exists.
func (c *Client) GetExerciseHistory(ctx context.Context, exerciseTemplateID string) (res []ExerciseHistoryEntry, err error) {
	path := fmt.Sprintf("/exercise_history/%s", exerciseTemplateID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

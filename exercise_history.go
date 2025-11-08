package hevy

import (
	"context"
	"fmt"
	"net/http"
)

// GetExerciseHistory retrieves exercise history for a specific exercise template
func (c *Client) GetExerciseHistory(ctx context.Context, exerciseTemplateID string) (res []ExerciseHistoryEntry, err error) {
	path := fmt.Sprintf("/exercise_history/%s", exerciseTemplateID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

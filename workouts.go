package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetWorkouts retrieves a paginated list of workouts, ordered by most recently updated first.
// Use PaginationParams to specify the page (1-based) and page size.
// Iterate until the returned PageCount is reached to retrieve all workouts.
// Returns *APIError on failure (e.g. StatusCode 401 for invalid API key, 429 for rate limiting).
func (c *Client) GetWorkouts(ctx context.Context, params PaginationParams) (res *PaginatedWorkoutsResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/workouts?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetWorkout retrieves a single workout by its unique ID.
// Returns *APIError with StatusCode 404 if no workout with the given ID exists.
func (c *Client) GetWorkout(ctx context.Context, workoutID string) (res *Workout, err error) {
	path := fmt.Sprintf("/workouts/%s", workoutID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateWorkout creates a new workout and returns the server-assigned record (including ID and timestamps).
// Populate Exercises and their Sets to record performance data. The ID field of the input is ignored.
func (c *Client) CreateWorkout(ctx context.Context, workout Workout) (res *Workout, err error) {
	err = c.request(ctx, http.MethodPost, "/workouts", workout, &res)
	return
}

// UpdateWorkout replaces the workout identified by workoutID with the provided data.
// This is a full replacement (PUT), not a partial update — all exercises and sets must be included.
// Returns *APIError with StatusCode 404 if no workout with the given ID exists.
func (c *Client) UpdateWorkout(ctx context.Context, workoutID string, workout Workout) (res *Workout, err error) {
	path := fmt.Sprintf("/workouts/%s", workoutID)

	err = c.request(ctx, http.MethodPut, path, workout, &res)
	return
}

// GetWorkoutsCount returns the total number of workouts logged on the authenticated account.
func (c *Client) GetWorkoutsCount(ctx context.Context) (res *WorkoutCountResponse, err error) {
	err = c.request(ctx, http.MethodGet, "/workouts/count", nil, &res)
	return
}

// GetWorkoutEvents retrieves a paginated stream of workout change events (creates, updates, deletions)
// that occurred after params.Since. Use this for incremental sync: poll periodically and process
// each event's Type to determine whether to upsert or delete a local record.
// Call GetWorkout with the ID from UpdatedWorkout to fetch full workout data after an "updated" event.
// Note: workout deletion events can appear even though the library has no DeleteWorkout endpoint —
// deletions originate from the Hevy app. Passing a zero-value Since returns all events (per Hevy API docs).
func (c *Client) GetWorkoutEvents(ctx context.Context, params WorkoutEventsParams) (res *PaginatedWorkoutEvents, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))
	urlParams.Add("since", params.Since.Format(time.RFC3339))

	path := fmt.Sprintf("/workouts/events?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

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
// Populate Exercises and their Sets to record performance data. Set Workout.IsPrivate for the POST body.
func (c *Client) CreateWorkout(ctx context.Context, workout Workout) (res *Workout, err error) {
	body := workoutToPostBody(workout)
	err = c.request(ctx, http.MethodPost, "/workouts", body, &res)
	return
}

// UpdateWorkout replaces the workout identified by workoutID with the provided data (OpenAPI PostWorkoutsRequestBody).
// Returns *APIError with StatusCode 404 if no workout with the given ID exists.
func (c *Client) UpdateWorkout(ctx context.Context, workoutID string, workout Workout) (res *Workout, err error) {
	path := fmt.Sprintf("/workouts/%s", workoutID)
	body := workoutToPostBody(workout)
	err = c.request(ctx, http.MethodPut, path, body, &res)
	return
}

// GetWorkoutsCount returns the total number of workouts logged on the authenticated account.
func (c *Client) GetWorkoutsCount(ctx context.Context) (res *WorkoutCountResponse, err error) {
	err = c.request(ctx, http.MethodGet, "/workouts/count", nil, &res)
	return
}

// GetWorkoutEvents retrieves a paginated stream of workout change events (creates, updates, deletions).
// For Type=="updated", Event.Workout contains the full workout payload. For Type=="deleted", use Event.ID.
// Pass a zero-value Since to omit the since query parameter (server default applies).
func (c *Client) GetWorkoutEvents(ctx context.Context, params WorkoutEventsParams) (res *PaginatedWorkoutEvents, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))
	if !params.Since.IsZero() {
		urlParams.Add("since", params.Since.UTC().Format(time.RFC3339))
	}

	path := fmt.Sprintf("/workouts/events?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

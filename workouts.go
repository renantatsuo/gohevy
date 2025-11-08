package hevy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetWorkouts retrieves a paginated list of workouts
func (c *Client) GetWorkouts(ctx context.Context, params PaginationParams) (res *PaginatedWorkoutsResponse, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))

	path := fmt.Sprintf("/workouts?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// GetWorkout retrieves a single workout by ID
func (c *Client) GetWorkout(ctx context.Context, workoutID string) (res *Workout, err error) {
	path := fmt.Sprintf("/workouts/%s", workoutID)

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

// CreateWorkout creates a new workout
func (c *Client) CreateWorkout(ctx context.Context, workout Workout) (res *Workout, err error) {
	err = c.request(ctx, http.MethodPost, "/workouts", workout, &res)
	return
}

// UpdateWorkout updates an existing workout
func (c *Client) UpdateWorkout(ctx context.Context, workoutID string, workout Workout) (res *Workout, err error) {
	path := fmt.Sprintf("/workouts/%s", workoutID)

	err = c.request(ctx, http.MethodPut, path, workout, &res)
	return
}

// GetWorkoutsCount retrieves the total number of workouts on the account
func (c *Client) GetWorkoutsCount(ctx context.Context) (res *WorkoutCountResponse, err error) {
	err = c.request(ctx, http.MethodGet, "/workouts/count", nil, &res)
	return
}

// GetWorkoutEvents retrieves a paged list of workout events (updates or deletes) since a given date
func (c *Client) GetWorkoutEvents(ctx context.Context, params WorkoutEventsParams) (res *PaginatedWorkoutEvents, err error) {
	urlParams := url.Values{}
	urlParams.Add("page", strconv.Itoa(params.Page))
	urlParams.Add("pageSize", strconv.Itoa(params.PageSize))
	urlParams.Add("since", params.Since.Format(time.RFC3339))

	path := fmt.Sprintf("/workouts/events?%s", urlParams.Encode())

	err = c.request(ctx, http.MethodGet, path, nil, &res)
	return
}

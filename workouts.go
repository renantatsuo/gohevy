package hevy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetWorkouts retrieves a paginated list of workouts
func (c *Client) GetWorkouts(params PaginationParams) (*PaginatedWorkoutsResponse, error) {
	url := fmt.Sprintf("%s/workouts?page=%d&pageSize=%d", c.baseURL, params.Page, params.PageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workouts: %s", resp.Status)
	}

	var result PaginatedWorkoutsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWorkout retrieves a single workout by ID
func (c *Client) GetWorkout(workoutID string) (*Workout, error) {
	url := fmt.Sprintf("%s/workouts/%s", c.baseURL, workoutID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workout: %s", resp.Status)
	}

	var workout Workout
	err = json.NewDecoder(resp.Body).Decode(&workout)
	if err != nil {
		return nil, err
	}

	return &workout, nil
}

// CreateWorkout creates a new workout
func (c *Client) CreateWorkout(workout Workout) (*Workout, error) {
	url := fmt.Sprintf("%s/workouts", c.baseURL)

	body, err := json.Marshal(workout)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create workout: %s", resp.Status)
	}

	var result Workout
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateWorkout updates an existing workout
func (c *Client) UpdateWorkout(workoutID string, workout Workout) (*Workout, error) {
	url := fmt.Sprintf("%s/workouts/%s", c.baseURL, workoutID)

	body, err := json.Marshal(workout)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update workout: %s", resp.Status)
	}

	var result Workout
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWorkoutsCount retrieves the total number of workouts on the account
func (c *Client) GetWorkoutsCount() (*WorkoutCountResponse, error) {
	url := fmt.Sprintf("%s/workouts/count", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workouts count: %s", resp.Status)
	}

	var result WorkoutCountResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// WorkoutEventsParams holds parameters for retrieving workout events
type WorkoutEventsParams struct {
	PaginationParams
	Since time.Time `json:"since"`
}

// GetWorkoutEvents retrieves a paged list of workout events (updates or deletes) since a given date
func (c *Client) GetWorkoutEvents(params WorkoutEventsParams) (*PaginatedWorkoutEvents, error) {
	url := fmt.Sprintf("%s/workouts/events?page=%d&pageSize=%d&since=%s",
		c.baseURL, params.Page, params.PageSize, params.Since.Format(time.RFC3339))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workout events: %s", resp.Status)
	}

	var result PaginatedWorkoutEvents
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

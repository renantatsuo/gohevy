package hevy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRoutines retrieves a paginated list of routines
func (c *Client) GetRoutines(params PaginationParams) (*PaginatedRoutinesResponse, error) {
	url := fmt.Sprintf("%s/routines?page=%d&pageSize=%d", c.baseURL, params.Page, params.PageSize)
	
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
		return nil, fmt.Errorf("failed to get routines: %s", resp.Status)
	}

	var result PaginatedRoutinesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRoutine retrieves a single routine by ID
func (c *Client) GetRoutine(routineID string) (*Routine, error) {
	url := fmt.Sprintf("%s/routines/%s", c.baseURL, routineID)
	
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
		return nil, fmt.Errorf("failed to get routine: %s", resp.Status)
	}

	var routine Routine
	err = json.NewDecoder(resp.Body).Decode(&routine)
	if err != nil {
		return nil, err
	}

	return &routine, nil
}

// CreateRoutine creates a new routine
func (c *Client) CreateRoutine(routine Routine) (*Routine, error) {
	url := fmt.Sprintf("%s/routines", c.baseURL)
	
	body, err := json.Marshal(routine)
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
		return nil, fmt.Errorf("failed to create routine: %s", resp.Status)
	}

	var result Routine
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateRoutine updates an existing routine
func (c *Client) UpdateRoutine(routineID string, routine Routine) (*Routine, error) {
	url := fmt.Sprintf("%s/routines/%s", c.baseURL, routineID)
	
	body, err := json.Marshal(routine)
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
		return nil, fmt.Errorf("failed to update routine: %s", resp.Status)
	}

	var result Routine
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


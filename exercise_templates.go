package hevy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExerciseTemplates retrieves a paginated list of exercise templates
func (c *Client) GetExerciseTemplates(params PaginationParams) (*PaginatedExerciseTemplatesResponse, error) {
	url := fmt.Sprintf("%s/exercise_templates?page=%d&pageSize=%d", c.baseURL, params.Page, params.PageSize)
	
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
		return nil, fmt.Errorf("failed to get exercise templates: %s", resp.Status)
	}

	var result PaginatedExerciseTemplatesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExerciseTemplate retrieves a single exercise template by ID
func (c *Client) GetExerciseTemplate(templateID string) (*ExerciseTemplate, error) {
	url := fmt.Sprintf("%s/exercise_templates/%s", c.baseURL, templateID)
	
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
		return nil, fmt.Errorf("failed to get exercise template: %s", resp.Status)
	}

	var template ExerciseTemplate
	err = json.NewDecoder(resp.Body).Decode(&template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// CreateExerciseTemplate creates a new custom exercise template
func (c *Client) CreateExerciseTemplate(template ExerciseTemplate) (*ExerciseTemplate, error) {
	url := fmt.Sprintf("%s/exercise_templates", c.baseURL)
	
	body, err := json.Marshal(template)
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
		return nil, fmt.Errorf("failed to create exercise template: %s", resp.Status)
	}

	var result ExerciseTemplate
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


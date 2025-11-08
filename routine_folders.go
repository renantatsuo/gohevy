package hevy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRoutineFolders retrieves a paginated list of routine folders
func (c *Client) GetRoutineFolders(params PaginationParams) (*PaginatedRoutineFoldersResponse, error) {
	url := fmt.Sprintf("%s/routine_folders?page=%d&pageSize=%d", c.baseURL, params.Page, params.PageSize)
	
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
		return nil, fmt.Errorf("failed to get routine folders: %s", resp.Status)
	}

	var result PaginatedRoutineFoldersResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRoutineFolder retrieves a single routine folder by ID
func (c *Client) GetRoutineFolder(folderID int) (*RoutineFolder, error) {
	url := fmt.Sprintf("%s/routine_folders/%d", c.baseURL, folderID)
	
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
		return nil, fmt.Errorf("failed to get routine folder: %s", resp.Status)
	}

	var folder RoutineFolder
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

// CreateRoutineFolder creates a new routine folder
func (c *Client) CreateRoutineFolder(folder RoutineFolder) (*RoutineFolder, error) {
	url := fmt.Sprintf("%s/routine_folders", c.baseURL)
	
	body, err := json.Marshal(folder)
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
		return nil, fmt.Errorf("failed to create routine folder: %s", resp.Status)
	}

	var result RoutineFolder
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


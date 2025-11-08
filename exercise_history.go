package hevy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetExerciseHistory retrieves exercise history for a specific exercise template
func (c *Client) GetExerciseHistory(exerciseTemplateID string) ([]ExerciseHistoryEntry, error) {
	url := fmt.Sprintf("%s/exercise_history/%s", c.baseURL, exerciseTemplateID)
	
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
		return nil, fmt.Errorf("failed to get exercise history: %s", resp.Status)
	}

	var history []ExerciseHistoryEntry
	err = json.NewDecoder(resp.Body).Decode(&history)
	if err != nil {
		return nil, err
	}

	return history, nil
}


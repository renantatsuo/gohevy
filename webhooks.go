package hevy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateWebhookSubscription creates a new webhook subscription
func (c *Client) CreateWebhookSubscription(subscription WebhookSubscription) (*WebhookSubscription, error) {
	url := fmt.Sprintf("%s/webhook-subscription", c.baseURL)
	
	body, err := json.Marshal(subscription)
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
		return nil, fmt.Errorf("failed to create webhook subscription: %s", resp.Status)
	}

	var result WebhookSubscription
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWebhookSubscription retrieves the current webhook subscription
func (c *Client) GetWebhookSubscription() (*WebhookSubscription, error) {
	url := fmt.Sprintf("%s/webhook-subscription", c.baseURL)
	
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
		return nil, fmt.Errorf("failed to get webhook subscription: %s", resp.Status)
	}

	var result WebhookSubscription
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteWebhookSubscription deletes the webhook subscription
func (c *Client) DeleteWebhookSubscription() error {
	url := fmt.Sprintf("%s/webhook-subscription", c.baseURL)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("api-key", c.apiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete webhook subscription: %s", resp.Status)
	}

	return nil
}


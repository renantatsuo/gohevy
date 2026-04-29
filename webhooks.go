package hevy

import (
	"context"
	"net/http"
)

// CreateWebhookSubscription registers a new webhook endpoint to receive Hevy events.
// Only one subscription is allowed per API key; creating a second will replace the existing one.
// Set subscription.URL to an HTTPS endpoint and subscription.Events to the event types to receive
// (e.g. "workout.created", "workout.updated", "workout.deleted").
// Returns the server-assigned subscription record including its ID.
func (c *Client) CreateWebhookSubscription(ctx context.Context, subscription WebhookSubscription) (res *WebhookSubscription, err error) {
	err = c.request(ctx, http.MethodPost, "/webhook-subscription", subscription, &res)
	return
}

// GetWebhookSubscription retrieves the active webhook subscription for the authenticated API key.
// Returns *APIError with StatusCode 404 if no subscription exists.
func (c *Client) GetWebhookSubscription(ctx context.Context) (res *WebhookSubscription, err error) {
	err = c.request(ctx, http.MethodGet, "/webhook-subscription", nil, &res)
	return
}

// DeleteWebhookSubscription removes the active webhook subscription for the authenticated API key.
// Returns nil on success. Returns *APIError with StatusCode 404 if no subscription exists.
func (c *Client) DeleteWebhookSubscription(ctx context.Context) error {
	return c.request(ctx, http.MethodDelete, "/webhook-subscription", nil, nil)
}

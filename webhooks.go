package hevy

import (
	"context"
	"net/http"
)

// CreateWebhookSubscription creates a new webhook subscription
func (c *Client) CreateWebhookSubscription(ctx context.Context, subscription WebhookSubscription) (res *WebhookSubscription, err error) {
	err = c.request(ctx, http.MethodPost, "/webhook-subscription", subscription, &res)
	return
}

// GetWebhookSubscription retrieves the current webhook subscription
func (c *Client) GetWebhookSubscription(ctx context.Context) (res *WebhookSubscription, err error) {
	err = c.request(ctx, http.MethodGet, "/webhook-subscription", nil, &res)
	return
}

// DeleteWebhookSubscription deletes the webhook subscription
func (c *Client) DeleteWebhookSubscription(ctx context.Context) error {
	return c.request(ctx, http.MethodDelete, "/webhook-subscription", nil, nil)
}

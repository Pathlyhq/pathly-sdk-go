package pathly

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
)

func newIdempotencyKey() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// ---------------------------------------------------------------------------
// Scenarios (HTTP monitoring)
// ---------------------------------------------------------------------------

// CreateScenario creates an HTTP scenario. An Idempotency-Key is always sent.
func (c *Client) CreateScenario(ctx context.Context, in ScenarioInput, idempotencyKey string) (*Scenario, error) {
	if idempotencyKey == "" {
		idempotencyKey = newIdempotencyKey()
	}
	out := &Scenario{}
	err := c.do(ctx, request{
		method:         http.MethodPost,
		path:           "/v1/scenarios",
		body:           in,
		idempotencyKey: idempotencyKey,
		out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetScenario reads one scenario by id.
func (c *Client) GetScenario(ctx context.Context, id string) (*Scenario, error) {
	out := &Scenario{}
	err := c.do(ctx, request{
		method: http.MethodGet,
		path:   "/v1/scenarios/" + url.PathEscape(id),
		out:    out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateScenario patches an existing scenario.
func (c *Client) UpdateScenario(ctx context.Context, id string, in ScenarioInput) (*Scenario, error) {
	out := &Scenario{}
	err := c.do(ctx, request{
		method: http.MethodPatch,
		path:   "/v1/scenarios/" + url.PathEscape(id),
		body:   in,
		out:    out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteScenario removes a scenario.
func (c *Client) DeleteScenario(ctx context.Context, id string) error {
	return c.do(ctx, request{
		method: http.MethodDelete,
		path:   "/v1/scenarios/" + url.PathEscape(id),
	})
}

// ListScenarios walks every page of HTTP scenarios.
func (c *Client) ListScenarios(ctx context.Context) ([]Scenario, error) {
	return listPaged[Scenario](ctx, c, "/v1/scenarios")
}

// MuteScenario mutes or unmutes. A nil until wakes the scenario up.
func (c *Client) MuteScenario(ctx context.Context, id string, until *string) error {
	body := map[string]any{"mutedUntil": until}
	return c.do(ctx, request{
		method: http.MethodPost,
		path:   "/v1/scenarios/" + url.PathEscape(id) + "/mute",
		body:   body,
	})
}

// ---------------------------------------------------------------------------
// Maintenance windows
// ---------------------------------------------------------------------------

// CreateMaintenanceWindow creates a maintenance window.
func (c *Client) CreateMaintenanceWindow(ctx context.Context, in MaintenanceWindowInput, idempotencyKey string) (*MaintenanceWindow, error) {
	if idempotencyKey == "" {
		idempotencyKey = newIdempotencyKey()
	}
	out := &MaintenanceWindow{}
	err := c.do(ctx, request{
		method:         http.MethodPost,
		path:           "/v1/maintenance-windows",
		body:           in,
		idempotencyKey: idempotencyKey,
		out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetMaintenanceWindow reads one window by id.
func (c *Client) GetMaintenanceWindow(ctx context.Context, id string) (*MaintenanceWindow, error) {
	out := &MaintenanceWindow{}
	err := c.do(ctx, request{
		method: http.MethodGet,
		path:   "/v1/maintenance-windows/" + url.PathEscape(id),
		out:    out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteMaintenanceWindow removes a window.
func (c *Client) DeleteMaintenanceWindow(ctx context.Context, id string) error {
	return c.do(ctx, request{
		method: http.MethodDelete,
		path:   "/v1/maintenance-windows/" + url.PathEscape(id),
	})
}

// ListMaintenanceWindows lists every page of maintenance windows.
func (c *Client) ListMaintenanceWindows(ctx context.Context) ([]MaintenanceWindow, error) {
	return listPaged[MaintenanceWindow](ctx, c, "/v1/maintenance-windows")
}

// ---------------------------------------------------------------------------
// Outbound webhooks
// ---------------------------------------------------------------------------

// CreateWebhook creates a signed outbound webhook. Store Secret immediately.
func (c *Client) CreateWebhook(ctx context.Context, in WebhookInput, idempotencyKey string) (*Webhook, error) {
	if idempotencyKey == "" {
		idempotencyKey = newIdempotencyKey()
	}
	out := &Webhook{}
	err := c.do(ctx, request{
		method:         http.MethodPost,
		path:           "/v1/webhooks",
		body:           in,
		idempotencyKey: idempotencyKey,
		out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetWebhook reads one webhook by id (URL never returned, only fingerprint).
func (c *Client) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	out := &Webhook{}
	err := c.do(ctx, request{
		method: http.MethodGet,
		path:   "/v1/webhooks/" + url.PathEscape(id),
		out:    out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWebhook removes a webhook.
func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	return c.do(ctx, request{
		method: http.MethodDelete,
		path:   "/v1/webhooks/" + url.PathEscape(id),
	})
}

// ListWebhooks lists every page of webhooks.
func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	return listPaged[Webhook](ctx, c, "/v1/webhooks")
}

// ---------------------------------------------------------------------------
// SLA targets
// ---------------------------------------------------------------------------

// UpsertSlaTarget creates or replaces an SLA target (PUT).
func (c *Client) UpsertSlaTarget(ctx context.Context, in SlaTargetInput, idempotencyKey string) (*SlaTarget, error) {
	if idempotencyKey == "" {
		idempotencyKey = newIdempotencyKey()
	}
	out := &SlaTarget{}
	err := c.do(ctx, request{
		method:         http.MethodPut,
		path:           "/v1/sla-targets",
		body:           in,
		idempotencyKey: idempotencyKey,
		out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetSlaTarget reads one SLA target by id.
func (c *Client) GetSlaTarget(ctx context.Context, id string) (*SlaTarget, error) {
	out := &SlaTarget{}
	err := c.do(ctx, request{
		method: http.MethodGet,
		path:   "/v1/sla-targets/" + url.PathEscape(id),
		out:    out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteSlaTarget removes an SLA target.
func (c *Client) DeleteSlaTarget(ctx context.Context, id string) error {
	return c.do(ctx, request{
		method: http.MethodDelete,
		path:   "/v1/sla-targets/" + url.PathEscape(id),
	})
}

// ListSlaTargets lists every page of SLA targets.
func (c *Client) ListSlaTargets(ctx context.Context) ([]SlaTarget, error) {
	return listPaged[SlaTarget](ctx, c, "/v1/sla-targets")
}

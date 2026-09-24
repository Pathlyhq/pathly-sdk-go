package pathly

// Scenario mirrors the public projection of an HTTP monitoring scenario.
type Scenario struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Type                string   `json:"type"`
	URL                 *string  `json:"url"`
	Enabled             *bool    `json:"enabled"`
	IntervalSec         *int64   `json:"intervalSec"`
	Method              *string  `json:"method"`
	ExpectedStatus      *int64   `json:"expectedStatus"`
	MaxLatencyMs        *int64   `json:"maxLatencyMs"`
	ExpectText          *string  `json:"expectText"`
	Runbook             *string  `json:"runbook"`
	Cron                *string  `json:"cron"`
	LastStatus          *string  `json:"lastStatus"`
	Regions             []string `json:"regions"`
	Tags                []string `json:"tags"`
	Folder              *string  `json:"folder"`
	Severity            *string  `json:"severity"`
	MutedUntil          *string  `json:"mutedUntil"`
	ScenarioFingerprint *string  `json:"scenarioFingerprint"`
	CreatedAt           *string  `json:"createdAt"`
}

// ScenarioInput serves both create and update. Absent pointer fields are not
// sent so the API keeps the existing value.
type ScenarioInput struct {
	Name           *string `json:"name,omitempty"`
	Type           *string `json:"type,omitempty"`
	URL            *string `json:"url,omitempty"`
	IntervalSec    *int64  `json:"intervalSec,omitempty"`
	Method         *string `json:"method,omitempty"`
	ExpectedStatus *int64  `json:"expectedStatus,omitempty"`
	MaxLatencyMs   *int64  `json:"maxLatencyMs,omitempty"`
	ExpectText     *string `json:"expectText,omitempty"`
	Regions        []string `json:"regions,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Folder         *string `json:"folder,omitempty"`
	Severity       *string `json:"severity,omitempty"`
	Runbook        *string `json:"runbook,omitempty"`
	Cron           *string `json:"cron,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
}

// MaintenanceWindow is a planned silence window for a scenario.
type MaintenanceWindow struct {
	ID          string  `json:"id"`
	MonitorID   *string `json:"monitorId"`
	StartsAt    *string `json:"startsAt"`
	EndsAt      *string `json:"endsAt"`
	Reason      *string `json:"reason"`
	Weekday     *int64  `json:"weekday"`
	StartMinute *int64  `json:"startMinute"`
	DurationMin *int64  `json:"durationMin"`
}

// MaintenanceWindowInput is the write payload for a maintenance window.
type MaintenanceWindowInput struct {
	MonitorID   *string `json:"monitorId,omitempty"`
	StartsAt    *string `json:"startsAt,omitempty"`
	EndsAt      *string `json:"endsAt,omitempty"`
	Reason      *string `json:"reason,omitempty"`
	Weekday     *int64  `json:"weekday,omitempty"`
	StartMinute *int64  `json:"startMinute,omitempty"`
	DurationMin *int64  `json:"durationMin,omitempty"`
}

// Webhook is an outbound signed webhook. The destination URL is never returned
// on read — only urlFingerprint. Secret is filled only on create.
type Webhook struct {
	ID             string   `json:"id"`
	Events         []string `json:"events"`
	Enabled        *bool    `json:"enabled"`
	HasSecret      *bool    `json:"hasSecret"`
	URLFingerprint *string  `json:"urlFingerprint"`
	CreatedAt      *string  `json:"createdAt"`
	Secret         *string  `json:"secret"`
}

// WebhookInput creates an outbound webhook.
type WebhookInput struct {
	URL    string   `json:"url"`
	Events []string `json:"events,omitempty"`
}

// SlaTarget is an availability objective bound to a scenario.
type SlaTarget struct {
	ID                 string   `json:"id"`
	MonitorID          *string  `json:"monitorId"`
	Name               *string  `json:"name"`
	ObjectivePct       *float64 `json:"objectivePct"`
	WindowDays         *int64   `json:"windowDays"`
	ExcludeMaintenance *bool    `json:"excludeMaintenance"`
	WarnAtBudgetRatio  *float64 `json:"warnAtBudgetRatio"`
	Enabled            *bool    `json:"enabled"`
}

// SlaTargetInput creates or replaces an SLA target (PUT upsert).
type SlaTargetInput struct {
	MonitorID          *string  `json:"monitorId,omitempty"`
	Name               *string  `json:"name,omitempty"`
	ObjectivePct       float64  `json:"objectivePct"`
	WindowDays         int64    `json:"windowDays"`
	ExcludeMaintenance *bool    `json:"excludeMaintenance,omitempty"`
	WarnAtBudgetRatio  *float64 `json:"warnAtBudgetRatio,omitempty"`
	Enabled            *bool    `json:"enabled,omitempty"`
}

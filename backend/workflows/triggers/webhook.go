package triggers

import (
	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterTrigger("WEBHOOK", &WebhookTrigger{})
}

// WebhookTrigger handles webhook-based workflow triggers
type WebhookTrigger struct {
	URL    string `json:"url" validate:"required,url" desc:"Webhook URL that will receive HTTP POST requests when the workflow is triggered. Must be a valid HTTP/HTTPS URL."`
	Secret string `json:"secret,omitempty" validate:"omitempty,min=8" desc:"Optional: Secret key for webhook signature verification. If provided, must be at least 8 characters long. Used to validate incoming webhook requests."`
}

func (t *WebhookTrigger) Type() string {
	return "WEBHOOK"
}

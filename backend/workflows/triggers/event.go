package triggers

import (
	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterTrigger("EVENT", &EventTrigger{})
}

// EventTrigger handles event-based workflow triggers
type EventTrigger struct {
	TriggerEvent string `json:"trigger_event" validate:"required" desc:"Name of the event that will trigger this workflow (e.g., 'user_signup', 'purchase_completed')."`
	SiteIDs      []int  `json:"site_ids,omitempty" validate:"omitempty,min=1" desc:"Optional: List of site IDs to filter events. If provided, only events from these sites will trigger the workflow."`
}

func (t *EventTrigger) Type() string {
	return "EVENT"
}

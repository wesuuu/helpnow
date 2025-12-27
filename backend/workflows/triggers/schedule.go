package triggers

import (
	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterTrigger("SCHEDULE", &ScheduleTrigger{})
}

// ScheduleTrigger handles scheduled/cron-based workflow triggers
type ScheduleTrigger struct {
	Cron string `json:"cron" validate:"required" desc:"Cron expression defining when the workflow should run (e.g., '0 9 * * *' for daily at 9am). Uses standard cron format: minute hour day month weekday."`
}

func (t *ScheduleTrigger) Type() string {
	return "SCHEDULE"
}

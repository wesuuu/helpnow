package actions

import (
	"context"
	"fmt"

	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterAction("Delay", &DelayAction{})
}

// DelayAction pauses workflow execution
type DelayAction struct {
	DelayMinutes int `json:"delay_minutes" validate:"required,min=1" desc:"Number of minutes to delay execution."`
}

func (a *DelayAction) Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error) {
	// Stub implementation
	return fmt.Sprintf("Simulated delay of %d minutes", a.DelayMinutes), nil
}

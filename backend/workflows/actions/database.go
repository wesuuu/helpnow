package actions

import (
	"context"
	"fmt"

	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterAction("Update DB", &UpdateDBAction{})
}

// UpdateDBAction updates a database record
type UpdateDBAction struct {
	Table    string `json:"table" validate:"required,oneof=users sites events" desc:"Database table to update."`
	RecordID string `json:"record_id" validate:"required" desc:"ID of the record to update (can use {{variables}})."`
	Data     string `json:"data" validate:"required" desc:"JSON string of data to update."`
}

func (a *UpdateDBAction) Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error) {
	// Stub implementation
	return fmt.Sprintf("Simulated update to %s:%s", a.Table, a.RecordID), nil
}

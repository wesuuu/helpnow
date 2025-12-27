package logic

import (
	"context"
	"math/rand"

	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterLogic("Condition", &ConditionLogic{})
}

// ConditionLogic evaluates a condition based on context data
type ConditionLogic struct {
	Force string `json:"force,omitempty" validate:"omitempty,oneof=true false" desc:"Optional: Set to 'true' or 'false' to force a specific result. If not provided, the condition will evaluate randomly (50/50)."`
}

func (l *ConditionLogic) Evaluate(ctx context.Context, contextData map[string]interface{}) (bool, string, error) {
	result := false

	// If Force is set on the struct, use that
	if l.Force == "true" {
		result = true
	} else if l.Force == "false" {
		result = false
	} else {
		// Otherwise, random 50/50
		result = rand.Float32() > 0.5
	}

	output := "Condition: False"
	if result {
		output = "Condition: True"
	}

	return result, output, nil
}

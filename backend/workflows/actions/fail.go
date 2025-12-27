package actions

import (
	"context"
	"errors"

	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterAction("FAIL", &FailAction{})
}

type FailAction struct{}

func (a *FailAction) Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error) {
	return "Manual Fail", errors.New("workflow manually failed by FAIL action")
}

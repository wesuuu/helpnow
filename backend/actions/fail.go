package actions

import (
	"context"
	"errors"
)

type FailAction struct{}

func (a *FailAction) Execute(ctx context.Context, actionCtx ActionContext) (output string, err error) {
	return "Simulated Failure", errors.New("simulated failure")
}

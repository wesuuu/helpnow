package actions

import (
	"context"
	"sync"
)

// ActionContext holds data available to an action during execution
type ActionContext struct {
	ExecutionID int
	NodeData    map[string]interface{}
	ContextData map[string]interface{} // Event context (e.g. email, person_id)
}

// Action interface for all workflow actions
type Action interface {
	Execute(ctx context.Context, actionCtx ActionContext) (output string, err error)
}

var (
	registry = make(map[string]Action)
	mu       sync.RWMutex
)

// Register adds an action to the registry
func Register(name string, action Action) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = action
}

// Get retrieves an action from the registry
func Get(name string) (Action, bool) {
	mu.RLock()
	defer mu.RUnlock()
	action, ok := registry[name]
	return action, ok
}

// RegisterStandardActions registers the default set of actions
func RegisterStandardActions() {
	Register("Send Email", &SendEmailAction{})
	Register("FAIL", &FailAction{})
}

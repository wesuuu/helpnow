package actions

import (
	"context"
	"fmt"

	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterAction("HTTP Request", &HTTPRequestAction{})
}

// HTTPRequestAction sends an HTTP request
type HTTPRequestAction struct {
	Method  string `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH" desc:"HTTP Method to use for the request."`
	URL     string `json:"url" validate:"required,url" desc:"Target URL for the request."`
	Headers string `json:"headers,omitempty" desc:"Optional JSON string of request headers."`
	Body    string `json:"body,omitempty" desc:"Optional request body payload."`
}

func (a *HTTPRequestAction) Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error) {
	// Stub implementation for now
	return fmt.Sprintf("Simulated %s to %s", a.Method, a.URL), nil
}

package workflows

import (
	"context"
	"sync"
	"time"
)

// --- Definitions (Moved from handlers) ---

type EventDefinition struct {
	ID          int       `json:"id"`
	SiteID      int       `json:"site_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// WorkflowTrigger represents a workflow trigger (event, schedule, webhook, etc.)
// Note: This is different from the Trigger interface which defines trigger behavior
type WorkflowTrigger struct {
	ID         int                    `json:"id"`
	WorkflowID int                    `json:"workflow_id"`
	NodeID     string                 `json:"node_id"`               // Links to graph node
	Type       string                 `json:"type"`                  // EVENT, SCHEDULE, WEBHOOK
	Config     map[string]interface{} `json:"config"`                // Trigger-specific configuration
	NextRunAt  *time.Time             `json:"next_run_at,omitempty"` // For scheduled triggers
	CreatedAt  time.Time              `json:"created_at"`
}

type Workflow struct {
	ID             int       `json:"id"`
	OrganizationID *int      `json:"organization_id"`
	SiteID         *int      `json:"site_id"`
	SiteName       string    `json:"site_name,omitempty"`
	AudienceID     *int      `json:"audience_id"`
	Name           string    `json:"name"`
	Steps          string    `json:"steps"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`

	// Multiple triggers support
	Triggers []WorkflowTrigger `json:"triggers,omitempty"`
}

type WorkflowExecution struct {
	ID          int        `json:"id"`
	WorkflowID  int        `json:"workflow_id"`
	SubjectID   int        `json:"subject_id"`
	CurrentStep int        `json:"current_step"`
	Status      string     `json:"status"`
	NextRunAt   time.Time  `json:"next_run_at"`
	Context     string     `json:"context"` // JSON string
	CreatedAt   time.Time  `json:"created_at"`
	FinishedAt  *time.Time `json:"finished_at"`
}

// --- Graph Models (Moved from models.go) ---

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"` // "TRIGGER", "ACTION", "CONDITION"
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
	Position   NodePosition           `json:"position"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Handle string `json:"handle"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// --- Interfaces & Registry (Moved from types.go) ---

// ActionContext is deprecated - kept for backward compatibility
// New pattern: actions have properties as struct fields and receive only runtime context
type ActionContext struct {
	ExecutionID int
	NodeData    map[string]interface{}
	ContextData map[string]interface{}
}

// Action interface - properties should be struct fields on the implementing type
type Action interface {
	Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error)
}

// Logic interface - properties should be struct fields on the implementing type
type Logic interface {
	Evaluate(ctx context.Context, contextData map[string]interface{}) (result bool, output string, err error)
}

type Trigger interface {
	Type() string
}

var (
	actionRegistry  = make(map[string]Action)
	logicRegistry   = make(map[string]Logic)
	triggerRegistry = make(map[string]Trigger)
	mu              sync.RWMutex
)

func RegisterAction(name string, action Action) {
	mu.Lock()
	defer mu.Unlock()
	actionRegistry[name] = action
}

func GetAction(name string) (Action, bool) {
	mu.RLock()
	defer mu.RUnlock()
	action, ok := actionRegistry[name]
	return action, ok
}

func RegisterLogic(name string, logic Logic) {
	mu.Lock()
	defer mu.Unlock()
	logicRegistry[name] = logic
}

func GetLogic(name string) (Logic, bool) {
	mu.RLock()
	defer mu.RUnlock()
	logic, ok := logicRegistry[name]
	return logic, ok
}

func RegisterTrigger(name string, trigger Trigger) {
	mu.Lock()
	defer mu.Unlock()
	triggerRegistry[name] = trigger
}

func GetTrigger(name string) (Trigger, bool) {
	mu.RLock()
	defer mu.RUnlock()
	trigger, ok := triggerRegistry[name]
	return trigger, ok
}

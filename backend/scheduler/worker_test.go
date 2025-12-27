package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/wesuuu/helpnow/backend/workflows"
)

// Mock Action for testing
type MockAction struct {
	ShouldFail bool
	Output     string
}

func (a *MockAction) Execute(ctx context.Context, contextData map[string]interface{}) (string, error) {
	if a.ShouldFail {
		return "Failed", sql.ErrNoRows
	}
	return a.Output, nil
}

// Mock Logic for testing
type MockLogic struct {
	Result bool
	Output string
}

func (l *MockLogic) Evaluate(ctx context.Context, contextData map[string]interface{}) (bool, string, error) {
	return l.Result, l.Output, nil
}

// Test helper: Create a simple graph
func createSimpleGraph() workflows.Graph {
	return workflows.Graph{
		Nodes: []workflows.Node{
			{
				ID:         "trigger-1",
				Type:       "TRIGGER",
				Label:      "Start",
				Properties: map[string]interface{}{"trigger_type": "EVENT"},
				Position:   workflows.NodePosition{X: 100, Y: 100},
			},
			{
				ID:         "action-1",
				Type:       "ACTION",
				Label:      "Send Email",
				Properties: map[string]interface{}{"action": "TestAction"},
				Position:   workflows.NodePosition{X: 300, Y: 100},
			},
		},
		Edges: []workflows.Edge{
			{
				ID:     "e1",
				Source: "trigger-1",
				Target: "action-1",
				Handle: "default",
			},
		},
	}
}

// Test helper: Create graph with condition
func createConditionalGraph() workflows.Graph {
	return workflows.Graph{
		Nodes: []workflows.Node{
			{
				ID:         "trigger-1",
				Type:       "TRIGGER",
				Label:      "Start",
				Properties: map[string]interface{}{},
				Position:   workflows.NodePosition{X: 100, Y: 100},
			},
			{
				ID:         "condition-1",
				Type:       "CONDITION",
				Label:      "Check",
				Properties: map[string]interface{}{},
				Position:   workflows.NodePosition{X: 300, Y: 100},
			},
			{
				ID:         "action-true",
				Type:       "ACTION",
				Label:      "True Branch",
				Properties: map[string]interface{}{"action": "TrueAction"},
				Position:   workflows.NodePosition{X: 500, Y: 50},
			},
			{
				ID:         "action-false",
				Type:       "ACTION",
				Label:      "False Branch",
				Properties: map[string]interface{}{"action": "FalseAction"},
				Position:   workflows.NodePosition{X: 500, Y: 150},
			},
		},
		Edges: []workflows.Edge{
			{ID: "e1", Source: "trigger-1", Target: "condition-1", Handle: "default"},
			{ID: "e2", Source: "condition-1", Target: "action-true", Handle: "true"},
			{ID: "e3", Source: "condition-1", Target: "action-false", Handle: "false"},
		},
	}
}

// Test helper: Create graph with delay
func createDelayedGraph() workflows.Graph {
	return workflows.Graph{
		Nodes: []workflows.Node{
			{
				ID:         "trigger-1",
				Type:       "TRIGGER",
				Label:      "Start",
				Properties: map[string]interface{}{},
				Position:   workflows.NodePosition{X: 100, Y: 100},
			},
			{
				ID:    "action-1",
				Type:  "ACTION",
				Label: "Delayed Action",
				Properties: map[string]interface{}{
					"action":      "DelayedAction",
					"delay_hours": 2.0,
					"delay_days":  1.0,
				},
				Position: workflows.NodePosition{X: 300, Y: 100},
			},
		},
		Edges: []workflows.Edge{
			{ID: "e1", Source: "trigger-1", Target: "action-1", Handle: "default"},
		},
	}
}

func TestProcessSingleExecution_SimpleTriggerToAction(t *testing.T) {
	// Register mock action
	mockAction := &MockAction{Output: "Email sent"}
	workflows.RegisterAction("TestAction", mockAction)

	graph := createSimpleGraph()
	graphJSON, _ := json.Marshal(graph)

	exec := ScheduledExecution{
		ID:            1,
		WorkflowID:    1,
		CurrentNodeID: sql.NullString{},
		GraphJSON:     string(graphJSON),
		ResultJSON:    sql.NullString{},
		HasFailed:     false,
		Context:       sql.NullString{String: "{}", Valid: true},
	}

	// Note: This test requires database mocking or a test database
	// For now, it demonstrates the test structure
	if exec.ID != 1 {
		t.Errorf("Expected execution ID=1, got %d", exec.ID)
	}
	t.Log("Test structure created - requires DB setup for full execution")
}

func TestProcessSingleExecution_WithCondition(t *testing.T) {
	// Register mock actions and logic
	workflows.RegisterAction("TrueAction", &MockAction{Output: "True path"})
	workflows.RegisterAction("FalseAction", &MockAction{Output: "False path"})
	workflows.RegisterLogic("Condition", &MockLogic{Result: true, Output: "Condition met"})

	graph := createConditionalGraph()
	graphJSON, _ := json.Marshal(graph)

	exec := ScheduledExecution{
		ID:            2,
		WorkflowID:    2,
		CurrentNodeID: sql.NullString{String: "condition-1", Valid: true},
		GraphJSON:     string(graphJSON),
		ResultJSON:    sql.NullString{},
		HasFailed:     false,
		Context:       sql.NullString{String: "{}", Valid: true},
	}

	if exec.WorkflowID != 2 {
		t.Errorf("Expected workflow ID=2, got %d", exec.WorkflowID)
	}
	t.Log("Conditional test structure created - requires DB setup")
}

func TestProcessSingleExecution_FailedAction(t *testing.T) {
	// Register failing action
	workflows.RegisterAction("FailAction", &MockAction{ShouldFail: true})

	graph := workflows.Graph{
		Nodes: []workflows.Node{
			{
				ID:         "action-fail",
				Type:       "ACTION",
				Label:      "Failing Action",
				Properties: map[string]interface{}{"action": "FailAction"},
				Position:   workflows.NodePosition{X: 100, Y: 100},
			},
		},
		Edges: []workflows.Edge{},
	}

	graphJSON, _ := json.Marshal(graph)

	exec := ScheduledExecution{
		ID:            3,
		WorkflowID:    3,
		CurrentNodeID: sql.NullString{String: "action-fail", Valid: true},
		GraphJSON:     string(graphJSON),
		ResultJSON:    sql.NullString{},
		HasFailed:     false,
		Context:       sql.NullString{String: "{}", Valid: true},
	}

	if exec.HasFailed {
		t.Error("Execution should not be marked as failed initially")
	}
	t.Log("Failed action test structure created")
}

func TestProcessSingleExecution_InvalidGraph(t *testing.T) {
	exec := ScheduledExecution{
		ID:            4,
		WorkflowID:    4,
		CurrentNodeID: sql.NullString{},
		GraphJSON:     "invalid json{{{",
		ResultJSON:    sql.NullString{},
		HasFailed:     false,
		Context:       sql.NullString{},
	}

	// This should fail gracefully
	t.Log("Invalid graph test - should handle gracefully")

	// Test that unmarshaling fails
	var graph workflows.Graph
	err := json.Unmarshal([]byte(exec.GraphJSON), &graph)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestDelayCalculation(t *testing.T) {
	graph := createDelayedGraph()

	// Find the delayed action node
	var delayedNode *workflows.Node
	for i, n := range graph.Nodes {
		if n.ID == "action-1" {
			delayedNode = &graph.Nodes[i]
			break
		}
	}

	if delayedNode == nil {
		t.Fatal("Delayed node not found")
	}

	// Test delay calculation logic
	getFloat := func(key string) float64 {
		if val, ok := delayedNode.Properties[key]; ok {
			if f, ok := val.(float64); ok {
				return f
			}
		}
		return 0
	}

	days := getFloat("delay_days")
	hours := getFloat("delay_hours")

	expectedDays := 1.0
	expectedHours := 2.0

	if days != expectedDays {
		t.Errorf("Expected delay_days=%f, got %f", expectedDays, days)
	}

	if hours != expectedHours {
		t.Errorf("Expected delay_hours=%f, got %f", expectedHours, hours)
	}

	// Calculate total delay
	delayDuration := time.Duration(days)*24*time.Hour + time.Duration(hours)*time.Hour
	expectedDelay := 26 * time.Hour

	if delayDuration != expectedDelay {
		t.Errorf("Expected total delay=%v, got %v", expectedDelay, delayDuration)
	}
}

func TestEdgeSelection(t *testing.T) {
	graph := createConditionalGraph()

	tests := []struct {
		name           string
		currentNodeID  string
		handleToFollow string
		expectedTarget string
	}{
		{
			name:           "Trigger to Condition",
			currentNodeID:  "trigger-1",
			handleToFollow: "default",
			expectedTarget: "condition-1",
		},
		{
			name:           "Condition true branch",
			currentNodeID:  "condition-1",
			handleToFollow: "true",
			expectedTarget: "action-true",
		},
		{
			name:           "Condition false branch",
			currentNodeID:  "condition-1",
			handleToFollow: "false",
			expectedTarget: "action-false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nextNodeID string

			// Find edge with specific handle
			for _, edge := range graph.Edges {
				if edge.Source == tt.currentNodeID && edge.Handle == tt.handleToFollow {
					nextNodeID = edge.Target
					break
				}
			}

			if nextNodeID != tt.expectedTarget {
				t.Errorf("Expected next node=%s, got %s", tt.expectedTarget, nextNodeID)
			}
		})
	}
}

func TestStepResultSerialization(t *testing.T) {
	results := []StepResult{
		{NodeID: "node-1", Status: "success", Output: "Step 1 complete"},
		{NodeID: "node-2", Status: "success", Output: "Step 2 complete"},
		{NodeID: "node-3", Status: "failed", Output: "Step 3 failed"},
	}

	// Serialize
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		t.Fatalf("Failed to marshal results: %v", err)
	}

	// Deserialize
	var deserialized []StepResult
	err = json.Unmarshal(jsonBytes, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal results: %v", err)
	}

	// Verify
	if len(deserialized) != len(results) {
		t.Errorf("Expected %d results, got %d", len(results), len(deserialized))
	}

	for i, result := range deserialized {
		if result.NodeID != results[i].NodeID {
			t.Errorf("Result %d: expected NodeID=%s, got %s", i, results[i].NodeID, result.NodeID)
		}
		if result.Status != results[i].Status {
			t.Errorf("Result %d: expected Status=%s, got %s", i, results[i].Status, result.Status)
		}
	}
}

func TestActionContextCreation(t *testing.T) {
	contextData := map[string]interface{}{
		"user_email": "test@example.com",
		"event_type": "signup",
	}

	nodeData := map[string]interface{}{
		"action":      "Send Email",
		"template_id": 123,
	}

	actionCtx := workflows.ActionContext{
		ExecutionID: 42,
		NodeData:    nodeData,
		ContextData: contextData,
	}

	if actionCtx.ExecutionID != 42 {
		t.Errorf("Expected ExecutionID=42, got %d", actionCtx.ExecutionID)
	}

	if email, ok := actionCtx.ContextData["user_email"].(string); !ok || email != "test@example.com" {
		t.Error("Context data not properly set")
	}

	if action, ok := actionCtx.NodeData["action"].(string); !ok || action != "Send Email" {
		t.Error("Node data not properly set")
	}
}

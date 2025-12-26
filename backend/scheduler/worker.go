package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/wesuuu/helpnow/backend/actions"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

// --- Graph Data Structures ---

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
	Handle string `json:"handle"` // "default", "true", "false"
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type ScheduledExecution struct {
	ID            int
	WorkflowID    int
	CurrentNodeID sql.NullString // Replaces CurrentStep
	GraphJSON     string         // Replaces StepsJSON
	ResultJSON    sql.NullString
	HasFailed     bool
	Context       sql.NullString
}

type StepResult struct {
	NodeID string `json:"node_id"`
	Status string `json:"status"` // "success", "failed", "skipped"
	Output string `json:"output"`
}

func StartWorker() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			processPendingExecutions()
		}
	}()
}

func processPendingExecutions() {
	// Query for executions that are PENDING and due
	// Note: We're reading `steps` into `GraphJSON`.
	rows, err := db.GetDB().Query(`
		SELECT we.id, we.workflow_id, we.current_node_id, w.steps, we.step_results, we.has_failed, we.context 
		FROM workflow_executions we
		JOIN workflows w ON we.workflow_id = w.id
		WHERE we.status = 'PENDING' AND we.next_run_at <= NOW()
	`)
	if err != nil {
		Logger.Error("Scheduler error querying executions:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exec ScheduledExecution
		if err := rows.Scan(&exec.ID, &exec.WorkflowID, &exec.CurrentNodeID, &exec.GraphJSON, &exec.ResultJSON, &exec.HasFailed, &exec.Context); err != nil {
			Logger.Error("Scheduler scan error:", err)
			continue
		}
		processSingleExecution(exec)
	}
}

func processSingleExecution(exec ScheduledExecution) {
	var graph Graph
	if err := json.Unmarshal([]byte(exec.GraphJSON), &graph); err != nil {
		Logger.Errorf("[Worker] Failed to unmarshal graph for execution %d: %v", exec.ID, err)
		// Fallback for legacy array-based steps?
		// For now, fail if not graph. In a real migration we'd check `[` vs `{`.
		markExecutionFinal(exec.ID, "FAILED", "Invalid graph JSON")
		return
	}

	currentNodeID := ""
	if exec.CurrentNodeID.Valid {
		currentNodeID = exec.CurrentNodeID.String
	} else {
		// New execution: Find Start Node (Trigger)
		// Assumption: The Trigger node has no incoming edges, or is explicitly Type="TRIGGER"
		for _, n := range graph.Nodes {
			if n.Type == string(models.NodeTypeTrigger) {
				currentNodeID = n.ID
				break
			}
		}
		if currentNodeID == "" && len(graph.Nodes) > 0 {
			// Fallback: take first node
			currentNodeID = graph.Nodes[0].ID
		}
	}

	if currentNodeID == "" {
		Logger.Errorf("[Worker] No start node found for execution %d", exec.ID)
		markExecutionFinal(exec.ID, "FAILED", "No start node found")
		return
	}

	// Find the Node Object
	var node *Node
	for i, n := range graph.Nodes {
		if n.ID == currentNodeID {
			node = &graph.Nodes[i]
			break
		}
	}

	if node == nil {
		Logger.Errorf("[Worker] Node %s not found in graph", currentNodeID)
		markExecutionFinal(exec.ID, "FAILED", "Node not found")
		return
	}

	Logger.Infof("[Worker] Executing Node: %s (%s)", node.Label, node.Type)

	// --- EXECUTE NODE LOGIC ---
	output := ""
	status := "success"
	handleToFollow := "default" // Used for branching

	switch models.NodeType(node.Type) {
	case models.NodeTypeTrigger:
		output = "Triggered"
		handleToFollow = "default"

	case models.NodeTypeAction:
		actionType, _ := node.Properties["action"].(string)
		output = "Executed " + actionType

		// Convert Context JSON to map
		var ctxData map[string]interface{}
		if exec.Context.Valid {
			json.Unmarshal([]byte(exec.Context.String), &ctxData)
		}

		actionCtx := actions.ActionContext{
			ExecutionID: exec.ID,
			NodeData:    node.Properties,
			ContextData: ctxData,
		}

		// Look up action
		if action, ok := actions.Get(actionType); ok {
			out, err := action.Execute(context.Background(), actionCtx)
			output = out
			if err != nil {
				status = "failed"
				exec.HasFailed = true
				Logger.Warnf("[Worker] Action %s failed: %v", actionType, err)
				// In new design, Action failure doesn't branch. It just marks execution as failed (or continues if we want execution to flow)?
				// User said: "There should be only a single handle of input and output for an action"
				// Assuming we still want to stop on failure or continue on 'default'?
				// If the user removed the fail handle, we likely treat it as linear.
				// Let's keep status='failed' but handleToFollow='default' so strictly linear.
				// But we still stop if status == "failed" later in the function (line 216).
				// Implementation note: If we want it to Continue on error, we shouldn't return at line 219.
				// For now, let's assume Failure STOPS the workflow (standard behavior) unless caught.
				// Since we removed 'fail' handle, we probably just stop.
			}
		} else {
			Logger.Warnf("[Worker] Unknown action type: %s", actionType)
			output = "Unknown action type"
			status = "failed"
			exec.HasFailed = true
		}
		// ALWAYS follow default for Action
		handleToFollow = "default"

	case models.NodeTypeCondition:
		// Logic: Check properties or random
		// For MVP, if property "force" is set, use that, else random 50/50
		force, ok := node.Properties["force"].(string)
		result := false

		if ok {
			if force == "true" {
				result = true
			} else {
				result = false
			}
		} else {
			// Random
			if rand.Float32() > 0.5 {
				result = true
			} else {
				result = false
			}
		}

		if result {
			output = "Condition: True"
			handleToFollow = "true"
		} else {
			output = "Condition: False"
			handleToFollow = "false"
		}
	}

	// Record Result
	recordStepResult(exec.ID, exec.ResultJSON, StepResult{
		NodeID: currentNodeID,
		Status: status,
		Output: output,
	}, exec.HasFailed)

	if status == "failed" {
		// For now, stop on failure unless we implement error handling edges later
		markExecutionFinal(exec.ID, "FAILED", "Node execution failed")
		return
	}

	// --- FIND NEXT NODE ---
	// Look for edge starting at currentNodeID with handle == handleToFollow
	// If handleToFollow is "default" but no default edge exists, take ANY edge? No, be strict.
	var nextNodeID string

	// First try specific handle
	for _, edge := range graph.Edges {
		if edge.Source == currentNodeID && edge.Handle == handleToFollow {
			nextNodeID = edge.Target
			break
		}
	}

	// If action/trigger and no "default" edge found, maybe allow any edge (if UI doesn't use handles for simple flows)
	if nextNodeID == "" && handleToFollow == "default" {
		for _, edge := range graph.Edges {
			if edge.Source == currentNodeID {
				nextNodeID = edge.Target
				break
			}
		}
	}

	if nextNodeID != "" {
		// Find next node object to check for delays
		var nextNode *Node
		for i, n := range graph.Nodes {
			if n.ID == nextNodeID {
				nextNode = &graph.Nodes[i]
				break
			}
		}

		delayDuration := time.Duration(0)
		if nextNode != nil {
			// Extract delay properties
			// properties map[string]interface{}, so values might be float64 (JSON default) or string
			// We need to safely convert.

			getFloat := func(key string) float64 {
				if val, ok := nextNode.Properties[key]; ok {
					switch v := val.(type) {
					case float64:
						return v
					case int:
						return float64(v)
					case string:
						if f, err := strconv.ParseFloat(v, 64); err == nil {
							return f
						}
					}
				}
				return 0
			}

			days := getFloat("delay_days")
			hours := getFloat("delay_hours")

			if days > 0 {
				delayDuration += time.Duration(days) * 24 * time.Hour
			}
			if hours > 0 {
				delayDuration += time.Duration(hours) * time.Hour
			}
		}

		// Schedule next node
		nextRunAt := time.Now().Add(delayDuration)
		Logger.Infof("[Worker] Scheduling next node %s to run at %s (Delay: %s)", nextNodeID, nextRunAt.Format(time.RFC3339), delayDuration)

		updateExecutionNode(exec.ID, nextNodeID, nextRunAt)
	} else {
		// End of flow
		markExecutionFinal(exec.ID, "COMPLETED", "")
	}
}

func recordStepResult(executionID int, currentJSON sql.NullString, result StepResult, hasFailed bool) error {
	var results []StepResult
	if currentJSON.Valid && currentJSON.String != "" {
		json.Unmarshal([]byte(currentJSON.String), &results)
	}
	results = append(results, result)

	jsonBytes, _ := json.Marshal(results)

	_, err := db.GetDB().Exec(`
		UPDATE workflow_executions 
		SET step_results = $1, has_failed = $2
		WHERE id = $3
	`, string(jsonBytes), hasFailed, executionID)

	return err
}

func updateExecutionNode(executionID int, nextNodeID string, nextRunAt time.Time) {
	_, err := db.GetDB().Exec(`
		UPDATE workflow_executions 
		SET current_node_id = $1, next_run_at = $2 
		WHERE id = $3
	`, nextNodeID, nextRunAt, executionID)
	if err != nil {
		Logger.Error("Failed to update execution node:", err)
	}
}

func markExecutionFinal(executionID int, status string, resultReason string) {
	_, err := db.GetDB().Exec(`
		UPDATE workflow_executions 
		SET status = $2, result = $3, finished_at = NOW() 
		WHERE id = $1
	`, executionID, status, resultReason)
	if err != nil {
		Logger.Error("Failed to mark execution final:", err)
	}
}

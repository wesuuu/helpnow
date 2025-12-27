package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
	"github.com/wesuuu/helpnow/backend/workflows"
)

// ScheduledExecution tracks an execution in progress
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
	var graph workflows.Graph
	if err := json.Unmarshal([]byte(exec.GraphJSON), &graph); err != nil {
		Logger.Errorf("[Worker] Failed to unmarshal graph for execution %d: %v", exec.ID, err)
		markExecutionFinal(exec.ID, "FAILED", "Invalid graph JSON")
		return
	}

	currentNodeID := ""
	if exec.CurrentNodeID.Valid {
		currentNodeID = exec.CurrentNodeID.String
	} else {
		// New execution: Find Start Node (Trigger)
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
	var node *workflows.Node
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

	// Prepare Context
	var ctxData map[string]interface{}
	if exec.Context.Valid {
		json.Unmarshal([]byte(exec.Context.String), &ctxData)
	}
	if ctxData == nil {
		ctxData = make(map[string]interface{})
	}

	switch models.NodeType(node.Type) {
	case models.NodeTypeTrigger:
		output = "Triggered"
		handleToFollow = "default"

	case models.NodeTypeAction:
		actionType, _ := node.Properties["action"].(string)
		output = "Executed " + actionType

		// Look up action template
		if actionTemplate, ok := workflows.GetAction(actionType); ok {
			// Create a new instance of the action (to avoid shared state)
			action := reflect.New(reflect.TypeOf(actionTemplate).Elem()).Interface().(workflows.Action)

			// Unmarshal node properties into the action struct
			if propBytes, err := json.Marshal(node.Properties); err == nil {
				if err := json.Unmarshal(propBytes, action); err != nil {
					Logger.Warnf("[Worker] Failed to unmarshal action properties: %v", err)
				}
			}

			// Execute with only context data
			out, err := action.Execute(context.Background(), ctxData)
			output = out
			if err != nil {
				status = "failed"
				exec.HasFailed = true
				Logger.Warnf("[Worker] Action %s failed: %v", actionType, err)
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
		// Logic Support
		logicType := "Condition" // Hardcoded for now, or property?

		if logicTemplate, ok := workflows.GetLogic(logicType); ok {
			// Create a new instance of the logic
			logic := reflect.New(reflect.TypeOf(logicTemplate).Elem()).Interface().(workflows.Logic)

			// Unmarshal node properties into the logic struct
			if propBytes, err := json.Marshal(node.Properties); err == nil {
				if err := json.Unmarshal(propBytes, logic); err != nil {
					Logger.Warnf("[Worker] Failed to unmarshal logic properties: %v", err)
				}
			}

			// Evaluate with only context data
			res, out, err := logic.Evaluate(context.Background(), ctxData)
			if err != nil {
				Logger.Warnf("[Worker] Logic failed: %v", err)
				status = "failed"
				exec.HasFailed = true
			}
			output = out
			if res {
				handleToFollow = "true"
			} else {
				handleToFollow = "false"
			}
		} else {
			// Fallback (shouldn't happen if registered)
			output = "Logic type not found"
			status = "failed"
		}
	}

	// Record Result
	recordStepResult(exec.ID, exec.ResultJSON, StepResult{
		NodeID: currentNodeID,
		Status: status,
		Output: output,
	}, exec.HasFailed)

	if status == "failed" {
		markExecutionFinal(exec.ID, "FAILED", "Node execution failed")
		return
	}

	// --- FIND NEXT NODE ---
	var nextNodeID string

	// First try specific handle
	for _, edge := range graph.Edges {
		if edge.Source == currentNodeID && edge.Handle == handleToFollow {
			nextNodeID = edge.Target
			break
		}
	}

	// Fallback to default if no specific handle edge found (for actions)
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
		var nextNode *workflows.Node
		for i, n := range graph.Nodes {
			if n.ID == nextNodeID {
				nextNode = &graph.Nodes[i]
				break
			}
		}

		delayDuration := time.Duration(0)
		if nextNode != nil {
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

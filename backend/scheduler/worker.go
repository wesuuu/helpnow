package scheduler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/outbound"
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
		log.Println("Scheduler error querying executions:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exec ScheduledExecution
		if err := rows.Scan(&exec.ID, &exec.WorkflowID, &exec.CurrentNodeID, &exec.GraphJSON, &exec.ResultJSON, &exec.HasFailed, &exec.Context); err != nil {
			log.Println("Scheduler scan error:", err)
			continue
		}
		processSingleExecution(exec)
	}
}

func processSingleExecution(exec ScheduledExecution) {
	var graph Graph
	if err := json.Unmarshal([]byte(exec.GraphJSON), &graph); err != nil {
		log.Printf("[Worker] Failed to unmarshal graph for execution %d: %v\n", exec.ID, err)
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
			if n.Type == "TRIGGER" {
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
		log.Printf("[Worker] No start node found for execution %d\n", exec.ID)
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
		log.Printf("[Worker] Node %s not found in graph\n", currentNodeID)
		markExecutionFinal(exec.ID, "FAILED", "Node not found")
		return
	}

	log.Printf("[Worker] Executing Node: %s (%s)\n", node.Label, node.Type)

	// --- EXECUTE NODE LOGIC ---
	output := ""
	status := "success"
	handleToFollow := "default" // Used for branching

	switch node.Type {
	case "TRIGGER":
		output = "Triggered"
		handleToFollow = "default"

	case "ACTION":
		actionType, _ := node.Properties["action"].(string)
		output = "Executed " + actionType

		if actionType == "Send Email" {
			// Extract template ID
			var templateID int
			if tid, ok := node.Properties["template_id"]; ok {
				switch v := tid.(type) {
				case float64:
					templateID = int(v)
				case int:
					templateID = v
				case string:
					// simple atoi if needed, or 0
				}
			}

			if templateID > 0 {
				// Fetch Template
				var subject, body string
				err := db.GetDB().QueryRow("SELECT subject, body FROM email_templates WHERE id = $1", templateID).Scan(&subject, &body)
				if err != nil {
					log.Printf("[Worker] Failed to fetch template %d: %v", templateID, err)
					status = "failed"
					output = "Failed to fetch template"
					exec.HasFailed = true
				} else {
					// Extract Recipient from Context
					recipient := ""
					if exec.Context.Valid {
						var ctxData map[string]interface{}
						if err := json.Unmarshal([]byte(exec.Context.String), &ctxData); err == nil {
							if email, ok := ctxData["email"].(string); ok {
								recipient = email
							} else if _, ok := ctxData["person_id"].(float64); ok {
								// TODO: Fetch person email if ID provided?
								// For now MVP expects email in context
							}
						}
					}

					// Fallback: If no recipient, check if we want to default to something (or fail)
					if recipient == "" {
						// MVP: Log warning, ideally fail
						log.Printf("[Worker] No recipient email found in context for execution %d", exec.ID)
						// For testing convenience, we might skip fail if it's a test run
						status = "failed"
						output = "No recipient email in context"
						exec.HasFailed = true
					} else {
						// Send Email
						err := outbound.SendEmail(nil, outbound.Email{
							To:      recipient,
							From:    "notifications@helpnow.ai", // Default sender
							Subject: subject,
							Body:    body,
						})
						if err != nil {
							log.Printf("[Worker] Failed to send email: %v", err)
							status = "failed"
							output = "Failed to send email: " + err.Error()
							exec.HasFailed = true
						} else {
							output = fmt.Sprintf("Email sent to %s (Template %d)", recipient, templateID)
						}
					}
				}
			} else {
				output = "No template selected"
			}
		}

		// Mock Fail
		if actionType == "FAIL" {
			status = "failed"
			exec.HasFailed = true
		}

	case "CONDITION":
		// logic: e.g. "random" or property check
		// For MVP, if property "force" is set, use that, else random 50/50
		force, ok := node.Properties["force"].(string)
		if ok {
			if force == "true" {
				output = "Condition met (forced)"
				handleToFollow = "true"
			} else {
				output = "Condition failed (forced)"
				handleToFollow = "false"
			}
		} else {
			// Random
			if rand.Float32() > 0.5 {
				output = "Condition met (random)"
				handleToFollow = "true"
			} else {
				output = "Condition failed (random)"
				handleToFollow = "false"
			}
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
						// Try parsing if string? For now assume valid JSON number types
						return 0
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
		log.Printf("[Worker] Scheduling next node %s to run at %s (Delay: %s)\n", nextNodeID, nextRunAt.Format(time.RFC3339), delayDuration)
		
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
		log.Println("Failed to update execution node:", err)
	}
}

func markExecutionFinal(executionID int, status string, resultReason string) {
	_, err := db.GetDB().Exec(`
		UPDATE workflow_executions 
		SET status = $2, result = $3, finished_at = NOW() 
		WHERE id = $1
	`, executionID, status, resultReason)
	if err != nil {
		log.Println("Failed to mark execution final:", err)
	}
}

package scheduler

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/wesuuu/helpnow/backend/db"
)

var Logger *log.Logger

func init() {
	Logger = log.New("scheduler")
	Logger.SetHeader("${time_rfc3339} | ${level} | ${prefix} |")
}

func Start() {
	ticker := time.NewTicker(1 * time.Minute)
	// ... existing code ...
	go func() {
		for range ticker.C {
			runDueCampaigns()
			runScheduledWorkflows()
		}
	}()
	Logger.Info("Scheduler started")
	StartWorker() // Start workflow worker
}

func runDueCampaigns() {
	// ... no changes to this function ...
	dbConn := db.GetDB()
	// Find Active campaigns where next_run_at <= now
	rows, err := dbConn.Query("SELECT id, name, schedule_interval FROM email_campaigns WHERE status = 'ACTIVE' AND next_run_at <= NOW()")
	if err != nil {
		Logger.Error("Scheduler error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, interval string
		if err := rows.Scan(&id, &name, &interval); err != nil {
			continue
		}

		Logger.Infof("Executing Campaign: %s (ID: %d)", name, id)

		// Simulate Execution
		// Mock stats
		sent := rand.Intn(1000) + 100
		success := 0.8 + (rand.Float64() * 0.2) // 80-100% success

		// Record Run
		_, err = dbConn.Exec("INSERT INTO campaign_runs (campaign_id, sent_count, success_rate, executed_at) VALUES ($1, $2, $3, NOW())", id, sent, success)
		if err != nil {
			Logger.Warn("Failed to record run:", err)
		}

		// Update Next Run
		var nextRun time.Time
		if interval == "DAILY" {
			nextRun = time.Now().Add(24 * time.Hour)
		} else if interval == "WEEKLY" {
			nextRun = time.Now().Add(7 * 24 * time.Hour)
		} else {
			// One time or unknown, just push it far future or set to null/completed?
			// For now, push it a day
			nextRun = time.Now().Add(24 * time.Hour)
		}

		_, err = dbConn.Exec("UPDATE email_campaigns SET next_run_at = $1 WHERE id = $2", nextRun, id)
		if err != nil {
			Logger.Warn("Failed to update next run:", err)
		}
	}
}

func runScheduledWorkflows() {
	dbConn := db.GetDB()
	// Query due scheduled triggers
	rows, err := dbConn.Query(`
		SELECT wt.id, wt.workflow_id, wt.node_id, w.audience_id, wt.config 
		FROM workflow_triggers wt
		JOIN workflows w ON wt.workflow_id = w.id 
		WHERE wt.type = 'SCHEDULE' AND w.status='ACTIVE' AND wt.next_run_at <= NOW()
	`)
	if err != nil {
		Logger.Error("Scheduler error checking triggers:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var triggerID, workflowID, audienceID int
		var nodeID string
		var audIDPtr *int
		var configStr sql.NullString

		if err := rows.Scan(&triggerID, &workflowID, &nodeID, &audIDPtr, &configStr); err != nil {
			Logger.Error("Scan error:", err)
			continue
		}
		if audIDPtr != nil {
			audienceID = *audIDPtr
		}

		Logger.Infof("Triggering Scheduled Workflow ID: %d (Trigger: %d, Node: %s)", workflowID, triggerID, nodeID)

		// Create Execution
		contextData := map[string]interface{}{}

		// 1. Legacy Audience ID
		if audienceID != 0 {
			contextData["audience_id"] = audienceID
		}

		// 2. Trigger Config Audience IDs
		if configStr.Valid && configStr.String != "" {
			var config struct {
				AudienceIDs []int `json:"audience_ids"`
			}
			if err := json.Unmarshal([]byte(configStr.String), &config); err == nil {
				if len(config.AudienceIDs) > 0 {
					contextData["audience_ids"] = config.AudienceIDs
				}
			}
		}
		contextJSON, _ := json.Marshal(contextData)

		// Note: We set current_node_id to the trigger node ID directly
		_, err = dbConn.Exec(`
			INSERT INTO workflow_executions (workflow_id, current_node_id, status, next_run_at, created_at, context) 
			VALUES ($1, $2, 'PENDING', NOW(), NOW(), $3)`,
			workflowID, nodeID, string(contextJSON))

		if err != nil {
			Logger.Error("Failed to create execution:", err)
			continue
		}

		// Calculate Next Run (MVP: +24h)
		nextRun := time.Now().Add(24 * time.Hour)

		_, err = dbConn.Exec("UPDATE workflow_triggers SET next_run_at = $1 WHERE id = $2", nextRun, triggerID)
		if err != nil {
			Logger.Error("Failed to update trigger next run:", err)
		}
	}
}

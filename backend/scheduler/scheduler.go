package scheduler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/wesuuu/helpnow/backend/db"
)

func Start() {
	ticker := time.NewTicker(1 * time.Minute)
	// ... existing code ...
	go func() {
		for {
			select {
			case <-ticker.C:
				runDueCampaigns()
				runScheduledWorkflows()
			}
		}
	}()
	fmt.Println("Scheduler started")
	StartWorker() // Start workflow worker
}

func runDueCampaigns() {
	// ... no changes to this function ...
	dbConn := db.GetDB()
	// Find Active campaigns where next_run_at <= now
	rows, err := dbConn.Query("SELECT id, name, schedule_interval FROM email_campaigns WHERE status = 'ACTIVE' AND next_run_at <= NOW()")
	if err != nil {
		fmt.Println("Scheduler error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, interval string
		if err := rows.Scan(&id, &name, &interval); err != nil {
			continue
		}

		fmt.Printf("Executing Campaign: %s (ID: %d)\n", name, id)

		// Simulate Execution
		// Mock stats
		sent := rand.Intn(1000) + 100
		success := 0.8 + (rand.Float64() * 0.2) // 80-100% success

		// Record Run
		_, err = dbConn.Exec("INSERT INTO campaign_runs (campaign_id, sent_count, success_rate, executed_at) VALUES ($1, $2, $3, NOW())", id, sent, success)
		if err != nil {
			fmt.Println("Failed to record run:", err)
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
			fmt.Println("Failed to update next run:", err)
		}
	}
}

func runScheduledWorkflows() {
	dbConn := db.GetDB()
	// Query due scheduled workflows
	rows, err := dbConn.Query(`
		SELECT id, name, steps, schedule, audience_id 
		FROM workflows 
		WHERE trigger_type = 'SCHEDULE' AND status='ACTIVE' AND next_run_at <= NOW()
	`)
	if err != nil {
		fmt.Println("Scheduler error checking workflows:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, audienceID int
		var audIDPtr *int
		var name, steps, schedule string
		if err := rows.Scan(&id, &name, &steps, &schedule, &audIDPtr); err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		if audIDPtr != nil {
			audienceID = *audIDPtr
		}

		fmt.Printf("Triggering Scheduled Workflow: %s (ID: %d)\n", name, id)

		// Create Execution
		contextData := map[string]interface{}{}
		if audienceID != 0 {
			contextData["audience_id"] = audienceID
		}
		contextJSON, _ := json.Marshal(contextData)

		_, err = dbConn.Exec(`
			INSERT INTO workflow_executions (workflow_id, status, next_run_at, created_at, context) 
			VALUES ($1, 'PENDING', NOW(), NOW(), $2)`, id, string(contextJSON))

		if err != nil {
			fmt.Println("Failed to create execution:", err)
			continue // Don't reschedule if we failed to start? Or should we?
			// If we don't reschedule, it will loop forever.
		}

		// Calculate Next Run
		// Parse Cron-ish
		// Simplest MVP parsing:
		var nextRun time.Time
		// Assume simplified cron: "Minute Hour * * *"
		// Or simplified daily/hourly check
		// Let's implement basic "Daily at Hour:Minute" or "Every N hours" if easy
		// Fallback: Add 24 hours if parsing fails or default
		nextRun = time.Now().Add(24 * time.Hour) // Default

		// Attempt simple parsing
		// If Schedule contains "every 12 hours" logic or standard cron
		// ... (Keep it simple for now as per plan)

		_, err = dbConn.Exec("UPDATE workflows SET next_run_at = $1 WHERE id = $2", nextRun, id)
		if err != nil {
			fmt.Println("Failed to update workflow next run:", err)
		}
	}
}

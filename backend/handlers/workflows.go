package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
)

type EventDefinition struct {
	ID          int       `json:"id"`
	SiteID      int       `json:"site_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Workflow struct {
	ID             int        `json:"id"`
	OrganizationID *int       `json:"organization_id"`
	SiteID         *int       `json:"site_id"`
	SiteName       string     `json:"site_name,omitempty"`
	AudienceID     *int       `json:"audience_id"`
	Name           string     `json:"name"`
	TriggerType    string     `json:"trigger_type"`
	TriggerEvent   *string    `json:"trigger_event"`
	Steps          string     `json:"steps"`
	Schedule       *string    `json:"schedule"`
	NextRunAt      *time.Time `json:"next_run_at"` // Added
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
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

// Workflows

func CreateWorkflow(c echo.Context) error {
	var wf Workflow
	if err := c.Bind(&wf); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Calculate Initial Next Run if Schedule
	if wf.TriggerType == "SCHEDULE" && wf.Schedule != nil && *wf.Schedule != "" {
		// MVP: If simple cron "* * * * *", current logic in updates handles parsing.
		// Here we just set it to NOW() to run immediately or soon, or parse properly.
		// For MVP simplicity, let's set it to NOW() so scheduler picks it up and calculates real next run.
		now := time.Now()
		wf.NextRunAt = &now
	}

	// Default Org ID to 1 for MVP if not set
	if wf.OrganizationID == nil {
		orgID := 1
		wf.OrganizationID = &orgID
	}

	query := `INSERT INTO workflows (organization_id, site_id, audience_id, name, trigger_type, trigger_event, steps, schedule, next_run_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at`
	err := db.GetDB().QueryRow(query, wf.OrganizationID, wf.SiteID, wf.AudienceID, wf.Name, wf.TriggerType, wf.TriggerEvent, wf.Steps, wf.Schedule, wf.NextRunAt, "ACTIVE").Scan(&wf.ID, &wf.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to create workflow: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create workflow"})
	}
	wf.Status = "ACTIVE"
	return c.JSON(http.StatusCreated, wf)
}

func ListWorkflows(c echo.Context) error {
	siteID := c.QueryParam("site_id")
	var rows *sql.Rows
	var err error

	if siteID != "" && siteID != "null" {
		rows, err = db.GetDB().Query(`
			SELECT w.id, w.organization_id, w.site_id, s.name, w.audience_id, w.name, w.trigger_type, w.trigger_event, w.steps, w.schedule, w.next_run_at, w.status, w.created_at 
			FROM workflows w
			LEFT JOIN sites s ON w.site_id = s.id
			WHERE w.site_id = $1 
			ORDER BY w.created_at DESC`, siteID)
	} else {
		rows, err = db.GetDB().Query(`
			SELECT w.id, w.organization_id, w.site_id, s.name, w.audience_id, w.name, w.trigger_type, w.trigger_event, w.steps, w.schedule, w.next_run_at, w.status, w.created_at 
			FROM workflows w
			LEFT JOIN sites s ON w.site_id = s.id
			ORDER BY w.created_at DESC`)
	}

	if err != nil {
		c.Logger().Error("Failed to list workflows: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list workflows"})
	}
	defer rows.Close()

	workflows := []Workflow{}
	for rows.Next() {
		var w Workflow
		var siteName sql.NullString // Handle Join NULLs
		if err := rows.Scan(&w.ID, &w.OrganizationID, &w.SiteID, &siteName, &w.AudienceID, &w.Name, &w.TriggerType, &w.TriggerEvent, &w.Steps, &w.Schedule, &w.NextRunAt, &w.Status, &w.CreatedAt); err == nil {
			if siteName.Valid {
				w.SiteName = siteName.String
			}
			workflows = append(workflows, w)
		} else {
			c.Logger().Error("Scan error: ", err)
		}
	}
	return c.JSON(http.StatusOK, workflows)
}

// Event Definitions

func CreateEventDefinition(c echo.Context) error {
	var ed EventDefinition
	if err := c.Bind(&ed); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	query := `INSERT INTO event_definitions (site_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.GetDB().QueryRow(query, ed.SiteID, ed.Name, ed.Description).Scan(&ed.ID, &ed.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to create event definition: ", err)
		// Usually unique constraint violation
		return c.JSON(http.StatusConflict, map[string]string{"error": "Event already exists"})
	}
	return c.JSON(http.StatusCreated, ed)
}

func ListEventDefinitions(c echo.Context) error {
	siteID := c.QueryParam("site_id")
	rows, err := db.GetDB().Query(`SELECT id, site_id, name, description, created_at FROM event_definitions WHERE site_id = $1 ORDER BY name ASC`, siteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list events"})
	}
	defer rows.Close()

	events := []EventDefinition{}
	for rows.Next() {
		var e EventDefinition
		if err := rows.Scan(&e.ID, &e.SiteID, &e.Name, &e.Description, &e.CreatedAt); err == nil {
			events = append(events, e)
		}
	}
	return c.JSON(http.StatusOK, events)
}

// Trigger Logic for Internal Use

func TriggerWorkflow(siteID int, eventName string, contextData map[string]interface{}) error {
	// 1. Find active workflows for this site and event
	rows, err := db.GetDB().Query(`SELECT id FROM workflows WHERE site_id = $1 AND trigger_event = $2 AND status = 'ACTIVE'`, siteID, eventName)
	if err != nil {
		return err
	}
	defer rows.Close()

	contextJSON, _ := json.Marshal(contextData)

	for rows.Next() {
		var workflowID int
		if err := rows.Scan(&workflowID); err == nil {
			// 2. Create Execution
			_, err := db.GetDB().Exec(`
				INSERT INTO workflow_executions (workflow_id, status, context, next_run_at) 
				VALUES ($1, 'PENDING', $2, NOW())`,
				workflowID, string(contextJSON))
			if err != nil {
				// Log but continue
				// fmt.Println("Failed to trigger execution:", err)
			}
		}
	}
	return nil
}

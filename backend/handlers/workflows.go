package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
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
	if wf.TriggerType == string(models.TriggerTypeSchedule) && wf.Schedule != nil && *wf.Schedule != "" {
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

	// Create Workflow Record
	query := `INSERT INTO workflows (organization_id, site_id, audience_id, name, trigger_type, trigger_event, steps, schedule, next_run_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at`
	// For legacy support/display, we might still want to populate trigger_type/schedule in the main table if we have a primary trigger?
	// But let's assume we just save them blank or as "multiple" if we move fully?
	// For now, keep saving them as is (binding form input) BUT ALSO parse graph.
	err := db.GetDB().QueryRow(query, wf.OrganizationID, wf.SiteID, wf.AudienceID, wf.Name, wf.TriggerType, wf.TriggerEvent, wf.Steps, wf.Schedule, wf.NextRunAt, "ACTIVE").Scan(&wf.ID, &wf.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to create workflow: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create workflow"})
	}
	wf.Status = "ACTIVE"

	// Parse Graph to extract Triggers
	type GraphNode struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	}
	type GraphStruct struct {
		Nodes []GraphNode `json:"nodes"`
	}

	var graph GraphStruct
	if err := json.Unmarshal([]byte(wf.Steps), &graph); err == nil {
		for _, node := range graph.Nodes {
			if node.Type == string(models.NodeTypeTrigger) {
				// Determine Type
				tType := string(models.TriggerTypeEvent) // Default
				if val, ok := node.Properties["trigger_type"].(string); ok {
					tType = val
				} else if node.Properties["cron"] != nil {
					tType = string(models.TriggerTypeSchedule)
				}

				// Build Config JSON
				configBytes, _ := json.Marshal(node.Properties)
				configJSON := string(configBytes)

				// Determine Next Run (if schedule)
				var nextRun sql.NullTime
				if tType == string(models.TriggerTypeSchedule) {
					nextRun.Time = time.Now() // Run immediately/soon
					nextRun.Valid = true
				}

				_, err := db.GetDB().Exec(`
					INSERT INTO workflow_triggers (workflow_id, node_id, type, config, next_run_at)
					VALUES ($1, $2, $3, $4, $5)
				`, wf.ID, node.ID, tType, configJSON, nextRun)

				if err != nil {
					c.Logger().Error("Failed to save trigger:", err)
				}
			}
		}
	}

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

func GetWorkflow(c echo.Context) error {
	id := c.Param("id")
	var w Workflow
	var siteName sql.NullString

	err := db.GetDB().QueryRow(`
		SELECT w.id, w.organization_id, w.site_id, s.name, w.audience_id, w.name, w.trigger_type, w.trigger_event, w.steps, w.schedule, w.next_run_at, w.status, w.created_at 
		FROM workflows w
		LEFT JOIN sites s ON w.site_id = s.id
		WHERE w.id = $1`, id).Scan(&w.ID, &w.OrganizationID, &w.SiteID, &siteName, &w.AudienceID, &w.Name, &w.TriggerType, &w.TriggerEvent, &w.Steps, &w.Schedule, &w.NextRunAt, &w.Status, &w.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Workflow not found"})
		}
		c.Logger().Error("Failed to get workflow: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get workflow"})
	}

	if siteName.Valid {
		w.SiteName = siteName.String
	}

	return c.JSON(http.StatusOK, w)
}

func UpdateWorkflow(c echo.Context) error {
	id := c.Param("id")
	var wf Workflow
	if err := c.Bind(&wf); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Calculate Initial Next Run if Schedule
	if wf.TriggerType == string(models.TriggerTypeSchedule) && wf.Schedule != nil && *wf.Schedule != "" {
		now := time.Now()
		wf.NextRunAt = &now
	}

	// Update Workflow Record
	query := `UPDATE workflows SET site_id=$1, audience_id=$2, name=$3, trigger_type=$4, trigger_event=$5, steps=$6, schedule=$7, next_run_at=$8 WHERE id=$9`
	_, err := db.GetDB().Exec(query, wf.SiteID, wf.AudienceID, wf.Name, wf.TriggerType, wf.TriggerEvent, wf.Steps, wf.Schedule, wf.NextRunAt, id)
	if err != nil {
		c.Logger().Error("Failed to update workflow: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update workflow"})
	}

	// Re-create Triggers
	// 1. Delete existing triggers for this workflow
	_, err = db.GetDB().Exec("DELETE FROM workflow_triggers WHERE workflow_id = $1", id)
	if err != nil {
		c.Logger().Error("Failed to clear triggers: ", err)
	}

	// 2. Parse Graph and create new triggers
	type GraphNode struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
	}
	type GraphStruct struct {
		Nodes []GraphNode `json:"nodes"`
	}

	var graph GraphStruct
	if err := json.Unmarshal([]byte(wf.Steps), &graph); err == nil {
		for _, node := range graph.Nodes {
			if node.Type == string(models.NodeTypeTrigger) {
				// Determine Type
				tType := string(models.TriggerTypeEvent) // Default
				if val, ok := node.Properties["trigger_type"].(string); ok {
					tType = val
				} else if node.Properties["cron"] != nil {
					tType = string(models.TriggerTypeSchedule)
				}

				// Build Config JSON
				configBytes, _ := json.Marshal(node.Properties)
				configJSON := string(configBytes)

				// Determine Next Run (if schedule)
				var nextRun sql.NullTime
				if tType == string(models.TriggerTypeSchedule) {
					nextRun.Time = time.Now()
					nextRun.Valid = true
				}

				_, err := db.GetDB().Exec(`
					INSERT INTO workflow_triggers (workflow_id, node_id, type, config, next_run_at)
					VALUES ($1, $2, $3, $4, $5)
				`, id, node.ID, tType, configJSON, nextRun)

				if err != nil {
					c.Logger().Error("Failed to save trigger:", err)
				}
			}
		}
	}

	return c.JSON(http.StatusOK, wf)
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
	var rows *sql.Rows
	var err error

	if siteID != "" {
		rows, err = db.GetDB().Query(`SELECT id, site_id, name, description, created_at FROM event_definitions WHERE site_id = $1 ORDER BY name ASC`, siteID)
	} else {
		rows, err = db.GetDB().Query(`SELECT id, site_id, name, description, created_at FROM event_definitions ORDER BY name ASC`)
	}

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
	// query both legacy columns (if we want to support old workflows) and new table.
	// For now, let's assume we migrated or only support new table
	rows, err := db.GetDB().Query(`
		SELECT wt.workflow_id, wt.node_id, wt.config 
		FROM workflow_triggers wt 
		JOIN workflows w ON wt.workflow_id = w.id 
		WHERE wt.type = '` + string(models.TriggerTypeEvent) + `' AND w.status = 'ACTIVE'
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	contextJSON, _ := json.Marshal(contextData)

	type TriggerConfig struct {
		TriggerEvent string `json:"trigger_event"`
		SiteIDs      []int  `json:"site_ids"`
		AudienceIDs  []int  `json:"audience_ids"`
	}

	for rows.Next() {
		var workflowID int
		var nodeID, configStr string
		if err := rows.Scan(&workflowID, &nodeID, &configStr); err == nil {
			// Check Config
			var config TriggerConfig
			if err := json.Unmarshal([]byte(configStr), &config); err != nil {
				continue
			}

			// 1. Check Event Name
			if config.TriggerEvent != eventName {
				continue
			}

			// 2. Check Site ID (if SiteIDs is not empty)
			// If SiteIDs is empty, maybe run for all? Or run for none?
			// Frontend: "Select logic: OR (Run if event happens on any selected site)".
			// If empty, probably shouldn't run unless it means "All"?
			// Let's assume empty means "None selected" so don't run.
			// Unless we add an "All Sites" flag. For now, strict check.
			siteMatch := false
			if len(config.SiteIDs) > 0 {
				siteMatch = false
				for _, sid := range config.SiteIDs {
					if sid == siteID {
						siteMatch = true
						break
					}
				}
				if !siteMatch {
					continue
				}
			}

			// Check Audience Filter
			if len(config.AudienceIDs) > 0 {
				email, ok := contextData["email"].(string)
				if !ok || email == "" {
					email, ok = contextData["user_email"].(string)
				}

				if !ok || email == "" {
					// Cannot verify audience without email
					continue
				}

				var count int
				err := db.GetDB().QueryRow(`
					SELECT COUNT(*) 
					FROM audience_memberships am
					JOIN people p ON am.person_id = p.id
					WHERE p.email = $1 AND am.audience_id = ANY($2)
				`, email, pq.Array(config.AudienceIDs)).Scan(&count)

				if err != nil || count == 0 {
					continue
				}
				if err != nil || count == 0 {
					continue
				}
			}

			// 2. Create Execution
			_, err := db.GetDB().Exec(`
				INSERT INTO workflow_executions (workflow_id, current_node_id, status, context, next_run_at) 
				VALUES ($1, $2, 'PENDING', $3, NOW())`,
				workflowID, nodeID, string(contextJSON))
			if err != nil {
				// fmt.Println("Failed to trigger execution:", err)
			}
		}
	}

	return nil
}

package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateCampaign(c echo.Context) error {
	var campaign models.EmailCampaign
	if err := c.Bind(&campaign); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Status defaults to DRAFT
	campaign.Status = "DRAFT"

	query := `INSERT INTO email_campaigns (organization_id, audience_segment_id, name, prompt, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, campaign.OrganizationID, campaign.AudienceSegmentID, campaign.Name, campaign.Prompt, campaign.Status).Scan(&campaign.ID, &campaign.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create campaign"})
	}

	return c.JSON(http.StatusCreated, campaign)
}

func GenerateCampaignContent(c echo.Context) error {
	id := c.Param("id")

	// Fetch campaign prompt
	var prompt string
	var orgID int
	err := db.GetDB().QueryRow("SELECT prompt, organization_id FROM email_campaigns WHERE id = $1", id).Scan(&prompt, &orgID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Campaign not found"})
	}

	if prompt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Campaign has no prompt"})
	}

	// Call AI Service (Mocked via direct routine execution logic usually, but here we can just ask a "Marketing Agent")
	// For simplicity, we'll construct a direct prompt to the AI service if possible, or trigger a specific "Copywriting Routine"
	// Let's assume we have a "Marketing Agent" available or can spin one up.
	// For MVP, lets just call the ExecuteRoutine logic but for a special "Generate Email" routine we assume exists,
	// OR we just use the AI Client directly if it supported raw generation.
	// Given the constraints, let's pretend we are triggering a generic "Content Generation" routine.
	// We'll fake the AI response for now since setting up a specific agent dynamically is complex in this step.

	// Real implementation would: find Marketing Agent -> Send Prompt -> Get Result.
	generatedContent := "Subject: Special Offer for You!\n\nHi there,\n\n" + prompt + "\n\nBest,\nHelpNow Team"

	// Update campaign content
	_, err = db.GetDB().Exec("UPDATE email_campaigns SET content = $1 WHERE id = $2", generatedContent, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save generated content"})
	}

	return c.JSON(http.StatusOK, map[string]string{"content": generatedContent})
}

func UpdateCampaign(c echo.Context) error {
	id := c.Param("id")
	var req models.EmailCampaign
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if req.ScheduleInterval != "" && req.Status == "ACTIVE" {
		// Calculate Next Run
		// Simple logic: Run in 1 minute
		now := time.Now().Add(1 * time.Minute)
		req.NextRunAt = &now
	}

	query := `UPDATE email_campaigns SET content = $1, schedule_interval = $2, status = $3, next_run_at = $4 WHERE id = $5`
	_, err := db.GetDB().Exec(query, req.Content, req.ScheduleInterval, req.Status, req.NextRunAt, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update campaign"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func ListCampaigns(c echo.Context) error {
	orgID := c.QueryParam("organization_id")

	rows, err := db.GetDB().Query("SELECT id, organization_id, audience_segment_id, name, prompt, content, schedule_interval, next_run_at, status, created_at FROM email_campaigns WHERE organization_id = $1", orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch campaigns"})
	}
	defer rows.Close()

	campaigns := []models.EmailCampaign{}
	for rows.Next() {
		var c models.EmailCampaign
		var content, interval sql.NullString // Handle nulls
		if err := rows.Scan(&c.ID, &c.OrganizationID, &c.AudienceSegmentID, &c.Name, &c.Prompt, &content, &interval, &c.NextRunAt, &c.Status, &c.CreatedAt); err != nil {
			continue
		}
		c.Content = content.String
		c.ScheduleInterval = interval.String
		campaigns = append(campaigns, c)
	}

	return c.JSON(http.StatusOK, campaigns)
}

func ListCampaignRuns(c echo.Context) error {
	campaignID := c.Param("id")
	rows, err := db.GetDB().Query("SELECT id, campaign_id, sent_count, success_rate, executed_at FROM campaign_runs WHERE campaign_id = $1 ORDER BY executed_at DESC", campaignID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch runs"})
	}
	defer rows.Close()

	runs := []models.CampaignRun{}
	for rows.Next() {
		var r models.CampaignRun
		if err := rows.Scan(&r.ID, &r.CampaignID, &r.SentCount, &r.SuccessRate, &r.ExecutedAt); err != nil {
			continue
		}
		runs = append(runs, r)
	}
	return c.JSON(http.StatusOK, runs)
}

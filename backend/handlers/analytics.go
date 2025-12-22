package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
)

type AnalyticsEvent struct {
	Token     string                 `json:"token"`
	EventType string                 `json:"event_type"` // impression, click, conversion
	Meta      map[string]interface{} `json:"meta"`
	Value     float64                `json:"value"` // For conversions
}

func getClickHouseDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("tcp://192.168.0.16:9001?username=%s&password=%s",
		os.Getenv("CLICKHOUSE_USER"),
		os.Getenv("CLICKHOUSE_PASSWORD"))

	return sql.Open("clickhouse", dsn)
}

func TrackEvent(c echo.Context) error {
	var event AnalyticsEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 1. Determine Identity (Site vs Signup ID)
	var organizationID int
	var campaignID int // Used for generic grouping ID for now

	if strings.HasPrefix(event.Token, "HN-") {
		// It's a Site Tracking ID
		// In ClickHouse tables, we have campaign_id. We might need to map Site ID to Campaign ID or just use Site ID as Campaign ID
		// For now, let's treat Site ID as the "Campaign ID" dimension in the metrics table, or we need to add site_id to clickhouse.
		// User requirement "metrics... for sales page". treating site_id as campaign_id is a reasonable MVP short to avoid schema change in ClickHouse if not requested.
		// However, user asked to "add site ID ... similar to google analytics".
		// Let's lookup the site.
		var siteID int
		query := `SELECT organization_id, id FROM sites WHERE tracking_id = $1`
		err := db.GetDB().QueryRow(query, event.Token).Scan(&organizationID, &siteID)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid site tracking ID"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		campaignID = siteID // Storing site_id in campaign_id column for now
	} else {
		// Fallback to Signup Campaign Token
		query := `SELECT organization_id, id FROM signup_campaigns WHERE token = $1`
		err := db.GetDB().QueryRow(query, event.Token).Scan(&organizationID, &campaignID)
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid token"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	}

	// Trigger Workflow (Async or Sync? For MVP Sync is fine or fire goroutine)
	if strings.HasPrefix(event.Token, "HN-") {
		// Lookup siteID if we haven't already (we have, in determining organizationID)
		// We can reuse the lookup logic or just pass the campaignID as siteID since we mapped it.
		// Earlier code: campaignID = siteID

		// Run in goroutine to not block analytics response
		go func(siteID int, eventType string, meta map[string]interface{}) {
			// Trigger defined workflows
			// We use the event.EventType. If user sends custom event type, it will be handled.
			_ = TriggerWorkflow(siteID, eventType, meta)
		}(campaignID, event.EventType, event.Meta)
	}

	// 2. Connect to ClickHouse
	ch, err := getClickHouseDB()
	if err != nil {
		c.Logger().Error("ClickHouse connection error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Analytics storage error"})
	}
	defer ch.Close()

	// 3. Insert into appropriate table
	metaJSON, _ := json.Marshal(event.Meta)
	metaStr := string(metaJSON)

	var insertQuery string
	switch event.EventType {
	case "impression":
		insertQuery = `INSERT INTO impressions (organization_id, campaign_id, count, meta) VALUES (?, ?, 1, ?)`
		_, err = ch.Exec(insertQuery, uint64(organizationID), uint64(campaignID), metaStr)
	case "click":
		insertQuery = `INSERT INTO click_rate (organization_id, campaign_id, clicks, meta) VALUES (?, ?, 1, ?)`
		_, err = ch.Exec(insertQuery, uint64(organizationID), uint64(campaignID), metaStr)
	case "conversion":
		insertQuery = `INSERT INTO cpa (organization_id, campaign_id, cost, meta) VALUES (?, ?, ?, ?)`
		_, err = ch.Exec(insertQuery, uint64(organizationID), uint64(campaignID), event.Value, metaStr)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unknown event type"})
	}

	if err != nil {
		c.Logger().Error("ClickHouse write error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to record event"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "recorded"})
}

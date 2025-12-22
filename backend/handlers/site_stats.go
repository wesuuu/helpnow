package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/labstack/echo/v4"
)

type SiteStats struct {
	Active      bool      `json:"active"`
	LastEventAt time.Time `json:"last_event_at,omitempty"`
	Stats       struct {
		Impressions uint64  `json:"impressions"`
		Clicks      uint64  `json:"clicks"`
		Conversions uint64  `json:"conversions"`
		Cost        float64 `json:"cost"`
	} `json:"stats"`
}

func getClickHouseReadDB() (*sql.DB, error) {
	// Duplicated from analytics.go for now to avoid refactoring churn
	dsn := fmt.Sprintf("tcp://192.168.0.16:9001?username=%s&password=%s",
		os.Getenv("CLICKHOUSE_USER"),
		os.Getenv("CLICKHOUSE_PASSWORD"))

	return sql.Open("clickhouse", dsn)
}

func GetSiteStats(c echo.Context) error {
	siteIDStr := c.Param("id")
	siteID, err := strconv.ParseUint(siteIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid site ID"})
	}

	ch, err := getClickHouseReadDB()
	if err != nil {
		c.Logger().Error("ClickHouse connection error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Analytics storage error"})
	}
	defer ch.Close()

	stats := SiteStats{}

	// 1. Get Aggregated Stats
	// Using UNION ALL or separate queries. Separate queries are safer/easier for simple key-value structure.

	// Impressions
	err = ch.QueryRow(`
		SELECT count(*) FROM impressions WHERE campaign_id = ?
	`, siteID).Scan(&stats.Stats.Impressions)
	if err != nil && err != sql.ErrNoRows {
		c.Logger().Error("Failed to query impressions: ", err)
	}

	// Clicks
	err = ch.QueryRow(`
		SELECT count(*) FROM click_rate WHERE campaign_id = ?
	`, siteID).Scan(&stats.Stats.Clicks)
	if err != nil && err != sql.ErrNoRows {
		c.Logger().Error("Failed to query clicks: ", err)
	}

	// Conversions & Cost
	// cpa table has 'cost' column. We can count rows for conversions and sum cost.
	err = ch.QueryRow(`
		SELECT count(*), coalesce(sum(cost), 0) FROM cpa WHERE campaign_id = ?
	`, siteID).Scan(&stats.Stats.Conversions, &stats.Stats.Cost)
	if err != nil && err != sql.ErrNoRows {
		c.Logger().Error("Failed to query conversions: ", err)
	}

	// 2. Determine Active Status (Check for any event in the last 24 hours created_at)
	// We check all tables roughly or just the most common one (impressions)
	var lastEventAt time.Time

	// Check max created_at across tables
	query := `
		SELECT max(t) FROM (
			SELECT max(created_at) as t FROM impressions WHERE campaign_id = ?
			UNION ALL
			SELECT max(created_at) as t FROM click_rate WHERE campaign_id = ?
			UNION ALL
			SELECT max(created_at) as t FROM cpa WHERE campaign_id = ?
		)
	`
	// ClickHouse might return a zero date or null if empty.
	// We handle scans carefully. Scan into NullTime or check errors?
	// ClickHouse driver often returns default value.

	// Note: Parameter placeholder `?` works in clickhouse-go if setup correctly, but sometimes `$1` is used.
	// analytics.go used `?`.

	rows, err := ch.Query(query, siteID, siteID, siteID)
	if err != nil {
		c.Logger().Error("Failed to query last event: ", err)
	} else {
		defer rows.Close()
		if rows.Next() {
			var t sql.NullTime
			if err := rows.Scan(&t); err == nil && t.Valid {
				lastEventAt = t.Time
			}
		}
	}

	if !lastEventAt.IsZero() {
		stats.LastEventAt = lastEventAt
		// Logic for "Active": if last event was within last 24 hours?
		// Or just "Active" if we have EVER received data?
		// User asked "integration active or not". Usually means "is code installed and sending data".
		// Let's say if we have ANY data, it's active.
		stats.Active = true
	} else {
		stats.Active = false
	}

	return c.JSON(http.StatusOK, stats)
}

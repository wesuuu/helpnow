package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
)

type DataSource struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Config         string    `json:"config"` // JSON string
	CreatedAt      time.Time `json:"created_at"`
}

type DataSync struct {
	ID         int        `json:"id"`
	SourceID   int        `json:"source_id"`
	AudienceID int        `json:"audience_id"`
	SyncType   string     `json:"sync_type"`
	Schedule   string     `json:"schedule"`
	Query      string     `json:"query"`
	Status     string     `json:"status"`
	LastRunAt  *time.Time `json:"last_run_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Data Sources

func CreateDataSource(c echo.Context) error {
	var ds DataSource
	if err := c.Bind(&ds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Hardcoded org for now, or get from context
	ds.OrganizationID = 1

	query := `INSERT INTO data_sources (organization_id, name, type, config) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := db.GetDB().QueryRow(query, ds.OrganizationID, ds.Name, ds.Type, ds.Config).Scan(&ds.ID, &ds.CreatedAt)
	if err != nil {
		log.Println("Error creating data source:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create data source"})
	}
	return c.JSON(http.StatusCreated, ds)
}

func ListDataSources(c echo.Context) error {
	// orgID := c.QueryParam("organization_id")
	orgID := 1 // Hardcoded
	rows, err := db.GetDB().Query(`SELECT id, organization_id, name, type, config, created_at FROM data_sources WHERE organization_id = $1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list sources"})
	}
	defer rows.Close()

	sources := []DataSource{}
	for rows.Next() {
		var ds DataSource
		if err := rows.Scan(&ds.ID, &ds.OrganizationID, &ds.Name, &ds.Type, &ds.Config, &ds.CreatedAt); err == nil {
			sources = append(sources, ds)
		}
	}
	return c.JSON(http.StatusOK, sources)
}

// Data Syncs

func CreateDataSync(c echo.Context) error {
	var sync DataSync
	if err := c.Bind(&sync); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	sync.Status = "PENDING"

	query := `INSERT INTO data_syncs (source_id, audience_id, sync_type, schedule, query, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	err := db.GetDB().QueryRow(query, sync.SourceID, sync.AudienceID, sync.SyncType, sync.Schedule, sync.Query, sync.Status).Scan(&sync.ID, &sync.CreatedAt)
	if err != nil {
		log.Println("Error creating data sync:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create sync job"})
	}
	return c.JSON(http.StatusCreated, sync)
}

func ListDataSyncs(c echo.Context) error {
	sourceID := c.QueryParam("source_id")
	if sourceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Source ID required"})
	}

	rows, err := db.GetDB().Query(`SELECT id, source_id, audience_id, sync_type, schedule, query, status, last_run_at, created_at FROM data_syncs WHERE source_id = $1 ORDER BY created_at DESC`, sourceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list syncs"})
	}
	defer rows.Close()

	syncs := []DataSync{}
	for rows.Next() {
		var s DataSync
		if err := rows.Scan(&s.ID, &s.SourceID, &s.AudienceID, &s.SyncType, &s.Schedule, &s.Query, &s.Status, &s.LastRunAt, &s.CreatedAt); err == nil {
			syncs = append(syncs, s)
		}
	}
	return c.JSON(http.StatusOK, syncs)
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateAudience(c echo.Context) error {
	var audience models.Audience
	if err := c.Bind(&audience); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	query := `INSERT INTO audiences (organization_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, audience.OrganizationID, audience.Name, audience.Description).Scan(&audience.ID, &audience.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create audience"})
	}

	return c.JSON(http.StatusCreated, audience)
}

func CreateAudienceSegment(c echo.Context) error {
	var segment models.AudienceSegment
	if err := c.Bind(&segment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	query := `INSERT INTO audience_segments (audience_id, name, filters) VALUES ($1, $2, $3) RETURNING id, created_at`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, segment.AudienceID, segment.Name, segment.Filters).Scan(&segment.ID, &segment.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create segment"})
	}

	return c.JSON(http.StatusCreated, segment)
}

func ListAudiences(c echo.Context) error {
	orgID := c.QueryParam("organization_id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}

	dbConn := db.GetDB()

	// Fetch Audiences
	rows, err := dbConn.Query(`SELECT id, organization_id, name, description, created_at FROM audiences WHERE organization_id = $1`, orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch audiences"})
	}
	defer rows.Close()

	type AudienceWithSegments struct {
		models.Audience
		Segments []models.AudienceSegment `json:"segments"`
	}

	audiences := []AudienceWithSegments{}
	for rows.Next() {
		var a AudienceWithSegments
		if err := rows.Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Description, &a.CreatedAt); err != nil {
			continue
		}
		a.Segments = []models.AudienceSegment{}
		audiences = append(audiences, a)
	}

	// Fetch Segments for each audience (not efficient for large datasets but fine for MVP)
	for i := range audiences {
		segRows, err := dbConn.Query(`SELECT id, audience_id, name, filters, created_at FROM audience_segments WHERE audience_id = $1`, audiences[i].ID)
		if err == nil {
			defer segRows.Close()
			for segRows.Next() {
				var s models.AudienceSegment
				if err := segRows.Scan(&s.ID, &s.AudienceID, &s.Name, &s.Filters, &s.CreatedAt); err == nil {
					audiences[i].Segments = append(audiences[i].Segments, s)
				}
			}
		}
	}

	return c.JSON(http.StatusOK, audiences)
}

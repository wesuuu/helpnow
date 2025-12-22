package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func generateTrackingID() string {
	bytes := make([]byte, 4) // 8 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "HN-ERROR"
	}
	return "HN-" + strings.ToUpper(hex.EncodeToString(bytes))
}

func CreateSite(c echo.Context) error {
	var site models.Site
	if err := c.Bind(&site); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	site.TrackingID = generateTrackingID()

	query := `INSERT INTO sites (organization_id, name, url, tracking_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := db.GetDB().QueryRow(query, site.OrganizationID, site.Name, site.URL, site.TrackingID).Scan(&site.ID, &site.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to create site: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create site"})
	}

	return c.JSON(http.StatusCreated, site)
}

func ListSites(c echo.Context) error {
	orgID := c.QueryParam("organization_id")

	rows, err := db.GetDB().Query(`SELECT id, organization_id, name, url, tracking_id, created_at FROM sites WHERE organization_id = $1`, orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list sites"})
	}
	defer rows.Close()

	sites := []models.Site{}
	for rows.Next() {
		var s models.Site
		// Handle potential NULL tracking_id for old records
		var trackingID sql.NullString

		if err := rows.Scan(&s.ID, &s.OrganizationID, &s.Name, &s.URL, &trackingID, &s.CreatedAt); err == nil {
			if trackingID.Valid {
				s.TrackingID = trackingID.String
			}
			sites = append(sites, s)
		}
	}
	return c.JSON(http.StatusOK, sites)
}

func GetSite(c echo.Context) error {
	id := c.Param("id")
	// For production, we should verify the site belongs to the authenticated organization.
	// organizationID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["organization_id"]

	var s models.Site
	var trackingID sql.NullString
	query := `SELECT id, organization_id, name, url, tracking_id, created_at FROM sites WHERE id = $1`

	err := db.GetDB().QueryRow(query, id).Scan(&s.ID, &s.OrganizationID, &s.Name, &s.URL, &trackingID, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Site not found"})
	} else if err != nil {
		c.Logger().Error("Database error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if trackingID.Valid {
		s.TrackingID = trackingID.String
	}

	return c.JSON(http.StatusOK, s)
}

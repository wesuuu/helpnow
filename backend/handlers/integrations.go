package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateIntegration(c echo.Context) error {
	var integration models.Integration
	if err := c.Bind(&integration); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate configuration is valid JSON
	if integration.Configuration != "" {
		if !json.Valid([]byte(integration.Configuration)) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Configuration must be valid JSON"})
		}
	}

	query := `INSERT INTO integrations (organization_id, name, type, configuration) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, integration.OrganizationID, integration.Name, integration.Type, integration.Configuration).Scan(&integration.ID, &integration.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create integration"})
	}

	return c.JSON(http.StatusCreated, integration)
}

func ListIntegrations(c echo.Context) error {
	// TODO: Get OrganizationID from context/auth
	// For now, assume it's passed as a query param or fetch all
	orgID := c.QueryParam("organization_id")

	dbConn := db.GetDB()
	var query string
	var rows *sql.Rows
	var err error

	if orgID != "" {
		query = `SELECT id, organization_id, name, type, configuration, created_at FROM integrations WHERE organization_id = $1`
		rows, err = dbConn.Query(query, orgID)
	} else {
		query = `SELECT id, organization_id, name, type, configuration, created_at FROM integrations`
		rows, err = dbConn.Query(query)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch integrations"})
	}
	defer rows.Close()

	var integrations []models.Integration
	for rows.Next() {
		var i models.Integration
		if err := rows.Scan(&i.ID, &i.OrganizationID, &i.Name, &i.Type, &i.Configuration, &i.CreatedAt); err != nil {
			continue
		}
		integrations = append(integrations, i)
	}

	return c.JSON(http.StatusOK, integrations)
}

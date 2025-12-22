package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func GetOrganization(c echo.Context) error {
	id := c.Param("id")

	var org models.Organization
	var systemPrompt sql.NullString

	query := `SELECT id, name, system_prompt, created_at FROM organizations WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&org.ID, &org.Name, &systemPrompt, &org.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Organization not found"})
	}

	org.SystemPrompt = systemPrompt.String
	return c.JSON(http.StatusOK, org)
}

func UpdateOrganization(c echo.Context) error {
	id := c.Param("id")
	var req models.Organization
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Update only name and system_prompt for now
	query := `UPDATE organizations SET name = $1, system_prompt = $2 WHERE id = $3 returning id, name, system_prompt, created_at`

	var org models.Organization
	var systemPrompt sql.NullString

	err := db.GetDB().QueryRow(query, req.Name, req.SystemPrompt, id).Scan(&org.ID, &org.Name, &systemPrompt, &org.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to update organization: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update organization"})
	}

	org.SystemPrompt = systemPrompt.String
	return c.JSON(http.StatusOK, org)
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateAgent(c echo.Context) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Perform insert
	query := `INSERT INTO agents (organization_id, name, type, model_config) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, agent.OrganizationID, agent.Name, agent.Type, agent.ModelConfig).Scan(&agent.ID, &agent.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create agent"})
	}

	return c.JSON(http.StatusCreated, agent)
}

func ListAgents(c echo.Context) error {
	orgID := c.QueryParam("org_id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_id is required"})
	}

	rows, err := db.GetDB().Query("SELECT id, organization_id, name, type, model_config, created_at FROM agents WHERE organization_id = $1", orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch agents"})
	}
	defer rows.Close()

	agents := []models.Agent{}
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(&a.ID, &a.OrganizationID, &a.Name, &a.Type, &a.ModelConfig, &a.CreatedAt); err != nil {
			continue
		}
		agents = append(agents, a)
	}

	return c.JSON(http.StatusOK, agents)
}

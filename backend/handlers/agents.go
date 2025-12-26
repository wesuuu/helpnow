package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
	"github.com/wesuuu/helpnow/backend/secrets"
)

func CreateAgent(c echo.Context) error {
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Transaction for ID + Secret
	tx, err := db.GetDB().Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer tx.Rollback()

	query := `INSERT INTO agents (organization_id, name, description, model_config) VALUES ($1, $2, $3, 'PENDING_SECRET') RETURNING id, created_at`
	err = tx.QueryRow(query, agent.OrganizationID, agent.Name, agent.Description).Scan(&agent.ID, &agent.CreatedAt)
	if err != nil {
		log.Println("Error creating agent:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create agent"})
	}

	// Store config in Vault
	if secrets.GlobalSecretStore == nil {
		log.Println("Error: Secret Store not configured")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal resource not available"})
	}

	secretData := map[string]interface{}{
		"model_config": agent.ModelConfig,
	}
	path := fmt.Sprintf("agents/%d", agent.ID)
	if err := secrets.GlobalSecretStore.Write(c.Request().Context(), path, secretData); err != nil {
		log.Println("Error writing agent secret:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to protect agent config"})
	}

	// Update DB marker
	_, err = tx.Exec(`UPDATE agents SET model_config = 'SECRET_IN_VAULT' WHERE id = $1`, agent.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to finalize agent creation"})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Transaction failed"})
	}

	// Reset config in response object if we want to mask it, but usually strictly matching input is fine for Create.
	// But List/Get will load it properly.
	return c.JSON(http.StatusCreated, agent)
}

func ListAgents(c echo.Context) error {
	orgID := c.QueryParam("org_id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "org_id is required"})
	}

	rows, err := db.GetDB().Query("SELECT id, organization_id, name, description, model_config, created_at FROM agents WHERE organization_id = $1", orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch agents"})
	}
	defer rows.Close()

	agents := []models.Agent{}
	for rows.Next() {
		var agent models.Agent
		if err := rows.Scan(&agent.ID, &agent.OrganizationID, &agent.Name, &agent.Description, &agent.ModelConfig, &agent.CreatedAt); err == nil {
			if agent.ModelConfig == "SECRET_IN_VAULT" {
				agent.ModelConfig = "" // Mask
			}
			agents = append(agents, agent)
		}
	}
	return c.JSON(http.StatusOK, agents)
}

func GetAgent(c echo.Context) error {
	id := c.Param("id")
	var agent models.Agent
	row := db.GetDB().QueryRow(`SELECT id, organization_id, name, description, model_config, created_at FROM agents WHERE id = $1`, id)
	err := row.Scan(&agent.ID, &agent.OrganizationID, &agent.Name, &agent.Description, &agent.ModelConfig, &agent.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Agent not found"})
	}

	// Fetch secret if needed
	if secrets.GlobalSecretStore != nil && agent.ModelConfig == "SECRET_IN_VAULT" {
		path := fmt.Sprintf("agents/%d", agent.ID)
		data, err := secrets.GlobalSecretStore.Read(c.Request().Context(), path)
		if err != nil {
			log.Println("Error reading agent secret:", err)
			// Don't fail the request, just return empty? or error?
		} else if data != nil {
			if cfg, ok := data["model_config"].(string); ok {
				agent.ModelConfig = cfg
			}
		}
	}

	return c.JSON(http.StatusOK, agent)
}

func UpdateAgent(c echo.Context) error {
	id := c.Param("id")
	var agent models.Agent
	if err := c.Bind(&agent); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 1. Update basic info
	_, err := db.GetDB().Exec(`UPDATE agents SET name = $1, description = $2 WHERE id = $3`,
		agent.Name, agent.Description, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update agent"})
	}

	// 2. Update Secret if provided
	// We assume if ModelConfig is populated it's an update.
	// If the frontend sends back "SECRET_IN_VAULT" we shouldn't overwrite unless we want to keep it.
	// But usually Update sends full state or partial.
	// If it's "SECRET_IN_VAULT", we skip writing.
	if agent.ModelConfig != "" && agent.ModelConfig != "SECRET_IN_VAULT" {
		if secrets.GlobalSecretStore == nil {
			log.Println("Error: Secret Store not configured")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal resource not available"})
		}

		secretData := map[string]interface{}{
			"model_config": agent.ModelConfig,
		}
		path := fmt.Sprintf("agents/%s", id) // id param is string
		if err := secrets.GlobalSecretStore.Write(c.Request().Context(), path, secretData); err != nil {
			log.Println("Error updating agent secret:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update config"})
		}

		// Ensure marker
		db.GetDB().Exec(`UPDATE agents SET model_config = 'SECRET_IN_VAULT' WHERE id = $1`, id)
	}

	return c.JSON(http.StatusOK, agent)
}

func DeleteAgent(c echo.Context) error {
	id := c.Param("id")
	_, err := db.GetDB().Exec("DELETE FROM agents WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete agent"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

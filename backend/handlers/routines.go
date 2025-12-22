package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/clients"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateRoutine(c echo.Context) error {
	var routine models.Routine
	if err := c.Bind(&routine); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	query := `INSERT INTO routines (agent_id, name, description, workflow) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := db.GetDB().QueryRow(query, routine.AgentID, routine.Name, routine.Description, routine.Workflow).Scan(&routine.ID, &routine.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create routine"})
	}

	return c.JSON(http.StatusCreated, routine)
}

func ListRoutines(c echo.Context) error {
	agentID := c.QueryParam("agent_id")
	if agentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "agent_id is required"})
	}

	rows, err := db.GetDB().Query("SELECT id, agent_id, name, description, workflow, created_at FROM routines WHERE agent_id = $1", agentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch routines"})
	}
	defer rows.Close()

	routines := []models.Routine{}
	for rows.Next() {
		var r models.Routine
		if err := rows.Scan(&r.ID, &r.AgentID, &r.Name, &r.Description, &r.Workflow, &r.CreatedAt); err != nil {
			continue
		}
		routines = append(routines, r)
	}

	return c.JSON(http.StatusOK, routines)
}

func ExecuteRoutine(c echo.Context) error {
	// Request should contain routine_id
	type ExecuteRequest struct {
		RoutineID int               `json:"routine_id"`
		Inputs    map[string]string `json:"inputs"`
	}
	var req ExecuteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 1. Fetch Routine metadata to get Agent ID and Config
	var agentID int
	var agentType, modelConfig string
	err := db.GetDB().QueryRow("SELECT a.id, a.type, a.model_config FROM routines r JOIN agents a ON r.agent_id = a.id WHERE r.id = $1", req.RoutineID).Scan(&agentID, &agentType, &modelConfig)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Routine or Agent not found"})
	}

	// 2. Prepare inputs with model config
	inputs := req.Inputs
	if inputs == nil {
		inputs = make(map[string]string)
	}
	inputs["model_config"] = modelConfig
	inputs["agent_type"] = agentType

	// 3. Call AI Service via gRPC
	// TODO: Get real user ID from context/auth
	userID := "1" 

	resp, err := clients.GlobalAIClient.ExecuteRoutine(c.Request().Context(), strconv.Itoa(req.RoutineID), strconv.Itoa(agentID), inputs, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to trigger AI service: " + err.Error()})
	}

	// 3. Create Run record
	var runID int
	query := `INSERT INTO routine_runs (routine_id, user_id, status, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id`
	// Note: We might want to use the ID returned by AI service or map them. For now, we create our own run ID.
	err = db.GetDB().QueryRow(query, req.RoutineID, userID, "PENDING").Scan(&runID)
	
	if err != nil {
		// Log error but we already triggered the AI...
		fmt.Println("Error recording run:", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Routine execution started",
		"run_id":  runID,
		"ai_execution_id": resp.ExecutionId,
		"initial_status": resp.Status,
	})
}

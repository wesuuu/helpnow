package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/workflows"
)

// ListActions returns all available workflow actions with their schemas
func ListActions(c echo.Context) error {
	schemas := workflows.ListActionSchemas()
	return c.JSON(http.StatusOK, schemas)
}

// GetAction returns the schema for a specific action
func GetAction(c echo.Context) error {
	name := c.Param("name")

	schema, ok := workflows.GetActionSchema(name)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Action not found"})
	}

	return c.JSON(http.StatusOK, schema)
}

// ListLogic returns all available workflow logic with their schemas
func ListLogic(c echo.Context) error {
	schemas := workflows.ListLogicSchemas()
	return c.JSON(http.StatusOK, schemas)
}

// GetLogic returns the schema for a specific logic
func GetLogic(c echo.Context) error {
	name := c.Param("name")

	schema, ok := workflows.GetLogicSchema(name)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Logic not found"})
	}

	return c.JSON(http.StatusOK, schema)
}

// ListTriggers returns all available workflow triggers with their schemas
func ListTriggers(c echo.Context) error {
	schemas := workflows.ListTriggerSchemas()
	return c.JSON(http.StatusOK, schemas)
}

// GetTrigger returns the schema for a specific trigger
func GetTrigger(c echo.Context) error {
	name := c.Param("name")

	schema, ok := workflows.GetTriggerSchema(name)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Trigger not found"})
	}

	return c.JSON(http.StatusOK, schema)
}

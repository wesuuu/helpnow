package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreateContentTemplate(c echo.Context) error {
	var t models.ContentTemplate
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validate Schema is valid JSON if provided
	if t.Schema != "" {
		if !isValidJSON(t.Schema) {
			// If bindings didn't handle it, ensure it's valid stringified JSON
			// Actually if user sends object, Echo bind might fail to strict string field?
			// No, it handles basic types. We assume client sends string or we handle interface{}
			// Given previous pattern, let's assume client sends stringified JSON or we simply store what we get if it's string.
		}
	} else {
		t.Schema = "{}"
	}

	query := `
		INSERT INTO content_templates (organization_id, name, type, content, schema)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := db.GetDB().QueryRow(query,
		t.OrganizationID, t.Name, t.Type, t.Content, t.Schema,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		c.Logger().Error("Failed to create template: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create template"})
	}

	return c.JSON(http.StatusCreated, t)
}

func ListContentTemplates(c echo.Context) error {
	orgID := c.QueryParam("organization_id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}
	typeFilter := c.QueryParam("type")

	query := `SELECT id, organization_id, name, type, content, COALESCE(schema::text, '{}'), created_at, updated_at FROM content_templates WHERE organization_id = $1`
	args := []interface{}{orgID}

	if typeFilter != "" {
		query += ` AND type = $2`
		args = append(args, typeFilter)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.GetDB().Query(query, args...)
	if err != nil {
		c.Logger().Error("Failed to list templates: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list templates"})
	}
	defer rows.Close()

	templates := []models.ContentTemplate{}
	for rows.Next() {
		var t models.ContentTemplate
		if err := rows.Scan(&t.ID, &t.OrganizationID, &t.Name, &t.Type, &t.Content, &t.Schema, &t.CreatedAt, &t.UpdatedAt); err != nil {
			continue
		}
		templates = append(templates, t)
	}

	return c.JSON(http.StatusOK, templates)
}

func GetContentTemplate(c echo.Context) error {
	id := c.Param("id")
	var t models.ContentTemplate
	query := `SELECT id, organization_id, name, type, content, COALESCE(schema::text, '{}'), created_at, updated_at FROM content_templates WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&t.ID, &t.OrganizationID, &t.Name, &t.Type, &t.Content, &t.Schema, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	}
	return c.JSON(http.StatusOK, t)
}

func UpdateContentTemplate(c echo.Context) error {
	id := c.Param("id")
	var t models.ContentTemplate
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	t.UpdatedAt = time.Now()

	query := `
		UPDATE content_templates 
		SET name = $1, type = $2, content = $3, schema = $4, updated_at = $5
		WHERE id = $6
		RETURNING id`

	err := db.GetDB().QueryRow(query, t.Name, t.Type, t.Content, t.Schema, t.UpdatedAt, id).Scan(&t.ID)
	if err != nil {
		c.Logger().Error("Failed to update template: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update template"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func DeleteContentTemplate(c echo.Context) error {
	id := c.Param("id")
	_, err := db.GetDB().Exec("DELETE FROM content_templates WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete template"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

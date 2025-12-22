package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
)

type EmailTemplate struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ListEmailTemplates(c echo.Context) error {
	orgIDStr := c.QueryParam("organization_id")
	if orgIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}

	rows, err := db.GetDB().Query("SELECT id, organization_id, name, subject, body, created_at, updated_at FROM email_templates WHERE organization_id = $1 ORDER BY updated_at DESC", orgIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	templates := []EmailTemplate{}
	for rows.Next() {
		var tmpl EmailTemplate
		if err := rows.Scan(&tmpl.ID, &tmpl.OrganizationID, &tmpl.Name, &tmpl.Subject, &tmpl.Body, &tmpl.CreatedAt, &tmpl.UpdatedAt); err != nil {
			continue
		}
		templates = append(templates, tmpl)
	}

	return c.JSON(http.StatusOK, templates)
}

func CreateEmailTemplate(c echo.Context) error {
	var tmpl EmailTemplate
	if err := c.Bind(&tmpl); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if tmpl.Name == "" || tmpl.OrganizationID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and organization_id are required"})
	}

	err := db.GetDB().QueryRow(`
		INSERT INTO email_templates (organization_id, name, subject, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		tmpl.OrganizationID, tmpl.Name, tmpl.Subject, tmpl.Body).Scan(&tmpl.ID, &tmpl.CreatedAt, &tmpl.UpdatedAt)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, tmpl)
}

func UpdateEmailTemplate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var tmpl EmailTemplate
	if err := c.Bind(&tmpl); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	_, err = db.GetDB().Exec(`
		UPDATE email_templates 
		SET name = $1, subject = $2, body = $3, updated_at = NOW()
		WHERE id = $4`,
		tmpl.Name, tmpl.Subject, tmpl.Body, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Fetch updated
	db.GetDB().QueryRow("SELECT created_at, updated_at FROM email_templates WHERE id = $1", id).Scan(&tmpl.CreatedAt, &tmpl.UpdatedAt)
	tmpl.ID = id

	return c.JSON(http.StatusOK, tmpl)
}

func DeleteEmailTemplate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	_, err = db.GetDB().Exec("DELETE FROM email_templates WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

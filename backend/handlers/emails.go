package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
)

type EmailDomain struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Domain         string    `json:"domain"`
	DKIMRecord     *string   `json:"dkim_record"`
	SPFRecord      *string   `json:"spf_record"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
}

func ListEmailDomains(c echo.Context) error {
	orgIDStr := c.QueryParam("organization_id")
	if orgIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}

	rows, err := db.GetDB().Query("SELECT id, organization_id, domain, dkim_record, spf_record, is_verified, created_at FROM email_domains WHERE organization_id = $1 ORDER BY created_at DESC", orgIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	domains := []EmailDomain{}
	for rows.Next() {
		var domain EmailDomain
		if err := rows.Scan(&domain.ID, &domain.OrganizationID, &domain.Domain, &domain.DKIMRecord, &domain.SPFRecord, &domain.IsVerified, &domain.CreatedAt); err != nil {
			continue
		}
		domains = append(domains, domain)
	}

	return c.JSON(http.StatusOK, domains)
}

func CreateEmailDomain(c echo.Context) error {
	var domain EmailDomain
	if err := c.Bind(&domain); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if domain.Domain == "" || domain.OrganizationID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "domain and organization_id are required"})
	}

	// Generate Mock Records
	dkim := fmt.Sprintf("v=DKIM1; k=rsa; p=%s-mock-public-key", domain.Domain)
	spf := "v=spf1 include:_spf.helpnow.ai ~all"

	err := db.GetDB().QueryRow(`
		INSERT INTO email_domains (organization_id, domain, dkim_record, spf_record)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, is_verified`,
		domain.OrganizationID, domain.Domain, dkim, spf).Scan(&domain.ID, &domain.CreatedAt, &domain.IsVerified)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	domain.DKIMRecord = &dkim
	domain.SPFRecord = &spf

	return c.JSON(http.StatusCreated, domain)
}

func VerifyEmailDomain(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Mock Verification Logic
	// In real world, do net.LookupTXT()
	isVerified := true

	_, err = db.GetDB().Exec("UPDATE email_domains SET is_verified = $1 WHERE id = $2", isVerified, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"is_verified": isVerified})
}

func DeleteEmailDomain(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	_, err = db.GetDB().Exec("DELETE FROM email_domains WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

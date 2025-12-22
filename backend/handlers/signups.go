package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func CreateSignupCampaign(c echo.Context) error {
	var campaign models.SignupCampaign
	if err := c.Bind(&campaign); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	campaign.Token = generateToken()

	query := `INSERT INTO signup_campaigns (organization_id, name, target_audience_id, token) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := db.GetDB().QueryRow(query, campaign.OrganizationID, campaign.Name, campaign.TargetAudienceID, campaign.Token).Scan(&campaign.ID, &campaign.CreatedAt)
	if err != nil {
		c.Logger().Error("Failed to create signup campaign: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create signup campaign"})
	}

	return c.JSON(http.StatusCreated, campaign)
}

func ListSignupCampaigns(c echo.Context) error {
	orgID := c.QueryParam("organization_id")

	rows, err := db.GetDB().Query(`SELECT id, organization_id, name, target_audience_id, token, created_at FROM signup_campaigns WHERE organization_id = $1`, orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list campaigns"})
	}
	defer rows.Close()

	campaigns := []models.SignupCampaign{}
	for rows.Next() {
		var cc models.SignupCampaign
		if err := rows.Scan(&cc.ID, &cc.OrganizationID, &cc.Name, &cc.TargetAudienceID, &cc.Token, &cc.CreatedAt); err == nil {
			campaigns = append(campaigns, cc)
		}
	}
	return c.JSON(http.StatusOK, campaigns)
}

// Public Endpoint
func CaptureLead(c echo.Context) error {
	token := c.Param("token")

	type LeadInput struct {
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		Age       int    `json:"age"`
		Gender    string `json:"gender"`
		Ethnicity string `json:"ethnicity"`
		Location  string `json:"location"`
	}

	var input LeadInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 1. Validate Campaign
	var campaignID, orgID, audienceID int
	err := db.GetDB().QueryRow("SELECT id, organization_id, target_audience_id FROM signup_campaigns WHERE token = $1", token).Scan(&campaignID, &orgID, &audienceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid campaign token"})
	}

	// 2. Upsert Person
	// Simple assumption: Email matches unique person.
	var personID int

	// Try to find existing
	err = db.GetDB().QueryRow("SELECT id FROM people WHERE organization_id = $1 AND email = $2", orgID, input.Email).Scan(&personID)

	if err != nil {
		// Create new
		err = db.GetDB().QueryRow(
			`INSERT INTO people (organization_id, full_name, email, age, gender, ethnicity, location, last_interaction_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`,
			orgID, input.FullName, input.Email, input.Age, input.Gender, input.Ethnicity, input.Location,
		).Scan(&personID)
		if err != nil {
			c.Logger().Error("Failed to create person from lead: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save lead"})
		}
	} else {
		// Update existing (optional, but good for lead capture to refresh data)
		// For MVP, we skip updating fields to avoid overwriting good data with partial data
	}

	// 3. Add to Audience
	_, err = db.GetDB().Exec("INSERT INTO audience_memberships (audience_id, person_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", audienceID, personID)
	if err != nil {
		c.Logger().Error("Failed to add lead to audience: ", err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success", "message": "Lead captured"})
}

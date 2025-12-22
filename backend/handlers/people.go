package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func CreatePerson(c echo.Context) error {
	var person models.Person
	if err := c.Bind(&person); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	query := `
		INSERT INTO people (organization_id, full_name, email, age, ethnicity, gender, location, last_interaction_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err := db.GetDB().QueryRow(query,
		person.OrganizationID, person.FullName, person.Email,
		person.Age, person.Ethnicity, person.Gender, person.Location, person.LastInteractionAt,
	).Scan(&person.ID, &person.CreatedAt)

	if err != nil {
		c.Logger().Error("Failed to create person: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create person"})
	}

	return c.JSON(http.StatusCreated, person)
}

func ListPeople(c echo.Context) error {
	orgID := c.QueryParam("organization_id")
	if orgID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "organization_id is required"})
	}

	// Optional filter by audience
	audienceID := c.QueryParam("audience_id")

	query := `SELECT id, organization_id, full_name, email, age, ethnicity, gender, location, last_interaction_at, created_at FROM people WHERE organization_id = $1`
	args := []interface{}{orgID}

	if audienceID != "" {
		query = `SELECT p.id, p.organization_id, p.full_name, p.email, p.age, p.ethnicity, p.gender, p.location, p.last_interaction_at, p.created_at 
				 FROM people p
				 JOIN audience_memberships am ON p.id = am.person_id
				 WHERE p.organization_id = $1 AND am.audience_id = $2`
		args = append(args, audienceID)
	}

	rows, err := db.GetDB().Query(query, args...)
	if err != nil {
		c.Logger().Error("Failed to fetch people: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch people"})
	}
	defer rows.Close()

	people := []models.Person{}
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.FullName, &p.Email, &p.Age, &p.Ethnicity, &p.Gender, &p.Location, &p.LastInteractionAt, &p.CreatedAt); err != nil {
			continue
		}
		people = append(people, p)
	}

	return c.JSON(http.StatusOK, people)
}

func AddPersonToAudience(c echo.Context) error {
	audienceID := c.Param("id")
	type Request struct {
		PersonID int `json:"person_id"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	_, err := db.GetDB().Exec("INSERT INTO audience_memberships (audience_id, person_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", audienceID, req.PersonID)
	if err != nil {
		c.Logger().Error("Failed to add person to audience: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add person to audience"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "added"})
}

func GetAudienceMembers(c echo.Context) error {
	// Re-using logic via ListPeople with audience_id param, but exposed as specific route
	// Or we can just use ListPeople? Let's just use ListPeople for now on frontend?
	// Actually, following RESTful pattern:
	c.QueryParams().Set("audience_id", c.Param("id"))
	return ListPeople(c)
}

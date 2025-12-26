package handlers

import (
	"encoding/json"
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
		INSERT INTO people (organization_id, full_name, email, age, ethnicity, gender, location, last_interaction_at, score)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := db.GetDB().QueryRow(query,
		person.OrganizationID, person.FullName, person.Email,
		person.Age, person.Ethnicity, person.Gender, person.Location, person.LastInteractionAt, person.Score,
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

	query := `SELECT id, organization_id, full_name, email, age, ethnicity, gender, location, last_interaction_at, score, COALESCE(event_history::text, '[]'), created_at FROM people WHERE organization_id = $1`
	args := []interface{}{orgID}

	if audienceID != "" {
		query = `SELECT p.id, p.organization_id, p.full_name, p.email, p.age, p.ethnicity, p.gender, p.location, p.last_interaction_at, p.score, COALESCE(p.event_history::text, '[]'), p.created_at 
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
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.FullName, &p.Email, &p.Age, &p.Ethnicity, &p.Gender, &p.Location, &p.LastInteractionAt, &p.Score, &p.EventHistory, &p.CreatedAt); err != nil {
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

func AppendPersonEvent(c echo.Context) error {
	id := c.Param("id")
	type Request struct {
		Event interface{} `json:"event"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	eventJSON, err := json.Marshal(req.Event)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event structure"})
	}

	// Append jsonb to event_history array. We wrap the single event in an array '[]' -> to_jsonb -> concatenation?
	// Postgres: event_history || '[eventJSON]'::jsonb
	// Actually simpler: passing the json string of the object, casting to jsonb, and || operator works if the existing is an array.
	query := `UPDATE people SET event_history = COALESCE(event_history, '[]'::jsonb) || $1::jsonb WHERE id = $2`

	// We need to pass the event as a single-item array string: "[{...}]" so that || appends it to the array.
	// Otherwise "array || object" might not be standard behavior depending on PG version (usually works, but safer to append array to array).
	// Let's try appending array to array.
	eventArrayJSON := "[" + string(eventJSON) + "]"

	_, err = db.GetDB().Exec(query, eventArrayJSON, id)
	if err != nil {
		c.Logger().Error("Failed to append person event: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update history"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func UpdatePersonHistory(c echo.Context) error {
	id := c.Param("id")
	type Request struct {
		History interface{} `json:"event_history"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	historyJSON, err := json.Marshal(req.History)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid history structure"})
	}

	query := `UPDATE people SET event_history = $1 WHERE id = $2`
	_, err = db.GetDB().Exec(query, string(historyJSON), id)
	if err != nil {
		c.Logger().Error("Failed to update person history: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update history"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

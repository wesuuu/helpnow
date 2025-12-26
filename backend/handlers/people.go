package handlers

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func getIP(c echo.Context) string {
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request().RemoteAddr
	}
	return ip
}

func getUserAgent(c echo.Context) string {
	return c.Request().UserAgent()
}

func updateMeta(existing *models.PersonMeta, c echo.Context) {
	ip := getIP(c)
	ua := getUserAgent(c)
	now := time.Now()

	// Update IP Addresses
	foundIP := false
	for i := range existing.IPAddresses {
		if existing.IPAddresses[i].Address == ip {
			existing.IPAddresses[i].Used = now
			existing.IPAddresses[i].Count++
			existing.IPAddresses[i].UserAgent = ua // Update latest UA for this IP
			foundIP = true
			break
		}
	}
	if !foundIP {
		existing.IPAddresses = append(existing.IPAddresses, models.IPAddress{
			Address:   ip,
			Used:      now,
			Count:     1,
			UserAgent: ua,
			// Location: to be populated using a geoip service ideally, handled by frontend or separate service?
			// For now leaving empty or simple
		})
	}

	// Update User Agents list
	foundUA := false
	for _, a := range existing.UserAgents {
		if a == ua {
			foundUA = true
			break
		}
	}
	if !foundUA {
		existing.UserAgents = append(existing.UserAgents, ua)
	}

	// Update Devices (naive implementation based on UA string)
	// In reality this needs a UA parser library. For now, adding the UA string itself to devices if unique?
	// Or simplistic check. User said "list of devices used".
	// Let's just store the UA string in Devices if it's not there, or assume "Device" == UA string if undefined.
	// Actually user separated devices and user_agent.
	// I'll leave devices empty or just copy UA for now unless I add a parser.
}

func CreatePerson(c echo.Context) error {
	var person models.Person
	if err := c.Bind(&person); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Initialize Meta
	person.Meta = models.PersonMeta{
		IPAddresses: []models.IPAddress{},
		Devices:     []string{},
		UserAgents:  []string{},
	}
	updateMeta(&person.Meta, c)

	metaJSON, _ := json.Marshal(person.Meta)

	query := `
		INSERT INTO people (organization_id, first_name, last_name, email, age, ethnicity, gender, location, last_interaction_at, score, interests, meta)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	err := db.GetDB().QueryRow(query,
		person.OrganizationID, person.FirstName, person.LastName, person.Email,
		person.Age, person.Ethnicity, person.Gender, person.Location, person.LastInteractionAt, person.Score, pq.Array(person.Interests), metaJSON,
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

	query := `SELECT id, organization_id, first_name, last_name, email, age, ethnicity, gender, location, last_interaction_at, score, interests, COALESCE(meta::text, '{}'), created_at FROM people WHERE organization_id = $1`
	args := []interface{}{orgID}

	if audienceID != "" {
		query = `SELECT p.id, p.organization_id, p.first_name, p.last_name, p.email, p.age, p.ethnicity, p.gender, p.location, p.last_interaction_at, p.score, p.interests, COALESCE(p.meta::text, '{}'), p.created_at 
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
		var metaString string
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.FirstName, &p.LastName, &p.Email, &p.Age, &p.Ethnicity, &p.Gender, &p.Location, &p.LastInteractionAt, &p.Score, pq.Array(&p.Interests), &metaString, &p.CreatedAt); err != nil {
			continue
		}
		if metaString != "" {
			json.Unmarshal([]byte(metaString), &p.Meta)
		}
		people = append(people, p)
	}

	return c.JSON(http.StatusOK, people)
}

func GetPerson(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var p models.Person
	var metaString string
	query := `SELECT id, organization_id, first_name, last_name, email, age, ethnicity, gender, location, last_interaction_at, score, interests, COALESCE(meta::text, '{}'), created_at FROM people WHERE id = $1`

	err := db.GetDB().QueryRow(query, id).Scan(&p.ID, &p.OrganizationID, &p.FirstName, &p.LastName, &p.Email, &p.Age, &p.Ethnicity, &p.Gender, &p.Location, &p.LastInteractionAt, &p.Score, pq.Array(&p.Interests), &metaString, &p.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
	}
	if metaString != "" {
		json.Unmarshal([]byte(metaString), &p.Meta)
	}

	// Fetch Events
	eventRows, err := db.GetDB().Query("SELECT id, person_id, event, created_at FROM person_events WHERE person_id = $1 ORDER BY created_at DESC", p.ID)
	if err == nil {
		defer eventRows.Close()
		p.Events = []models.PersonEvent{}
		for eventRows.Next() {
			var evt models.PersonEvent
			var eventJSON string
			if err := eventRows.Scan(&evt.ID, &evt.PersonID, &eventJSON, &evt.CreatedAt); err == nil {
				json.Unmarshal([]byte(eventJSON), &evt.Event)
				p.Events = append(p.Events, evt)
			}
		}
	}

	return c.JSON(http.StatusOK, p)
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

	// Insert into person_events table
	query := `INSERT INTO person_events (person_id, event) VALUES ($1, $2)`
	_, err = db.GetDB().Exec(query, id, string(eventJSON))
	if err != nil {
		c.Logger().Error("Failed to append person event: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update history"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

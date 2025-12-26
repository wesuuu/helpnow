package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/models"
)

func Register(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// TODO: Hash password properly
	// TODO: Check if org exists or create one

	query := `INSERT INTO users (email, password_hash, first_name, middle_name, last_name, organization_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, user.Email, user.PasswordHash, user.FirstName, user.MiddleName, user.LastName, user.OrganizationID).Scan(&user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// TODO: Validate password and generate JWT
	// For now, fetch user to get details (and pretend password is valid)
	var user models.User
	// We only need basic details for the session
	query := `SELECT id, organization_id, email, first_name, last_name FROM users WHERE email = $1`
	err := db.GetDB().QueryRow(query, req.Email).Scan(&user.ID, &user.OrganizationID, &user.Email, &user.FirstName, &user.LastName)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": "mock-jwt-token",
		"user":  user,
	})
}

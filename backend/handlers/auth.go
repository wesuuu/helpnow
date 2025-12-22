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

	query := `INSERT INTO users (email, password_hash, full_name, organization_id) VALUES ($1, $2, $3, $4) RETURNING id`
	dbConn := db.GetDB()
	err := dbConn.QueryRow(query, user.Email, user.PasswordHash, user.FullName, user.OrganizationID).Scan(&user.ID)
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
	// For now, return a mock token

	return c.JSON(http.StatusOK, map[string]string{"token": "mock-jwt-token", "user_email": req.Email})
}

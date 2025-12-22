package outbound

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/wesuuu/helpnow/backend/db"
)

// Email represents an email to be sent.
type Email struct {
	To      string
	From    string // Sender email address (e.g., "marketing@mybrand.com")
	Subject string
	Body    string
}

// SendEmail attempts to send an email.
// It verifies that the sender domain is allowed (either "helpnow.ai" or a verified domain).
func SendEmail(ctx context.Context, email Email) error {
	// 1. Extract domain from "From" address
	parts := strings.Split(email.From, "@")
	if len(parts) != 2 {
		return errors.New("invalid 'From' email address format")
	}
	domain := parts[1]

	// 2. Verify Domain Authorization
	if strings.ToLower(domain) == "helpnow.ai" {
		// System domain is always allowed
		log.Printf("[Outbound] Allowed system domain: %s", domain)
	} else {
		// Check against verified user domains
		verified, err := isDomainVerified(ctx, domain)
		if err != nil {
			return fmt.Errorf("failed to check domain verification: %w", err)
		}
		if !verified {
			return fmt.Errorf("sender domain '%s' is not verified", domain)
		}
		log.Printf("[Outbound] Verified user domain: %s", domain)
	}

	// 3. Mock Send Logic (Replace with actual SMTP/API later)
	log.Printf("--- SENDING EMAIL ---\nTo: %s\nFrom: %s\nSubject: %s\nBody: %s\n---------------------",
		email.To, email.From, email.Subject, email.Body)

	return nil
}

// isDomainVerified checks the database for a verified domain.
func isDomainVerified(ctx context.Context, domain string) (bool, error) {
	conn := db.GetDB()
	if conn == nil {
		return false, errors.New("database connection not available")
	}

	var isVerified bool
	// Case-insensitive check for the domain where is_verified = true
	query := `SELECT is_verified FROM email_domains WHERE LOWER(domain) = LOWER($1) AND is_verified = true`
	err := conn.QueryRowContext(ctx, query, domain).Scan(&isVerified)

	if err == sql.ErrNoRows {
		// Domain not found or not verified
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return isVerified, nil
}

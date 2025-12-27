package actions

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/outbound"
	"github.com/wesuuu/helpnow/backend/workflows"
)

func init() {
	workflows.RegisterAction("Send Email", &SendEmailAction{})
}

// SendEmailAction sends an email using a template
type SendEmailAction struct {
	TemplateID int `json:"template_id" validate:"required,min=1" desc:"ID of the email template to use. Must reference an existing template in the email_templates table."`
}

func (a *SendEmailAction) Execute(ctx context.Context, contextData map[string]interface{}) (output string, err error) {
	// Fetch Template from database
	var subject, body string
	row := db.GetDB().QueryRow("SELECT subject, body FROM email_templates WHERE id = $1", a.TemplateID)
	if err := row.Scan(&subject, &body); err != nil {
		return "Failed to fetch template", fmt.Errorf("failed to fetch template %d: %w", a.TemplateID, err)
	}

	// Extract Recipient from context
	recipient := ""
	if email, ok := contextData["email"].(string); ok {
		recipient = email
	} else if _, ok := contextData["person_id"].(float64); ok {
		// MVP: Not implemented - would fetch person's email from DB
	}

	if recipient == "" {
		return "No recipient email in context", errors.New("no recipient email found in context")
	}

	// Send Email
	sendErr := outbound.SendEmail(ctx, outbound.Email{
		To:      recipient,
		From:    "notifications@helpnow.ai",
		Subject: subject,
		Body:    body,
	})
	if sendErr != nil {
		return "Failed to send email: " + sendErr.Error(), sendErr
	}

	return fmt.Sprintf("Email sent to %s (Template %d)", recipient, a.TemplateID), nil
}

// Helper function to parse template_id from various types (for backward compatibility during migration)
func parseTemplateID(val interface{}) int {
	switch v := val.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			return id
		}
	}
	return 0
}

package actions

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/outbound"
)

type SendEmailAction struct{}

func (a *SendEmailAction) Execute(ctx context.Context, actionCtx ActionContext) (output string, err error) {
	// Extract template ID
	var templateID int
	if tid, ok := actionCtx.NodeData["template_id"]; ok {
		switch v := tid.(type) {
		case float64:
			templateID = int(v)
		case int:
			templateID = v
		case string:
			if id, err := strconv.Atoi(v); err == nil {
				templateID = id
			}
		}
	}

	if templateID <= 0 {
		return "No template selected", nil // Not strictly an error, just no op? Or maybe error? Worker code returned output "No template selected" and success status.
	}

	// Fetch Template
	var subject, body string
	row := db.GetDB().QueryRow("SELECT subject, body FROM email_templates WHERE id = $1", templateID)
	if err := row.Scan(&subject, &body); err != nil {
		return "Failed to fetch template", fmt.Errorf("failed to fetch template %d: %w", templateID, err)
	}

	// Extract Recipient
	recipient := ""
	if email, ok := actionCtx.ContextData["email"].(string); ok {
		recipient = email
	} else if _, ok := actionCtx.ContextData["person_id"].(float64); ok {
		// MVP: Not implemented in worker yet
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

	return fmt.Sprintf("Email sent to %s (Template %d)", recipient, templateID), nil
}

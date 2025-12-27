package workflows

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateActionNode validates an action node's properties
func ValidateActionNode(actionType string, properties map[string]interface{}) error {
	// Get the registered action template
	actionTemplate, ok := GetAction(actionType)
	if !ok {
		return fmt.Errorf("unknown action type: %s", actionType)
	}

	// Create a new instance
	action := reflect.New(reflect.TypeOf(actionTemplate).Elem()).Interface()

	// Unmarshal properties into the action
	propBytes, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	if err := json.Unmarshal(propBytes, action); err != nil {
		return fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	// Validate using struct tags
	if err := validate.Struct(action); err != nil {
		return formatValidationError(err, actionType)
	}

	return nil
}

// ValidateLogicNode validates a logic node's properties
func ValidateLogicNode(logicType string, properties map[string]interface{}) error {
	// Get the registered logic template
	logicTemplate, ok := GetLogic(logicType)
	if !ok {
		return fmt.Errorf("unknown logic type: %s", logicType)
	}

	// Create a new instance
	logic := reflect.New(reflect.TypeOf(logicTemplate).Elem()).Interface()

	// Unmarshal properties into the logic
	propBytes, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	if err := json.Unmarshal(propBytes, logic); err != nil {
		return fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	// Validate using struct tags
	if err := validate.Struct(logic); err != nil {
		return formatValidationError(err, logicType)
	}

	return nil
}

// ValidateTriggerNode validates a trigger node's properties
func ValidateTriggerNode(triggerType string, properties map[string]interface{}) error {
	// Get the registered trigger template
	triggerTemplate, ok := GetTrigger(triggerType)
	if !ok {
		return fmt.Errorf("unknown trigger type: %s", triggerType)
	}

	// Create a new instance
	trigger := reflect.New(reflect.TypeOf(triggerTemplate).Elem()).Interface()

	// Unmarshal properties into the trigger
	propBytes, err := json.Marshal(properties)
	if err != nil {
		return fmt.Errorf("failed to marshal properties: %w", err)
	}

	if err := json.Unmarshal(propBytes, trigger); err != nil {
		return fmt.Errorf("failed to unmarshal properties: %w", err)
	}

	// Validate using struct tags
	if err := validate.Struct(trigger); err != nil {
		return formatValidationError(err, triggerType)
	}

	return nil
}

// ValidateWorkflowGraph validates all nodes in a workflow graph
func ValidateWorkflowGraph(graph Graph) error {
	for _, node := range graph.Nodes {
		switch node.Type {
		case "ACTION":
			actionType, ok := node.Properties["action"].(string)
			if !ok {
				return fmt.Errorf("node %s: missing 'action' property", node.ID)
			}
			if err := ValidateActionNode(actionType, node.Properties); err != nil {
				return fmt.Errorf("node %s (%s): %w", node.ID, actionType, err)
			}

		case "CONDITION":
			logicType := "Condition" // Default logic type
			if err := ValidateLogicNode(logicType, node.Properties); err != nil {
				return fmt.Errorf("node %s (Condition): %w", node.ID, err)
			}

		case "TRIGGER":
			triggerType, ok := node.Properties["trigger_type"].(string)
			if !ok {
				triggerType = "EVENT" // Default
			}
			if err := ValidateTriggerNode(triggerType, node.Properties); err != nil {
				return fmt.Errorf("node %s (%s): %w", node.ID, triggerType, err)
			}
		}
	}

	return nil
}

// formatValidationError converts validator errors to user-friendly messages
func formatValidationError(err error, nodeType string) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, fmt.Sprintf("%s: %s", e.Field(), getErrorMessage(e)))
		}
		return fmt.Errorf("validation failed for %s: %v", nodeType, messages)
	}
	return err
}

// getErrorMessage returns a user-friendly error message for a validation error
func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", e.Param())
	case "email":
		return "must be a valid email"
	case "url":
		return "must be a valid URL"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	default:
		return fmt.Sprintf("failed validation: %s", e.Tag())
	}
}

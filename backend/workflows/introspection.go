package workflows

import (
	"reflect"
	"strings"
)

// FieldSchema describes a field in an action/logic/trigger
type FieldSchema struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Description string   `json:"description,omitempty"`
	Validations []string `json:"validations,omitempty"`
}

// ComponentSchema describes an action, logic, or trigger
type ComponentSchema struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"` // "action", "logic", or "trigger"
	Description string        `json:"description,omitempty"`
	Fields      []FieldSchema `json:"fields"`
}

// GetActionSchema returns the schema for a registered action
func GetActionSchema(name string) (*ComponentSchema, bool) {
	action, ok := GetAction(name)
	if !ok {
		return nil, false
	}

	return buildComponentSchema(name, "action", action), true
}

// GetLogicSchema returns the schema for a registered logic
func GetLogicSchema(name string) (*ComponentSchema, bool) {
	logic, ok := GetLogic(name)
	if !ok {
		return nil, false
	}

	return buildComponentSchema(name, "logic", logic), true
}

// GetTriggerSchema returns the schema for a registered trigger
func GetTriggerSchema(name string) (*ComponentSchema, bool) {
	trigger, ok := GetTrigger(name)
	if !ok {
		return nil, false
	}

	return buildComponentSchema(name, "trigger", trigger), true
}

// ListActionSchemas returns schemas for all registered actions
func ListActionSchemas() []ComponentSchema {
	mu.RLock()
	defer mu.RUnlock()

	schemas := make([]ComponentSchema, 0, len(actionRegistry))
	for name, action := range actionRegistry {
		schemas = append(schemas, *buildComponentSchema(name, "action", action))
	}
	return schemas
}

// ListLogicSchemas returns schemas for all registered logic
func ListLogicSchemas() []ComponentSchema {
	mu.RLock()
	defer mu.RUnlock()

	schemas := make([]ComponentSchema, 0, len(logicRegistry))
	for name, logic := range logicRegistry {
		schemas = append(schemas, *buildComponentSchema(name, "logic", logic))
	}
	return schemas
}

// ListTriggerSchemas returns schemas for all registered triggers
func ListTriggerSchemas() []ComponentSchema {
	mu.RLock()
	defer mu.RUnlock()

	schemas := make([]ComponentSchema, 0, len(triggerRegistry))
	for name, trigger := range triggerRegistry {
		schemas = append(schemas, *buildComponentSchema(name, "trigger", trigger))
	}
	return schemas
}

// buildComponentSchema builds a schema from a component using reflection
func buildComponentSchema(name, componentType string, component interface{}) *ComponentSchema {
	schema := &ComponentSchema{
		Name:   name,
		Type:   componentType,
		Fields: []FieldSchema{},
	}

	// Use reflection to inspect the struct
	t := reflect.TypeOf(component)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldSchema := FieldSchema{
			Name:        getJSONFieldName(field),
			Type:        field.Type.String(),
			Validations: []string{},
		}

		// Parse validation tags
		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			validations := parseValidationTag(validateTag)
			fieldSchema.Validations = validations
			fieldSchema.Required = contains(validations, "required")
		}

		// Parse description from comment (if we had it)
		// For now, we can add descriptions manually or via a custom tag
		if descTag := field.Tag.Get("desc"); descTag != "" {
			fieldSchema.Description = descTag
		}

		schema.Fields = append(schema.Fields, fieldSchema)
	}

	return schema
}

// getJSONFieldName extracts the JSON field name from struct tags
func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}

	// Handle "field_name,omitempty" format
	parts := strings.Split(jsonTag, ",")
	if parts[0] == "-" {
		return field.Name
	}
	return parts[0]
}

// parseValidationTag parses a validation tag into individual rules
func parseValidationTag(tag string) []string {
	rules := strings.Split(tag, ",")
	result := []string{}

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule != "" && rule != "omitempty" {
			result = append(result, rule)
		}
	}

	return result
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetValidationErrorMessage returns a user-friendly error message for a validation rule
func GetValidationErrorMessage(rule string) string {
	parts := strings.Split(rule, "=")
	tag := parts[0]

	switch tag {
	case "required":
		return "This field is required"
	case "min":
		if len(parts) > 1 {
			return "Must be at least " + parts[1]
		}
		return "Must meet minimum requirement"
	case "max":
		if len(parts) > 1 {
			return "Must be at most " + parts[1]
		}
		return "Must not exceed maximum"
	case "email":
		return "Must be a valid email address"
	case "url":
		return "Must be a valid URL"
	case "oneof":
		if len(parts) > 1 {
			return "Must be one of: " + strings.ReplaceAll(parts[1], " ", ", ")
		}
		return "Must be a valid option"
	default:
		return "Must satisfy: " + rule
	}
}

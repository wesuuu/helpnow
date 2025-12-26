package models

import (
	"time"
)

type Organization struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SystemPrompt string    `json:"system_prompt"`
	CreatedAt    time.Time `json:"created_at"`
}

type Team struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

type CampaignType string

const (
	CampaignTypeBrandAwareness CampaignType = "brand_awareness"
	CampaignTypeEmail          CampaignType = "email"
)

type User struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"` // Optional: Users might belong to org directly or via teams
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	FirstName      *string   `json:"first_name"`
	MiddleName     *string   `json:"middle_name"`
	LastName       *string   `json:"last_name"`
	CreatedAt      time.Time `json:"created_at"`
}

type TeamMember struct {
	TeamID int  `json:"team_id"`
	UserID int  `json:"user_id"`
	Role   Role `json:"role"`
}

type Agent struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ModelConfig    string    `json:"model_config"` // JSON string for model settings
	CreatedAt      time.Time `json:"created_at"`
}

type Routine struct {
	ID          int       `json:"id"`
	AgentID     int       `json:"agent_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Workflow    string    `json:"workflow"` // JSON string defining the steps
	CreatedAt   time.Time `json:"created_at"`
}

type RunStatus string

const (
	StatusPending         RunStatus = "PENDING"
	StatusRunning         RunStatus = "RUNNING"
	StatusCompleted       RunStatus = "COMPLETED"
	StatusFailed          RunStatus = "FAILED"
	StatusWaitingForHuman RunStatus = "WAITING_FOR_HUMAN"
)

type RoutineRun struct {
	ID         int       `json:"id"`
	RoutineID  int       `json:"routine_id"`
	UserID     int       `json:"user_id"`
	Status     RunStatus `json:"status"`
	Result     string    `json:"result"` // JSON result
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}

type EventLog struct {
	ID        int       `json:"id"`
	AgentID   int       `json:"agent_id"`
	RunID     *int      `json:"run_id,omitempty"`
	EventType string    `json:"event_type"` // e.g., "EMAIL_SENT", "POST_CREATED"
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

type Integration struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Configuration  string    `json:"configuration"` // JSON string
	CreatedAt      time.Time `json:"created_at"`
}

type Audience struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type AudienceSegment struct {
	ID         int       `json:"id"`
	AudienceID int       `json:"audience_id"`
	Name       string    `json:"name"`
	Filters    string    `json:"filters"`
	CreatedAt  time.Time `json:"created_at"`
}

type Campaign struct {
	ID                int          `json:"id"`
	OrganizationID    int          `json:"organization_id"`
	AudienceSegmentID *int         `json:"audience_segment_id,omitempty"` // Made pointer as it might be optional for some types? Or keep as is.
	Type              CampaignType `json:"type"`
	WorkflowID        *int         `json:"workflow_id,omitempty"`
	Name              string       `json:"name"`
	PrimaryGoal       string       `json:"primary_goal"`
	KPIs              []string     `json:"kpis"` // JSON array of strings
	Prompt            string       `json:"prompt"`
	Content           string       `json:"content"`
	ScheduleInterval  string       `json:"schedule_interval"`
	NextRunAt         *time.Time   `json:"next_run_at"`
	Status            string       `json:"status"`
	CreatedAt         time.Time    `json:"created_at"`
}

type CampaignRun struct {
	ID          int       `json:"id"`
	CampaignID  int       `json:"campaign_id"`
	SentCount   int       `json:"sent_count"`
	SuccessRate float64   `json:"success_rate"`
	ExecutedAt  time.Time `json:"executed_at"`
}

type Location struct {
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type IPAddress struct {
	Address   string    `json:"address"`
	Used      time.Time `json:"used"`
	Count     int       `json:"count"`
	UserAgent string    `json:"user_agent"`
	Location  *Location `json:"location,omitempty"`
}

type PersonMeta struct {
	IPAddresses []IPAddress `json:"ip_addresses,omitempty"`
	Devices     []string    `json:"devices,omitempty"`
	UserAgents  []string    `json:"user_agents,omitempty"`
	Location    *Location   `json:"location,omitempty"`
}

type PersonEvent struct {
	ID        int         `json:"id"`
	PersonID  int         `json:"person_id"`
	Event     interface{} `json:"event"` // JSONB
	CreatedAt time.Time   `json:"created_at"`
}

type Person struct {
	ID                int           `json:"id"`
	OrganizationID    int           `json:"organization_id"`
	FirstName         *string       `json:"first_name"`
	LastName          *string       `json:"last_name"`
	Interests         []string      `json:"interests"`
	Email             string        `json:"email"`
	Age               *int          `json:"age"`
	Ethnicity         string        `json:"ethnicity"`
	Gender            string        `json:"gender"`
	Location          string        `json:"location"` // Kept as string to match TEXT column
	LastInteractionAt *time.Time    `json:"last_interaction_at"`
	Score             *int          `json:"score"`
	Events            []PersonEvent `json:"events,omitempty"` // populated on detail view
	Meta              PersonMeta    `json:"meta"`
	CreatedAt         time.Time     `json:"created_at"`
}

type SignupCampaign struct {
	ID               int       `json:"id"`
	OrganizationID   int       `json:"organization_id"`
	Name             string    `json:"name"`
	TargetAudienceID int       `json:"target_audience_id"`
	Token            string    `json:"token"`
	CreatedAt        time.Time `json:"created_at"`
}

type Site struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	TrackingID     string    `json:"tracking_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContentTemplate struct {
	ID             int       `json:"id"`
	OrganizationID int       `json:"organization_id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Content        string    `json:"content"`
	Schema         string    `json:"schema"` // JSON string
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wesuuu/helpnow/backend/clients"
	"github.com/wesuuu/helpnow/backend/db"
	"github.com/wesuuu/helpnow/backend/handlers"
	"github.com/wesuuu/helpnow/backend/scheduler"
	"github.com/wesuuu/helpnow/backend/secrets"
	_ "github.com/wesuuu/helpnow/backend/workflows/actions"
	_ "github.com/wesuuu/helpnow/backend/workflows/logic"
	_ "github.com/wesuuu/helpnow/backend/workflows/triggers"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		// Try root .env if ../.env fails (depending on where we run from)
		_ = godotenv.Load()
	}

	// Initialize Database
	db.InitDB()

	// Register Workflow Components (handled by init in subpackages)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${latency_human} | ${method} ${uri}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HelpNow AI Management API is running")
	})

	e.GET("/health", func(c echo.Context) error {
		if err := db.GetDB().Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "db_down"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Initialize AI Client
	clients.InitAIClient()
	defer clients.GlobalAIClient.Close()

	// Initialize Secret Store
	if err := secrets.InitSecretStore(); err != nil {
		log.Printf("Warning: Failed to initialize secret store: %v", err)
	} else {
		log.Println("Secret Store initialized")
	}

	// Auth
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	// Agents
	e.POST("/agents", handlers.CreateAgent)
	e.GET("/agents", handlers.ListAgents)
	e.GET("/agents/:id", handlers.GetAgent)
	e.PUT("/agents/:id", handlers.UpdateAgent)
	e.DELETE("/agents/:id", handlers.DeleteAgent)

	// Routines
	e.POST("/routines", handlers.CreateRoutine)

	// Migration fix/init
	db.GetDB().Exec("ALTER TABLE organizations ADD COLUMN IF NOT EXISTS system_prompt TEXT")
	db.GetDB().Exec("ALTER TABLE sites ADD COLUMN IF NOT EXISTS tracking_id TEXT UNIQUE")

	// Workflow Migrations
	db.GetDB().Exec("ALTER TABLE workflows ADD COLUMN IF NOT EXISTS organization_id INTEGER REFERENCES organizations(id)")
	// Wait, let's stick to plan: 'schedule'
	db.GetDB().Exec("ALTER TABLE workflows ADD COLUMN IF NOT EXISTS schedule TEXT")
	db.GetDB().Exec("ALTER TABLE workflows ADD COLUMN IF NOT EXISTS audience_id INTEGER REFERENCES audiences(id)")
	db.GetDB().Exec("ALTER TABLE workflows ADD COLUMN IF NOT EXISTS next_run_at TIMESTAMP WITH TIME ZONE")
	db.GetDB().Exec("ALTER TABLE workflows ALTER COLUMN site_id DROP NOT NULL")
	db.GetDB().Exec("ALTER TABLE workflows ALTER COLUMN trigger_event DROP NOT NULL")

	// Triggers
	db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS workflow_triggers (
		id SERIAL PRIMARY KEY,
		workflow_id INTEGER REFERENCES workflows(id),
		node_id TEXT NOT NULL,
		type TEXT NOT NULL,
		config TEXT,
		next_run_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)

	db.GetDB().Exec("ALTER TABLE workflow_executions ADD COLUMN IF NOT EXISTS current_node_id TEXT")
	db.GetDB().Exec("ALTER TABLE workflow_executions ADD COLUMN IF NOT EXISTS step_results TEXT")
	db.GetDB().Exec("ALTER TABLE workflow_executions ADD COLUMN IF NOT EXISTS has_failed BOOLEAN DEFAULT FALSE")
	db.GetDB().Exec("ALTER TABLE workflow_executions ADD COLUMN IF NOT EXISTS result TEXT")
	db.GetDB().Exec("ALTER TABLE workflow_executions ADD COLUMN IF NOT EXISTS context TEXT")

	// Create tables if they didn't exist
	// Tables are handled by schema.sql migration logic below

	// People & Audience Members
	e.POST("/people", handlers.CreatePerson)
	e.GET("/people", handlers.ListPeople)
	e.GET("/people/:id", handlers.GetPerson)
	e.POST("/audiences/:id/members", handlers.AddPersonToAudience)
	e.GET("/audiences/:id/members", handlers.GetAudienceMembers)
	e.POST("/people/:id/events", handlers.AppendPersonEvent)

	// Signup Campaigns
	e.POST("/signup-campaigns", handlers.CreateSignupCampaign)
	e.GET("/signup-campaigns", handlers.ListSignupCampaigns)
	e.POST("/public/capture/:token", handlers.CaptureLead)
	e.POST("/public/analytics", handlers.TrackEvent)

	// Sites
	e.POST("/sites", handlers.CreateSite)
	e.GET("/sites", handlers.ListSites)
	e.GET("/sites/:id", handlers.GetSite)
	e.GET("/sites/:id/stats", handlers.GetSiteStats)

	// Workflows & Events
	e.POST("/workflows", handlers.CreateWorkflow)
	e.GET("/workflows", handlers.ListWorkflows)
	e.GET("/workflows/:id", handlers.GetWorkflow)
	e.PUT("/workflows/:id", handlers.UpdateWorkflow)
	e.POST("/events/definitions", handlers.CreateEventDefinition)
	e.GET("/events/definitions", handlers.ListEventDefinitions)

	// Data Sources & Syncs
	e.POST("/sources", handlers.CreateDataSource)
	e.GET("/sources", handlers.ListDataSources)
	e.POST("/syncs", handlers.CreateDataSync)
	e.GET("/syncs", handlers.ListDataSyncs)

	// Seed Organization 1 if not exists
	var count int
	db.GetDB().QueryRow("SELECT COUNT(*) FROM organizations WHERE id = 1").Scan(&count)
	if count == 0 {
		_, err := db.GetDB().Exec("INSERT INTO organizations (id, name, system_prompt) VALUES (1, 'Acme Corp', '')")
		if err != nil {
			e.Logger.Fatal("Failed to seed organization:", err)
		}
	}

	// Seed "Do not call" Audience
	var dncCount int
	db.GetDB().QueryRow("SELECT COUNT(*) FROM audiences WHERE organization_id = 1 AND name = 'Do not call'").Scan(&dncCount)
	if dncCount == 0 {
		_, err := db.GetDB().Exec("INSERT INTO audiences (organization_id, name, description) VALUES (1, 'Do not call', 'People who should not be contacted')")
		if err != nil {
			log.Println("Failed to seed Do not call audience:", err)
		}
	}

	// Campaigns Migration (Simplified - Manual Rename assumed or done by schema if fresh)
	// We just ensure columns exist on 'campaigns' table.
	// If 'campaigns' table doesn't exist (because it's still email_campaigns), these will fail safely (logged error but continued execution if I didn't check err... wait Exec returns err).
	// db.GetDB().Exec ignores return values if not assigned?
	// In Go, "_, _ = ..." or just "db.GetDB().Exec(...)" is valid statement value discard? Yes.

	db.GetDB().Exec("ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS type TEXT")
	db.GetDB().Exec("ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS workflow_id INTEGER REFERENCES workflows(id)")
	db.GetDB().Exec("ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS primary_goal TEXT")
	db.GetDB().Exec("ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS kpis JSONB")

	// Seed Site 1 if not exists
	var siteCount int
	db.GetDB().QueryRow("SELECT COUNT(*) FROM sites WHERE id = 1").Scan(&siteCount)
	if siteCount == 0 {
		_, err := db.GetDB().Exec("INSERT INTO sites (id, organization_id, name, url) VALUES (1, 1, 'Default Site', 'https://example.com')")
		if err != nil {
			log.Println("Failed to seed site:", err)
		}
	}

	// Create Metric Tables
	db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS metric_signups (
		id SERIAL PRIMARY KEY,
		site_id INTEGER,
		user_email TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS metric_purchases (
		id SERIAL PRIMARY KEY,
		site_id INTEGER,
		amount DECIMAL,
		currency TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS metric_page_views (
		id SERIAL PRIMARY KEY,
		site_id INTEGER,
		path TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)

	db.GetDB().Exec(`CREATE TABLE IF NOT EXISTS content_templates (
		id SERIAL PRIMARY KEY,
		organization_id INTEGER REFERENCES organizations(id),
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		content TEXT,
		schema JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)

	// Seed Event Definitions
	events := []string{"signup", "purchase", "page_view"}
	for _, eventName := range events {
		_, _ = db.GetDB().Exec("INSERT INTO event_definitions (site_id, name, description) VALUES (1, $1, 'System metric event') ON CONFLICT DO NOTHING", eventName)
	}

	// Routes
	e.GET("/routines", handlers.ListRoutines)
	e.POST("/routines/execute", handlers.ExecuteRoutine)

	// Organizations
	e.GET("/organizations/:id", handlers.GetOrganization)
	e.PUT("/organizations/:id", handlers.UpdateOrganization)

	// Integrations
	e.POST("/integrations", handlers.CreateIntegration)
	e.GET("/integrations", handlers.ListIntegrations)

	// Audiences
	e.POST("/audiences", handlers.CreateAudience)
	e.GET("/audiences", handlers.ListAudiences)
	e.POST("/audiences/segments", handlers.CreateAudienceSegment)

	// Campaigns
	e.POST("/campaigns", handlers.CreateCampaign)
	e.GET("/email-domains", handlers.ListEmailDomains)
	e.POST("/email-domains", handlers.CreateEmailDomain)
	e.POST("/email-domains/:id/verify", handlers.VerifyEmailDomain)
	e.DELETE("/email-domains/:id", handlers.DeleteEmailDomain)

	// Email Templates
	e.GET("/email-templates", handlers.ListEmailTemplates)
	e.POST("/email-templates", handlers.CreateEmailTemplate)
	e.PUT("/email-templates/:id", handlers.UpdateEmailTemplate)
	e.DELETE("/email-templates/:id", handlers.DeleteEmailTemplate)

	e.GET("/campaigns", handlers.ListCampaigns)
	e.POST("/campaigns/:id/content", handlers.GenerateCampaignContent) // POST to trigger generation
	e.PUT("/campaigns/:id", handlers.UpdateCampaign)
	e.GET("/campaigns/:id/runs", handlers.ListCampaignRuns)

	// Content Templates
	e.POST("/templates/generate", handlers.GenerateTemplateContent)
	e.POST("/templates", handlers.CreateContentTemplate)
	e.GET("/templates", handlers.ListContentTemplates)
	e.GET("/templates/:id", handlers.GetContentTemplate)
	e.PUT("/templates/:id", handlers.UpdateContentTemplate)
	e.DELETE("/templates/:id", handlers.DeleteContentTemplate)

	// Start Scheduler
	scheduler.Start()

	// Simple Migration (MVP)
	schema, err := os.ReadFile("schema.sql")
	if err == nil {
		_, err = db.GetDB().Exec(string(schema))
		if err != nil {
			log.Println("Schema migration warning:", err)
		} else {
			log.Println("Schema migrated successfully")
		}
	} else {
		log.Println("schema.sql not found, skipping migration")
	}

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on port " + port)
	e.Logger.Fatal(e.Start(":" + port))
}

// Workflow Component Introspection
e.GET("/workflow-components/actions", handlers.ListActions)
e.GET("/workflow-components/actions/:name", handlers.GetAction)
e.GET("/workflow-components/logic", handlers.ListLogic)
e.GET("/workflow-components/logic/:name", handlers.GetLogic)
e.GET("/workflow-components/triggers", handlers.ListTriggers)
e.GET("/workflow-components/triggers/:name", handlers.GetTrigger)

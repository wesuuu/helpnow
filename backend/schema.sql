CREATE TABLE IF NOT EXISTS organizations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    system_prompt TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    middle_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS agents (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    description TEXT,
    model_config TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS routines (
    id SERIAL PRIMARY KEY,
    agent_id INTEGER REFERENCES agents(id),
    name TEXT NOT NULL,
    description TEXT,
    workflow TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS routine_runs (
    id SERIAL PRIMARY KEY,
    routine_id INTEGER REFERENCES routines(id),
    user_id INTEGER REFERENCES users(id),
    status TEXT NOT NULL,
    result TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS integrations (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- e.g. "BLOG_WORDPRESS", "SOCIAL_TWITTER"
    configuration TEXT, -- JSON string storing API keys, URLs, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audiences (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    age INTEGER,
    ethnicity TEXT,
    gender TEXT,
    location TEXT,
    last_interaction_at TIMESTAMP WITH TIME ZONE,
    score INTEGER DEFAULT 0,
    event_history JSONB,
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audience_memberships (
    audience_id INTEGER REFERENCES audiences(id),
    person_id INTEGER REFERENCES people(id),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (audience_id, person_id)
);

CREATE TABLE IF NOT EXISTS audience_segments (
    id SERIAL PRIMARY KEY,
    audience_id INTEGER REFERENCES audiences(id),
    type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_campaigns (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    audience_segment_id INTEGER REFERENCES audience_segments(id),
    name TEXT NOT NULL,
    prompt TEXT,
    content TEXT,
    schedule_interval TEXT, -- 'DAILY', 'WEEKLY', etc.
    next_run_at TIMESTAMP WITH TIME ZONE,
    status TEXT DEFAULT 'DRAFT', -- DRAFT, ACTIVE, PAUSED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS campaign_runs (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES email_campaigns(id),
    sent_count INTEGER DEFAULT 0,
    success_rate FLOAT DEFAULT 0,
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS signup_campaigns (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    target_audience_id INTEGER REFERENCES audiences(id),
    token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sites (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    tracking_id TEXT UNIQUE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_definitions (
    id SERIAL PRIMARY KEY,
    site_id INTEGER REFERENCES sites(id),
    name TEXT NOT NULL, -- e.g. "purchased_item"
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(site_id, name)
);

CREATE TABLE IF NOT EXISTS workflows (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id), -- New: Owner org
    site_id INTEGER REFERENCES sites(id), -- Nullable now if org-level
    audience_id INTEGER REFERENCES audiences(id), -- New: Context for schedule
    name TEXT NOT NULL,
    trigger_type TEXT NOT NULL, -- 'EVENT', 'SCHEDULE'
    trigger_event TEXT, -- Name of event if type=EVENT
    steps TEXT NOT NULL, -- JSON array of steps
    schedule TEXT, -- New: Cron expression or interval
    next_run_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- New: For scheduling
    status TEXT DEFAULT 'ACTIVE', -- ACTIVE, PAUSED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(organization_id, name)
);

CREATE TABLE IF NOT EXISTS workflow_executions (
    id SERIAL PRIMARY KEY,
    workflow_id INTEGER REFERENCES workflows(id),
    subject_id INTEGER, -- Optional: Link to a 'people' record if applicable
    current_step INTEGER DEFAULT 0, -- Deprecated in favor of current_node_id for Graph workflows
    current_node_id TEXT, -- ID of the current node in the graph
    status TEXT NOT NULL, -- PENDING, COMPLETED, FAILED
    next_run_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    context TEXT, -- JSON blob of event data
    step_results TEXT, -- JSON array of results from each step
    has_failed BOOLEAN DEFAULT FALSE, -- Track if any step has failed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS data_sources (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'POSTGRES', 'DATABRICKS', 'WEBHOOK'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS data_syncs (
    id SERIAL PRIMARY KEY,
    source_id INTEGER REFERENCES data_sources(id),
    audience_id INTEGER REFERENCES audiences(id),
    sync_type TEXT NOT NULL, -- 'ONE_OFF', 'CRON'
    schedule TEXT, -- cron expression
    query TEXT, -- SQL query or mapping config
    status TEXT NOT NULL, -- 'PENDING', 'RUNNING', 'COMPLETED', 'FAILED'
    last_run_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS email_domains (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    domain TEXT NOT NULL,
    dkim_record TEXT,
    spf_record TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(organization_id, domain)
);

CREATE TABLE IF NOT EXISTS email_templates (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS content_templates (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id),
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- 'AD', 'FORM', 'LANDING_PAGE'
    content TEXT, -- HTML, Markdown, or JSON structure
    schema JSONB, -- JSON schema for variables
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

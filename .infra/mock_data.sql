-- Mock Users compatible with current schema.sql
-- Organizations: name only
-- Users: email, password_hash, full_name, org_id

INSERT INTO organizations (name) VALUES ('HelpNow Demo Org') ON CONFLICT DO NOTHING;

INSERT INTO users (email, password_hash, full_name, organization_id) 
VALUES 
('admin@helpnow.ai', '$2a$10$dummyhash', 'Admin User', (SELECT id FROM organizations WHERE name='HelpNow Demo Org' LIMIT 1)),
('user@helpnow.ai', '$2a$10$dummyhash', 'Demo User', (SELECT id FROM organizations WHERE name='HelpNow Demo Org' LIMIT 1))
ON CONFLICT (email) DO NOTHING;

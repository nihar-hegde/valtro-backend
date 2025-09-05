-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    -- Unique identifier for the project, using UUID.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key linking this project to the organization that owns it.
    -- ON DELETE CASCADE means if an organization is deleted, all its projects are also deleted.
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- The user-provided name for the project (e.g., "Web App - Production").
    name VARCHAR(255) NOT NULL,

    -- The public API key used by the SDK to send logs.
    -- This MUST be unique across all projects in the system.
    api_key VARCHAR(255) NOT NULL UNIQUE,

    -- Standard timestamps managed by PostgreSQL.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Soft delete timestamp for GORM soft delete functionality
    deleted_at TIMESTAMP
);

-- Create an index on the organization_id for quickly fetching all projects for an org.
CREATE INDEX IF NOT EXISTS idx_projects_organization_id ON projects(organization_id);

-- Create an index on the api_key for fast lookups during SDK requests.
CREATE INDEX IF NOT EXISTS idx_projects_api_key ON projects(api_key);

-- Create an index on deleted_at for efficient soft delete filtering
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);

-- Add comments for documentation
COMMENT ON TABLE projects IS 'Projects belonging to organizations for log collection';
COMMENT ON COLUMN projects.organization_id IS 'Organization that owns this project';
COMMENT ON COLUMN projects.name IS 'User-provided name for the project';
COMMENT ON COLUMN projects.api_key IS 'Unique API key for SDK authentication';

-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    -- Unique identifier for the organization, using UUID for global uniqueness.
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- The user-provided name for the organization (e.g., "Acme Inc.").
    name VARCHAR(255) NOT NULL,

    -- Foreign key linking to the user who created and "owns" the organization.
    -- This is crucial for billing and top-level permissions.
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Standard timestamps managed by PostgreSQL.
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Soft delete timestamp for GORM soft delete functionality
    deleted_at TIMESTAMP
);

-- Create an index on the owner_id for fast lookups of organizations by a specific user.
CREATE INDEX IF NOT EXISTS idx_organizations_owner_id ON organizations(owner_id);

-- Create an index on deleted_at for efficient soft delete filtering
CREATE INDEX IF NOT EXISTS idx_organizations_deleted_at ON organizations(deleted_at);

-- Add comments for documentation
COMMENT ON TABLE organizations IS 'Organizations owned by users for project management';
COMMENT ON COLUMN organizations.name IS 'User-provided name for the organization';
COMMENT ON COLUMN organizations.owner_id IS 'User who created and owns the organization';

-- Rollback email_verified column reordering
-- This moves email_verified back to the end of the table

BEGIN;

-- Create table with email_verified at the end (original position)
CREATE TABLE users_rollback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    username VARCHAR(50),
    image_url VARCHAR(500),
    active BOOLEAN DEFAULT true,
    last_sign_in TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    email_verified BOOLEAN DEFAULT false
);

-- Copy data from current table
INSERT INTO users_rollback (id, clerk_user_id, email, full_name, username, image_url, active, last_sign_in, created_at, updated_at, deleted_at, email_verified)
SELECT id, clerk_user_id, email, full_name, username, image_url, active, last_sign_in, created_at, updated_at, deleted_at, email_verified
FROM users;

-- Drop current table
DROP TABLE users;

-- Rename rollback table
ALTER TABLE users_rollback RENAME TO users;

-- Recreate indexes
CREATE UNIQUE INDEX idx_clerk_user ON users(clerk_user_id);
CREATE UNIQUE INDEX idx_user_email ON users(email);
CREATE UNIQUE INDEX idx_username ON users(username) WHERE username IS NOT NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Add comments for documentation
COMMENT ON TABLE users IS 'User accounts integrated with Clerk authentication';
COMMENT ON COLUMN users.clerk_user_id IS 'Unique identifier from Clerk auth provider';
COMMENT ON COLUMN users.email IS 'User email address, must be unique';
COMMENT ON COLUMN users.username IS 'Optional unique username';

COMMIT;

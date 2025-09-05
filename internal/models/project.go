package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents the project model
type Project struct {
	// ID is the primary key for the project record, automatically generated as a UUID
	// Uses PostgreSQL's gen_random_uuid() function for generation
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// OrganizationID is a foreign key reference to the organization that owns this project
	// Required field with CASCADE delete behavior (if organization is deleted, project is deleted)
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_projects_organization_id"`

	// Organization is the relationship to the Organization model
	// This allows GORM to handle the foreign key relationship
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`

	// Name stores the project's name
	// Required field with maximum length of 255 characters
	Name string `gorm:"type:varchar(255);not null"`

	// APIKey is the unique API key used by the SDK to send logs
	// Must be unique across all projects in the system
	// Has a unique index for fast lookups during SDK requests
	APIKey string `gorm:"type:varchar(255);not null;uniqueIndex:idx_projects_api_key"`

	// Standard timestamp fields

	// CreatedAt is automatically managed by GORM
	// Records when the project record was created
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// UpdatedAt is automatically managed by GORM
	// Records when the project record was last updated
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// DeletedAt enables soft deletion in GORM
	// When a record is "deleted", this field is set instead of removing the record
	// Has an index for efficient filtering of deleted records
	// GORM automatically excludes soft-deleted records from queries
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

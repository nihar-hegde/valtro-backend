package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents the organization model
type Organization struct {
	// ID is the primary key for the organization record, automatically generated as a UUID
	// Uses PostgreSQL's gen_random_uuid() function for generation
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Name stores the organization's name
	// Required field with maximum length of 255 characters
	Name string `gorm:"type:varchar(255);not null"`

	// OwnerID is a foreign key reference to the user who owns this organization
	// Required field with CASCADE delete behavior (if user is deleted, organization is deleted)
	OwnerID uuid.UUID `gorm:"type:uuid;not null;index:idx_organizations_owner_id"`

	// Owner is the relationship to the User model
	// This allows GORM to handle the foreign key relationship
	Owner User `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`

	// Projects is the relationship to Project models
	// One organization can have many projects
	Projects []Project `gorm:"foreignKey:OrganizationID"`

	// Standard timestamp fields

	// CreatedAt is automatically managed by GORM
	// Records when the organization record was created
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// UpdatedAt is automatically managed by GORM
	// Records when the organization record was last updated
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// DeletedAt enables soft deletion in GORM
	// When a record is "deleted", this field is set instead of removing the record
	// Has an index for efficient filtering of deleted records
	// GORM automatically excludes soft-deleted records from queries
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

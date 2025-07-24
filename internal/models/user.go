package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user model, storing details from the auth provider (Clerk).
type User struct {
	// ID is the primary key for the user record, automatically generated as a UUID
	// Uses PostgreSQL's gen_random_uuid() function for generation
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// ClerkUserID stores the unique identifier from Clerk auth provider
	// Has a named unique index (idx_clerk_user) for efficient lookups
	// Limited to 255 characters for storage efficiency
	ClerkUserID string `gorm:"type:varchar(255);uniqueIndex:idx_clerk_user;not null"`

	// Email stores the user's email address
	// Has a named unique index (idx_user_email) for efficient lookups and ensuring uniqueness
	// Required field (not null constraint)
	Email string `gorm:"type:varchar(255);uniqueIndex:idx_user_email;not null"`

	// FullName stores the user's full name
	// Optional field, can be empty
	FullName string `gorm:"type:varchar(255)"`

	// Username stores the user's chosen username
	// Pointer type (*string) makes it nullable for users who haven't set a username
	// Has a unique index (idx_username) to ensure username uniqueness
	// Limited to 50 characters for UI/UX considerations
	Username *string `gorm:"type:varchar(50);uniqueIndex:idx_username"`

	// ImageURL stores the URL to the user's profile image
	// Allows longer strings (500 chars) to accommodate various URL formats
	ImageURL string `gorm:"type:varchar(500)"`
	
	// Additional useful fields for Clerk integration
	
	// EmailVerified indicates whether the user's email has been verified
	// Defaults to false until verification is confirmed
	EmailVerified bool `gorm:"default:false"`

	// Active indicates whether the user account is currently active
	// Defaults to true for new accounts
	// Can be used for temporary deactivation without deletion
	Active bool `gorm:"default:true"`

	// LastSignIn tracks when the user last authenticated
	// Nullable (pointer type) for users who have never signed in
	LastSignIn *time.Time
	
	// Standard timestamp fields
	
	// CreatedAt is automatically managed by GORM
	// Records when the user record was created
	CreatedAt time.Time

	// UpdatedAt is automatically managed by GORM
	// Records when the user record was last updated
	UpdatedAt time.Time

	// DeletedAt enables soft deletion in GORM
	// When a record is "deleted", this field is set instead of removing the record
	// Has an index for efficient filtering of deleted records
	// GORM automatically excludes soft-deleted records from queries
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

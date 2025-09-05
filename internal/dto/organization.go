package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateOrganizationRequest represents the request payload for creating an organization
type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// UpdateOrganizationRequest represents the request payload for updating an organization
type UpdateOrganizationRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
}

// OrganizationResponse represents the response structure for organization data
type OrganizationResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrganizationWithProjectsResponse represents organization data with its projects
type OrganizationWithProjectsResponse struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	OwnerID   uuid.UUID         `json:"owner_id"`
	Projects  []ProjectResponse `json:"projects"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// UserOrganizationCheckResponse represents the response for checking user's organization membership
type UserOrganizationCheckResponse struct {
	HasOrganization bool                  `json:"has_organization"`
	Organization    *OrganizationResponse `json:"organization,omitempty"`
}

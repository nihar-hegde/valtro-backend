package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateProjectRequest represents the request payload for creating a project
type CreateProjectRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	Name           string    `json:"name" validate:"required,min=2,max=255"`
}

// UpdateProjectRequest represents the request payload for updating a project
type UpdateProjectRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
}

// ProjectResponse represents the response structure for project data
type ProjectResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	APIKey         string    `json:"api_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// OnboardingRequest represents the request payload for the complete onboarding flow
type OnboardingRequest struct {
	OrganizationName string `json:"organization_name" validate:"required,min=2,max=255"`
	ProjectName      string `json:"project_name" validate:"required,min=2,max=255"`
}

// OnboardingResponse represents the response for the complete onboarding flow
type OnboardingResponse struct {
	Organization OrganizationResponse `json:"organization"`
	Project      ProjectResponse      `json:"project"`
}

package organization

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	"github.com/nihar-hegde/valtro-backend/internal/errors"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	orgRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/organization"
)

// Service handles organization business logic
type Service struct {
	orgRepo *orgRepo.Repository
}

// NewService creates a new organization service
func NewService(orgRepo *orgRepo.Repository) *Service {
	return &Service{
		orgRepo: orgRepo,
	}
}

// CreateOrganization creates a new organization with business logic validation
func (s *Service) CreateOrganization(req dto.CreateOrganizationRequest, ownerID uuid.UUID) (*dto.OrganizationResponse, error) {
	// Validate business rules
	if err := s.validateCreateOrganization(req); err != nil {
		return nil, err
	}

	// Check if organization name already exists for this owner
	nameExists, err := s.orgRepo.NameExistsForOwner(req.Name, ownerID)
	if err != nil {
		return nil, err // Repository now returns structured errors
	}
	if nameExists {
		return nil, errors.NewConflictError("Organization with this name already exists for user", "Name: "+req.Name)
	}

	// Create organization model
	organization := &models.Organization{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(req.Name),
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.orgRepo.Create(organization); err != nil {
		return nil, err // Repository now returns structured errors
	}

	// Convert to response DTO
	return s.toOrganizationResponse(organization), nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *Service) GetOrganizationByID(id uuid.UUID) (*dto.OrganizationResponse, error) {
	organization, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toOrganizationResponse(organization), nil
}

// GetOrganizationsByOwner retrieves all organizations for a specific owner
func (s *Service) GetOrganizationsByOwner(ownerID uuid.UUID) ([]*dto.OrganizationResponse, error) {
	organizations, err := s.orgRepo.GetByOwnerID(ownerID)
	if err != nil {
		return nil, err // Repository now returns structured errors
	}

	// Convert to response DTOs
	var responses []*dto.OrganizationResponse
	for _, org := range organizations {
		responses = append(responses, s.toOrganizationResponse(org))
	}

	return responses, nil
}

// GetOrganizationWithProjects retrieves an organization with its projects
func (s *Service) GetOrganizationWithProjects(ownerID uuid.UUID) (*dto.OrganizationWithProjectsResponse, error) {
	organization, err := s.orgRepo.GetByOwnerIDWithProjects(ownerID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO with projects
	response := &dto.OrganizationWithProjectsResponse{
		ID:        organization.ID,
		Name:      organization.Name,
		OwnerID:   organization.OwnerID,
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
		Projects:  make([]dto.ProjectResponse, len(organization.Projects)),
	}

	// Convert projects to DTOs
	for i, project := range organization.Projects {
		response.Projects[i] = dto.ProjectResponse{
			ID:             project.ID,
			OrganizationID: project.OrganizationID,
			Name:           project.Name,
			APIKey:         project.APIKey,
			CreatedAt:      project.CreatedAt,
			UpdatedAt:      project.UpdatedAt,
		}
	}

	return response, nil
}

// UpdateOrganization updates an organization
func (s *Service) UpdateOrganization(id uuid.UUID, req dto.UpdateOrganizationRequest, ownerID uuid.UUID) (*dto.OrganizationResponse, error) {
	// Get existing organization
	organization, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user owns this organization
	if organization.OwnerID != ownerID {
		return nil, errors.NewForbiddenError("Unauthorized: you don't own this organization", "Organization ID: "+id.String())
	}

	// Update fields if provided
	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, errors.NewValidationError("Organization name cannot be empty")
		}

		// Check if new name already exists for this owner (excluding current organization)
		nameExists, err := s.orgRepo.NameExistsForOwner(trimmedName, ownerID)
		if err != nil {
			return nil, err // Repository now returns structured errors
		}
		if nameExists && organization.Name != trimmedName {
			return nil, errors.NewConflictError("Organization with this name already exists for user", "Name: "+trimmedName)
		}

		organization.Name = trimmedName
	}

	organization.UpdatedAt = time.Now()

	// Save changes
	if err := s.orgRepo.Update(organization); err != nil {
		return nil, err // Repository now returns structured errors
	}

	return s.toOrganizationResponse(organization), nil
}

// DeleteOrganization soft deletes an organization
func (s *Service) DeleteOrganization(id uuid.UUID, ownerID uuid.UUID) error {
	// Get existing organization
	organization, err := s.orgRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user owns this organization
	if organization.OwnerID != ownerID {
		return errors.NewForbiddenError("Unauthorized: you don't own this organization", "Organization ID: "+id.String())
	}

	// Soft delete
	if err := s.orgRepo.Delete(id); err != nil {
		return err // Repository now returns structured errors
	}

	return nil
}

// CheckUserOrganization checks if a user has an organization
func (s *Service) CheckUserOrganization(ownerID uuid.UUID) (*dto.UserOrganizationCheckResponse, error) {
	hasOrg, err := s.orgRepo.HasOrganization(ownerID)
	if err != nil {
		return nil, err // Repository now returns structured errors
	}

	response := &dto.UserOrganizationCheckResponse{
		HasOrganization: hasOrg,
	}

	if hasOrg {
		organizations, err := s.orgRepo.GetByOwnerID(ownerID)
		if err != nil {
			return nil, err // Repository now returns structured errors
		}
		if len(organizations) > 0 {
			response.Organization = s.toOrganizationResponse(organizations[0])
		}
	}

	return response, nil
}

// validateCreateOrganization validates the create organization request
func (s *Service) validateCreateOrganization(req dto.CreateOrganizationRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.NewValidationError("Organization name is required")
	}
	if len(req.Name) < 2 {
		return errors.NewValidationError("Organization name must be at least 2 characters")
	}
	if len(req.Name) > 255 {
		return errors.NewValidationError("Organization name must be less than 255 characters")
	}
	return nil
}

// toOrganizationResponse converts an organization model to response DTO
func (s *Service) toOrganizationResponse(organization *models.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		ID:        organization.ID,
		Name:      organization.Name,
		OwnerID:   organization.OwnerID,
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
	}
}

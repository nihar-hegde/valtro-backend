package project

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/constants"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	projectRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/project"
)

// Service handles project business logic
type Service struct {
	projectRepo *projectRepo.Repository
}

// NewService creates a new project service
func NewService(projectRepo *projectRepo.Repository) *Service {
	return &Service{
		projectRepo: projectRepo,
	}
}

// CreateProject creates a new project with business logic validation
func (s *Service) CreateProject(req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	// Validate business rules
	if err := s.validateCreateProject(req); err != nil {
		return nil, err
	}

	// Check if project name already exists for this organization
	nameExists, err := s.projectRepo.NameExistsForOrganization(req.Name, req.OrganizationID)
	if err != nil {
		return nil, errors.New("failed to check project name existence")
	}
	if nameExists {
		return nil, errors.New("project with this name already exists in organization")
	}

	// Generate unique API key
	apiKey, err := s.generateUniqueAPIKey()
	if err != nil {
		return nil, errors.New("failed to generate API key")
	}

	// Create project model
	project := &models.Project{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		Name:           strings.TrimSpace(req.Name),
		APIKey:         apiKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.New("failed to create project")
	}

	// Convert to response DTO
	return s.toProjectResponse(project), nil
}

// GetProjectByID retrieves a project by ID
func (s *Service) GetProjectByID(id uuid.UUID) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toProjectResponse(project), nil
}

// GetProjectByAPIKey retrieves a project by API key
func (s *Service) GetProjectByAPIKey(apiKey string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	return s.toProjectResponse(project), nil
}

// GetProjectsByOrganization retrieves all projects for a specific organization
func (s *Service) GetProjectsByOrganization(organizationID uuid.UUID) ([]*dto.ProjectResponse, error) {
	projects, err := s.projectRepo.GetByOrganizationID(organizationID)
	if err != nil {
		return nil, errors.New("failed to retrieve projects")
	}

	// Convert to response DTOs
	var responses []*dto.ProjectResponse
	for _, project := range projects {
		responses = append(responses, s.toProjectResponse(project))
	}

	return responses, nil
}

// UpdateProject updates a project
func (s *Service) UpdateProject(id uuid.UUID, req dto.UpdateProjectRequest, organizationID uuid.UUID) (*dto.ProjectResponse, error) {
	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if project belongs to the specified organization
	if project.OrganizationID != organizationID {
		return nil, errors.New("unauthorized: project doesn't belong to this organization")
	}

	// Update fields if provided
	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return nil, errors.New("project name cannot be empty")
		}

		// Check if new name already exists for this organization (excluding current project)
		nameExists, err := s.projectRepo.NameExistsForOrganization(trimmedName, organizationID)
		if err != nil {
			return nil, errors.New("failed to check project name existence")
		}
		if nameExists && project.Name != trimmedName {
			return nil, errors.New("project with this name already exists in organization")
		}

		project.Name = trimmedName
	}

	project.UpdatedAt = time.Now()

	// Save changes
	if err := s.projectRepo.Update(project); err != nil {
		return nil, errors.New("failed to update project")
	}

	return s.toProjectResponse(project), nil
}

// DeleteProject soft deletes a project
func (s *Service) DeleteProject(id uuid.UUID, organizationID uuid.UUID) error {
	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if project belongs to the specified organization
	if project.OrganizationID != organizationID {
		return errors.New("unauthorized: project doesn't belong to this organization")
	}

	// Soft delete
	if err := s.projectRepo.Delete(id); err != nil {
		return errors.New("failed to delete project")
	}

	return nil
}

// RegenerateAPIKey generates a new API key for a project
func (s *Service) RegenerateAPIKey(id uuid.UUID, organizationID uuid.UUID) (*dto.ProjectResponse, error) {
	// Get existing project
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if project belongs to the specified organization
	if project.OrganizationID != organizationID {
		return nil, errors.New("unauthorized: project doesn't belong to this organization")
	}

	// Generate new unique API key
	apiKey, err := s.generateUniqueAPIKey()
	if err != nil {
		return nil, errors.New("failed to generate new API key")
	}

	project.APIKey = apiKey
	project.UpdatedAt = time.Now()

	// Save changes
	if err := s.projectRepo.Update(project); err != nil {
		return nil, errors.New("failed to update project")
	}

	return s.toProjectResponse(project), nil
}

// generateUniqueAPIKey generates a unique API key
func (s *Service) generateUniqueAPIKey() (string, error) {
	for attempt := 0; attempt < constants.APIKeyMaxAttempts; attempt++ {
		// Generate random bytes and encode as hex
		bytes := make([]byte, constants.APIKeyByteSize)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		
		apiKey := constants.APIKeyPrefix + hex.EncodeToString(bytes)
		
		// Check if this API key already exists
		exists, err := s.projectRepo.APIKeyExists(apiKey)
		if err != nil {
			return "", err
		}
		
		if !exists {
			return apiKey, nil
		}
	}
	
	return "", errors.New("failed to generate unique API key after multiple attempts")
}

// validateCreateProject validates the create project request
func (s *Service) validateCreateProject(req dto.CreateProjectRequest) error {
	if req.OrganizationID == uuid.Nil {
		return errors.New("organization ID is required")
	}
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("project name is required")
	}
	if len(req.Name) < 2 {
		return errors.New("project name must be at least 2 characters")
	}
	if len(req.Name) > 255 {
		return errors.New("project name must be less than 255 characters")
	}
	return nil
}

// toProjectResponse converts a project model to response DTO
func (s *Service) toProjectResponse(project *models.Project) *dto.ProjectResponse {
	return &dto.ProjectResponse{
		ID:             project.ID,
		OrganizationID: project.OrganizationID,
		Name:           project.Name,
		APIKey:         project.APIKey,
		CreatedAt:      project.CreatedAt,
		UpdatedAt:      project.UpdatedAt,
	}
}

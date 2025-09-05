package project

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	"gorm.io/gorm"
)

// Repository handles project data access operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new project repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new project in the database
func (r *Repository) Create(project *models.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a project by its ID
func (r *Repository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &project, nil
}

// GetByAPIKey retrieves a project by its API key
func (r *Repository) GetByAPIKey(apiKey string) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, "api_key = ?", apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &project, nil
}

// GetByOrganizationID retrieves all projects for a specific organization
func (r *Repository) GetByOrganizationID(organizationID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	if err := r.db.Where("organization_id = ?", organizationID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Update updates a project in the database
func (r *Repository) Update(project *models.Project) error {
	if err := r.db.Save(project).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes a project from the database
func (r *Repository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Project{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// APIKeyExists checks if an API key already exists
func (r *Repository) APIKeyExists(apiKey string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("api_key = ?", apiKey).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// NameExistsForOrganization checks if a project name already exists for a specific organization
func (r *Repository) NameExistsForOrganization(name string, organizationID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("name = ? AND organization_id = ?", name, organizationID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

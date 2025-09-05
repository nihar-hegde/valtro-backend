package organization

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	"gorm.io/gorm"
)

// Repository handles organization data access operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new organization repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new organization in the database
func (r *Repository) Create(organization *models.Organization) error {
	if err := r.db.Create(organization).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves an organization by its ID
func (r *Repository) GetByID(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &organization, nil
}

// GetByOwnerID retrieves organizations owned by a specific user
func (r *Repository) GetByOwnerID(ownerID uuid.UUID) ([]*models.Organization, error) {
	var organizations []*models.Organization
	if err := r.db.Where("owner_id = ? AND deleted_at IS NULL", ownerID).Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

// GetByOwnerIDPaginated retrieves organizations owned by a specific user with pagination
func (r *Repository) GetByOwnerIDPaginated(ownerID uuid.UUID, offset, limit int) ([]*models.Organization, int64, error) {
	var organizations []*models.Organization
	var totalCount int64
	
	// Get total count
	if err := r.db.Model(&models.Organization{}).
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	if err := r.db.Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&organizations).Error; err != nil {
		return nil, 0, err
	}
	
	return organizations, totalCount, nil
}

// GetByOwnerIDWithProjects retrieves an organization by owner ID with its projects
func (r *Repository) GetByOwnerIDWithProjects(ownerID uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	// Optimize preload to only load non-deleted projects and essential fields
	if err := r.db.Preload("Projects", "deleted_at IS NULL").
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &organization, nil
}

// GetByIDWithProjects retrieves an organization by ID with its projects
func (r *Repository) GetByIDWithProjects(id uuid.UUID) (*models.Organization, error) {
	var organization models.Organization
	// Optimize preload to only load non-deleted projects
	if err := r.db.Preload("Projects", "deleted_at IS NULL").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &organization, nil
}

// Update updates an organization in the database
func (r *Repository) Update(organization *models.Organization) error {
	if err := r.db.Save(organization).Error; err != nil {
		return err
	}
	return nil
}

// Delete soft deletes an organization from the database
func (r *Repository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Organization{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// NameExistsForOwner checks if an organization name already exists for a specific owner
func (r *Repository) NameExistsForOwner(name string, ownerID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Organization{}).
		Where("name = ? AND owner_id = ? AND deleted_at IS NULL", name, ownerID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasOrganization checks if a user has any organization
func (r *Repository) HasOrganization(ownerID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Organization{}).
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

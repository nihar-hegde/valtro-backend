package user

import (
	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/errors"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	"gorm.io/gorm"
)

// Repository handles user data access operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new user in the database
func (r *Repository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.NewInternalError("Failed to create user", err.Error())
	}
	return nil
}

// GetByID retrieves a user by their ID
func (r *Repository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewNotFoundError("User", "User with ID "+id.String()+" not found")
		}
		return nil, errors.NewInternalError("Failed to retrieve user", err.Error())
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *Repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewNotFoundError("User", "User with email "+email+" not found")
		}
		return nil, errors.NewInternalError("Failed to retrieve user by email", err.Error())
	}
	return &user, nil
}

// GetByClerkUserID retrieves a user by their Clerk user ID
func (r *Repository) GetByClerkUserID(clerkUserID string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "clerk_user_id = ?", clerkUserID).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewNotFoundError("User", "User with Clerk ID "+clerkUserID+" not found")
		}
		return nil, errors.NewInternalError("Failed to retrieve user by Clerk ID", err.Error())
	}
	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *Repository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, errors.NewNotFoundError("User", "User with username "+username+" not found")
		}
		return nil, errors.NewInternalError("Failed to retrieve user by username", err.Error())
	}
	return &user, nil
}

// GetAll retrieves all active users with pagination
func (r *Repository) GetAll(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Where("active = ?", true).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, errors.NewInternalError("Failed to retrieve users", err.Error())
	}
	return users, nil
}

// Update updates a user in the database
func (r *Repository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.NewInternalError("Failed to update user", err.Error())
	}
	return nil
}

// Delete soft deletes a user from the database
func (r *Repository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return errors.NewInternalError("Failed to delete user", err.Error())
	}
	return nil
}

// EmailExists checks if an email already exists
func (r *Repository) EmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, errors.NewInternalError("Failed to check email existence", err.Error())
	}
	return count > 0, nil
}

// UsernameExists checks if a username already exists
func (r *Repository) UsernameExists(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errors.NewInternalError("Failed to check username existence", err.Error())
	}
	return count > 0, nil
}

// ClerkUserIDExists checks if a Clerk user ID already exists
func (r *Repository) ClerkUserIDExists(clerkUserID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("clerk_user_id = ?", clerkUserID).Count(&count).Error; err != nil {
		return false, errors.NewInternalError("Failed to check Clerk user ID existence", err.Error())
	}
	return count > 0, nil
}

package user

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	"github.com/nihar-hegde/valtro-backend/internal/errors"
	"github.com/nihar-hegde/valtro-backend/internal/models"
	userRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/user"
)

// Service handles user business logic
type Service struct {
	userRepo *userRepo.Repository
}

// NewService creates a new user service
func NewService(userRepo *userRepo.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business logic validation
func (s *Service) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Validate business rules
	if err := s.validateCreateUser(req); err != nil {
		return nil, err
	}

	// Check if user already exists by Clerk user ID
	exists, err := s.userRepo.ClerkUserIDExists(req.ClerkUserID)
	if err != nil {
		return nil, err // Repository now returns structured errors
	}
	if exists {
		return nil, errors.NewConflictError("User already exists with this Clerk ID", "Clerk ID: "+req.ClerkUserID)
	}

	// Create user model
	user := &models.User{
		ID:          uuid.New(),
		ClerkUserID: req.ClerkUserID,
		Email:       strings.TrimSpace(req.Email),
		FullName:    strings.TrimSpace(req.FullName),
		Username:    req.Username,
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err // Repository now returns structured errors
	}

	// Convert to response DTO
	return s.toUserResponse(user), nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

// GetUserByClerkID retrieves a user by Clerk user ID
func (s *Service) GetUserByClerkID(clerkUserID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByClerkUserID(clerkUserID)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

// GetAllUsers retrieves all users with pagination
func (s *Service) GetAllUsers(limit, offset int) ([]*dto.UserResponse, error) {
	// Set default pagination limits
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err // Repository now returns structured errors
	}

	// Convert to response DTOs
	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, s.toUserResponse(user))
	}

	return responses, nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FullName != nil {
		user.FullName = strings.TrimSpace(*req.FullName)
	}
	if req.Username != nil {
		// Check if username is already taken by another user
		if *req.Username != "" {
			existingUser, err := s.userRepo.GetByUsername(*req.Username)
			if err == nil && existingUser.ID != user.ID {
				return nil, errors.NewConflictError("Username already taken", "Username: "+*req.Username)
			}
		}
		user.Username = req.Username
	}
	if req.ImageURL != nil {
		user.ImageURL = *req.ImageURL
	}

	user.UpdatedAt = time.Now()

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, err // Repository now returns structured errors
	}

	return s.toUserResponse(user), nil
}

// DeleteUser soft deletes a user
func (s *Service) DeleteUser(id uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete
	if err := s.userRepo.Delete(id); err != nil {
		return err // Repository now returns structured errors
	}

	return nil
}

// UpdateLastSignIn updates the user's last sign-in time
func (s *Service) UpdateLastSignIn(clerkUserID string) error {
	user, err := s.userRepo.GetByClerkUserID(clerkUserID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.LastSignIn = &now
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

// validateCreateUser validates the create user request
func (s *Service) validateCreateUser(req dto.CreateUserRequest) error {
	if strings.TrimSpace(req.ClerkUserID) == "" {
		return errors.NewValidationError("Clerk user ID is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.NewValidationError("Email is required")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return errors.NewValidationError("Full name is required")
	}
	if len(req.FullName) < 2 {
		return errors.NewValidationError("Full name must be at least 2 characters")
	}
	if req.Username != nil && *req.Username != "" && len(*req.Username) < 3 {
		return errors.NewValidationError("Username must be at least 3 characters")
	}
	return nil
}

// toUserResponse converts a user model to response DTO
func (s *Service) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:            user.ID,
		ClerkUserID:   user.ClerkUserID,
		Email:         user.Email,
		FullName:      user.FullName,
		Username:      user.Username,
		ImageURL:      user.ImageURL,
		EmailVerified: user.EmailVerified,
		Active:        user.Active,
		LastSignIn:    user.LastSignIn,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

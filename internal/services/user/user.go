package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
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

	// Check if user already exists
	exists, err := s.userRepo.ClerkUserIDExists(req.ClerkUserID)
	if err != nil {
		return nil, errors.New("failed to check user existence")
	}
	if exists {
		return nil, errors.New("user with this Clerk ID already exists")
	}

	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, errors.New("failed to check email existence")
	}
	if emailExists {
		return nil, errors.New("user with this email already exists")
	}

	// Check if username already exists (if provided)
	if req.Username != nil && *req.Username != "" {
		usernameExists, err := s.userRepo.UsernameExists(*req.Username)
		if err != nil {
			return nil, errors.New("failed to check username existence")
		}
		if usernameExists {
			return nil, errors.New("username already taken")
		}
	}

	// Create user model
	user := &models.User{
		ID:            uuid.New(), // Generate UUID
		ClerkUserID:   req.ClerkUserID,
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		FullName:      strings.TrimSpace(req.FullName),
		Username:      req.Username,
		ImageURL:      req.ImageURL,
		EmailVerified: req.EmailVerified,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
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
		return nil, errors.New("failed to retrieve users")
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
				return nil, errors.New("username already taken")
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
		return nil, errors.New("failed to update user")
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
		return errors.New("failed to delete user")
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
		return errors.New("clerk user ID is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(req.FullName) == "" {
		return errors.New("full name is required")
	}
	if len(req.FullName) < 2 {
		return errors.New("full name must be at least 2 characters")
	}
	if req.Username != nil && *req.Username != "" && len(*req.Username) < 3 {
		return errors.New("username must be at least 3 characters")
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

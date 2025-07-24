package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	ClerkUserID   string  `json:"clerk_user_id" validate:"required"`
	Email         string  `json:"email" validate:"required,email"`
	FullName      string  `json:"full_name" validate:"required,min=2,max=255"`
	Username      *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	ImageURL      string  `json:"image_url,omitempty" validate:"omitempty,url"`
	EmailVerified bool    `json:"email_verified"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	ImageURL *string `json:"image_url,omitempty" validate:"omitempty,url"`
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	ClerkUserID   string     `json:"clerk_user_id"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name"`
	Username      *string    `json:"username"`
	ImageURL      string     `json:"image_url"`
	EmailVerified bool       `json:"email_verified"`
	Active        bool       `json:"active"`
	LastSignIn    *time.Time `json:"last_sign_in"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

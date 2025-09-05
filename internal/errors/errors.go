package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents different types of application errors
type ErrorType string

const (
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeServiceUnavailable ErrorType = "SERVICE_UNAVAILABLE"
)

// AppError represents a structured application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Code    string    `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// HTTPStatus returns the appropriate HTTP status code for the error type
func (e *AppError) HTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Error constructors

// NewValidationError creates a validation error
func NewValidationError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: detail,
		Code:    "VAL_001",
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: detail,
		Code:    "NF_001",
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Details: detail,
		Code:    "AUTH_001",
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Details: detail,
		Code:    "FORB_001",
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Details: detail,
		Code:    "CONF_001",
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Details: detail,
		Code:    "INT_001",
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Details: detail,
		Code:    "BR_001",
	}
}

// Helper functions to check error types

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsUnauthorizedError checks if error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeUnauthorized
	}
	return false
}

// IsConflictError checks if error is a conflict error
func IsConflictError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == ErrorTypeConflict
	}
	return false
}

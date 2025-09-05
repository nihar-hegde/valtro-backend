package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nihar-hegde/valtro-backend/internal/constants"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// ValidateRequired checks if a string field is not empty after trimming
func ValidateRequired(field, fieldName string) error {
	if strings.TrimSpace(field) == "" {
		return ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length requirement
func ValidateMinLength(field, fieldName string, minLength int) error {
	if len(strings.TrimSpace(field)) < minLength {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, minLength),
		}
	}
	return nil
}

// ValidateMaxLength checks if a string meets maximum length requirement
func ValidateMaxLength(field, fieldName string, maxLength int) error {
	if len(strings.TrimSpace(field)) > maxLength {
		return ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("%s must be at most %d characters long", fieldName, maxLength),
		}
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return ValidationError{Field: "email", Message: "email is required"}
	}
	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "email must be a valid email address"}
	}
	return nil
}

// ValidateOptionalMinLength validates optional fields that have minimum length when provided
func ValidateOptionalMinLength(field *string, fieldName string, minLength int) error {
	if field != nil && *field != "" {
		return ValidateMinLength(*field, fieldName, minLength)
	}
	return nil
}

// ValidateCreateUserRequest validates user creation request
func ValidateCreateUserRequest(req interface{}) error {
	// For now, return nil to maintain compatibility
	// TODO: Implement proper struct validation when go-playground/validator is available
	return nil
}

// ValidateCreateOrganizationRequest validates organization creation request
func ValidateCreateOrganizationRequest(name string) error {
	var errs []error
	
	if err := ValidateRequired(name, "organization name"); err != nil {
		errs = append(errs, err)
	}
	
	if err := ValidateMinLength(name, "organization name", constants.MinOrganizationNameLength); err != nil {
		errs = append(errs, err)
	}
	
	if err := ValidateMaxLength(name, "organization name", constants.MaxOrganizationNameLength); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		var messages []string
		for _, err := range errs {
			messages = append(messages, err.Error())
		}
		return errors.New(strings.Join(messages, "; "))
	}
	
	return nil
}

// ValidateCreateProjectRequest validates project creation request
func ValidateCreateProjectRequest(name string) error {
	var errs []error
	
	if err := ValidateRequired(name, "project name"); err != nil {
		errs = append(errs, err)
	}
	
	if err := ValidateMinLength(name, "project name", constants.MinProjectNameLength); err != nil {
		errs = append(errs, err)
	}
	
	if err := ValidateMaxLength(name, "project name", constants.MaxProjectNameLength); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		var messages []string
		for _, err := range errs {
			messages = append(messages, err.Error())
		}
		return errors.New(strings.Join(messages, "; "))
	}
	
	return nil
}

// SanitizeString trims whitespace
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeEmail trims whitespace and converts email to lowercase
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
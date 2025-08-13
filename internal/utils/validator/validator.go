package validator

import (
    "errors"
    "regexp"
    "strings"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateRequired checks if a string field is not empty after trimming
func ValidateRequired(field, fieldName string) error {
    if strings.TrimSpace(field) == "" {
        return errors.New(fieldName + " is required")
    }
    return nil
}

// ValidateMinLength checks if a string meets minimum length requirement
func ValidateMinLength(field, fieldName string, minLength int) error {
    if len(strings.TrimSpace(field)) < minLength {
        return errors.New(fieldName + " must be at least " + string(rune(minLength+'0')) + " characters")
    }
    return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" {
        return errors.New("email is required")
    }
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
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

// SanitizeString trims whitespace and converts to lowercase if needed
func SanitizeString(s string) string {
    return strings.TrimSpace(s)
}

// SanitizeEmail trims whitespace and converts email to lowercase
func SanitizeEmail(email string) string {
    return strings.ToLower(strings.TrimSpace(email))
}
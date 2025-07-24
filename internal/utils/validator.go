package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidURL checks if a URL is valid
func IsValidURL(url string) bool {
	if url == "" {
		return true // Empty URL is valid (optional field)
	}
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// SanitizeString trims whitespace and converts to lowercase if needed
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// SanitizeEmail trims whitespace and converts to lowercase
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

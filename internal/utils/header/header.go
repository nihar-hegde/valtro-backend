package header

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/errors"
)

// ExtractUserID extracts and validates user ID from X-User-ID header
func ExtractUserID(r *http.Request) (uuid.UUID, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil, errors.NewUnauthorizedError("User ID required", "X-User-ID header is missing")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.NewValidationError("Invalid user ID format", err.Error())
	}

	return userID, nil
}

// ExtractClerkUserID extracts Clerk user ID from X-Clerk-User-ID header
func ExtractClerkUserID(r *http.Request) (string, error) {
	clerkUserID := r.Header.Get("X-Clerk-User-ID")
	if clerkUserID == "" {
		return "", errors.NewUnauthorizedError("Clerk User ID required", "X-Clerk-User-ID header is missing")
	}
	return clerkUserID, nil
}

// ExtractBearerToken extracts and validates Bearer token from Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.NewUnauthorizedError("Authorization header required", "Authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.NewValidationError("Invalid authorization format", "Expected 'Bearer <token>' format")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.NewValidationError("Empty token", "Bearer token cannot be empty")
	}

	return token, nil
}

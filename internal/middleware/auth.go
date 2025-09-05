package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
)

// RequireUserID middleware extracts and validates X-User-ID header
func RequireUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			response.SendUnauthorized(w, "User ID required")
			return
		}

		// Validate UUID format
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response.SendValidationError(w, "Invalid user ID: ID must be a valid UUID")
			return
		}

		// Add to context for use in handlers
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireOrganizationID middleware extracts and validates X-Organization-ID header
func RequireOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgIDStr := r.Header.Get("X-Organization-ID")
		if orgIDStr == "" {
			response.SendUnauthorized(w, "Organization ID required")
			return
		}

		// Validate UUID format
		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			response.SendValidationError(w, "Invalid organization ID: ID must be a valid UUID")
			return
		}

		// Add to context for use in handlers
		ctx := context.WithValue(r.Context(), "organizationID", orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireClerkUserID middleware extracts and validates X-Clerk-User-ID header
func RequireClerkUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clerkUserID := r.Header.Get("X-Clerk-User-ID")
		if clerkUserID == "" {
			response.SendUnauthorized(w, "User not authenticated")
			return
		}

		// Add to context for use in handlers
		ctx := context.WithValue(r.Context(), "clerkUserID", clerkUserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	return userID, ok
}

// GetOrganizationIDFromContext extracts organization ID from request context
func GetOrganizationIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	orgID, ok := ctx.Value("organizationID").(uuid.UUID)
	return orgID, ok
}

// GetClerkUserIDFromContext extracts Clerk user ID from request context
func GetClerkUserIDFromContext(ctx context.Context) (string, bool) {
	clerkUserID, ok := ctx.Value("clerkUserID").(string)
	return clerkUserID, ok
}

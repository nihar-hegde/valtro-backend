package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/onboarding"
	"github.com/nihar-hegde/valtro-backend/internal/middleware"
	"gorm.io/gorm"
)

// RegisterOnboardingRoutes registers all onboarding-related routes
func RegisterOnboardingRoutes(r chi.Router, db *gorm.DB, onboardingHandler *onboarding.Handler) {
	r.Route("/onboarding", func(r chi.Router) {
		// Apply Clerk JWT authentication to onboarding routes
		r.Use(middleware.ClerkJWTMiddleware(db))
		
		r.Post("/", onboardingHandler.CompleteOnboarding) // POST /api/v1/onboarding
	})
}

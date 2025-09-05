package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/onboarding"
)

// RegisterOnboardingRoutes registers all onboarding-related routes
func RegisterOnboardingRoutes(r chi.Router, onboardingHandler *onboarding.Handler) {
	r.Post("/onboarding", onboardingHandler.CompleteOnboarding) // POST /api/v1/onboarding
}

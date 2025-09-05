package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nihar-hegde/valtro-backend/internal/server/routes"
)

// corsMiddleware handles CORS headers for all requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-User-ID, X-Organization-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes sets up all the routes for the server.
func (s *Server) RegisterRoutes() {
	// Configure CORS middleware (must be first)
	s.router.Use(corsMiddleware)
	
	// Use standard chi middleware
	s.router.Use(middleware.Logger)    // Logs request details
	s.router.Use(middleware.Recoverer) // Recovers from panics

	// Root welcome endpoint
	s.router.Get("/", s.welcomeHandler)

	// Health check endpoint
	s.router.Get("/health", s.healthHandler.HealthCheck)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// User routes
		routes.RegisterUserRoutes(r, s.db, s.userHandler)
		
		// Organization routes
		routes.RegisterOrganizationRoutes(r, s.db, s.orgHandler)
		
		// Project routes
		routes.RegisterProjectRoutes(r, s.db, s.projectHandler)
		
		// Onboarding routes
		routes.RegisterOnboardingRoutes(r, s.db, s.onboardingHandler)
	})

	// Webhook routes (outside of API versioning as they're called by external services)
	s.router.Route("/api", func(r chi.Router) {
		routes.RegisterWebhookRoutes(r, s.webhookHandler)
	})
}

// welcomeHandler handles the root route.
func (s *Server) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the Valtro API!"})
}
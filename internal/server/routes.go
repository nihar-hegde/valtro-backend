package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nihar-hegde/valtro-backend/internal/server/routes"
)

// RegisterRoutes sets up all the routes for the server.
func (s *Server) RegisterRoutes() {
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
		routes.RegisterUserRoutes(r, s.userHandler)
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
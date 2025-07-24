package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		r.Route("/users", func(r chi.Router) {
			r.Post("/", s.userHandler.Create)        // POST /api/v1/users
			r.Get("/", s.userHandler.GetAll)         // GET /api/v1/users
			r.Get("/profile", s.userHandler.GetProfile) // GET /api/v1/users/profile
			r.Get("/{id}", s.userHandler.GetByID)     // GET /api/v1/users/{id}
			r.Put("/{id}", s.userHandler.Update)      // PUT /api/v1/users/{id}
			r.Delete("/{id}", s.userHandler.Delete)   // DELETE /api/v1/users/{id}
		})

		// Organization routes
		r.Route("/organization", func(r chi.Router) {
			r.Post("/", s.orgHandler.Create)         // POST /api/v1/organization
			r.Get("/", s.orgHandler.GetAll)          // GET /api/v1/organization
			r.Get("/{id}", s.orgHandler.GetByID)     // GET /api/v1/organization/{id}
			r.Put("/{id}", s.orgHandler.Update)      // PUT /api/v1/organization/{id}
			r.Delete("/{id}", s.orgHandler.Delete)   // DELETE /api/v1/organization/{id}
		})

		// Project routes
		r.Route("/projects", func(r chi.Router) {
			r.Post("/", s.projectHandler.Create)     // POST /api/v1/projects
			r.Get("/", s.projectHandler.GetAll)      // GET /api/v1/projects
			r.Get("/{id}", s.projectHandler.GetByID) // GET /api/v1/projects/{id}
			r.Put("/{id}", s.projectHandler.Update)  // PUT /api/v1/projects/{id}
			r.Delete("/{id}", s.projectHandler.Delete) // DELETE /api/v1/projects/{id}
		})
	})
}

// welcomeHandler handles the root route.
func (s *Server) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the Valtro API!"})
}
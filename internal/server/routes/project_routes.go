package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/project"
	"github.com/nihar-hegde/valtro-backend/internal/middleware"
	"gorm.io/gorm"
)

// RegisterProjectRoutes registers all project-related routes
func RegisterProjectRoutes(r chi.Router, db *gorm.DB, projectHandler *project.Handler) {
	r.Route("/projects", func(r chi.Router) {
		// Apply Clerk JWT authentication to all project routes
		r.Use(middleware.ClerkJWTMiddleware(db))
		
		r.Post("/", projectHandler.Create)                                           // POST /api/v1/projects
		r.Get("/{id}", projectHandler.GetByID)                                       // GET /api/v1/projects/{id}
		r.Get("/by-api-key/{apiKey}", projectHandler.GetByAPIKey)                    // GET /api/v1/projects/by-api-key/{apiKey}
		r.Get("/organization/{organizationId}", projectHandler.GetByOrganization)    // GET /api/v1/projects/organization/{organizationId}
		r.Put("/{id}", projectHandler.Update)                                        // PUT /api/v1/projects/{id}
		r.Delete("/{id}", projectHandler.Delete)                                     // DELETE /api/v1/projects/{id}
		r.Post("/{id}/regenerate-api-key", projectHandler.RegenerateAPIKey)          // POST /api/v1/projects/{id}/regenerate-api-key
	})
}

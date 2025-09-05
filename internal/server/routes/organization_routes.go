package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/organization"
	"github.com/nihar-hegde/valtro-backend/internal/middleware"
	"gorm.io/gorm"
)

// RegisterOrganizationRoutes registers all organization-related routes
func RegisterOrganizationRoutes(r chi.Router, db *gorm.DB, orgHandler *organization.Handler) {
	r.Route("/organizations", func(r chi.Router) {
		// Apply Clerk JWT authentication to all organization routes
		r.Use(middleware.ClerkJWTMiddleware(db))
		
		r.Post("/", orgHandler.Create)                      // POST /api/v1/organizations
		r.Get("/", orgHandler.GetAll)                       // GET /api/v1/organizations
		r.Get("/check", orgHandler.CheckUserOrganization)   // GET /api/v1/organizations/check
		r.Get("/with-projects", orgHandler.GetWithProjects) // GET /api/v1/organizations/with-projects
		r.Get("/{id}", orgHandler.GetByID)                  // GET /api/v1/organizations/{id}
		r.Put("/{id}", orgHandler.Update)                   // PUT /api/v1/organizations/{id}
		r.Delete("/{id}", orgHandler.Delete)                // DELETE /api/v1/organizations/{id}
	})
}

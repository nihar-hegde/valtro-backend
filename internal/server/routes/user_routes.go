package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/user"
)

func RegisterUserRoutes(r chi.Router, userHandler *user.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)           // POST /api/v1/users
		r.Get("/", userHandler.GetAll)            // GET /api/v1/users
		r.Get("/profile", userHandler.GetProfile) // GET /api/v1/users/profile
		r.Get("/{id}", userHandler.GetByID)       // GET /api/v1/users/{id}
		r.Put("/{id}", userHandler.Update)        // PUT /api/v1/users/{id}
		r.Delete("/{id}", userHandler.Delete)     // DELETE /api/v1/users/{id}
	})
}

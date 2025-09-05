package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/health"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/onboarding"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/organization"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/project"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/user"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/webhook"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server holds the dependencies for our HTTP server.
type Server struct {
	db                 *gorm.DB
	router             *chi.Mux
	healthHandler      *health.Handler
	userHandler        *user.Handler
	orgHandler         *organization.Handler
	projectHandler     *project.Handler
	webhookHandler     *webhook.Handler
	onboardingHandler  *onboarding.Handler
}

// NewServer creates a new Server instance.
func NewServer(db *gorm.DB) *Server {
	server := &Server{
		db:                db,
		router:            chi.NewRouter(),
		healthHandler:     health.NewHandler(db),
		userHandler:       user.NewHandler(db),
		orgHandler:        organization.NewHandler(db),
		projectHandler:    project.NewHandler(db),
		webhookHandler:    webhook.NewHandler(db),
		onboardingHandler: onboarding.NewHandler(db),
	}

	// Register all the application routes.
	server.RegisterRoutes()

	return server
}



// Start runs the HTTP server.
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s", addr)

	return http.ListenAndServe(addr, s.router)
}

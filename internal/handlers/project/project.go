package project

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// Handler handles project-related HTTP requests.
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new project handler.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Create handles POST /api/v1/projects
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Create project endpoint"})
}

// GetAll handles GET /api/v1/projects
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Get all projects endpoint"})
}

// GetByID handles GET /api/v1/projects/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Get project by ID endpoint"})
}

// Update handles PUT /api/v1/projects/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Update project endpoint"})
}

// Delete handles DELETE /api/v1/projects/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete project endpoint"})
}

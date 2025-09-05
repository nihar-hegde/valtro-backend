package project

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	projectRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/project"
	projectService "github.com/nihar-hegde/valtro-backend/internal/services/project"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"gorm.io/gorm"
)

// Handler handles project-related HTTP requests
type Handler struct {
	projectService *projectService.Service
}

// NewHandler creates a new project handler
func NewHandler(db *gorm.DB) *Handler {
	projectRepository := projectRepo.NewRepository(db)
	projectSvc := projectService.NewService(projectRepository)

	return &Handler{
		projectService: projectSvc,
	}
}

// Create handles POST /api/v1/projects
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Create project through service
	project, err := h.projectService.CreateProject(req)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to create project", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusCreated, "Project created successfully", project)
}

// GetByID handles GET /api/v1/projects/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Get project through service
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.SendNotFound(w, "Project")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project retrieved successfully", project)
}

// GetByAPIKey handles GET /api/v1/projects/by-api-key/{apiKey}
func (h *Handler) GetByAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get API key from URL
	apiKey := chi.URLParam(r, "apiKey")
	if apiKey == "" {
		response.SendValidationError(w, "API key is required")
		return
	}

	// Get project through service
	project, err := h.projectService.GetProjectByAPIKey(apiKey)
	if err != nil {
		response.SendNotFound(w, "Project")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project retrieved successfully", project)
}

// GetByOrganization handles GET /api/v1/projects/organization/{organizationId}
func (h *Handler) GetByOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get organization ID from URL
	orgIDStr := chi.URLParam(r, "organizationId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Get projects for organization
	projects, err := h.projectService.GetProjectsByOrganization(orgID)
	if err != nil {
		response.SendInternalError(w, "Failed to retrieve projects: "+err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Projects retrieved successfully", projects)
}

// Update handles PUT /api/v1/projects/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get organization ID from header (for authorization)
	orgIDStr := r.Header.Get("X-Organization-ID")
	if orgIDStr == "" {
		response.SendUnauthorized(w, "Organization ID required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Parse request body
	var req dto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Update project through service
	project, err := h.projectService.UpdateProject(projectID, req, orgID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to update project", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project updated successfully", project)
}

// Delete handles DELETE /api/v1/projects/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get organization ID from header (for authorization)
	orgIDStr := r.Header.Get("X-Organization-ID")
	if orgIDStr == "" {
		response.SendUnauthorized(w, "Organization ID required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Delete project through service
	if err := h.projectService.DeleteProject(projectID, orgID); err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to delete project", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project deleted successfully", nil)
}

// RegenerateAPIKey handles POST /api/v1/projects/{id}/regenerate-api-key
func (h *Handler) RegenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get organization ID from header (for authorization)
	orgIDStr := r.Header.Get("X-Organization-ID")
	if orgIDStr == "" {
		response.SendUnauthorized(w, "Organization ID required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Regenerate API key through service
	project, err := h.projectService.RegenerateAPIKey(projectID, orgID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to regenerate API key", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "API key regenerated successfully", project)
}


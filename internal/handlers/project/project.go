package project

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	orgRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/organization"
	projectRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/project"
	orgService "github.com/nihar-hegde/valtro-backend/internal/services/organization"
	projectService "github.com/nihar-hegde/valtro-backend/internal/services/project"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"gorm.io/gorm"
)

// Handler handles project-related HTTP requests
type Handler struct {
	projectService *projectService.Service
	orgService     *orgService.Service
}

// NewHandler creates a new project handler
func NewHandler(db *gorm.DB) *Handler {
	projectRepository := projectRepo.NewRepository(db)
	projectSvc := projectService.NewService(projectRepository)
	
	orgRepository := orgRepo.NewRepository(db)
	orgSvc := orgService.NewService(orgRepository)

	return &Handler{
		projectService: projectSvc,
		orgService:     orgSvc,
	}
}

// validateProjectOwnership is a DRY helper function to validate if user owns the project's organization
func (h *Handler) validateProjectOwnership(w http.ResponseWriter, r *http.Request, projectID uuid.UUID) (uuid.UUID, bool) {
	// Get current user ID from JWT middleware
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return uuid.Nil, false
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid current user ID: "+err.Error())
		return uuid.Nil, false
	}

	// Get project to find its organization
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.SendNotFound(w, "Project")
		return uuid.Nil, false
	}

	// Verify user owns the organization
	organization, err := h.orgService.GetOrganizationByID(project.OrganizationID)
	if err != nil {
		response.SendNotFound(w, "Organization")
		return uuid.Nil, false
	}

	if organization.OwnerID != currentUserID {
		response.SendForbidden(w, "You can only access projects for organizations you own")
		return uuid.Nil, false
	}

	return currentUserID, true
}

// validateOrganizationOwnership is a DRY helper function to validate if user owns the organization
func (h *Handler) validateOrganizationOwnership(w http.ResponseWriter, r *http.Request, orgID uuid.UUID) (uuid.UUID, bool) {
	// Get current user ID from JWT middleware
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return uuid.Nil, false
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid current user ID: "+err.Error())
		return uuid.Nil, false
	}

	// Verify user owns the organization
	organization, err := h.orgService.GetOrganizationByID(orgID)
	if err != nil {
		response.SendNotFound(w, "Organization")
		return uuid.Nil, false
	}

	if organization.OwnerID != currentUserID {
		response.SendForbidden(w, "You can only access organizations you own")
		return uuid.Nil, false
	}

	return currentUserID, true
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

	// Authorization: Verify user owns the organization using DRY helper
	_, valid := h.validateOrganizationOwnership(w, r, req.OrganizationID)
	if !valid {
		return // Response already sent by helper
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

	// Authorization: Verify user owns the project's organization using DRY helper
	_, valid := h.validateProjectOwnership(w, r, projectID)
	if !valid {
		return // Response already sent by helper
	}

	// Get project through service (we know it exists from validation)
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

	// Authorization: Verify user owns the project's organization using DRY helper
	_, valid := h.validateProjectOwnership(w, r, project.ID)
	if !valid {
		return // Response already sent by helper
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

	// Authorization: Verify user owns the organization using DRY helper
	_, valid := h.validateOrganizationOwnership(w, r, orgID)
	if !valid {
		return // Response already sent by helper
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

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Authorization: Verify user owns the project's organization using DRY helper
	_, valid := h.validateProjectOwnership(w, r, projectID)
	if !valid {
		return // Response already sent by helper
	}

	// Parse request body
	var req dto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Get project to find organization ID for service call
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.SendNotFound(w, "Project")
		return
	}

	// Update project through service
	updatedProject, err := h.projectService.UpdateProject(projectID, req, project.OrganizationID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to update project", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project updated successfully", updatedProject)
}

// Delete handles DELETE /api/v1/projects/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Authorization: Verify user owns the project's organization using DRY helper
	_, valid := h.validateProjectOwnership(w, r, projectID)
	if !valid {
		return // Response already sent by helper
	}

	// Get project to find organization ID for service call
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.SendNotFound(w, "Project")
		return
	}

	// Delete project through service
	if err := h.projectService.DeleteProject(projectID, project.OrganizationID); err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to delete project", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Project deleted successfully", nil)
}

// RegenerateAPIKey handles POST /api/v1/projects/{id}/regenerate-api-key
func (h *Handler) RegenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid project ID: "+err.Error())
		return
	}

	// Authorization: Verify user owns the project's organization using DRY helper
	_, valid := h.validateProjectOwnership(w, r, projectID)
	if !valid {
		return // Response already sent by helper
	}

	// Get project to find organization ID for service call
	project, err := h.projectService.GetProjectByID(projectID)
	if err != nil {
		response.SendNotFound(w, "Project")
		return
	}

	// Regenerate API key through service
	updatedProject, err := h.projectService.RegenerateAPIKey(projectID, project.OrganizationID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to regenerate API key", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "API key regenerated successfully", updatedProject)
}


package organization

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nihar-hegde/valtro-backend/internal/dto"
	orgRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/organization"
	orgService "github.com/nihar-hegde/valtro-backend/internal/services/organization"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"github.com/nihar-hegde/valtro-backend/internal/utils/validator"
)

// Handler handles organization-related HTTP requests
type Handler struct {
	orgService *orgService.Service
}

// NewHandler creates a new organization handler
func NewHandler(db *gorm.DB) *Handler {
	orgRepository := orgRepo.NewRepository(db)
	orgSvc := orgService.NewService(orgRepository)

	return &Handler{
		orgService: orgSvc,
	}
}

// Create handles POST /api/v1/organizations
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For now, we'll get the user ID from header (in real app, this would come from JWT middleware)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Parse request body
	var req dto.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Validate request
	if err := validator.ValidateCreateOrganizationRequest(req.Name); err != nil {
		response.SendValidationError(w, err.Error())
		return
	}

	// Sanitize input
	req.Name = validator.SanitizeString(req.Name)

	// Create organization through service
	organization, err := h.orgService.CreateOrganization(req, userID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to create organization", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusCreated, "Organization created successfully", organization)
}

// GetAll handles GET /api/v1/organizations
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Get organizations for user
	organizations, err := h.orgService.GetOrganizationsByOwner(userID)
	if err != nil {
		response.SendInternalError(w, "Failed to retrieve organizations: "+err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Organizations retrieved successfully", organizations)
}

// GetByID handles GET /api/v1/organizations/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get organization ID from URL
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Get organization through service
	organization, err := h.orgService.GetOrganizationByID(orgID)
	if err != nil {
		response.SendNotFound(w, "Organization")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Organization retrieved successfully", organization)
}

// GetWithProjects handles GET /api/v1/organizations/with-projects
func (h *Handler) GetWithProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Get organization with projects
	organization, err := h.orgService.GetOrganizationWithProjects(userID)
	if err != nil {
		response.SendNotFound(w, "Organization")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Organization with projects retrieved successfully", organization)
}

// Update handles PUT /api/v1/organizations/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Get organization ID from URL
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Parse request body
	var req dto.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Validate request
	if req.Name != nil {
		if err := validator.ValidateCreateOrganizationRequest(*req.Name); err != nil {
			response.SendValidationError(w, err.Error())
			return
		}
		// Sanitize input
		sanitized := validator.SanitizeString(*req.Name)
		req.Name = &sanitized
	}

	// Update organization through service
	organization, err := h.orgService.UpdateOrganization(orgID, req, userID)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to update organization", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Organization updated successfully", organization)
}

// Delete handles DELETE /api/v1/organizations/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Get organization ID from URL
	orgIDStr := chi.URLParam(r, "id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid organization ID: "+err.Error())
		return
	}

	// Delete organization through service
	if err := h.orgService.DeleteOrganization(orgID, userID); err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to delete organization", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Organization deleted successfully", nil)
}

// CheckUserOrganization handles GET /api/v1/organizations/check
func (h *Handler) CheckUserOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from header
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: "+err.Error())
		return
	}

	// Check user organization
	result, err := h.orgService.CheckUserOrganization(userID)
	if err != nil {
		response.SendInternalError(w, "Failed to check user organization: "+err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "User organization check completed", result)
}


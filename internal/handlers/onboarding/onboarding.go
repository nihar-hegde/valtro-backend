package onboarding

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	orgRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/organization"
	projectRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/project"
	orgService "github.com/nihar-hegde/valtro-backend/internal/services/organization"
	projectService "github.com/nihar-hegde/valtro-backend/internal/services/project"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"gorm.io/gorm"
)

// Handler handles onboarding-related HTTP requests
type Handler struct {
	db             *gorm.DB
	orgService     *orgService.Service
	projectService *projectService.Service
}

// NewHandler creates a new onboarding handler
func NewHandler(db *gorm.DB) *Handler {
	orgRepository := orgRepo.NewRepository(db)
	projectRepository := projectRepo.NewRepository(db)
	
	orgSvc := orgService.NewService(orgRepository)
	projectSvc := projectService.NewService(projectRepository)

	return &Handler{
		db:             db,
		orgService:     orgSvc,
		projectService: projectSvc,
	}
}

// CompleteOnboarding handles POST /api/v1/onboarding
// This endpoint creates both an organization and a project in a single transaction
func (h *Handler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var req dto.OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Start database transaction
	tx := h.db.Begin()
	if tx.Error != nil {
		response.SendInternalError(w, "Failed to start transaction: "+tx.Error.Error())
		return
	}

	// Rollback transaction if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create organization repositories and services with transaction
	orgTxRepo := orgRepo.NewRepository(tx)
	projectTxRepo := projectRepo.NewRepository(tx)
	
	orgTxService := orgService.NewService(orgTxRepo)
	projectTxService := projectService.NewService(projectTxRepo)

	// Create organization
	orgRequest := dto.CreateOrganizationRequest{
		Name: req.OrganizationName,
	}

	organization, err := orgTxService.CreateOrganization(orgRequest, userID)
	if err != nil {
		tx.Rollback()
		response.SendError(w, http.StatusBadRequest, "Failed to create organization", err.Error())
		return
	}

	// Create project
	projectRequest := dto.CreateProjectRequest{
		OrganizationID: organization.ID,
		Name:           req.ProjectName,
	}

	project, err := projectTxService.CreateProject(projectRequest)
	if err != nil {
		tx.Rollback()
		response.SendError(w, http.StatusBadRequest, "Failed to create project", err.Error())
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response.SendInternalError(w, "Failed to complete onboarding: "+err.Error())
		return
	}

	// Prepare response
	onboardingResponse := dto.OnboardingResponse{
		Organization: *organization,
		Project:      *project,
	}

	// Send success response
	response.SendSuccess(w, http.StatusCreated, "Onboarding completed successfully", onboardingResponse)
}


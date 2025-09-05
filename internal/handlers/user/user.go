package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
	userRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/user"
	userService "github.com/nihar-hegde/valtro-backend/internal/services/user"
	"github.com/nihar-hegde/valtro-backend/internal/utils/response"
	"gorm.io/gorm"
)

// Handler handles user-related HTTP requests.
type Handler struct {
	userService *userService.Service
}

// NewHandler creates a new user handler.
func NewHandler(db *gorm.DB) *Handler {
	userRepository := userRepo.NewRepository(db)
	userSvc := userService.NewService(userRepository)

	return &Handler{
		userService: userSvc,
	}
}

// Create handles POST /api/v1/users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Create user through service
	user, err := h.userService.CreateUser(req)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusCreated, "User created successfully", user)
}

// GetAll handles GET /api/v1/users
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters for pagination
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	offset := 0 // default

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Get users through service
	users, err := h.userService.GetAllUsers(limit, offset)
	if err != nil {
		response.SendInternalError(w, "Failed to retrieve users: "+err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Users retrieved successfully", users)
}

// GetByID handles GET /api/v1/users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: ID must be a valid UUID")
		return
	}

	// Get user through service
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.SendNotFound(w, "User")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "User retrieved successfully", user)
}

// Update handles PUT /api/v1/users/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user ID from JWT middleware (set by ClerkJWTMiddleware)
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid current user ID: "+err.Error())
		return
	}

	// Parse target user ID from URL
	targetUserIDStr := chi.URLParam(r, "id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: ID must be a valid UUID")
		return
	}

	// Authorization: Users can only update their own profile
	if currentUserID != targetUserID {
		response.SendForbidden(w, "You can only update your own profile")
		return
	}

	// Parse request body
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendValidationError(w, "Invalid request body: "+err.Error())
		return
	}

	// Update user through service
	user, err := h.userService.UpdateUser(targetUserID, req)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "User updated successfully", user)
}

// Delete handles DELETE /api/v1/users/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user ID from JWT middleware (set by ClerkJWTMiddleware)
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid current user ID: "+err.Error())
		return
	}

	// Parse target user ID from URL
	targetUserIDStr := chi.URLParam(r, "id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid user ID: ID must be a valid UUID")
		return
	}

	// Authorization: Users can only delete their own account
	if currentUserID != targetUserID {
		response.SendForbidden(w, "You can only delete your own account")
		return
	}

	// Delete user through service
	if err := h.userService.DeleteUser(targetUserID); err != nil {
		response.SendNotFound(w, "User")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "User deleted successfully", nil)
}

// GetProfile handles GET /api/v1/users/profile
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get current user ID from JWT middleware (set by ClerkJWTMiddleware)
	currentUserIDStr := r.Header.Get("X-User-ID")
	if currentUserIDStr == "" {
		response.SendUnauthorized(w, "User ID required")
		return
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		response.SendValidationError(w, "Invalid current user ID: "+err.Error())
		return
	}

	// Get user profile
	user, err := h.userService.GetUserByID(currentUserID)
	if err != nil {
		response.SendNotFound(w, "User")
		return
	}

	// Send success response
	response.SendSuccess(w, http.StatusOK, "Profile retrieved successfully", user)
}


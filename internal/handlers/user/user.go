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
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Create user through service
	user, err := h.userService.CreateUser(req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Failed to create user", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusCreated, "User created successfully", user)
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
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

// GetByID handles GET /api/v1/users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a valid UUID")
		return
	}

	// Get user through service
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// Update handles PUT /api/v1/users/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a valid UUID")
		return
	}

	// Parse request body
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Update user through service
	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusOK, "User updated successfully", user)
}

// Delete handles DELETE /api/v1/users/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a valid UUID")
		return
	}

	// Delete user through service
	if err := h.userService.DeleteUser(id); err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "Failed to delete user", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// GetProfile handles GET /api/v1/users/profile
// This would typically get the current user's profile from JWT token
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For now, this is a placeholder
	// In a real app, you'd extract the user ID from JWT token
	clerkUserID := r.Header.Get("X-Clerk-User-ID") // Example header
	if clerkUserID == "" {
		h.sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized", "User not authenticated")
		return
	}

	// Get user by Clerk ID
	user, err := h.userService.GetUserByClerkID(clerkUserID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Send success response
	h.sendSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user)
}

// Helper methods for consistent response formatting

func (h *Handler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	response := dto.SuccessResponse{
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	w.WriteHeader(statusCode)
	response := dto.ErrorResponse{
		Error:   error,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

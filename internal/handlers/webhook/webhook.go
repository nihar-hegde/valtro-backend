package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nihar-hegde/valtro-backend/internal/dto"
	userRepo "github.com/nihar-hegde/valtro-backend/internal/repositories/user"
	userService "github.com/nihar-hegde/valtro-backend/internal/services/user"
	svix "github.com/svix/svix-webhooks/go"
	"gorm.io/gorm"
)

// Handler handles webhook-related HTTP requests
type Handler struct {
	userService *userService.Service
	svix        *svix.Webhook
}

// NewHandler creates a new webhook handler
func NewHandler(db *gorm.DB) *Handler {
	userRepository := userRepo.NewRepository(db)
	userSvc := userService.NewService(userRepository)

	// Initialize Svix webhook with signing secret
	signingSecret := os.Getenv("CLERK_WEBHOOK_SIGNING_SECRET")
	if signingSecret == "" {
		panic("CLERK_WEBHOOK_SIGNING_SECRET environment variable is required")
	}

	webhook, err := svix.NewWebhook(signingSecret)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize webhook verifier: %v", err))
	}

	return &Handler{
		userService: userSvc,
		svix:        webhook,
	}
}

// HandleClerkWebhook handles incoming Clerk webhook events
func (h *Handler) HandleClerkWebhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}

	// Get headers for signature verification
	headers := http.Header{}
	headers.Set("svix-id", r.Header.Get("svix-id"))
	headers.Set("svix-timestamp", r.Header.Get("svix-timestamp"))
	headers.Set("svix-signature", r.Header.Get("svix-signature"))

	// Verify webhook signature
	if err := h.svix.Verify(body, headers); err != nil {
		h.sendWebhookErrorResponse(w, http.StatusUnauthorized, "Invalid webhook signature", err.Error())
		return
	}

	// Parse webhook event
	var event dto.WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.sendWebhookErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	// Route to appropriate handler based on event type
	switch event.Type {
	case "user.created":
		h.handleUserCreated(w, event)
	case "user.updated":
		h.handleUserUpdated(w, event)
	case "user.deleted":
		h.handleUserDeleted(w, event)
	default:
		// For unsupported event types, return success but don't process
		h.sendWebhookSuccessResponse(w, "Event type not handled", nil)
	}
}

// handleUserCreated processes user.created events
func (h *Handler) handleUserCreated(w http.ResponseWriter, event dto.WebhookEvent) {
	fmt.Println("event: ", event)
	// Parse user data from the event
	userData, err := h.parseClerkUserData(event.Data)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusBadRequest, "Failed to parse user data", err.Error())
		return
	}

	// Convert Clerk user to CreateUserRequest
	createUserReq := h.clerkUserToCreateRequest(userData)
	fmt.Println("createUserReq: ", createUserReq)

	// Create user in database
	user, err := h.userService.CreateUser(createUserReq)

	if err != nil {
		// Check if it's a duplicate user error
		if strings.Contains(strings.ToLower(err.Error()), "already exists") {
			h.sendWebhookSuccessResponse(w, "User already exists", nil)
			return
		}
		h.sendWebhookErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	h.sendWebhookSuccessResponse(w, "User created successfully", user)
}

// handleUserUpdated processes user.updated events
func (h *Handler) handleUserUpdated(w http.ResponseWriter, event dto.WebhookEvent) {
	// Parse user data from the event
	userData, err := h.parseClerkUserData(event.Data)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusBadRequest, "Failed to parse user data", err.Error())
		return
	}

	// Get existing user by Clerk ID
	existingUser, err := h.userService.GetUserByClerkID(userData.ID)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Convert Clerk user to UpdateUserRequest
	updateUserReq := h.clerkUserToUpdateRequest(userData)

	// Update user in database
	updatedUser, err := h.userService.UpdateUser(existingUser.ID, updateUserReq)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	h.sendWebhookSuccessResponse(w, "User updated successfully", updatedUser)
}

// handleUserDeleted processes user.deleted events
func (h *Handler) handleUserDeleted(w http.ResponseWriter, event dto.WebhookEvent) {
	// Parse user data from the event
	userData, err := h.parseClerkUserData(event.Data)
	if err != nil {
		h.sendWebhookErrorResponse(w, http.StatusBadRequest, "Failed to parse user data", err.Error())
		return
	}

	// Get existing user by Clerk ID
	existingUser, err := h.userService.GetUserByClerkID(userData.ID)
	if err != nil {
		// User doesn't exist, consider it successful
		h.sendWebhookSuccessResponse(w, "User already deleted or doesn't exist", nil)
		return
	}

	// Delete user from database (soft delete)
	if err := h.userService.DeleteUser(existingUser.ID); err != nil {
		h.sendWebhookErrorResponse(w, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	h.sendWebhookSuccessResponse(w, "User deleted successfully", nil)
}

// parseClerkUserData parses the raw event data into a ClerkUser struct
func (h *Handler) parseClerkUserData(data interface{}) (*dto.ClerkUser, error) {
	// Convert interface{} to JSON and back to ClerkUser
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var clerkUser dto.ClerkUser
	if err := json.Unmarshal(jsonData, &clerkUser); err != nil {
		return nil, err
	}

	return &clerkUser, nil
}

// clerkUserToCreateRequest converts a Clerk user to CreateUserRequest
func (h *Handler) clerkUserToCreateRequest(clerkUser *dto.ClerkUser) dto.CreateUserRequest {
	var email string
	var emailVerified bool

	// Extract primary email
	for _, emailAddr := range clerkUser.EmailAddresses {
		if clerkUser.PrimaryEmailAddressID != nil && emailAddr.ID == *clerkUser.PrimaryEmailAddressID {
			email = emailAddr.EmailAddress
			emailVerified = emailAddr.Verification.Status == "verified"
			break
		}
	}

	// Fallback to first email if primary not found
	if email == "" && len(clerkUser.EmailAddresses) > 0 {
		email = clerkUser.EmailAddresses[0].EmailAddress
		emailVerified = clerkUser.EmailAddresses[0].Verification.Status == "verified"
	}

	// Fallback for test webhooks - generate email from user ID
	if email == "" {
		email = fmt.Sprintf("test+%s@clerk.dev", clerkUser.ID)
		emailVerified = false
	}

	// Construct full name
	fullName := ""
	if clerkUser.FirstName != nil && clerkUser.LastName != nil {
		fullName = strings.TrimSpace(*clerkUser.FirstName + " " + *clerkUser.LastName)
	} else if clerkUser.FirstName != nil {
		fullName = strings.TrimSpace(*clerkUser.FirstName)
	} else if clerkUser.LastName != nil {
		fullName = strings.TrimSpace(*clerkUser.LastName)
	}

	// If no name components, try to extract from email
	if fullName == "" && email != "" {
		fullName = strings.Split(email, "@")[0]
	}

	return dto.CreateUserRequest{
		ClerkUserID:   clerkUser.ID,
		Email:         email,
		FullName:      fullName,
		Username:      clerkUser.Username,
		ImageURL:      clerkUser.ImageURL,
		EmailVerified: emailVerified,
	}
}

// clerkUserToUpdateRequest converts a Clerk user to UpdateUserRequest
func (h *Handler) clerkUserToUpdateRequest(clerkUser *dto.ClerkUser) dto.UpdateUserRequest {
	// Construct full name
	var fullName *string
	if clerkUser.FirstName != nil || clerkUser.LastName != nil {
		name := ""
		if clerkUser.FirstName != nil && clerkUser.LastName != nil {
			name = strings.TrimSpace(*clerkUser.FirstName + " " + *clerkUser.LastName)
		} else if clerkUser.FirstName != nil {
			name = strings.TrimSpace(*clerkUser.FirstName)
		} else if clerkUser.LastName != nil {
			name = strings.TrimSpace(*clerkUser.LastName)
		}
		if name != "" {
			fullName = &name
		}
	}

	return dto.UpdateUserRequest{
		FullName: fullName,
		Username: clerkUser.Username,
		ImageURL: &clerkUser.ImageURL,
	}
}

// Helper methods for response formatting
func (h *Handler) sendWebhookSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.WriteHeader(http.StatusOK)
	response := dto.WebhookResponse{
		Success: true,
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) sendWebhookErrorResponse(w http.ResponseWriter, statusCode int, error, details string) {
	w.WriteHeader(statusCode)
	response := dto.WebhookResponse{
		Success: false,
		Error:   fmt.Sprintf("%s: %s", error, details),
	}
	json.NewEncoder(w).Encode(response)
}

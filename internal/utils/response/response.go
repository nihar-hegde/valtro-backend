package response

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/nihar-hegde/valtro-backend/internal/errors"
	"github.com/nihar-hegde/valtro-backend/internal/dto"
)

// SendSuccess sends a standardized success response
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := dto.SuccessResponse{
        Message: message,
        Data:    data,
    }
    json.NewEncoder(w).Encode(response)
}

// SendPaginatedSuccess sends a standardized paginated success response
func SendPaginatedSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}, pagination dto.PaginationMetadata) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := dto.PaginatedResponse{
        Message:    message,
        Data:       data,
        Pagination: pagination,
    }
    json.NewEncoder(w).Encode(response)
}

// SendError sends a standardized error response
func SendError(w http.ResponseWriter, statusCode int, message string, details string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := dto.ErrorResponse{
        Message: message,
        Error:   details,
    }
    
    json.NewEncoder(w).Encode(response)
}

// SendAppError sends a standardized error response using AppError
func SendAppError(w http.ResponseWriter, err *appErrors.AppError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.HTTPStatus())
    
    response := dto.ErrorResponse{
        Message: err.Message,
        Error:   err.Details,
    }
    
    json.NewEncoder(w).Encode(response)
}

// SendValidationError sends a 400 Bad Request with validation details
func SendValidationError(w http.ResponseWriter, message string) {
    SendError(w, http.StatusBadRequest, "Validation Error", message)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(w http.ResponseWriter, resource string) {
    SendError(w, http.StatusNotFound, "Not Found", resource+" not found")
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(w http.ResponseWriter, message string) {
    SendError(w, http.StatusUnauthorized, "Unauthorized", message)
}

// SendInternalError sends a 500 Internal Server Error response
func SendInternalError(w http.ResponseWriter, message string) {
    SendError(w, http.StatusInternalServerError, "Internal Server Error", message)
}
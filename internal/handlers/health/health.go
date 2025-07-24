package health

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// Handler handles health-related HTTP requests.
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new health handler.
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// HealthCheck performs a health check on the service and its dependencies.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		http.Error(w, "Failed to get DB instance", http.StatusInternalServerError)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		response := map[string]string{
			"status":   "error",
			"database": "unhealthy",
			"message":  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	// If everything is okay
	response := map[string]string{
		"status":   "ok",
		"database": "healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

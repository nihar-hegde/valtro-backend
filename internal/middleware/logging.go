package middleware

import (
	"net/http"
	"time"

	"github.com/nihar-hegde/valtro-backend/internal/utils/logger"
)

// RequestLogger middleware logs HTTP requests with structured logging
func RequestLogger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom ResponseWriter to capture status code
			wrapper := &responseWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Extract client IP
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				clientIP = realIP
			}

			// Call the next handler
			next.ServeHTTP(wrapper, r)

			// Log the request with structured fields
			duration := time.Since(start)
			log.LogRequest(
				r.Method,
				r.RequestURI,
				r.UserAgent(),
				clientIP,
				duration,
			)

			// Log additional request details
			fields := logger.Fields{
				"status_code": wrapper.statusCode,
				"content_length": r.ContentLength,
			}

			// Add user ID if available from JWT middleware
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				fields["user_id"] = userID
			}

			log.WithFields(fields).Debug("Request details")
		})
	}
}

// responseWrapper wraps http.ResponseWriter to capture status code
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

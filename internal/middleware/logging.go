package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger middleware logs HTTP requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture status code
		wrapper := &responseWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Log the request
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %v %s %s",
			r.Method,
			r.RequestURI,
			wrapper.statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
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

package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
)

// Logger wraps the standard slog.Logger with additional convenience methods
type Logger struct {
	*slog.Logger
}

// Fields represents structured log fields
type Fields map[string]interface{}

// New creates a new structured logger
func New() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	}

	// Use JSON handler for production-like structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)
	
	return &Logger{
		Logger: slog.New(handler),
	}
}

// NewWithLevel creates a new logger with specified log level
func NewWithLevel(level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	
	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithFields returns a logger with additional structured fields
func (l *Logger) WithFields(fields Fields) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// WithContext extracts common fields from context and returns logger with those fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := Fields{}
	
	// Extract request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}
	
	// Extract user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}
	
	// Extract trace ID if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}
	
	return l.WithFields(fields)
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(Fields{"request_id": requestID})
}

// WithUserID adds user ID to logger context
func (l *Logger) WithUserID(userID uuid.UUID) *Logger {
	return l.WithFields(Fields{"user_id": userID.String()})
}

// WithError adds error details to logger context
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(Fields{"error": err.Error()})
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(method, path, userAgent, clientIP string, duration time.Duration) {
	l.Info("HTTP request",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"client_ip", clientIP,
		"duration_ms", duration.Milliseconds(),
	)
}

// LogError logs structured error information
func (l *Logger) LogError(message string, err error, fields Fields) {
	logFields := Fields{"error": err.Error()}
	for k, v := range fields {
		logFields[k] = v
	}
	
	logger := l.WithFields(logFields)
	logger.Error(message)
}

// LogDatabaseOperation logs database operation details
func (l *Logger) LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	fields := Fields{
		"operation": operation,
		"table": table,
		"duration_ms": duration.Milliseconds(),
	}
	
	if err != nil {
		fields["error"] = err.Error()
		l.WithFields(fields).Error("Database operation failed")
	} else {
		l.WithFields(fields).Debug("Database operation completed")
	}
}

// LogAPICall logs external API call details
func (l *Logger) LogAPICall(service, endpoint string, statusCode int, duration time.Duration, err error) {
	fields := Fields{
		"service": service,
		"endpoint": endpoint,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}
	
	if err != nil {
		fields["error"] = err.Error()
		l.WithFields(fields).Error("External API call failed")
	} else {
		l.WithFields(fields).Info("External API call completed")
	}
}

// LogUserAction logs user-initiated actions
func (l *Logger) LogUserAction(userID uuid.UUID, action, resource string, fields Fields) {
	logFields := Fields{
		"user_id": userID.String(),
		"action": action,
		"resource": resource,
	}
	
	for k, v := range fields {
		logFields[k] = v
	}
	
	l.WithFields(logFields).Info("User action")
}

// LogSecurityEvent logs security-related events
func (l *Logger) LogSecurityEvent(event, userID, clientIP, details string) {
	l.WithFields(Fields{
		"event": event,
		"user_id": userID,
		"client_ip": clientIP,
		"details": details,
	}).Warn("Security event")
}

// LogBusinessEvent logs important business logic events
func (l *Logger) LogBusinessEvent(event, entityType, entityID string, fields Fields) {
	logFields := Fields{
		"event": event,
		"entity_type": entityType,
		"entity_id": entityID,
	}
	
	for k, v := range fields {
		logFields[k] = v
	}
	
	l.WithFields(logFields).Info("Business event")
}

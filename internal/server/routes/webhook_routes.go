package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nihar-hegde/valtro-backend/internal/handlers/webhook"
)

// RegisterWebhookRoutes registers all webhook-related routes
func RegisterWebhookRoutes(r chi.Router, webhookHandler *webhook.Handler) {
	// Webhook routes - no authentication needed as they're called by external services
	r.Post("/webhooks/clerk", webhookHandler.HandleClerkWebhook)
}

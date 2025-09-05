package constants

const (
	// API Key Generation Constants
	APIKeyMaxAttempts = 10
	APIKeyByteSize    = 32
	APIKeyPrefix      = "vltro_"
	
	// Validation Constants
	MinOrganizationNameLength = 2
	MaxOrganizationNameLength = 255
	MinProjectNameLength      = 2
	MaxProjectNameLength      = 255
	MinUserNameLength         = 2
	MaxUserNameLength         = 100
	
	// HTTP Status Messages
	UserIDRequired          = "User ID required"
	OrganizationIDRequired  = "Organization ID required"
	ProjectIDRequired       = "Project ID required"
	InvalidRequestBody      = "Invalid request body"
	
	// Database Query Limits
	DefaultPageSize = 20
	MaxPageSize     = 100
	
	// Error Messages
	OrganizationNotFound = "Organization not found"
	ProjectNotFound      = "Project not found"
	UserNotFound         = "User not found"
	UnauthorizedAccess   = "Unauthorized access"
	InternalServerError  = "Internal server error"
)

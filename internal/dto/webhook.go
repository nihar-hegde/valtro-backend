package dto

// WebhookEvent represents the top-level structure of a Clerk webhook event
type WebhookEvent struct {
	Data       interface{} `json:"data"`
	Object     string      `json:"object"`
	Type       string      `json:"type"`
	Timestamp  int64       `json:"timestamp"`
	InstanceID string      `json:"instance_id"`
}

// ClerkUser represents the Clerk user object structure from webhooks
type ClerkUser struct {
	ID                            string                 `json:"id"`
	Object                        string                 `json:"object"`
	Username                      *string                `json:"username"`
	FirstName                     *string                `json:"first_name"`
	LastName                      *string                `json:"last_name"`
	ImageURL                      string                 `json:"image_url"`
	HasImage                      bool                   `json:"has_image"`
	PrimaryEmailAddressID         *string                `json:"primary_email_address_id"`
	PrimaryPhoneNumberID          *string                `json:"primary_phone_number_id"`
	PrimaryWeb3WalletID           *string                `json:"primary_web3_wallet_id"`
	PasswordEnabled               bool                   `json:"password_enabled"`
	TwoFactorEnabled              bool                   `json:"two_factor_enabled"`
	TotpEnabled                   bool                   `json:"totp_enabled"`
	BackupCodeEnabled             bool                   `json:"backup_code_enabled"`
	EmailAddresses                []ClerkEmailAddress    `json:"email_addresses"`
	PhoneNumbers                  []ClerkPhoneNumber     `json:"phone_numbers"`
	Web3Wallets                   []interface{}          `json:"web3_wallets"`
	ExternalAccounts              []ClerkExternalAccount `json:"external_accounts"`
	PublicMetadata                map[string]interface{} `json:"public_metadata"`
	PrivateMetadata               map[string]interface{} `json:"private_metadata"`
	UnsafeMetadata                map[string]interface{} `json:"unsafe_metadata"`
	LastSignInAt                  *int64                 `json:"last_sign_in_at"`
	Banned                        bool                   `json:"banned"`
	Locked                        bool                   `json:"locked"`
	LockoutExpiresInMs            *int64                 `json:"lockout_expires_in_ms"`
	VerificationAttemptsRemaining int                    `json:"verification_attempts_remaining"`
	CreatedAt                     int64                  `json:"created_at"`
	UpdatedAt                     int64                  `json:"updated_at"`
}

// ClerkEmailAddress represents an email address in Clerk
type ClerkEmailAddress struct {
	ID           string            `json:"id"`
	Object       string            `json:"object"`
	EmailAddress string            `json:"email_address"`
	Reserved     bool              `json:"reserved"`
	Verification ClerkVerification `json:"verification"`
	LinkedTo     []ClerkLinkedTo   `json:"linked_to"`
	CreatedAt    int64             `json:"created_at"`
	UpdatedAt    int64             `json:"updated_at"`
}

// ClerkPhoneNumber represents a phone number in Clerk
type ClerkPhoneNumber struct {
	ID                  string            `json:"id"`
	Object              string            `json:"object"`
	PhoneNumber         string            `json:"phone_number"`
	Reserved            bool              `json:"reserved"`
	DefaultSecondFactor bool              `json:"default_second_factor"`
	Verification        ClerkVerification `json:"verification"`
	LinkedTo            []ClerkLinkedTo   `json:"linked_to"`
	BackupCodes         []interface{}     `json:"backup_codes"`
	CreatedAt           int64             `json:"created_at"`
	UpdatedAt           int64             `json:"updated_at"`
}

// ClerkExternalAccount represents an external account in Clerk
type ClerkExternalAccount struct {
	ID               string                 `json:"id"`
	Object           string                 `json:"object"`
	Provider         string                 `json:"provider"`
	IdentificationID string                 `json:"identification_id"`
	ProviderUserID   string                 `json:"provider_user_id"`
	ApprovedScopes   string                 `json:"approved_scopes"`
	EmailAddress     string                 `json:"email_address"`
	FirstName        string                 `json:"first_name"`
	LastName         string                 `json:"last_name"`
	ImageURL         string                 `json:"image_url"`
	Username         *string                `json:"username"`
	PublicMetadata   map[string]interface{} `json:"public_metadata"`
	Label            *string                `json:"label"`
	Verification     *ClerkVerification     `json:"verification"`
	CreatedAt        int64                  `json:"created_at"`
	UpdatedAt        int64                  `json:"updated_at"`
}

// ClerkVerification represents verification status in Clerk
type ClerkVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
	Attempts *int   `json:"attempts"`
	ExpireAt *int64 `json:"expire_at"`
}

// ClerkLinkedTo represents linked identification in Clerk
type ClerkLinkedTo struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// WebhookResponse represents the response structure for webhook endpoints
type WebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

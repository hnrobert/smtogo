package models

// EmailRequest represents an email sending request
type EmailRequest struct {
	RecipientEmail string `json:"recipient_email" form:"recipient_email" binding:"required,email" example:"recipient@example.com"`
	Subject        string `json:"subject" form:"subject" binding:"required" example:"Test Email Subject"`
	Body           string `json:"body" form:"body" binding:"required" example:"This is the email body content"`
	BodyType       string `json:"body_type" form:"body_type" example:"plain" enums:"plain,html"`
	Debug          bool   `json:"debug" form:"debug" example:"false"`
}

// EmailResult represents the result of an email sending operation
type EmailResult struct {
	EmailID       string            `json:"email_id"`
	Status        string            `json:"status"`
	Detail        string            `json:"detail"`
	Timestamp     string            `json:"timestamp"`
	ClientIP      string            `json:"client_ip"`
	Headers       map[string]string `json:"headers"`
	MessageLength int               `json:"message_length"`
}

// EmailResponse represents an email sending response
type EmailResponse struct {
	Message string `json:"message" example:"Email is being sent in the background"`
	EmailID string `json:"email_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid email address"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:""`
	Version   string `json:"version" example:"1.0.0"`
}

// APIResponse represents a standard API response (deprecated, use specific responses)
type APIResponse struct {
	Message string `json:"message"`
	EmailID string `json:"email_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

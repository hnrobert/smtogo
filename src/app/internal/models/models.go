package models

// EmailRequest represents an email sending request
type EmailRequest struct {
	RecipientEmail string `json:"recipient_email" form:"recipient_email" binding:"required,email"`
	Subject        string `json:"subject" form:"subject" binding:"required"`
	Body           string `json:"body" form:"body" binding:"required"`
	BodyType       string `json:"body_type" form:"body_type"`
	Debug          bool   `json:"debug" form:"debug"`
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

// APIResponse represents a standard API response
type APIResponse struct {
	Message string `json:"message"`
	EmailID string `json:"email_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

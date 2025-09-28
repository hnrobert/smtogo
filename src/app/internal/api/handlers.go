package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hnrobert/smtogo/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// sendEmail handles the JSON email sending endpoint
func (s *Server) sendEmail(c *gin.Context) {
	var req models.EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email request
	if err := s.validateEmailRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate email ID
	emailID := uuid.New().String()

	// Get client info
	clientIP := getClientIP(c)
	headers := getHeaders(c)

	// Send email asynchronously
	go func() {
		err := s.emailSender.SendEmail(&req, emailID, clientIP, headers, nil)
		if err != nil {
			fmt.Printf("Failed to send email %s: %v\n", emailID, err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Email is being sent in the background",
		"email_id": emailID,
	})
}

// validateEmailRequest validates the email request
func (s *Server) validateEmailRequest(req *models.EmailRequest) error {
	// Validate recipient email length
	if len(req.RecipientEmail) > s.config.MaxLenRecipientEmail {
		return fmt.Errorf("email address must be less than %d characters", s.config.MaxLenRecipientEmail)
	}

	// Validate subject length
	if len(req.Subject) > s.config.MaxLenSubject {
		return fmt.Errorf("subject must be less than %d characters", s.config.MaxLenSubject)
	}

	// Validate body length
	if len(req.Body) > s.config.MaxLenBody {
		return fmt.Errorf("body content must be less than %d characters", s.config.MaxLenBody)
	}

	// Validate body type
	if req.BodyType != "plain" && req.BodyType != "html" {
		return fmt.Errorf("body type must be either 'plain' or 'html'")
	}

	// Basic email validation (simple check)
	if !strings.Contains(req.RecipientEmail, "@") {
		return fmt.Errorf("invalid email address format")
	}

	return nil
}

// getClientIP extracts the client IP address
func getClientIP(c *gin.Context) string {
	// Check X-Real-IP header first
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	// Check X-Forwarded-For header
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		// Take the first IP from the list
		if parts := strings.Split(ip, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Fall back to remote address
	return c.ClientIP()
}

// getHeaders extracts request headers
func getHeaders(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			// Remove sensitive headers
			if strings.ToLower(key) != "x-api-key" {
				headers[key] = values[0]
			}
		}
	}
	return headers
}

// getFormValue safely extracts a form value
func getFormValue(values map[string][]string, key string) string {
	if vals, exists := values[key]; exists && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

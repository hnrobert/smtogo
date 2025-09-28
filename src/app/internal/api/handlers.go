package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"smtogo/internal/models"
	"smtogo/internal/storage"

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

// sendEmailWithAttachments handles the multipart form email sending endpoint
func (s *Server) sendEmailWithAttachments(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Extract form values
	req := models.EmailRequest{
		RecipientEmail: getFormValue(form.Value, "recipient_email"),
		Subject:        getFormValue(form.Value, "subject"),
		Body:           getFormValue(form.Value, "body"),
		BodyType:       getFormValue(form.Value, "body_type"),
	}

	// Set default body type
	if req.BodyType == "" {
		req.BodyType = "plain"
	}

	// Parse debug flag
	if debugStr := getFormValue(form.Value, "debug"); debugStr != "" {
		req.Debug, _ = strconv.ParseBool(debugStr)
	}

	// Validate email request
	if err := s.validateEmailRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle attachments
	var attachmentNames []string
	if files := form.File["attachments"]; len(files) > 0 {
		// Validate attachment count
		if len(files) > 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can only upload up to 2 attachments"})
			return
		}

		// Process each attachment
		for _, file := range files {
			// Validate file size (2MB limit)
			if file.Size > 2*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Attachments must be smaller than 2MB"})
				return
			}

			// Open file for reading
			fileReader, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read attachment"})
				return
			}
			defer fileReader.Close()

			// Generate unique object name
			objectName := fmt.Sprintf("%s_%s", uuid.New().String(), file.Filename)

			// Upload to MinIO
			contentType := storage.GetContentType(file.Filename)
			err = s.storage.UploadFile(c.Request.Context(), objectName, fileReader, file.Size, contentType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload attachment"})
				return
			}

			attachmentNames = append(attachmentNames, objectName)
		}
	}

	// Generate email ID
	emailID := uuid.New().String()

	// Get client info
	clientIP := getClientIP(c)
	headers := getHeaders(c)

	// Send email asynchronously
	go func() {
		err := s.emailSender.SendEmail(&req, emailID, clientIP, headers, attachmentNames)
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

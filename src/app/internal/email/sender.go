package email

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"smtogo/internal/config"
	"smtogo/internal/models"
	"smtogo/internal/storage"

	"gopkg.in/gomail.v2"
)

// Sender handles email sending operations
type Sender struct {
	config  *config.Config
	storage *storage.MinIOClient
}

// NewSender creates a new email sender
func NewSender(cfg *config.Config) *Sender {
	return &Sender{
		config:  cfg,
		storage: storage.NewMinIOClient(cfg),
	}
}

// SendEmail sends an email with optional attachments
func (s *Sender) SendEmail(req *models.EmailRequest, emailID, clientIP string, headers map[string]string, attachmentNames []string) error {
	// Create email message
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", s.config.GetDisplayEmail())
	m.SetHeader("To", req.RecipientEmail)
	m.SetHeader("Subject", req.Subject)
	m.SetHeader("Message-ID", fmt.Sprintf("<%s@%s>", emailID, s.config.SenderDomain))

	// Set body
	if req.BodyType == "html" {
		m.SetBody("text/html", req.Body)
	} else {
		m.SetBody("text/plain", req.Body)
	}

	// Add attachments if any
	if len(attachmentNames) > 0 {
		for _, objectName := range attachmentNames {
			if objectName != "" {
				if err := s.addAttachment(m, objectName); err != nil {
					s.saveEmailResult(emailID, "failure", fmt.Sprintf("Failed to add attachment: %v", err), clientIP, headers, 0)
					return err
				}
			}
		}
	}

	// Calculate message length (approximate)
	messageLength := len(req.Subject) + len(req.Body) + len(req.RecipientEmail)

	// Send email
	if err := s.sendMessage(m); err != nil {
		s.saveEmailResult(emailID, "failure", fmt.Sprintf("Failed to send email: %v", err), clientIP, headers, messageLength)
		return err
	}

	// Save success result
	s.saveEmailResult(emailID, "success", "Email sent successfully", clientIP, headers, messageLength)

	// Save debug email if requested
	if req.Debug {
		s.saveDebugEmail(emailID, m, req)
	}

	return nil
}

// sendMessage sends the email message via SMTP
func (s *Sender) sendMessage(m *gomail.Message) error {
	// Create SMTP dialer
	d := gomail.NewDialer(s.config.SMTPServer, s.config.SMTPPort, s.config.SenderEmail, s.config.SenderPassword)

	// Configure TLS/SSL
	if s.config.UseSSL {
		d.SSL = true
	}
	if s.config.UseTLS {
		d.TLSConfig = nil // Use default TLS config
	}

	// Disable authentication if not required
	if !s.config.UsePassword {
		d.Username = ""
		d.Password = ""
	}

	// Send the message
	return d.DialAndSend(m)
}

// addAttachment adds an attachment from MinIO to the email
func (s *Sender) addAttachment(m *gomail.Message, objectName string) error {
	// Download file from MinIO
	reader, err := s.storage.GetFile(objectName)
	if err != nil {
		return fmt.Errorf("failed to get file from storage: %w", err)
	}
	defer reader.Close()

	// Read file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Extract filename from object name (remove UUID prefix)
	filename := objectName
	if len(objectName) > 36 && objectName[36] == '_' {
		filename = objectName[37:] // Remove UUID prefix
	}

	// Attach file to message
	m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(content)
		return err
	}))

	return nil
}

// saveEmailResult saves the email sending result to a JSON file
func (s *Sender) saveEmailResult(emailID, status, detail, clientIP string, headers map[string]string, messageLength int) {
	result := models.EmailResult{
		EmailID:       emailID,
		Status:        status,
		Detail:        detail,
		Timestamp:     time.Now().Format(time.RFC3339),
		ClientIP:      clientIP,
		Headers:       headers,
		MessageLength: messageLength,
	}

	// Create directory structure
	dateStr := time.Now().Format("2006-01-02")
	statusDir := "success"
	if status != "success" {
		statusDir = "failure"
	}
	dirPath := filepath.Join("data", dateStr, statusDir)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", dirPath, err)
		return
	}

	// Save result to file
	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.json", emailID))
	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal email result: %v\n", err)
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		fmt.Printf("Failed to save email result to %s: %v\n", filePath, err)
	}
}

// saveDebugEmail saves the raw email message for debugging
func (s *Sender) saveDebugEmail(emailID string, m *gomail.Message, req *models.EmailRequest) {
	// Create debug directory
	dateStr := time.Now().Format("2006-01-02")
	dirPath := filepath.Join("data", dateStr, "debug")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("Failed to create debug directory %s: %v\n", dirPath, err)
		return
	}

	// Save message to file
	filePath := filepath.Join(dirPath, fmt.Sprintf("%s_email.txt", emailID))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create debug file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	// Write message to file
	if _, err := m.WriteTo(file); err != nil {
		fmt.Printf("Failed to write debug email to %s: %v\n", filePath, err)
	}
}

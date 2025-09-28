package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// API Configuration
	APIKey         string `json:"api_key"`
	APIName        string `json:"api_name"`
	APIDescription string `json:"api_description"`
	Port           int    `json:"port"`

	// SMTP Server Settings
	SMTPServer  string `json:"smtp_server"`
	SMTPPort    int    `json:"smtp_port"`
	UseSSL      bool   `json:"use_ssl"`
	UsePassword bool   `json:"use_password"`
	UseTLS      bool   `json:"use_tls"`

	// Email Limits
	MaxLenRecipientEmail int `json:"max_len_recipient_email"`
	MaxLenSubject        int `json:"max_len_subject"`
	MaxLenBody           int `json:"max_len_body"`

	// Sender Configuration
	SenderEmail        string `json:"sender_email"`
	SenderEmailDisplay string `json:"sender_email_display"`
	SenderDomain       string `json:"sender_domain"`
	SenderPassword     string `json:"sender_password"`
}

// Load reads configuration from file
func Load() (*Config, error) {
	configPaths := []string{"smtp_config.jsonc", "smtp_config.json"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path)
		}
	}

	return nil, fmt.Errorf("config file not found in paths: %v", configPaths)
}

// loadFromFile loads configuration from a specific file
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Remove comments for JSONC files
	if strings.HasSuffix(path, ".jsonc") {
		data = removeJSONComments(data)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Set defaults
	config.setDefaults()

	return &config, nil
}

// setDefaults sets default values for configuration
func (c *Config) setDefaults() {
	if c.APIName == "" {
		c.APIName = "High-Performance SMTP API"
	}
	if c.APIDescription == "" {
		c.APIDescription = "SMTP API mail dispatch with support for attachments."
	}
	if c.MaxLenRecipientEmail == 0 {
		c.MaxLenRecipientEmail = 64
	}
	if c.MaxLenSubject == 0 {
		c.MaxLenSubject = 255
	}
	if c.MaxLenBody == 0 {
		c.MaxLenBody = 50000
	}
}

// IsAPIKeyAuthEnabled returns true if API key authentication is enabled
func (c *Config) IsAPIKeyAuthEnabled() bool {
	return strings.TrimSpace(c.APIKey) != ""
}

// GetDisplayEmail returns the display email or falls back to sender email
func (c *Config) GetDisplayEmail() string {
	if strings.TrimSpace(c.SenderEmailDisplay) != "" {
		return c.SenderEmailDisplay
	}
	return c.SenderEmail
}

// removeJSONComments removes single-line comments from JSONC
func removeJSONComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	var cleanLines []string

	for _, line := range lines {
		// Find comment position
		commentPos := strings.Index(line, "//")
		if commentPos != -1 {
			// Keep only the part before the comment
			line = strings.TrimSpace(line[:commentPos])
			// Remove trailing comma if it exists after removing comment
			if strings.HasSuffix(line, ",") && (commentPos > 0) {
				// This is a simple approach; a more robust parser would handle edge cases
			}
		}
		cleanLines = append(cleanLines, line)
	}

	return []byte(strings.Join(cleanLines, "\n"))
}

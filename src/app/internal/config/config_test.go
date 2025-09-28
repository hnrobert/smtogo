package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `{
		"api_name": "Test API",
		"smtp_server": "smtp.test.com",
		"smtp_port": 587,
		"sender_email": "test@example.com",
		"minio_endpoint": "localhost:9000",
		"minio_bucket": "test-bucket"
	}`

	tmpFile, err := os.CreateTemp("", "test_config*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(configContent)
	assert.NoError(t, err)
	tmpFile.Close()

	// Test loading the config
	config, err := loadFromFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "Test API", config.APIName)
	assert.Equal(t, "smtp.test.com", config.SMTPServer)
	assert.Equal(t, 587, config.SMTPPort)
	assert.Equal(t, "test@example.com", config.SenderEmail)
	assert.Equal(t, "localhost:9000", config.MinIOEndpoint)
	assert.Equal(t, "test-bucket", config.MinioBucket)
}

func TestRemoveJSONComments(t *testing.T) {
	input := `{
		"field1": "value1", // This is a comment
		"field2": "value2"
	}`

	_ = `{
		"field1": "value1", 
		"field2": "value2"
	}`

	result := removeJSONComments([]byte(input))
	assert.Contains(t, string(result), "field1")
	assert.Contains(t, string(result), "value1")
	assert.NotContains(t, string(result), "// This is a comment")
}

func TestIsAPIKeyAuthEnabled(t *testing.T) {
	config := &Config{}

	// Test with empty API key
	assert.False(t, config.IsAPIKeyAuthEnabled())

	// Test with whitespace API key
	config.APIKey = "   "
	assert.False(t, config.IsAPIKeyAuthEnabled())

	// Test with valid API key
	config.APIKey = "test-api-key"
	assert.True(t, config.IsAPIKeyAuthEnabled())
}

func TestGetDisplayEmail(t *testing.T) {
	config := &Config{
		SenderEmail: "sender@example.com",
	}

	// Test without display email
	assert.Equal(t, "sender@example.com", config.GetDisplayEmail())

	// Test with display email
	config.SenderEmailDisplay = "Display Name <sender@example.com>"
	assert.Equal(t, "Display Name <sender@example.com>", config.GetDisplayEmail())

	// Test with empty display email
	config.SenderEmailDisplay = "   "
	assert.Equal(t, "sender@example.com", config.GetDisplayEmail())
}

func TestSetDefaults(t *testing.T) {
	config := &Config{}
	config.setDefaults()

	assert.Equal(t, "High-Performance SMTP API", config.APIName)
	assert.Equal(t, "SMTP API mail dispatch with support for attachments.", config.APIDescription)
	assert.Equal(t, 64, config.MaxLenRecipientEmail)
	assert.Equal(t, 255, config.MaxLenSubject)
	assert.Equal(t, 50000, config.MaxLenBody)
	assert.Equal(t, "email-attachments", config.MinioBucket)
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"smtogo/internal/api"
	"smtogo/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		APIName: "Test SMTP API",
		Port:    8000,
	}

	// Create test server
	server := api.NewServer(cfg)
	router := server.GetRouter()

	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body contains expected content
	assert.Contains(t, rr.Body.String(), "status")
}

func TestOpenAPIEndpoint(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		APIName:        "Test SMTP API",
		APIDescription: "Test description",
		Port:           8000,
	}

	// Create test server
	server := api.NewServer(cfg)
	router := server.GetRouter()

	// Create test request
	req, err := http.NewRequest("GET", "/openapi.json", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check content type
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	// Check response body contains OpenAPI spec
	assert.Contains(t, rr.Body.String(), "openapi")
	assert.Contains(t, rr.Body.String(), "Test SMTP API")
}

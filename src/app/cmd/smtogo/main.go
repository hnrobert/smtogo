package main

import (
	"log"
	_ "smtogo/docs" // This line is necessary for Swagger
	"smtogo/internal/api"
	"smtogo/internal/config"
)

// @title SMToGo API
// @version 1.0
// @description High-performance SMTP API server for reliable email sending
// @contact.name SMToGo Support
// @contact.email hnrobert@qq.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8000
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Start the API server
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"

	"github.com/hnrobert/smtogo/internal/api"
	"github.com/hnrobert/smtogo/internal/config"
)

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

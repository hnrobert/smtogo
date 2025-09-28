package api

import (
	"fmt"
	"net/http"
	"smtogo/internal/config"
	"smtogo/internal/email"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the API server
type Server struct {
	config      *config.Config
	emailSender *email.Sender
	router      *gin.Engine
}

// NewServer creates a new API server instance
func NewServer(cfg *config.Config) *Server {
	// Initialize email sender
	emailSender := email.NewSender(cfg)

	server := &Server{
		config:      cfg,
		emailSender: emailSender,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	s.router = gin.Default()

	// Health check endpoint
	s.router.GET("/health", s.getHealth)

	// Swagger documentation endpoints - mounted at root
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	// Email endpoints
	v1 := s.router.Group("/v1")
	{
		mail := v1.Group("/mail")
		{
			if s.config.IsAPIKeyAuthEnabled() {
				mail.Use(s.apiKeyAuthMiddleware())
			}
			mail.POST("/send", s.sendEmail)
		}
	}
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	port := ":8000"
	fmt.Printf("Starting server on port %s\n", port)
	return s.router.Run(port)
}

// apiKeyAuthMiddleware validates API key if authentication is enabled
func (s *Server) apiKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != s.config.APIKey {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Could not validate credentials",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getHealth handles health check requests
// @Summary Health Check
// @Description Check API health status
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse "API is healthy"
// @Router /health [get]
func (s *Server) getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": c.Request.Header.Get("X-Request-Time"),
		"version":   "1.0.0",
	})
}

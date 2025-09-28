package api

import (
	"fmt"
	"net/http"
	"smtogo/internal/config"
	"smtogo/internal/email"

	"github.com/gin-gonic/gin"
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

	// API documentation endpoints
	s.router.GET("/", s.getDocumentation)
	// s.router.GET("/docs", s.getDocumentation)
	s.router.GET("/openapi.json", s.getOpenAPISpec)

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
func (s *Server) getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": c.Request.Header.Get("X-Request-Time"),
		"version":   "1.0.0",
	})
}

// getDocumentation serves the Swagger UI documentation
func (s *Server) getDocumentation(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>` + s.config.APIName + ` - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/openapi.json',
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.standalone
            ]
        });
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html", []byte(html))
}

// getOpenAPISpec returns the OpenAPI specification
func (s *Server) getOpenAPISpec(c *gin.Context) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       s.config.APIName,
			"description": s.config.APIDescription,
			"version":     "1.0.0",
		},
		"paths": map[string]interface{}{
			"/v1/mail/send": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Send email",
					"description": "Send an email without attachments",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/EmailRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Email queued successfully",
						},
					},
				},
			},
			"/v1/mail/send-with-attachments": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Send email with attachments",
					"description": "Send an email with file attachments",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"multipart/form-data": map[string]interface{}{
								"schema": map[string]interface{}{
									"$ref": "#/components/schemas/EmailWithAttachmentsRequest",
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Email queued successfully",
						},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"EmailRequest": map[string]interface{}{
					"type":     "object",
					"required": []string{"recipient_email", "subject", "body"},
					"properties": map[string]interface{}{
						"recipient_email": map[string]interface{}{
							"type":        "string",
							"format":      "email",
							"description": "Recipient email address",
						},
						"subject": map[string]interface{}{
							"type":        "string",
							"description": "Email subject",
						},
						"body": map[string]interface{}{
							"type":        "string",
							"description": "Email body content",
						},
						"body_type": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"plain", "html"},
							"default":     "plain",
							"description": "Email body type",
						},
						"debug": map[string]interface{}{
							"type":        "boolean",
							"default":     false,
							"description": "Enable debug mode",
						},
					},
				},
				"EmailWithAttachmentsRequest": map[string]interface{}{
					"type":     "object",
					"required": []string{"recipient_email", "subject", "body"},
					"properties": map[string]interface{}{
						"recipient_email": map[string]interface{}{
							"type":        "string",
							"format":      "email",
							"description": "Recipient email address",
						},
						"subject": map[string]interface{}{
							"type":        "string",
							"description": "Email subject",
						},
						"body": map[string]interface{}{
							"type":        "string",
							"description": "Email body content",
						},
						"body_type": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"plain", "html"},
							"default":     "plain",
							"description": "Email body type",
						},
						"debug": map[string]interface{}{
							"type":        "boolean",
							"default":     false,
							"description": "Enable debug mode",
						},
						"attachments": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type":   "string",
								"format": "binary",
							},
							"description": "File attachments (max 2 files, 2MB each)",
						},
					},
				},
			},
			"securitySchemes": func() map[string]interface{} {
				if s.config.IsAPIKeyAuthEnabled() {
					return map[string]interface{}{
						"ApiKeyAuth": map[string]interface{}{
							"type": "apiKey",
							"in":   "header",
							"name": "X-API-Key",
						},
					}
				}
				return map[string]interface{}{}
			}(),
		},
		"security": func() []interface{} {
			if s.config.IsAPIKeyAuthEnabled() {
				return []interface{}{
					map[string]interface{}{
						"ApiKeyAuth": []interface{}{},
					},
				}
			}
			return []interface{}{}
		}(),
	}

	c.JSON(http.StatusOK, spec)
}

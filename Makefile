# SMToGo Makefile

.PHONY: all build test clean run docker-build docker-run dev format lint

# Variables
BINARY_NAME=smtogo
MAIN_PATH=./src/app/cmd/smtogo
DOCKER_TAG=smtogo:latest

# Default target
all: clean format lint test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="-w -s" -o $(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# Format code
format:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"; \
		go vet ./...; \
	fi

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8000:8000 -v $(PWD)/smtp_config.jsonc:/app/smtp_config.jsonc:ro $(DOCKER_TAG)

# Docker Compose operations
compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

compose-down:
	@echo "Stopping services..."
	@docker-compose down

compose-logs:
	@echo "Showing logs..."
	@docker-compose logs -f

# Generate OpenAPI documentation
docs:
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/smtogo/main.go -o docs/; \
	else \
		echo "swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Database migration (placeholder for future use)
migrate-up:
	@echo "Running database migrations..."
	@echo "No migrations to run yet"

migrate-down:
	@echo "Rolling back database migrations..."
	@echo "No migrations to rollback yet"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Development tools installed successfully"

# Help
help:
	@echo "Available targets:"
	@echo "  all           - Run clean, format, lint, test, and build"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  run           - Run the application"
	@echo "  dev           - Start development mode with hot reload"
	@echo "  format        - Format code"
	@echo "  lint          - Lint code"
	@echo "  security      - Run security scan"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  deps-update   - Update dependencies"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  compose-up    - Start services with Docker Compose"
	@echo "  compose-down  - Stop services"
	@echo "  compose-logs  - Show Docker Compose logs"
	@echo "  docs          - Generate API documentation"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message"
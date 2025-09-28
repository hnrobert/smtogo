# SMToGo Makefile

.PHONY: all build test clean run docker-build docker-run dev format lint

# Variables
BINARY_NAME=smtogo
MAIN_PATH=./src/app/cmd/smtogo
DOCKER_TAG=smtogo:latest

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

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html

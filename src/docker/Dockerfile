# Build stage
FROM golang:1.21-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o smtogo ./cmd/smtogo

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh smtogo

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/smtogo .

# Create data directory
RUN mkdir -p data && chown -R smtogo:smtogo /app

# Switch to non-root user
USER smtogo

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# Run the application
ENTRYPOINT ["./smtogo"]

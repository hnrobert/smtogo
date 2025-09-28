# SMToGo - High-Performance SMTP API Server

[![Build Status](https://github.com/YOUR_USERNAME/smtogo/workflows/Build%20and%20Test%20SMToGo/badge.svg)](https://github.com/YOUR_USERNAME/smtogo/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_USERNAME/smtogo)](https://goreportcard.com/report/github.com/YOUR_USERNAME/smtogo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance SMTP API server written in Go, designed for reliable email sending with attachment support via MinIO object storage.

## Features

- 🚀 **High Performance**: Built with Go and Gin framework for excellent performance
- 📧 **SMTP Support**: Full SMTP configuration with SSL/TLS support
- 📎 **File Attachments**: Seamless file attachment handling via MinIO object storage
- 🔐 **Optional Authentication**: API key-based authentication (optional)
- 📊 **OpenAPI Documentation**: Built-in Swagger/ReDoc documentation
- 🐳 **Docker Ready**: Complete Docker and Docker Compose setup
- 🔄 **CI/CD**: GitHub Actions workflow for testing and deployment
- 📝 **JSONC Configuration**: Support for JSON with comments configuration files
- 🏥 **Health Checks**: Built-in health check endpoints
- 📈 **Structured Logging**: Comprehensive logging and request tracking

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**:

   ```bash
   git clone https://github.com/YOUR_USERNAME/smtogo.git
   cd smtogo
   ```

2. **Configure SMTP settings**:

   ```bash
   cp smtp_config.jsonc.example smtp_config.jsonc
   # Edit smtp_config.jsonc with your SMTP server details
   ```

3. **Start the services**:

   ```bash
   docker-compose up -d
   ```

4. **Access the services**:
   - API Server: <http://localhost:8000>
   - API Documentation: <http://localhost:8000/docs>
   - MinIO Console: <http://localhost:9001> (admin/admin123)

### Manual Installation

1. **Prerequisites**:

   - Go 1.21 or later
   - MinIO server (for file attachments)

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Configure the application**:

   ```bash
   cp smtp_config.jsonc.example smtp_config.jsonc
   # Edit the configuration file
   ```

4. **Run the application**:

```bash
go run cmd/smtogo/main.go
```

## Configuration

The application uses a JSONC configuration file (`smtp_config.jsonc`) that supports comments:

```jsonc
{
  // API Configuration
  "api_key": "", // Optional: API key for authentication
  "api_name": "High-Performance SMTP API",
  "api_description": "SMTP API mail dispatch with support for attachments.",

  // SMTP Server Settings
  "smtp_server": "smtp.example.com",
  "smtp_port": 587,
  "use_ssl": false,
  "use_password": true,
  "use_tls": true,

  // Email Limits
  "max_len_recipient_email": 64,
  "max_len_subject": 255,
  "max_len_body": 50000,

  // Sender Configuration
  "sender_email": "sender@example.com",
  "sender_email_display": "Display Name <sender@example.com>",
  "sender_domain": "example.com",
  "sender_password": "your_smtp_password",

  // MinIO Object Storage Settings
  "minio_endpoint": "localhost:9000",
  "minio_access_key": "minioadmin",
  "minio_secret_key": "minioadmin",
  "minio_bucket": "email-attachments",
  "minio_use_ssl": false
}
```

### Environment Variables

You can also configure the application using environment variables:

- `SMTP_SERVER`: SMTP server hostname
- `SMTP_PORT`: SMTP server port
- `SENDER_EMAIL`: Sender email address
- `SENDER_PASSWORD`: SMTP password
- `API_KEY`: Optional API key for authentication
- `MINIO_ENDPOINT`: MinIO server endpoint
- `MINIO_ACCESS_KEY`: MinIO access key
- `MINIO_SECRET_KEY`: MinIO secret key

## API Usage

### Send Email (JSON)

```bash
curl -X POST http://localhost:8000/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "recipient_email": "recipient@example.com",
    "subject": "Test Email",
    "body": "This is a test email",
    "body_type": "plain"
  }'
```

### Send Email with Attachments (Multipart Form)

```bash
curl -X POST http://localhost:8000/send-form \
  -H "X-API-Key: your-api-key" \
  -F "recipient_email=recipient@example.com" \
  -F "subject=Test with Attachment" \
  -F "body=Email with attachment" \
  -F "body_type=plain" \
  -F "file1=@/path/to/attachment.pdf" \
  -F "file2=@/path/to/image.jpg"
```

### API Response

```json
{
  "success": true,
  "message": "Email sent successfully",
  "email_id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2024-01-15T10:30:45Z"
}
```

## Architecture

```mermaid
flowchart LR
    A[Nginx Proxy<br/>(Optional)] --> B[SMToGo API]
    B --> C[SMTP Server]
    B --> D[MinIO Store<br/>(Attachments)]
```

### Project Structure

```text
smtogo/
├── cmd/smtogo/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── config/          # Configuration management
│   ├── email/           # Email sending logic
│   ├── models/          # Data structures
│   └── storage/         # MinIO client
├── nginx/               # Nginx configuration
├── .github/workflows/   # CI/CD pipelines
├── docker-compose.yml   # Docker orchestration
├── Dockerfile          # Container definition
└── README.md
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

### Building

```bash
# Build for current platform
go build -o smtogo ./cmd/smtogo

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o smtogo-linux ./cmd/smtogo

# Build Docker image
docker build -t smtogo .
```

## Deployment

### Docker Deployment

1. **Build and push image**:

   ```bash
   docker build -t your-registry/smtogo:latest .
   docker push your-registry/smtogo:latest
   ```

2. **Deploy with docker-compose**:

   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Kubernetes Deployment

See the `k8s/` directory for Kubernetes manifests.

### Production Considerations

- Use strong API keys for authentication
- Configure TLS/SSL for SMTP connections
- Set up proper logging and monitoring
- Configure rate limiting in Nginx
- Use persistent volumes for MinIO data
- Set up backup strategies for email logs
- Monitor disk space for attachment storage

## Monitoring

### Health Checks

- `GET /health`: Basic health check
- `GET /ready`: Readiness check (includes dependencies)

### Metrics

The application exposes metrics endpoints for monitoring:

- Request/response times
- Email sending success/failure rates
- Storage usage statistics
- SMTP connection health

## Security

- Optional API key authentication
- Request rate limiting
- Input validation and sanitization
- Secure file upload handling
- SMTP credential protection
- Container security best practices

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📧 Email: <support@example.com>
- 💬 Issues: [GitHub Issues](https://github.com/YOUR_USERNAME/smtogo/issues)
- 📖 Documentation: [API Docs](http://localhost:8000/docs)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.

---

Made with ❤️ by [Your Name](https://github.com/YOUR_USERNAME)

package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"smtogo/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient handles MinIO operations
type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg *config.Config) *MinIOClient {
	// Create MinIO client
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create MinIO client: %v", err))
	}

	return &MinIOClient{
		client:     client,
		bucketName: cfg.MinioBucket,
	}
}

// EnsureBucket ensures the bucket exists
func (mc *MinIOClient) EnsureBucket(ctx context.Context) error {
	exists, err := mc.client.BucketExists(ctx, mc.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := mc.client.MakeBucket(ctx, mc.bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (mc *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := mc.client.PutObject(ctx, mc.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// GetFile downloads a file from MinIO
func (mc *MinIOClient) GetFile(objectName string) (io.ReadCloser, error) {
	object, err := mc.client.GetObject(context.Background(), mc.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return object, nil
}

// DeleteFile deletes a file from MinIO
func (mc *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := mc.client.RemoveObject(ctx, mc.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListFiles lists files in the bucket with optional prefix
func (mc *MinIOClient) ListFiles(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo

	for object := range mc.client.ListObjects(ctx, mc.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		objects = append(objects, object)
	}

	return objects, nil
}

// GetFileInfo gets information about a file
func (mc *MinIOClient) GetFileInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	info, err := mc.client.StatObject(ctx, mc.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get file info: %w", err)
	}

	return info, nil
}

// GeneratePresignedURL generates a presigned URL for file access
func (mc *MinIOClient) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := mc.client.PresignedGetObject(ctx, mc.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// CleanupOldFiles deletes files older than the specified duration
func (mc *MinIOClient) CleanupOldFiles(ctx context.Context, maxAge time.Duration) error {
	cutoffTime := time.Now().Add(-maxAge)

	objects, err := mc.ListFiles(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list files for cleanup: %w", err)
	}

	for _, object := range objects {
		if object.LastModified.Before(cutoffTime) {
			if err := mc.DeleteFile(ctx, object.Key); err != nil {
				// Log error but continue cleanup
				fmt.Printf("Failed to delete old file %s: %v\n", object.Key, err)
			}
		}
	}

	return nil
}

// ValidateObjectName validates the object name format
func ValidateObjectName(objectName string) error {
	if objectName == "" {
		return fmt.Errorf("object name cannot be empty")
	}

	if strings.Contains(objectName, "..") {
		return fmt.Errorf("object name cannot contain '..'")
	}

	if strings.HasPrefix(objectName, "/") || strings.HasSuffix(objectName, "/") {
		return fmt.Errorf("object name cannot start or end with '/'")
	}

	return nil
}

// GetContentType returns the MIME type based on file extension
func GetContentType(filename string) string {
	extension := strings.ToLower(filename)

	if strings.HasSuffix(extension, ".jpg") || strings.HasSuffix(extension, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(extension, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(extension, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(extension, ".pdf") {
		return "application/pdf"
	}
	if strings.HasSuffix(extension, ".txt") {
		return "text/plain"
	}
	if strings.HasSuffix(extension, ".doc") {
		return "application/msword"
	}
	if strings.HasSuffix(extension, ".docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	if strings.HasSuffix(extension, ".xls") {
		return "application/vnd.ms-excel"
	}
	if strings.HasSuffix(extension, ".xlsx") {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	if strings.HasSuffix(extension, ".zip") {
		return "application/zip"
	}

	return "application/octet-stream"
}

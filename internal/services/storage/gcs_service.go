package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/disintegration/imaging"
	"google.golang.org/api/option"
)

// GCSImageService handles image uploads and processing using Google Cloud Storage
type GCSImageService struct {
	client     *storage.Client
	bucketName string
	cdnBaseURL string
	bucket     *storage.BucketHandle
}

// NewGCSImageService creates a new Google Cloud Storage image service
func NewGCSImageService(bucketName, cdnBaseURL string, opts ...option.ClientOption) (*GCSImageService, error) {
	ctx := context.Background()

	// Create GCS client
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	bucket := client.Bucket(bucketName)

	// Check if bucket exists
	if _, err := bucket.Attrs(ctx); err != nil {
		return nil, fmt.Errorf("bucket %s not accessible: %w", bucketName, err)
	}

	return &GCSImageService{
		client:     client,
		bucketName: bucketName,
		cdnBaseURL: cdnBaseURL,
		bucket:     bucket,
	}, nil
}

// UploadTicketImage uploads and processes a ticket image
func (s *GCSImageService) UploadTicketImage(tokenID uint64, imageData []byte, contentType string) (string, error) {
	ctx := context.Background()
	config := DefaultUploadConfig()

	// Validate image
	if err := s.validateImage(imageData, contentType, config); err != nil {
		return "", err
	}

	// Process image (resize, optimize)
	processedData, err := s.processImage(imageData, contentType, config)
	if err != nil {
		return "", err
	}

	// Generate object paths
	originalPath := fmt.Sprintf("tickets/%d/original.jpg", tokenID)
	thumbnailPath := fmt.Sprintf("tickets/%d/thumbnail.jpg", tokenID)

	// Upload original
	if err := s.uploadToGCS(ctx, originalPath, processedData, "image/jpeg"); err != nil {
		return "", fmt.Errorf("failed to upload original: %w", err)
	}

	// Create and upload thumbnail
	thumbnailData, err := s.createThumbnail(processedData, 400, 400)
	if err != nil {
		return "", fmt.Errorf("failed to create thumbnail: %w", err)
	}

	if err := s.uploadToGCS(ctx, thumbnailPath, thumbnailData, "image/jpeg"); err != nil {
		return "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	// Return CDN URL or direct GCS URL
	if s.cdnBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.cdnBaseURL, originalPath), nil
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, originalPath), nil
}

// uploadToGCS uploads data to Google Cloud Storage
func (s *GCSImageService) uploadToGCS(ctx context.Context, objectPath string, data []byte, contentType string) error {
	obj := s.bucket.Object(objectPath)
	writer := obj.NewWriter(ctx)

	// Set metadata
	writer.ContentType = contentType
	writer.CacheControl = "public, max-age=86400" // Cache for 1 day
	writer.Metadata = map[string]string{
		"uploaded": time.Now().Format(time.RFC3339),
	}

	// Write data
	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write to GCS: %w", err)
	}

	// Close and commit
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	// Make object publicly readable (optional)
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		// Log warning but don't fail - bucket might have uniform access control
		fmt.Printf("Warning: Could not set public ACL: %v\n", err)
	}

	return nil
}

// GenerateSignedUploadURL generates a signed URL for direct frontend upload to GCS
func (s *GCSImageService) GenerateSignedUploadURL(tokenID uint64, contentType string) (string, error) {
	objectPath := fmt.Sprintf("tickets/%d/original.jpg", tokenID)

	// Generate signed URL valid for 15 minutes
	url, err := s.bucket.SignedURL(objectPath, &storage.SignedURLOptions{
		Method:      "PUT",
		ContentType: contentType,
		Expires:     time.Now().Add(15 * time.Minute),
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}

// GetImageURL returns the public URL for a ticket image
func (s *GCSImageService) GetImageURL(tokenID uint64) string {
	objectPath := fmt.Sprintf("tickets/%d/original.jpg", tokenID)

	if s.cdnBaseURL != "" {
		return fmt.Sprintf("%s/%s", s.cdnBaseURL, objectPath)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, objectPath)
}

// DeleteTicketImages deletes all images for a ticket
func (s *GCSImageService) DeleteTicketImages(tokenID uint64) error {
	ctx := context.Background()

	paths := []string{
		fmt.Sprintf("tickets/%d/original.jpg", tokenID),
		fmt.Sprintf("tickets/%d/thumbnail.jpg", tokenID),
	}

	for _, path := range paths {
		if err := s.bucket.Object(path).Delete(ctx); err != nil {
			// Ignore not found errors
			if err != storage.ErrObjectNotExist {
				return fmt.Errorf("failed to delete %s: %w", path, err)
			}
		}
	}

	return nil
}

// validateImage validates image data against configuration
func (s *GCSImageService) validateImage(data []byte, contentType string, config UploadConfig) error {
	// Check size
	if int64(len(data)) > config.MaxSizeBytes {
		return fmt.Errorf("image size %d exceeds maximum %d bytes", len(data), config.MaxSizeBytes)
	}

	// Check content type
	allowed := false
	for _, t := range config.AllowedTypes {
		if strings.EqualFold(contentType, t) {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("content type %s not allowed", contentType)
	}

	// Verify it's actually an image
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid image data: %w", err)
	}

	// Verify format matches content type
	expectedFormat := strings.TrimPrefix(contentType, "image/")
	if !strings.EqualFold(format, expectedFormat) &&
		!(expectedFormat == "jpeg" && format == "jpg") {
		return fmt.Errorf("content type %s doesn't match image format %s", contentType, format)
	}

	return nil
}

// processImage resizes and optimizes an image
func (s *GCSImageService) processImage(data []byte, _ string, config UploadConfig) ([]byte, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize if needed (maintain aspect ratio)
	bounds := img.Bounds()
	if bounds.Dx() > config.ResizeWidth || bounds.Dy() > config.ResizeHeight {
		img = imaging.Fit(img, config.ResizeWidth, config.ResizeHeight, imaging.Lanczos)
	}

	// Encode as JPEG with specified quality
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: config.Quality}); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// createThumbnail creates a thumbnail version of an image
func (s *GCSImageService) createThumbnail(data []byte, width, height int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image for thumbnail: %w", err)
	}

	// Create thumbnail
	thumbnail := imaging.Thumbnail(img, width, height, imaging.Lanczos)

	// Encode as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

// Close closes the GCS client
func (s *GCSImageService) Close() error {
	return s.client.Close()
}

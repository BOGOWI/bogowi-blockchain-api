package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

// ImageService handles image uploads and processing for NFT tickets
type ImageService struct {
	s3Client   *s3.Client
	bucketName string
	cdnBaseURL string
}

// NewImageService creates a new image service
func NewImageService(bucketName, cdnBaseURL string) (*ImageService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &ImageService{
		s3Client:   s3.NewFromConfig(cfg),
		bucketName: bucketName,
		cdnBaseURL: cdnBaseURL,
	}, nil
}

// UploadConfig contains configuration for image upload
type UploadConfig struct {
	MaxSizeBytes int64
	AllowedTypes []string
	ResizeWidth  int
	ResizeHeight int
	Quality      int
}

// DefaultUploadConfig returns default configuration
func DefaultUploadConfig() UploadConfig {
	return UploadConfig{
		MaxSizeBytes: 10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/webp"},
		ResizeWidth:  1200,
		ResizeHeight: 1200,
		Quality:      85,
	}
}

// GeneratePresignedUploadURL generates a presigned URL for direct frontend upload
func (s *ImageService) GeneratePresignedUploadURL(tokenID uint64, contentType string) (string, error) {
	key := fmt.Sprintf("tickets/%d/original.jpg", tokenID)

	presignClient := s3.NewPresignClient(s.s3Client)

	request, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * time.Minute)
	})

	if err != nil {
		return "", fmt.Errorf("failed to create presigned URL: %w", err)
	}

	return request.URL, nil
}

// ProcessAndUploadImage processes and uploads an image
func (s *ImageService) ProcessAndUploadImage(reader io.Reader, tokenID uint64, config UploadConfig) (string, error) {
	// Read image
	img, format, err := image.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image if needed
	processed := imaging.Fit(img, config.ResizeWidth, config.ResizeHeight, imaging.Lanczos)

	// Generate multiple sizes for responsive display
	sizes := map[string]image.Image{
		"original": processed,
		"large":    imaging.Resize(processed, 800, 0, imaging.Lanczos),
		"medium":   imaging.Resize(processed, 400, 0, imaging.Lanczos),
		"thumb":    imaging.Fill(processed, 200, 200, imaging.Center, imaging.Lanczos),
	}

	baseKey := fmt.Sprintf("tickets/%d", tokenID)

	for sizeName, sizedImg := range sizes {
		var buf bytes.Buffer

		// Encode based on format
		switch format {
		case "png":
			err = png.Encode(&buf, sizedImg)
		default:
			err = jpeg.Encode(&buf, sizedImg, &jpeg.Options{Quality: config.Quality})
		}

		if err != nil {
			return "", fmt.Errorf("failed to encode image: %w", err)
		}

		// Upload to S3
		key := fmt.Sprintf("%s/%s.jpg", baseKey, sizeName)
		_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:       aws.String(s.bucketName),
			Key:          aws.String(key),
			Body:         bytes.NewReader(buf.Bytes()),
			ContentType:  aws.String("image/jpeg"),
			CacheControl: aws.String("public, max-age=31536000"),
			Metadata: map[string]string{
				"token-id": fmt.Sprintf("%d", tokenID),
			},
		})

		if err != nil {
			return "", fmt.Errorf("failed to upload to S3: %w", err)
		}
	}

	// Return CDN URL for the original size
	return fmt.Sprintf("%s/tickets/%d/original.jpg", s.cdnBaseURL, tokenID), nil
}

// GenerateDefaultImage generates a default image for tickets without custom images
func (s *ImageService) GenerateDefaultImage(tokenID uint64, experienceType string) (string, error) {
	// Map experience types to default images
	defaultImages := map[string]string{
		"Wildlife Safari":     "defaults/wildlife-safari.jpg",
		"Marine Conservation": "defaults/marine-conservation.jpg",
		"Forest Trek":         "defaults/forest-trek.jpg",
		"Cultural Experience": "defaults/cultural-experience.jpg",
		"Eco Lodge":           "defaults/eco-lodge.jpg",
		"Adventure":           "defaults/adventure.jpg",
	}

	// Get default image path
	defaultPath, exists := defaultImages[experienceType]
	if !exists {
		defaultPath = "defaults/generic-experience.jpg"
	}

	// Copy default to token-specific path
	sourceKey := defaultPath
	destKey := fmt.Sprintf("tickets/%d/original.jpg", tokenID)

	_, err := s.s3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucketName, sourceKey)),
		Key:        aws.String(destKey),
		Metadata: map[string]string{
			"token-id": fmt.Sprintf("%d", tokenID),
			"type":     "default",
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to copy default image: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.cdnBaseURL, destKey), nil
}

// ValidateImage validates an image meets requirements
func (s *ImageService) ValidateImage(reader io.Reader, config UploadConfig) error {
	// Check size
	limitedReader := &io.LimitedReader{R: reader, N: config.MaxSizeBytes + 1}
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return fmt.Errorf("failed to read image: %w", err)
	}

	if int64(len(data)) > config.MaxSizeBytes {
		return fmt.Errorf("image size exceeds maximum of %d bytes", config.MaxSizeBytes)
	}

	// Check content type
	contentType := http.DetectContentType(data)
	validType := false
	for _, allowed := range config.AllowedTypes {
		if strings.HasPrefix(contentType, allowed) {
			validType = true
			break
		}
	}

	if !validType {
		return fmt.Errorf("invalid image type: %s", contentType)
	}

	// Try to decode to validate it's a real image
	_, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("invalid image format: %w", err)
	}

	return nil
}

// GetImageURLs returns all image URLs for a token
func (s *ImageService) GetImageURLs(tokenID uint64) map[string]string {
	baseURL := fmt.Sprintf("%s/tickets/%d", s.cdnBaseURL, tokenID)
	return map[string]string{
		"original": fmt.Sprintf("%s/original.jpg", baseURL),
		"large":    fmt.Sprintf("%s/large.jpg", baseURL),
		"medium":   fmt.Sprintf("%s/medium.jpg", baseURL),
		"thumb":    fmt.Sprintf("%s/thumb.jpg", baseURL),
	}
}

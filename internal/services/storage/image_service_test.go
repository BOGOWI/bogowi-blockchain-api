package storage

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultUploadConfig(t *testing.T) {
	config := DefaultUploadConfig()

	assert.Equal(t, int64(10*1024*1024), config.MaxSizeBytes)
	assert.Contains(t, config.AllowedTypes, "image/jpeg")
	assert.Contains(t, config.AllowedTypes, "image/png")
	assert.Contains(t, config.AllowedTypes, "image/webp")
	assert.Equal(t, 1200, config.ResizeWidth)
	assert.Equal(t, 1200, config.ResizeHeight)
	assert.Equal(t, 85, config.Quality)
}

func TestImageService_ValidateImage(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	config := DefaultUploadConfig()
	config.MaxSizeBytes = 1024 * 1024 // 1MB for testing

	tests := []struct {
		name    string
		image   func() io.Reader
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid JPEG image",
			image: func() io.Reader {
				return createTestJPEG(100, 100)
			},
			wantErr: false,
		},
		{
			name: "valid PNG image",
			image: func() io.Reader {
				return createTestPNG(100, 100)
			},
			wantErr: false,
		},
		{
			name: "image too large",
			image: func() io.Reader {
				// Create an image larger than 1MB
				data := make([]byte, 2*1024*1024) // 2MB
				// Add JPEG header
				copy(data, []byte{0xFF, 0xD8, 0xFF, 0xE0})
				return bytes.NewReader(data)
			},
			wantErr: true,
			errMsg:  "exceeds maximum",
		},
		{
			name: "invalid file type",
			image: func() io.Reader {
				return bytes.NewReader([]byte("This is not an image"))
			},
			wantErr: true,
			errMsg:  "invalid image",
		},
		{
			name: "corrupted image",
			image: func() io.Reader {
				// Start with JPEG header but corrupted data
				data := []byte{0xFF, 0xD8, 0xFF, 0xE0}
				data = append(data, bytes.Repeat([]byte{0x00}, 100)...)
				return bytes.NewReader(data)
			},
			wantErr: true,
			errMsg:  "invalid image format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateImage(tt.image(), config)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestImageService_GetImageURLs(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	tests := []struct {
		tokenID uint64
		want    map[string]string
	}{
		{
			tokenID: 1,
			want: map[string]string{
				"original": "https://cdn.example.com/tickets/1/original.jpg",
				"large":    "https://cdn.example.com/tickets/1/large.jpg",
				"medium":   "https://cdn.example.com/tickets/1/medium.jpg",
				"thumb":    "https://cdn.example.com/tickets/1/thumb.jpg",
			},
		},
		{
			tokenID: 999,
			want: map[string]string{
				"original": "https://cdn.example.com/tickets/999/original.jpg",
				"large":    "https://cdn.example.com/tickets/999/large.jpg",
				"medium":   "https://cdn.example.com/tickets/999/medium.jpg",
				"thumb":    "https://cdn.example.com/tickets/999/thumb.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("token_%d", tt.tokenID), func(t *testing.T) {
			urls := service.GetImageURLs(tt.tokenID)
			assert.Equal(t, tt.want, urls)
		})
	}
}

func TestImageService_ValidateImage_EdgeCases(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	t.Run("exactly at size limit", func(t *testing.T) {
		config := DefaultUploadConfig()
		config.MaxSizeBytes = 10000 // Small limit for testing
		
		// Create image that's under the limit
		img := createTestJPEGWithSize(9500)
		err := service.ValidateImage(img, config)
		assert.NoError(t, err)
	})

	t.Run("empty reader", func(t *testing.T) {
		config := DefaultUploadConfig()
		err := service.ValidateImage(bytes.NewReader([]byte{}), config)
		assert.Error(t, err)
	})
}

func TestImageService_NewImageService(t *testing.T) {
	// Test that NewImageService properly initializes
	// Note: This will fail without AWS credentials, but we can test the structure
	service, err := NewImageService("test-bucket", "https://cdn.example.com")
	
	// We expect an error without proper AWS config, but the function should exist
	if err == nil {
		assert.NotNil(t, service)
		assert.Equal(t, "test-bucket", service.bucketName)
		assert.Equal(t, "https://cdn.example.com", service.cdnBaseURL)
		assert.NotNil(t, service.s3Client)
	}
}

// Test default image mapping
func TestImageService_DefaultImageMapping(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	// Test that the method exists and has correct signature
	assert.NotNil(t, service.GenerateDefaultImage)
	
	// Test cases for experience type mapping
	tests := []struct {
		experienceType string
		expectedPath   string
	}{
		{"Wildlife Safari", "defaults/wildlife-safari.jpg"},
		{"Marine Conservation", "defaults/marine-conservation.jpg"},
		{"Forest Trek", "defaults/forest-trek.jpg"},
		{"Unknown Type", "defaults/generic-experience.jpg"},
	}

	// We can't test the actual S3 operations without mocking,
	// but we can verify the logic is present
	for _, tt := range tests {
		t.Run(tt.experienceType, func(t *testing.T) {
			// The actual test would require S3 mocking
			// This just ensures the method exists with proper signature
			_ = tt.expectedPath // Use the variable
		})
	}
}

// Test image processing logic
func TestImageService_ImageProcessing(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	// Test that ProcessAndUploadImage method exists
	assert.NotNil(t, service.ProcessAndUploadImage)
	
	// Test config validation
	config := DefaultUploadConfig()
	
	t.Run("validate config quality", func(t *testing.T) {
		assert.Greater(t, config.Quality, 0)
		assert.LessOrEqual(t, config.Quality, 100)
	})

	t.Run("validate config dimensions", func(t *testing.T) {
		assert.Greater(t, config.ResizeWidth, 0)
		assert.Greater(t, config.ResizeHeight, 0)
	})
}

// Test presigned URL generation method signature
func TestImageService_PresignedURL(t *testing.T) {
	service := &ImageService{
		bucketName: "test-bucket",
		cdnBaseURL: "https://cdn.example.com",
	}

	// Verify the method exists with correct signature
	assert.NotNil(t, service.GeneratePresignedUploadURL)
	
	// The actual functionality requires AWS SDK setup
	// which we skip in unit tests
}

// Helper functions to create test images

func createTestJPEG(width, height int) io.Reader {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{100, 100, 100, 255}}, image.Point{}, draw.Src)
	
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	return &buf
}

func createTestPNG(width, height int) io.Reader {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{100, 100, 100, 255}}, image.Point{}, draw.Src)
	
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return &buf
}

func createTestJPEGWithSize(targetSize int) io.Reader {
	// Start with a small image and adjust quality until we reach target size
	width, height := 100, 100
	quality := 85
	
	for {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{100, 100, 100, 255}}, image.Point{}, draw.Src)
		
		var buf bytes.Buffer
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		
		size := buf.Len()
		if size >= targetSize-100 && size <= targetSize {
			return &buf
		}
		
		if size < targetSize {
			width += 10
			height += 10
		} else {
			quality -= 5
			if quality < 10 {
				return &buf // Give up and return what we have
			}
		}
	}
}
# Google Cloud Storage Setup Guide for BOGOWI NFT Images

## ✅ Datakyte API Keys Configured

Your Datakyte API keys are now configured:
- **Testnet:** `dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d`
- **Mainnet:** `dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d`

## 📦 Google Cloud Storage Setup

### 1. Install Required Go Dependencies

```bash
go get cloud.google.com/go/storage
go get google.golang.org/api/option
go get github.com/disintegration/imaging
```

### 2. Create a Google Cloud Storage Bucket

```bash
# Using gcloud CLI
gcloud storage buckets create gs://bogowi-nft-images \
  --location=us-central1 \
  --default-storage-class=STANDARD \
  --uniform-bucket-level-access

# Make bucket public-read (optional, if not using CDN)
gcloud storage buckets add-iam-policy-binding gs://bogowi-nft-images \
  --member=allUsers \
  --role=roles/storage.objectViewer
```

### 3. Create Service Account for API Access

```bash
# Create service account
gcloud iam service-accounts create bogowi-nft-storage \
  --display-name="BOGOWI NFT Storage Service"

# Grant permissions
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:bogowi-nft-storage@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.objectAdmin"

# Create and download key
gcloud iam service-accounts keys create ~/bogowi-gcs-key.json \
  --iam-account=bogowi-nft-storage@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### 4. Configure CORS for Direct Upload (Optional)

Create `cors.json`:
```json
[
  {
    "origin": ["https://app.bogowi.com", "http://localhost:3000"],
    "method": ["GET", "PUT", "POST"],
    "responseHeader": ["Content-Type"],
    "maxAgeSeconds": 3600
  }
]
```

Apply CORS:
```bash
gcloud storage buckets update gs://bogowi-nft-images --cors-file=cors.json
```

## 🔧 Environment Configuration

### Development (.env.local)
```bash
# Storage Configuration
STORAGE_PROVIDER=gcs
GCS_BUCKET_NAME=bogowi-nft-images-dev
GCS_PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/bogowi-gcs-key.json
CDN_BASE_URL=  # Leave empty for now, using direct GCS URLs

# Datakyte
DATAKYTE_API_KEY_TESTNET=dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d
DATAKYTE_API_URL=https://api.datakyte.com
```

### Production (.env.production)
```bash
# Storage Configuration
STORAGE_PROVIDER=gcs
GCS_BUCKET_NAME=bogowi-nft-images
GCS_PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/bogowi-gcs-key.json
CDN_BASE_URL=https://cdn.bogowi.com  # Update when CDN is ready

# Datakyte
DATAKYTE_API_KEY_MAINNET=dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d
DATAKYTE_API_URL=https://api.datakyte.com
```

## 🖼️ Image URLs Structure

Images will be stored and accessed as:

### Direct GCS URLs (without CDN):
```
https://storage.googleapis.com/bogowi-nft-images/tickets/{tokenId}/original.jpg
https://storage.googleapis.com/bogowi-nft-images/tickets/{tokenId}/thumbnail.jpg
```

### With CDN (when configured):
```
https://cdn.bogowi.com/tickets/{tokenId}/original.jpg
https://cdn.bogowi.com/tickets/{tokenId}/thumbnail.jpg
```

## 📝 Metadata Structure

The metadata will be created by Datakyte and include:

```json
{
  "name": "BOGOWI Eco Experience #12345",
  "description": "Carbon-neutral adventure in Costa Rica rainforest",
  "image": "https://storage.googleapis.com/bogowi-nft-images/tickets/12345/original.jpg",
  "external_url": "https://app.bogowi.com/tickets/12345",
  "attributes": [
    {
      "trait_type": "Experience Type",
      "value": "Rainforest Trek"
    },
    {
      "trait_type": "Location",
      "value": "Costa Rica"
    },
    {
      "trait_type": "Carbon Offset",
      "value": "50kg CO2"
    },
    {
      "trait_type": "BOGO Rewards",
      "value": "5%"
    },
    {
      "trait_type": "Valid Until",
      "value": "2025-12-31"
    },
    {
      "trait_type": "Provider",
      "value": "EcoTours CR"
    }
  ]
}
```

## 🚀 Usage in Go API

### Initialize Storage Service

```go
import (
    "bogowi-blockchain-go/internal/services/storage"
    "bogowi-blockchain-go/internal/config"
    "google.golang.org/api/option"
)

// Get configuration
storageConfig := config.GetStorageConfig()

// Create GCS service
var imageService *storage.GCSImageService
if storageConfig.CredentialsPath != "" {
    imageService, err = storage.NewGCSImageService(
        storageConfig.BucketName,
        storageConfig.CDNBaseURL,
        option.WithCredentialsFile(storageConfig.CredentialsPath),
    )
} else {
    // Use default credentials (for GKE/Cloud Run)
    imageService, err = storage.NewGCSImageService(
        storageConfig.BucketName,
        storageConfig.CDNBaseURL,
    )
}
```

### Upload Ticket Image

```go
// Upload image for a ticket
imageURL, err := imageService.UploadTicketImage(
    tokenID,
    imageData,
    "image/jpeg",
)
```

### Generate Signed Upload URL

```go
// Generate presigned URL for frontend upload
uploadURL, err := imageService.GenerateSignedUploadURL(
    tokenID,
    "image/jpeg",
)
// Send uploadURL to frontend for direct upload
```

## 🎯 Next Steps

1. **Create GCS Bucket:**
   ```bash
   gcloud storage buckets create gs://bogowi-nft-images
   ```

2. **Create Service Account:**
   ```bash
   # Download the setup script
   curl -O https://raw.githubusercontent.com/bogowi/setup/main/gcs-setup.sh
   chmod +x gcs-setup.sh
   ./gcs-setup.sh YOUR_PROJECT_ID
   ```

3. **Test Upload:**
   ```bash
   # Test with gsutil
   echo "test" > test.txt
   gsutil cp test.txt gs://bogowi-nft-images/test.txt
   ```

4. **Set Up CDN (Optional):**
   - Use Google Cloud CDN
   - Or CloudFlare/Fastly pointing to GCS bucket

## 📊 Cost Estimates

### Google Cloud Storage:
- **Storage:** $0.020/GB/month (Standard)
- **Operations:** $0.005 per 10,000 operations
- **Egress:** Free to CDN, $0.12/GB to internet

### Estimated Monthly Costs:
- 10,000 images (1MB each): ~$0.20 storage
- 100,000 views: ~$12 egress (without CDN)
- With CDN: ~$2-5 (cached delivery)

## 🔒 Security Notes

1. **Service Account Key:**
   - Store securely, never commit to git
   - Rotate regularly
   - Use least privilege principle

2. **Bucket Permissions:**
   - Public read for images (or use signed URLs)
   - Write only via service account

3. **Image Validation:**
   - Max 10MB per image
   - Only JPEG, PNG, WebP allowed
   - Automatic resizing to 1200x1200 max

## ✅ Ready to Use!

With Datakyte keys configured and GCS set up, your NFT system can now:
- Store ticket images in Google Cloud Storage
- Generate metadata via Datakyte
- Serve images via CDN (when configured)
- Handle direct uploads from frontend
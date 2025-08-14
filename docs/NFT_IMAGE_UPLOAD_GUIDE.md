# NFT Ticket Image Upload Guide

## Overview

The BOGOWI NFT Tickets support custom images for each ticket. Images can be uploaded either:
1. **Via Pre-signed URLs** (recommended for frontend direct upload)
2. **Via Backend API** (for server-side processing)
3. **Using Default Images** (automatic fallback)

## Architecture

```
Frontend → Backend → S3/CDN → Datakyte
   ↓          ↓         ↓         ↓
Upload → Process → Store → Reference
```

## Image Requirements

- **Max Size**: 10MB
- **Formats**: JPEG, PNG, WebP
- **Recommended**: 1200x1200px minimum
- **Auto-resize**: Creates multiple sizes (original, large, medium, thumb)

## Option 1: Frontend Direct Upload (Recommended)

### Step 1: Get Pre-signed Upload URL

```javascript
// Frontend code
async function getUploadUrl(tokenId) {
  const response = await fetch(`/api/v2/nft/tickets/${tokenId}/upload-url`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'X-Network': 'testnet' // or 'mainnet'
    }
  });
  
  const data = await response.json();
  return data.uploadUrl;
}
```

### Step 2: Upload Image Directly to S3

```javascript
async function uploadImageToS3(file, uploadUrl) {
  const response = await fetch(uploadUrl, {
    method: 'PUT',
    body: file,
    headers: {
      'Content-Type': file.type
    }
  });
  
  return response.ok;
}
```

### Step 3: Complete Process

```javascript
async function uploadTicketImage(tokenId, imageFile) {
  try {
    // 1. Get pre-signed URL
    const uploadUrl = await getUploadUrl(tokenId);
    
    // 2. Upload to S3
    const uploaded = await uploadImageToS3(imageFile, uploadUrl);
    
    if (uploaded) {
      // 3. Image URL is automatically set
      const imageUrl = `https://storage.bogowi.com/tickets/${tokenId}/original.jpg`;
      
      // 4. Update UI
      updateTicketImage(tokenId, imageUrl);
      
      return imageUrl;
    }
  } catch (error) {
    console.error('Upload failed:', error);
  }
}
```

## Option 2: Backend Upload

```javascript
// Frontend code
async function uploadViaBackend(tokenId, imageFile) {
  const formData = new FormData();
  formData.append('image', imageFile);
  
  const response = await fetch(`/api/v2/nft/tickets/${tokenId}/image`, {
    method: 'POST',
    body: formData,
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'X-Network': 'testnet'
    }
  });
  
  const data = await response.json();
  return data.imageUrl;
}
```

## Option 3: Using Default Images

If no custom image is uploaded, the system automatically assigns a default image based on the experience type:

```javascript
const defaultImages = {
  'Wildlife Safari': 'wildlife-safari.jpg',
  'Marine Conservation': 'marine-conservation.jpg',
  'Forest Trek': 'forest-trek.jpg',
  'Cultural Experience': 'cultural-experience.jpg',
  'Eco Lodge': 'eco-lodge.jpg',
  'Adventure': 'adventure.jpg'
};
```

## React Component Example

```jsx
import React, { useState } from 'react';

function TicketImageUploader({ tokenId, onSuccess }) {
  const [uploading, setUploading] = useState(false);
  const [preview, setPreview] = useState(null);
  
  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      // Validate
      if (file.size > 10 * 1024 * 1024) {
        alert('File too large. Max 10MB.');
        return;
      }
      
      // Preview
      const reader = new FileReader();
      reader.onload = (e) => setPreview(e.target.result);
      reader.readAsDataURL(file);
      
      // Upload
      uploadImage(file);
    }
  };
  
  const uploadImage = async (file) => {
    setUploading(true);
    
    try {
      // Get pre-signed URL
      const urlResponse = await fetch(
        `/api/v2/nft/tickets/${tokenId}/upload-url`,
        {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'X-Network': 'testnet'
          }
        }
      );
      
      const { uploadUrl } = await urlResponse.json();
      
      // Upload to S3
      await fetch(uploadUrl, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type
        }
      });
      
      // Success
      const finalUrl = `https://storage.bogowi.com/tickets/${tokenId}/original.jpg`;
      onSuccess(finalUrl);
      
    } catch (error) {
      console.error('Upload failed:', error);
      alert('Upload failed. Please try again.');
    } finally {
      setUploading(false);
    }
  };
  
  return (
    <div className="image-uploader">
      {preview && (
        <img 
          src={preview} 
          alt="Preview" 
          style={{ maxWidth: '200px', marginBottom: '10px' }}
        />
      )}
      
      <input
        type="file"
        accept="image/jpeg,image/png,image/webp"
        onChange={handleFileSelect}
        disabled={uploading}
      />
      
      {uploading && <p>Uploading...</p>}
    </div>
  );
}
```

## Integration with Datakyte

After image upload, the metadata is automatically updated:

1. **Image uploaded to S3** → CDN URL generated
2. **Metadata updated** → Image URL added to NFT metadata
3. **Datakyte notified** → Metadata synced with Datakyte API
4. **Marketplace ready** → Image visible on OpenSea, etc.

## Image URLs Structure

For each uploaded image, multiple sizes are generated:

```json
{
  "original": "https://storage.bogowi.com/tickets/10001/original.jpg",
  "large": "https://storage.bogowi.com/tickets/10001/large.jpg",     // 800px wide
  "medium": "https://storage.bogowi.com/tickets/10001/medium.jpg",   // 400px wide
  "thumb": "https://storage.bogowi.com/tickets/10001/thumb.jpg"      // 200x200px
}
```

## Best Practices

1. **Validate on Frontend**: Check file size and type before upload
2. **Show Progress**: Display upload progress to users
3. **Handle Errors**: Implement retry logic for failed uploads
4. **Optimize Images**: Compress images before upload if possible
5. **Cache URLs**: Store generated URLs to avoid repeated API calls

## Security Considerations

1. **Pre-signed URLs expire** in 15 minutes
2. **File type validation** on both frontend and backend
3. **Size limits** enforced at multiple levels
4. **Token ownership** verified before allowing uploads
5. **Rate limiting** on upload endpoints

## Environment Variables

```bash
# S3 Configuration
TICKET_IMAGE_UPLOAD_BUCKET=bogowi-tickets
TICKET_IMAGE_BASE_URL=https://storage.bogowi.com/tickets

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret

# Feature Flags
ENABLE_IMAGE_UPLOAD=true
MAX_IMAGE_SIZE_MB=10
```

## API Endpoints

### Get Pre-signed Upload URL
```
GET /api/v2/nft/tickets/:tokenId/upload-url
Response: { uploadUrl, tokenId, contentType, expiresIn }
```

### Direct Upload
```
POST /api/v2/nft/tickets/:tokenId/image
Body: multipart/form-data with 'image' field
Response: { success, tokenId, imageUrl, sizes }
```

## Troubleshooting

### Upload Fails
- Check file size (max 10MB)
- Verify file format (JPEG, PNG, WebP)
- Ensure token exists and user owns it
- Check network connectivity

### Image Not Showing
- Verify CDN URL is accessible
- Check Datakyte metadata sync
- Ensure proper CORS headers
- Clear browser cache

### S3 Errors
- Verify AWS credentials
- Check bucket permissions
- Ensure bucket exists
- Review CORS configuration
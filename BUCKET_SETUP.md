# Bucket Setup Strategy for BOGOWI NFTs

## 🎯 Recommended: Separate Buckets per Network

### **Why Separate Buckets?**

1. **Clear Separation** - No risk of mixing test/production data
2. **Different Permissions** - Testnet can be more permissive
3. **Cost Tracking** - Easy to see costs per environment
4. **Easy Cleanup** - Delete entire testnet bucket when needed
5. **Different Retention** - Testnet images can auto-delete after 30 days

## 📦 Bucket Structure

### **Option 1: Separate Buckets (RECOMMENDED)**

```
bogowi-nft-images-testnet/     # Columbus testnet
├── tickets/
│   ├── 1/
│   │   ├── original.jpg
│   │   └── thumbnail.jpg
│   └── 2/
│       ├── original.jpg
│       └── thumbnail.jpg

bogowi-nft-images/              # Camino mainnet
├── tickets/
│   ├── 1/
│   │   ├── original.jpg
│   │   └── thumbnail.jpg
│   └── 2/
│       ├── original.jpg
│       └── thumbnail.jpg

bogowi-nft-images-local/        # Local development (optional)
└── tickets/
    └── ...
```

### **Option 2: Single Bucket with Folders**

```
bogowi-nft-images/
├── testnet/
│   └── tickets/
│       └── ...
├── mainnet/
│   └── tickets/
│       └── ...
└── local/
    └── tickets/
        └── ...
```

## 🛠️ Setup Commands

### Create Separate Buckets (Recommended)

```bash
# 1. Create testnet bucket (can be deleted/recreated freely)
gcloud storage buckets create gs://bogowi-nft-images-testnet \
  --location=us-central1 \
  --default-storage-class=STANDARD

# 2. Create mainnet bucket (production - be careful!)
gcloud storage buckets create gs://bogowi-nft-images \
  --location=us-central1 \
  --default-storage-class=STANDARD

# 3. Optional: Create local dev bucket
gcloud storage buckets create gs://bogowi-nft-images-local \
  --location=us-central1 \
  --default-storage-class=STANDARD

# 4. Set lifecycle rules for testnet (auto-delete after 30 days)
cat > lifecycle-testnet.json << EOF
{
  "lifecycle": {
    "rule": [
      {
        "action": {"type": "Delete"},
        "condition": {
          "age": 30,
          "matchesPrefix": ["tickets/"]
        }
      }
    ]
  }
}
EOF

gcloud storage buckets update gs://bogowi-nft-images-testnet \
  --lifecycle-file=lifecycle-testnet.json
```

### Make Buckets Public (for NFT viewing)

```bash
# Testnet - public read
gcloud storage buckets add-iam-policy-binding gs://bogowi-nft-images-testnet \
  --member=allUsers \
  --role=roles/storage.objectViewer

# Mainnet - public read
gcloud storage buckets add-iam-policy-binding gs://bogowi-nft-images \
  --member=allUsers \
  --role=roles/storage.objectViewer
```

## 💰 Cost Implications

### Separate Buckets:
- **Storage**: ~$0.02/GB/month per bucket
- **Operations**: Same cost, just split
- **Management**: Easier to track costs

### Example Monthly Costs:
```
Testnet Bucket:
- 1,000 test images (1MB each) = 1GB = $0.02
- Auto-deletion after 30 days = no accumulation

Mainnet Bucket:
- 10,000 production images (1MB each) = 10GB = $0.20
- Permanent storage

Total: ~$0.22/month
```

## 🔧 Configuration in Code

### Environment Variables

**.env.testnet**
```bash
NETWORK=testnet
GCS_BUCKET_NAME=bogowi-nft-images-testnet
DATAKYTE_API_KEY=dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d
```

**.env.mainnet**
```bash
NETWORK=mainnet
GCS_BUCKET_NAME=bogowi-nft-images
DATAKYTE_API_KEY=dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d
```

**.env.local**
```bash
NETWORK=local
GCS_BUCKET_NAME=bogowi-nft-images-local
# Or use testnet bucket for local
GCS_BUCKET_NAME=bogowi-nft-images-testnet
```

### Go Code Configuration

```go
// Determine bucket based on network
func GetBucketName(network string) string {
    switch network {
    case "mainnet":
        return os.Getenv("GCS_BUCKET_NAME_MAINNET")
    case "testnet":
        return os.Getenv("GCS_BUCKET_NAME_TESTNET")
    case "local":
        // Use testnet bucket for local, or separate local bucket
        return os.Getenv("GCS_BUCKET_NAME_LOCAL")
    default:
        return os.Getenv("GCS_BUCKET_NAME_TESTNET")
    }
}
```

## 🚨 Important Considerations

### For Testnet:
- ✅ Can delete and recreate anytime
- ✅ Can have auto-deletion policies
- ✅ More relaxed permissions
- ✅ Can test dangerous operations

### For Mainnet:
- ⚠️ NEVER delete production images
- ⚠️ Strict access controls
- ⚠️ Regular backups recommended
- ⚠️ Monitor for unauthorized access

## 📊 Monitoring

### Set up alerts:
```bash
# Alert if mainnet bucket exceeds 100GB
gcloud monitoring policies create \
  --notification-channels=YOUR_CHANNEL_ID \
  --display-name="Mainnet Bucket Size Alert" \
  --condition="storage.googleapis.com/storage/total_bytes > 100000000000"
```

## ✅ Recommendation

**Use SEPARATE BUCKETS** for:
- Clear separation of environments
- Easier cost tracking
- Safer testing on testnet
- No risk of mixing data

The small additional complexity is worth the safety and clarity!
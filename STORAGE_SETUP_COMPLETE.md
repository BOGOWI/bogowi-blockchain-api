# ✅ Storage Configuration Complete!

## 🪣 Bucket Setup

You've created separate buckets as recommended:

| Environment | Bucket Name | Service Account | Datakyte Key |
|------------|-------------|-----------------|--------------|
| **Local** | `bogowi-nft-images-testnet` | `testnet-service-account.json` | Testnet key |
| **Testnet** | `bogowi-nft-images-testnet` | `testnet-service-account.json` | `dk_d707e26c...` |
| **Mainnet** | `bogowi-nft-images` | `mainnet-service-account.json` | `dk_e2aad94d...` |

## 📁 Files Created

### Environment Files:
- `.env.local` - Local development (uses testnet bucket)
- `.env.testnet` - Columbus testnet 
- `.env.mainnet` - Camino mainnet

### Configuration:
- `/internal/config/storage_config.go` - Dynamic storage config
- `/internal/services/storage/gcs_service.go` - GCS implementation

### Testing:
- `/scripts/test-storage.sh` - Test storage setup

## 🔑 Service Account Setup

Place your service account JSON files in the project root:

```bash
bogowi-blockchain-go/
├── testnet-service-account.json   # For local & testnet
├── mainnet-service-account.json   # For mainnet only
└── .gitignore                      # Make sure *.json is ignored!
```

## 🚀 Quick Start

### 1. Add Service Account Files

Download from Google Cloud Console and save as:
- `testnet-service-account.json`
- `mainnet-service-account.json`

### 2. Test Your Setup

```bash
# Test testnet storage
./scripts/test-storage.sh testnet

# Test mainnet storage
./scripts/test-storage.sh mainnet

# Test local (uses testnet bucket)
./scripts/test-storage.sh local
```

### 3. Run API with Correct Environment

```bash
# Local development
source .env.local
go run cmd/api/main.go

# Testnet
source .env.testnet
go run cmd/api/main.go

# Mainnet (be careful!)
source .env.mainnet
go run cmd/api/main.go
```

## 🖼️ Image URLs

Your NFT images will be accessible at:

### Testnet/Local:
```
https://storage.googleapis.com/bogowi-nft-images-testnet/tickets/{tokenId}/original.jpg
https://storage.googleapis.com/bogowi-nft-images-testnet/tickets/{tokenId}/thumbnail.jpg
```

### Mainnet:
```
https://storage.googleapis.com/bogowi-nft-images/tickets/{tokenId}/original.jpg
https://storage.googleapis.com/bogowi-nft-images/tickets/{tokenId}/thumbnail.jpg
```

## 📝 Metadata Structure

When minting an NFT, the metadata will include:

```json
{
  "name": "BOGOWI Experience #123",
  "description": "Eco-friendly adventure experience",
  "image": "https://storage.googleapis.com/bogowi-nft-images/tickets/123/original.jpg",
  "external_url": "https://app.bogowi.com/experience/123",
  "attributes": [
    {"trait_type": "Experience Type", "value": "Rainforest Trek"},
    {"trait_type": "Location", "value": "Costa Rica"},
    {"trait_type": "Carbon Offset", "value": "50kg CO2"},
    {"trait_type": "BOGO Rewards", "value": "5%"}
  ]
}
```

## 🔒 Security Checklist

- [ ] Service account files added to `.gitignore`
- [ ] Never commit service account JSON files
- [ ] Use separate service accounts for testnet/mainnet
- [ ] Buckets are public-read (required for NFTs)
- [ ] Service accounts have minimal permissions

## 📊 Cost Estimates

- **Storage**: ~$0.02/GB/month
- **Operations**: ~$0.005 per 10,000 operations
- **Network**: Free within Google Cloud, ~$0.12/GB external

Example: 10,000 images (1MB each) = 10GB = $0.20/month

## ✅ You're Ready!

Your storage is configured with:
- ✅ Separate buckets for testnet/mainnet
- ✅ Datakyte API keys configured
- ✅ Google Cloud Storage ready
- ✅ Public access for NFT images
- ✅ Service account separation

Start minting NFTs with images! 🎨
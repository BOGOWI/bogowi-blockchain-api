# NFT Metadata Storage Implementation Plan

## Overview
BOGOWI Tickets NFT metadata needs reliable, decentralized storage with fast retrieval and update capabilities.

## Storage Options Comparison

### Option 1: IPFS (Recommended for Decentralization)
**Pros:**
- Fully decentralized
- Content-addressed (immutable)
- Censorship resistant
- Industry standard for NFTs

**Cons:**
- Requires pinning service (Pinata, Infura, etc.)
- Slower retrieval times
- No native update mechanism

**Implementation:**
```javascript
// Example IPFS upload
const metadata = {
  name: "BOGOWI Ticket #10001",
  description: "Eco-experience ticket",
  image: "ipfs://QmXxx...",
  attributes: [...]
};
const cid = await ipfs.add(JSON.stringify(metadata));
const metadataURI = `ipfs://${cid}`;
```

### Option 2: Arweave (Recommended for Permanence)
**Pros:**
- Permanent storage (pay once)
- Decentralized
- Fast retrieval via gateways

**Cons:**
- Higher upfront cost
- Cannot update metadata

### Option 3: Hybrid Approach (Recommended for Production)
**Best of both worlds:**
1. Store static assets (images) on IPFS/Arweave
2. Store dynamic metadata on cloud with IPFS backup
3. Use smart contract to point to current metadata location

```solidity
// Contract can update base URI if needed
string private _baseTokenURI = "https://api.bogowi.com/nft/metadata/";

function tokenURI(uint256 tokenId) returns (string memory) {
    return string(abi.encodePacked(_baseTokenURI, tokenId.toString()));
}
```

## Recommended Architecture

### 1. Metadata Service Components

```
┌─────────────────────────────────────────────────┐
│                 User Mints NFT                   │
└────────────────┬────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────────┐
│           BOGOWI Backend API                     │
│  - Generate metadata JSON                        │
│  - Create QR code                               │
│  - Generate ticket image                        │
└────────────────┬────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────────┐
│           Storage Pipeline                       │
│  1. Upload image to IPFS                        │
│  2. Upload metadata to IPFS                     │
│  3. Cache in PostgreSQL                         │
│  4. CDN distribution                            │
└────────────────┬────────────────────────────────┘
                 ▼
┌─────────────────────────────────────────────────┐
│         Smart Contract Minting                   │
│  - Store IPFS URI or API endpoint               │
│  - Emit events with metadata hash               │
└─────────────────────────────────────────────────┘
```

### 2. Database Schema

```sql
CREATE TABLE nft_metadata (
    token_id BIGINT PRIMARY KEY,
    contract_address VARCHAR(42) NOT NULL,
    metadata_uri VARCHAR(255) NOT NULL,
    ipfs_hash VARCHAR(66),
    metadata_json JSONB NOT NULL,
    image_url VARCHAR(255),
    qr_code_data TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    INDEX idx_contract_token (contract_address, token_id),
    INDEX idx_status (status)
);

CREATE TABLE nft_attributes (
    id SERIAL PRIMARY KEY,
    token_id BIGINT REFERENCES nft_metadata(token_id),
    trait_type VARCHAR(100) NOT NULL,
    value VARCHAR(255) NOT NULL,
    display_type VARCHAR(50),
    INDEX idx_token_traits (token_id, trait_type)
);
```

### 3. API Endpoints

```yaml
/api/v2/nft/metadata/{tokenId}:
  get:
    description: Get NFT metadata
    responses:
      200:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/NFTMetadata'

/api/v2/nft/metadata/generate:
  post:
    description: Generate metadata for new NFT
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            properties:
              bookingId: string
              eventId: string
              experienceDetails: object
              recipient: string

/api/v2/nft/metadata/{tokenId}/image:
  get:
    description: Get ticket image with QR code
    responses:
      200:
        content:
          image/png:
            schema:
              type: string
              format: binary
```

### 4. Image Generation Service

```go
package nft

import (
    "image"
    "image/png"
    "github.com/skip2/go-qrcode"
)

type TicketImageGenerator struct {
    templatePath string
    qrSize      int
}

func (g *TicketImageGenerator) GenerateTicketImage(
    tokenId uint64,
    bookingId string,
    experienceTitle string,
    location string,
    validUntil time.Time,
) ([]byte, error) {
    // 1. Load template image
    template := g.loadTemplate()
    
    // 2. Generate QR code with redemption data
    qrData := fmt.Sprintf("bogowi://redeem/%d/%s", tokenId, bookingId)
    qr, _ := qrcode.New(qrData, qrcode.Medium)
    
    // 3. Composite QR onto template
    // 4. Add text overlays (title, location, dates)
    // 5. Return PNG bytes
    
    return imageBytes, nil
}
```

### 5. Metadata Update Strategy

For dynamic attributes (like redemption status):

1. **On-chain events** - Emit events for status changes
2. **Off-chain indexing** - Index events and update database
3. **API serves latest** - API combines on-chain + off-chain data

```javascript
// Frontend can query both
const tokenURI = await contract.tokenURI(tokenId);
const metadata = await fetch(tokenURI).then(r => r.json());
const onChainStatus = await contract.isRedeemed(tokenId);
metadata.attributes.find(a => a.trait_type === "Status").value = 
    onChainStatus ? "Redeemed" : "Active";
```

## Implementation Steps

1. **Phase 1: Basic Metadata (Week 1)**
   - [ ] Set up IPFS node or Pinata account
   - [ ] Create metadata generation service
   - [ ] Implement basic image templates
   - [ ] Database schema setup

2. **Phase 2: Storage Pipeline (Week 2)**
   - [ ] IPFS upload service
   - [ ] Metadata caching layer
   - [ ] CDN configuration
   - [ ] API endpoints

3. **Phase 3: Dynamic Features (Week 3)**
   - [ ] QR code generation
   - [ ] Real-time status updates
   - [ ] Event indexing service
   - [ ] Metadata refresh mechanism

4. **Phase 4: Production Ready (Week 4)**
   - [ ] Load testing
   - [ ] Backup strategies
   - [ ] Monitoring and alerts
   - [ ] Documentation

## Cost Estimates

### IPFS (via Pinata)
- Free tier: 1GB storage, 100 pins
- Paid: $20/month for 50GB, unlimited pins
- Estimated need: ~10KB per NFT = 100,000 NFTs per GB

### Arweave
- ~$0.005 per NFT for permanent storage
- One-time cost, no recurring fees

### Hybrid (Recommended)
- Cloud storage: ~$50/month
- IPFS pinning: $20/month
- CDN: ~$100/month
- **Total: ~$170/month for 1M NFTs**

## Security Considerations

1. **Access Control**
   - Metadata generation requires authentication
   - Rate limiting on API endpoints
   - CORS configuration for allowed domains

2. **Data Integrity**
   - Hash metadata and store on-chain
   - Verify IPFS content matches hash
   - Backup critical data

3. **Privacy**
   - Don't expose personal data in public metadata
   - Separate public/private attributes
   - Encrypted storage for sensitive data

## Conclusion

Recommended approach:
1. Start with hybrid model (API + IPFS backup)
2. Use IPFS for images, API for dynamic metadata
3. Implement caching and CDN for performance
4. Plan migration to full IPFS when stable
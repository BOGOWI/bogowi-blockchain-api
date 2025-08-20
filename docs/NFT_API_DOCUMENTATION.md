# NFT API Documentation

## Overview
The BOGOWI NFT API provides endpoints for minting, managing, and querying NFT tickets on the blockchain. All NFT operations are performed on real blockchain networks (testnet/mainnet) with metadata stored in Datakyte and mappings persisted in a local SQLite database.

## Base URL
- Production: `https://web3.bogowi.com/api`
- Testnet: `https://testnet.web3.bogowi.com/api`
- Local: `http://localhost:8080/api`

## Authentication
Currently, the API uses network selection via headers:
- `X-Network-Type`: `testnet` or `mainnet` (defaults to `testnet`)

## Endpoints

### 1. Mint NFT Ticket
**POST** `/nft/tickets/mint`

Mints a new NFT ticket on the blockchain and creates metadata in Datakyte.

#### Request Headers
```
X-Network-Type: testnet
Content-Type: application/json
```

#### Request Body
```json
{
  "to": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
  "bookingId": "BOOK-2024-001",
  "eventId": "EVENT-ECO-001",
  "experienceTitle": "Amazon Rainforest Conservation Tour",
  "experienceType": "Eco-Adventure",
  "location": "Manaus, Brazil",
  "duration": "7 days",
  "maxParticipants": 20,
  "carbonOffset": 500,
  "conservationImpact": "Protecting 100 hectares of rainforest",
  "validUntil": "2024-12-31T23:59:59Z",
  "transferableAfter": "2024-07-01T00:00:00Z",
  "expiresAt": "2025-01-31T23:59:59Z",
  "rewardBasisPoints": 250,
  "recipientName": "John Doe",
  "providerName": "EcoTours Brazil",
  "providerContact": "contact@ecotours.br",
  "imageUrl": "https://storage.bogowi.com/tickets/amazon-tour.jpg"
}
```

#### Response
```json
{
  "success": true,
  "tokenId": 10001,
  "txHash": "0x123abc...",
  "metadataUri": "https://metadata.datakyte.io/nft/10001",
  "datakyteId": "dk_nft_abc123"
}
```

### 2. Batch Mint Tickets
**POST** `/nft/tickets/batch-mint`

Mints multiple NFT tickets in a single transaction (max 100).

#### Request Body
```json
{
  "tickets": [
    {
      "to": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
      "bookingId": "BOOK-2024-001",
      "eventId": "EVENT-ECO-001",
      // ... same fields as single mint
    },
    {
      "to": "0x123...",
      "bookingId": "BOOK-2024-002",
      // ... more tickets
    }
  ]
}
```

#### Response
```json
{
  "success": true,
  "message": "Successfully batch minted 2 tickets",
  "results": [
    {
      "success": true,
      "tokenId": 10001,
      "txHash": "0x123...",
      "metadataUri": "...",
      "datakyteId": "dk_nft_abc123"
    },
    {
      "success": true,
      "tokenId": 10002,
      "txHash": "0x123...",
      "metadataUri": "...",
      "datakyteId": "dk_nft_def456"
    }
  ],
  "txHash": "0xbatch123..."
}
```

### 3. Get Ticket Metadata
**GET** `/nft/tickets/{tokenId}/metadata`

Retrieves metadata for a specific NFT ticket from Datakyte.

#### Response
```json
{
  "name": "BOGOWI Eco-Adventure #10001",
  "description": "Amazon Rainforest Conservation Tour",
  "image": "https://storage.bogowi.com/tickets/10001.jpg",
  "attributes": [
    {
      "trait_type": "Experience Type",
      "value": "Eco-Adventure"
    },
    {
      "trait_type": "Location",
      "value": "Manaus, Brazil"
    },
    {
      "trait_type": "Conservation Impact",
      "value": "Protecting 100 hectares of rainforest"
    },
    {
      "trait_type": "BOGO Rewards",
      "value": 250
    }
  ]
}
```

### 4. Redeem Ticket
**POST** `/nft/tickets/redeem`

Redeems an NFT ticket using EIP-712 signature verification.

#### Request Body
```json
{
  "tokenId": 10001,
  "redeemer": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
  "nonce": 1234567890,
  "deadline": 1735689600,
  "signature": "0xabc123..."
}
```

#### Response
```json
{
  "success": true,
  "message": "Ticket redeemed successfully",
  "tokenId": 10001,
  "txHash": "0xredeem123..."
}
```

### 5. Get User's Tickets
**GET** `/nft/users/{address}/tickets`

Retrieves all NFT tickets owned by a specific address.

#### Response
```json
{
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
  "balance": "3",
  "tickets": [
    {
      "tokenId": 10001,
      "state": "active",
      "expiresAt": 1735689600,
      "transferUnlockAt": 1719792000,
      "bookingId": "0xabc123...",
      "eventId": "0xdef456...",
      "metadata": {
        "name": "BOGOWI Eco-Adventure #10001",
        "image": "https://storage.bogowi.com/tickets/10001.jpg"
      }
    }
  ]
}
```

### 6. Upload Ticket Image
**POST** `/nft/tickets/{tokenId}/image`

Uploads an image for an NFT ticket.

#### Request
- Content-Type: `multipart/form-data`
- Field: `image` (file)

#### Response
```json
{
  "success": true,
  "tokenId": 10001,
  "imageUrl": "https://storage.googleapis.com/bogowi-nft-images/10001.jpg",
  "sizes": {
    "original": "https://storage.googleapis.com/bogowi-nft-images/10001.jpg",
    "thumbnail": "https://storage.googleapis.com/bogowi-nft-images/10001_thumb.jpg",
    "medium": "https://storage.googleapis.com/bogowi-nft-images/10001_medium.jpg"
  },
  "message": "Image uploaded successfully"
}
```

### 7. Get Presigned Upload URL
**GET** `/nft/tickets/{tokenId}/upload-url`

Generates a presigned URL for direct image upload to storage.

#### Query Parameters
- `contentType`: Image content type (default: `image/jpeg`)

#### Response
```json
{
  "uploadUrl": "https://storage.googleapis.com/...",
  "tokenId": 10001,
  "contentType": "image/jpeg",
  "expiresIn": 900
}
```

### 8. Update Ticket Status
**PUT** `/nft/tickets/{tokenId}/status`

Updates the status of an NFT ticket.

#### Request Body
```json
{
  "status": "Redeemed"
}
```

Valid statuses: `Active`, `Redeemed`, `Expired`, `Burned`

#### Response
```json
{
  "success": true,
  "tokenId": 10001,
  "status": "Redeemed"
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid parameters)
- `404`: Not Found (token or resource not found)
- `500`: Internal Server Error

## Database Schema

The API uses SQLite to store NFT mappings between blockchain token IDs and Datakyte NFT IDs:

```sql
CREATE TABLE nft_token_mappings (
    id INTEGER PRIMARY KEY,
    token_id INTEGER NOT NULL,
    datakyte_nft_id TEXT NOT NULL,
    network TEXT NOT NULL,
    contract_address TEXT NOT NULL,
    owner_address TEXT NOT NULL,
    booking_id TEXT,
    event_id TEXT,
    status TEXT DEFAULT 'active',
    metadata_uri TEXT,
    image_url TEXT,
    tx_hash TEXT NOT NULL,
    minted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    redeemed_at DATETIME,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(token_id, network, contract_address)
);
```

## Smart Contract Addresses

### Testnet (Camino Columbus - Chain ID: 501)
- BOGOWITickets: `0x3Aa5ebB10DC797CAC828524e59A333d0A371443c`
- NFTRegistry: `0x68B1D87F95878fE05B998F19b66F4baba5De1aed`
- RoleManager: `0x9A676e781A523b5d0C0e43731313A708CB607508`

### Mainnet (Camino - Chain ID: 500)
- BOGOWITickets: *To be deployed*
- NFTRegistry: *To be deployed*
- RoleManager: *To be deployed*

## Testing

### Local Development
1. Run local blockchain: `npx hardhat node`
2. Deploy contracts: `npm run deploy:local`
3. Start API: `go run main.go`
4. Test endpoints: Use the provided Postman collection or curl commands

### Example Mint Request (curl)
```bash
curl -X POST http://localhost:8080/api/nft/tickets/mint \
  -H "Content-Type: application/json" \
  -H "X-Network-Type: testnet" \
  -d '{
    "to": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1",
    "bookingId": "BOOK-TEST-001",
    "eventId": "EVENT-TEST-001",
    "experienceTitle": "Test Experience",
    "experienceType": "Adventure",
    "location": "Test Location",
    "duration": "1 day",
    "validUntil": "2024-12-31T23:59:59Z",
    "transferableAfter": "2024-01-01T00:00:00Z",
    "expiresAt": "2025-12-31T23:59:59Z",
    "rewardBasisPoints": 100
  }'
```

## Integration

### Datakyte
- Testnet API Key: `dk_d707e26c919e72ab2bb3b81897566c393f4e2eba54d07ff680d765ee03d6cc5d`
- Mainnet API Key: `dk_e2aad94de12a2a7e7865a70b369e1eab69e2b5e2896577a5fbcbbb50d709bd3d`

### Google Cloud Storage
- Testnet Bucket: `bogowi-nft-images-testnet`
- Mainnet Bucket: `bogowi-nft-images-mainnet`

## Notes

1. **No More Hardcoded Data**: All NFT operations interact with real blockchain contracts
2. **Database Persistence**: SQLite database stores all NFT mappings locally
3. **Automatic Metadata**: Datakyte metadata is created automatically during minting
4. **Image Processing**: Images are automatically resized and optimized when uploaded
5. **Network Isolation**: Testnet and mainnet operations are completely isolated

## Support

For issues or questions, please contact the BOGOWI development team or create an issue in the GitHub repository.
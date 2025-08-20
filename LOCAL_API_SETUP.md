# Local API Configuration Guide

## ✅ Configuration Complete!

Your Go API is now configured to work with the locally deployed smart contracts on Hardhat network.

## 🔧 Configuration Files Created

1. **`.env.local`** - Environment variables for the Go API
2. **`config.local.yaml`** - YAML configuration file for structured config
3. **`run-local-api.sh`** - Startup script to run API with local config
4. **`test-local-connection.go`** - Test script to verify connectivity

## 📋 Contract Addresses (Local Deployment)

```
Network:        localhost
Chain ID:       501 (Camino testnet config)
RPC URL:        http://127.0.0.1:8545

Contracts:
- RoleManager:   0x9A676e781A523b5d0C0e43731313A708CB607508
- BOGOToken:     0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1
- NFTRegistry:   0x68B1D87F95878fE05B998F19b66F4baba5De1aed
- BOGOWITickets: 0x3Aa5ebB10DC797CAC828524e59A333d0A371443c
```

## 🚀 Quick Start

### 1. Start Local Blockchain (if not running)
```bash
cd contracts/v1
npx hardhat node
```

### 2. Deploy Contracts (if not deployed)
```bash
cd contracts/v1
npm run deploy-nft-local
```

### 3. Run Go API with Local Config

#### Option A: Using the startup script
```bash
./run-local-api.sh
```

#### Option B: Using environment variables
```bash
source .env.local
go run cmd/api/main.go
```

#### Option C: Using config file
```bash
go run cmd/api/main.go --config config.local.yaml
```

### 4. Test the Connection
```bash
go run test-local-connection.go
```

## 🧪 Testing the API

### Test Direct Contract Interaction
```bash
cd contracts/v1
node scripts/test-contracts-direct.js
```

### Test Go API Endpoints
```bash
cd contracts/v1
node scripts/test-go-api-local.js
```

### Manual Testing with curl
```bash
# Check API health
curl http://localhost:8080/health

# Get NFT ticket info (replace with actual ticket ID)
curl http://localhost:8080/api/v1/nft/tickets/10001 \
  -H "X-API-Key: test-api-key"

# Mint a ticket (requires proper signature)
curl -X POST http://localhost:8080/api/v1/nft/tickets/mint \
  -H "Content-Type: application/json" \
  -H "X-API-Key: test-api-key" \
  -d '{
    "recipient": "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
    "ticketId": "12345",
    "rewardBasisPoints": 500,
    "metadataURI": "https://api.bogowi.com/metadata/12345",
    "network": "localhost"
  }'
```

## 🔑 Test Accounts

### Admin Account
- Address: `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
- Private Key: `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`

### Minter Account
- Address: `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
- Private Key: `0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a`

### Backend Signer Account
- Address: `0x90F79bf6EB2c4f870365E785982E1f101E93b906`
- Private Key: `0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6`

### Test Users
- User1: `0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65`
- User2: `0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc`

## ⚠️ Important Notes

1. **These are Hardhat default accounts** - NEVER use these private keys in production!
2. **Chain ID 501** - The local Hardhat is configured to use Camino testnet chain ID
3. **Signatures** - When creating signatures for redemption, use chainId: 501
4. **Gas** - Local network has unlimited gas, but the API may have limits configured

## 📝 Environment Variables Reference

Key environment variables your Go API needs:

```bash
# Network
RPC_URL=http://127.0.0.1:8545
CHAIN_ID=501
NETWORK=localhost

# Contracts
ROLE_MANAGER_ADDRESS=0x9A676e781A523b5d0C0e43731313A708CB607508
BOGO_TOKEN_ADDRESS=0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1
NFT_REGISTRY_ADDRESS=0x68B1D87F95878fE05B998F19b66F4baba5De1aed
BOGOWI_TICKETS_ADDRESS=0x3Aa5ebB10DC797CAC828524e59A333d0A371443c

# Signers
BACKEND_PRIVATE_KEY=0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6
MINTER_PRIVATE_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

# API
API_KEY=test-api-key
API_PORT=8080
```

## 🐛 Troubleshooting

### "Contract not found" error
- Make sure contracts are deployed: `npm run deploy-nft-local`
- Check the deployment file exists: `contracts/v1/scripts/deployment-nft-localhost.json`

### "Invalid signature" error
- Ensure you're using chainId: 501 in signature creation
- Check the backend signer has the BACKEND_ROLE

### "Connection refused" error
- Verify Hardhat node is running: `npx hardhat node`
- Check RPC URL is correct: `http://127.0.0.1:8545`

### "Unauthorized" error
- Include API key header: `X-API-Key: test-api-key`
- Check the minter has NFT_MINTER_ROLE

## ✅ Ready to Test!

Your local environment is fully configured. The Go API can now:
- Connect to local blockchain ✅
- Interact with deployed contracts ✅
- Mint NFT tickets ✅
- Query ticket data ✅
- Perform batch operations ✅

Start testing with `./run-local-api.sh`!
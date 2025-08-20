#!/bin/bash

# Script to run BOGOWI Go API with local blockchain configuration
# Usage: ./run-local-api.sh

echo "🚀 Starting BOGOWI API with Local Configuration"
echo "================================================"

# Check if Hardhat node is running
echo "🔍 Checking local blockchain..."
if curl -s -X POST http://127.0.0.1:8545 \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' > /dev/null 2>&1; then
    echo "✅ Local blockchain is running on http://127.0.0.1:8545"
else
    echo "❌ Local blockchain is not running!"
    echo "   Please run: cd contracts/v1 && npx hardhat node"
    exit 1
fi

# Check if contracts are deployed
if [ -f "contracts/v1/scripts/deployment-nft-localhost.json" ]; then
    echo "✅ Contracts are deployed"
    echo ""
    echo "📋 Contract Addresses:"
    echo "   RoleManager:    0x9A676e781A523b5d0C0e43731313A708CB607508"
    echo "   BOGOToken:      0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1"
    echo "   NFTRegistry:    0x68B1D87F95878fE05B998F19b66F4baba5De1aed"
    echo "   BOGOWITickets:  0x3Aa5ebB10DC797CAC828524e59A333d0A371443c"
else
    echo "⚠️  Contracts may not be deployed"
    echo "   Run: cd contracts/v1 && npm run deploy-nft-local"
fi

echo ""
echo "🔧 Setting environment variables..."

# Export environment variables from .env.local
if [ -f ".env.local" ]; then
    export $(cat .env.local | grep -v '^#' | xargs)
    echo "✅ Environment variables loaded from .env.local"
else
    # Set inline if .env.local doesn't exist
    export NETWORK=localhost
    export CHAIN_ID=501
    export RPC_URL=http://127.0.0.1:8545
    export ROLE_MANAGER_ADDRESS=0x9A676e781A523b5d0C0e43731313A708CB607508
    export BOGO_TOKEN_ADDRESS=0x959922bE3CAee4b8Cd9a407cc3ac1C251C2007B1
    export NFT_REGISTRY_ADDRESS=0x68B1D87F95878fE05B998F19b66F4baba5De1aed
    export BOGOWI_TICKETS_ADDRESS=0x3Aa5ebB10DC797CAC828524e59A333d0A371443c
    export NFT_TICKETS_CONTRACT=0x3Aa5ebB10DC797CAC828524e59A333d0A371443c
    export NFT_TICKETS_TESTNET_CONTRACT=0x3Aa5ebB10DC797CAC828524e59A333d0A371443c
    export BACKEND_PRIVATE_KEY=0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6
    export MINTER_PRIVATE_KEY=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
    export API_KEY=test-api-key
    export LOG_LEVEL=debug
    echo "✅ Environment variables set inline"
fi

echo ""
echo "🌐 API Configuration:"
echo "   Network:  $NETWORK (Chain ID: $CHAIN_ID)"
echo "   RPC URL:  $RPC_URL"
echo "   API Port: 8080"
echo "   API Key:  $API_KEY"
echo ""

# Check if config file exists
CONFIG_FILE=""
if [ -f "config.local.yaml" ]; then
    CONFIG_FILE="--config config.local.yaml"
    echo "📝 Using config file: config.local.yaml"
fi

echo "================================================"
echo "🚀 Starting Go API..."
echo ""

# Run the API
if [ -f "cmd/api/main.go" ]; then
    go run cmd/api/main.go $CONFIG_FILE
elif [ -f "main.go" ]; then
    go run main.go $CONFIG_FILE
else
    echo "❌ Could not find main.go file"
    echo "   Please run this script from the project root"
    exit 1
fi
#!/bin/bash

# BOGOWI V1 Secure Deployment Script
# Uses macOS Keychain for private key management

set -e  # Exit on error

echo "🔐 BOGOWI V1 Secure Deployment"
echo "================================"

# Determine network
NETWORK=${1:-testnet}

if [ "$NETWORK" = "mainnet" ]; then
    echo "🌐 Deploying to MAINNET"
    KEYCHAIN_ACCOUNT="bogowi-mainnet"
    RPC_URL="https://api.camino.network/ext/bc/C/rpc"
    ADMIN_ADDRESS="0x444ddA4cA50765D3c0c0c662aAecF3b5D49761Ea"
elif [ "$NETWORK" = "testnet" ]; then
    echo "🧪 Deploying to TESTNET (Columbus)"
    KEYCHAIN_ACCOUNT="bogowi-testnet"
    RPC_URL="https://columbus.camino.network/ext/bc/C/rpc"
    ADMIN_ADDRESS="0xB34A822F735CDE477cbB39a06118267D00948ef7"
else
    echo "❌ Invalid network. Use: ./deploy-secure.sh [testnet|mainnet]"
    exit 1
fi

# Retrieve private key from keychain
echo "🔑 Retrieving private key from macOS Keychain..."
PRIVATE_KEY=$(security find-generic-password -a "$KEYCHAIN_ACCOUNT" -s "${KEYCHAIN_ACCOUNT}-pk" -w 2>/dev/null)

if [ -z "$PRIVATE_KEY" ]; then
    echo "❌ Private key not found in keychain!"
    echo "   Please add it first with:"
    echo "   security add-generic-password -a \"$KEYCHAIN_ACCOUNT\" -s \"${KEYCHAIN_ACCOUNT}-pk\" -w"
    exit 1
fi

# Export environment variables
export PRIVATE_KEY
export RPC_URL
export ADMIN_ADDRESS

echo "✅ Configuration loaded:"
echo "   Network: $NETWORK"
echo "   RPC URL: $RPC_URL"
echo "   Admin: $ADMIN_ADDRESS"
echo "   Private Key: [SECURED IN KEYCHAIN]"
echo ""

# Change to contracts directory
cd "$(dirname "$0")/.."
echo "📁 Working directory: $(pwd)"
echo ""

# Run deployment steps
echo "🚀 Starting deployment sequence..."
echo "================================"

echo "1️⃣ Deploying contracts..."
npx hardhat run scripts/deploy-all.js --network $NETWORK

echo ""
echo "2️⃣ Registering contracts..."
npx hardhat run scripts/register-contracts.js --network $NETWORK

echo ""
echo "3️⃣ Setting up roles..."
npx hardhat run scripts/setup-roles.js --network $NETWORK

echo ""
echo "4️⃣ Minting initial supply..."
npx hardhat run scripts/mint-initial-supply.js --network $NETWORK

echo ""
echo "5️⃣ Funding distributor..."
npx hardhat run scripts/fund-distributor.js --network $NETWORK

echo ""
echo "6️⃣ Verifying deployment..."
npx hardhat run scripts/verify-deployment.js --network $NETWORK

echo ""
echo "✅ Deployment complete!"
echo ""
echo "⚠️  IMPORTANT:"
echo "1. Save the deployment-${NETWORK}.json file"
echo "2. Update .env with the new contract addresses"
echo "3. Remove BUSINESS_ROLE from deployer for security"

# Clear sensitive variables
unset PRIVATE_KEY
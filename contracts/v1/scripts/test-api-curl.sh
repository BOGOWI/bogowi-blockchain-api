#!/bin/bash

# NFT API Testing with curl
# Make sure the API server is running: node scripts/test-nft-api-local.js

API_URL="http://localhost:3000"
USER1="0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65"
USER2="0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc"

echo "🚀 NFT API Testing with curl"
echo "================================"

# 1. Check status
echo -e "\n📊 Checking system status..."
curl -s "$API_URL/status" | jq '.'

# 2. Get statistics
echo -e "\n📈 Getting statistics..."
curl -s "$API_URL/stats" | jq '.'

# 3. Mint a single ticket
echo -e "\n🎫 Minting single ticket..."
TICKET_ID=$((RANDOM % 1000000))
curl -s -X POST "$API_URL/mint" \
  -H "Content-Type: application/json" \
  -d "{
    \"recipient\": \"$USER1\",
    \"ticketId\": $TICKET_ID,
    \"rewardBasisPoints\": 500,
    \"metadataURI\": \"https://api.bogowi.com/metadata/$TICKET_ID\"
  }" | jq '.'

# 4. Get ticket details
echo -e "\n🔍 Getting ticket details..."
curl -s "$API_URL/ticket/$TICKET_ID" | jq '.'

# 5. Batch mint
echo -e "\n🎫 Batch minting tickets..."
BASE_ID=$((RANDOM % 1000000))
curl -s -X POST "$API_URL/mint-batch" \
  -H "Content-Type: application/json" \
  -d "{
    \"recipients\": [\"$USER1\", \"$USER2\", \"$USER1\"],
    \"ticketIds\": [$BASE_ID, $((BASE_ID + 1)), $((BASE_ID + 2))],
    \"rewardBasisPoints\": [300, 500, 1000],
    \"metadataURIs\": [
      \"https://api.bogowi.com/metadata/$BASE_ID\",
      \"https://api.bogowi.com/metadata/$((BASE_ID + 1))\",
      \"https://api.bogowi.com/metadata/$((BASE_ID + 2))\"
    ]
  }" | jq '.'

# 6. Get user tickets
echo -e "\n👤 Getting User 1 tickets..."
curl -s "$API_URL/user/$USER1/tickets" | jq '.'

echo -e "\n👤 Getting User 2 tickets..."
curl -s "$API_URL/user/$USER2/tickets" | jq '.'

# 7. Get registry
echo -e "\n📚 Getting registry contents..."
curl -s "$API_URL/registry" | jq '.'

# 8. Redeem a ticket
echo -e "\n🔥 Redeeming ticket #$TICKET_ID..."
curl -s -X POST "$API_URL/redeem" \
  -H "Content-Type: application/json" \
  -d "{
    \"ticketId\": $TICKET_ID,
    \"recipient\": \"$USER1\",
    \"rewardBasisPoints\": 500,
    \"metadataURI\": \"https://api.bogowi.com/redeemed/$TICKET_ID\"
  }" | jq '.'

# 9. Verify ticket is burned
echo -e "\n🔍 Verifying ticket is burned (should fail)..."
curl -s "$API_URL/ticket/$TICKET_ID" 2>/dev/null | jq '.'

echo -e "\n✅ Testing complete!"
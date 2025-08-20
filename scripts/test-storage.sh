#!/bin/bash

# Test Google Cloud Storage Setup for BOGOWI NFTs
# Usage: ./scripts/test-storage.sh [local|testnet|mainnet]

NETWORK=${1:-testnet}
echo "🧪 Testing Storage Configuration for: $NETWORK"
echo "================================================"

# Load environment based on network
if [ "$NETWORK" = "mainnet" ]; then
    if [ -f ".env.mainnet" ]; then
        export $(cat .env.mainnet | grep -v '^#' | xargs)
        echo "✅ Loaded mainnet configuration"
    fi
elif [ "$NETWORK" = "local" ]; then
    if [ -f ".env.local" ]; then
        export $(cat .env.local | grep -v '^#' | xargs)
        echo "✅ Loaded local configuration"
    fi
else
    if [ -f ".env.testnet" ]; then
        export $(cat .env.testnet | grep -v '^#' | xargs)
        echo "✅ Loaded testnet configuration"
    fi
fi

# Display configuration
echo ""
echo "📋 Configuration:"
echo "  Network: $NETWORK"
echo "  Bucket: $GCS_BUCKET_NAME"
echo "  Credentials: $GOOGLE_APPLICATION_CREDENTIALS"
echo "  Base URL: $STORAGE_BASE_URL"
echo ""

# Check if service account file exists
if [ ! -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    echo "❌ Service account file not found: $GOOGLE_APPLICATION_CREDENTIALS"
    echo ""
    echo "📝 To fix this:"
    echo "  1. Download the service account JSON from Google Cloud Console"
    echo "  2. Save it as: $GOOGLE_APPLICATION_CREDENTIALS"
    echo ""
    exit 1
fi

echo "✅ Service account file found"

# Test authentication
echo ""
echo "🔐 Testing authentication..."
gcloud auth activate-service-account --key-file="$GOOGLE_APPLICATION_CREDENTIALS" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Authentication successful"
else
    echo "❌ Authentication failed"
    exit 1
fi

# Test bucket access
echo ""
echo "🪣 Testing bucket access..."
gsutil ls -b gs://$GCS_BUCKET_NAME > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Bucket accessible: gs://$GCS_BUCKET_NAME"
else
    echo "❌ Cannot access bucket: gs://$GCS_BUCKET_NAME"
    echo "   Make sure the bucket exists and service account has permissions"
    exit 1
fi

# Test upload
echo ""
echo "📤 Testing file upload..."
TEST_FILE="/tmp/bogowi-test-$(date +%s).txt"
echo "BOGOWI NFT Storage Test - $NETWORK" > $TEST_FILE

gsutil cp $TEST_FILE gs://$GCS_BUCKET_NAME/test/upload-test.txt > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Upload successful"
else
    echo "❌ Upload failed"
    exit 1
fi

# Test public access
echo ""
echo "🌐 Testing public access..."
PUBLIC_URL="$STORAGE_BASE_URL/test/upload-test.txt"
curl -s -o /dev/null -w "%{http_code}" $PUBLIC_URL | grep -q "200"
if [ $? -eq 0 ]; then
    echo "✅ Public access working: $PUBLIC_URL"
else
    echo "⚠️  Public access not configured (this is optional)"
    echo "   To enable: gcloud storage buckets add-iam-policy-binding gs://$GCS_BUCKET_NAME --member=allUsers --role=roles/storage.objectViewer"
fi

# Cleanup
echo ""
echo "🧹 Cleaning up test files..."
gsutil rm gs://$GCS_BUCKET_NAME/test/upload-test.txt > /dev/null 2>&1
rm -f $TEST_FILE

echo ""
echo "✅ Storage configuration test complete!"
echo ""
echo "📊 Summary:"
echo "  Network: $NETWORK"
echo "  Bucket: $GCS_BUCKET_NAME"
echo "  Status: READY"
echo ""
echo "🎯 Next steps:"
echo "  1. Your storage is configured and working!"
echo "  2. Images will be stored at: $STORAGE_BASE_URL/tickets/{tokenId}/original.jpg"
echo "  3. You can now mint NFTs with images"
#!/bin/bash

# Test script for MinIO document upload/download functionality
# Make sure docker-compose is running before executing this script

set -e

BASE_URL="http://localhost:8080/api/v1"
UPLOAD_ENDPOINT="$BASE_URL/documents/upload"
LIST_ENDPOINT="$BASE_URL/documents"

echo "==================================="
echo "MinIO Document Upload/Download Test"
echo "==================================="
echo ""

# Check if MinIO is running
echo "1. Checking MinIO health..."
if curl -f -s http://localhost:9000/minio/health/live > /dev/null; then
    echo "✅ MinIO is running"
else
    echo "❌ MinIO is not running. Start it with: docker-compose up -d minio"
    exit 1
fi

# Check if API is running
echo ""
echo "2. Checking API health..."
if curl -f -s http://localhost:8080/health > /dev/null; then
    echo "✅ API is running"
else
    echo "❌ API is not running. Start it with: go run cmd/api/main.go"
    exit 1
fi

# Create test files
echo ""
echo "3. Creating test files..."
mkdir -p /tmp/insurance-test

# Create a test PDF
echo "This is a test PDF document for insurance policy" > /tmp/insurance-test/policy.txt
echo "✅ Created test policy document"

# Create a test image (simple PNG)
if command -v convert &> /dev/null; then
    convert -size 100x100 xc:white /tmp/insurance-test/id-card.png
    echo "✅ Created test ID card image"
else
    # Create a dummy file if ImageMagick is not available
    echo "PNG" > /tmp/insurance-test/id-card.png
    echo "⚠️  ImageMagick not found, using dummy PNG file"
fi

# Test upload (without auth for now - you'll need to add JWT token later)
echo ""
echo "4. Testing document upload..."
echo "   Uploading policy document..."

RESPONSE=$(curl -s -X POST "$UPLOAD_ENDPOINT" \
  -F "file=@/tmp/insurance-test/policy.txt" \
  -F "document_type=policy_certificate" \
  -F "title=Test Policy Certificate" \
  -F "description=Test document for MinIO integration" \
  -H "Content-Type: multipart/form-data")

echo "   Response: $RESPONSE"

# Extract document ID from response (requires jq)
if command -v jq &> /dev/null; then
    DOC_ID=$(echo "$RESPONSE" | jq -r '.data.id // empty')
    
    if [ -z "$DOC_ID" ]; then
        echo "⚠️  Could not extract document ID. Response: $RESPONSE"
        echo "   This might be due to authentication requirement."
        echo ""
        echo "To test with authentication:"
        echo "1. Login and get JWT token:"
        echo "   curl -X POST $BASE_URL/auth/login -H 'Content-Type: application/json' -d '{\"email\":\"your@email.com\",\"password\":\"yourpass\"}'"
        echo ""
        echo "2. Use the token in upload request:"
        echo "   curl -X POST $UPLOAD_ENDPOINT -H 'Authorization: Bearer YOUR_TOKEN' -F 'file=@/tmp/insurance-test/policy.txt' ..."
    else
        echo "✅ Document uploaded successfully. ID: $DOC_ID"
        
        # Test download
        echo ""
        echo "5. Testing document download..."
        DOWNLOAD_RESPONSE=$(curl -s "$BASE_URL/documents/$DOC_ID/download")
        DOWNLOAD_URL=$(echo "$DOWNLOAD_RESPONSE" | jq -r '.download_url // empty')
        
        if [ -n "$DOWNLOAD_URL" ]; then
            echo "✅ Download URL generated: $DOWNLOAD_URL"
            echo "   URL valid for 1 hour"
            
            # Test actual download
            echo ""
            echo "6. Testing actual file download..."
            if curl -f -s "$DOWNLOAD_URL" -o /tmp/insurance-test/downloaded.txt; then
                echo "✅ File downloaded successfully"
                
                # Compare files
                if diff /tmp/insurance-test/policy.txt /tmp/insurance-test/downloaded.txt > /dev/null; then
                    echo "✅ Downloaded file matches original"
                else
                    echo "⚠️  Downloaded file differs from original"
                fi
            else
                echo "❌ Failed to download file"
            fi
        else
            echo "❌ Failed to generate download URL"
        fi
    fi
else
    echo "⚠️  jq not installed, skipping response parsing"
    echo "   Install jq to enable full testing: apt-get install jq"
fi

echo ""
echo "==================================="
echo "Test Summary"
echo "==================================="
echo ""
echo "MinIO Console: http://localhost:9001"
echo "   Username: minioadmin"
echo "   Password: minioadmin"
echo ""
echo "Check uploaded files in MinIO Console -> Buckets -> insurance-documents"
echo ""
echo "API Endpoints:"
echo "   POST   $UPLOAD_ENDPOINT"
echo "   GET    $LIST_ENDPOINT"
echo "   GET    $BASE_URL/documents/:id/download"
echo "   DELETE $BASE_URL/documents/:id"
echo ""
echo "Supported file types: JPEG, PNG, PDF (max 10MB)"
echo ""

# Cleanup
rm -rf /tmp/insurance-test
echo "✅ Cleaned up test files"

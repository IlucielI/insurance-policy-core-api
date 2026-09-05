#!/bin/bash
# Complete deployment script for semantic search integration

set -e

echo "🚀 Semantic Search Deployment Script"
echo "====================================="
echo ""

# Configuration
export LLM_BASE_URL="http://100.103.220.104:20128/v1"
export EMBEDDINGS_MODEL="bge-m3"

echo "📋 Configuration:"
echo "  LLM_BASE_URL: $LLM_BASE_URL"
echo "  EMBEDDINGS_MODEL: $EMBEDDINGS_MODEL"
echo ""

# Step 1: Build the application
echo "Step 1: Building application..."
cd /home/bayu/Project/insurance-policy-core-api
go build -o insurance-api cmd/api/main.go
echo "✅ Build successful"
echo ""

# Step 2: Apply database migration (if not already applied)
echo "Step 2: Database migration..."
echo "Note: Run this SQL manually if not already done:"
echo "ALTER TABLE product_embeddings ADD CONSTRAINT unique_product_chunk_type UNIQUE (product_id, chunk_type);"
echo ""
read -p "Press Enter after migration is complete..."

# Step 3: Generate embeddings
echo "Step 3: Generating embeddings for products..."
go run scripts/generate_embeddings.go
echo "✅ Embeddings generated"
echo ""

# Step 4: Start the API server (in background)
echo "Step 4: Starting API server..."
./insurance-api &
API_PID=$!
echo "✅ API server started (PID: $API_PID)"
sleep 5
echo ""

# Step 5: Test semantic search
echo "Step 5: Testing semantic search endpoint..."
bash scripts/test_search.sh
echo ""

# Step 6: Cleanup
echo "Step 6: Stopping API server..."
kill $API_PID 2>/dev/null || true
echo "✅ Deployment test complete"
echo ""
echo "🎉 Semantic search integration is ready!"
echo ""
echo "To start the server:"
echo "  ./insurance-api"
echo ""
echo "To test the endpoint:"
echo "  bash scripts/test_search.sh"

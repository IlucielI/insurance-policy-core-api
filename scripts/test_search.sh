#!/bin/bash
# Test semantic search endpoint

API_URL="http://localhost:8080/api/v1/products/search"

echo "🔍 Testing Semantic Search Endpoint"
echo "===================================="
echo ""

# Test query
QUERY="asuransi untuk keluarga"
echo "📤 Sending request..."
echo "Query: $QUERY"
echo "Limit: 5"
echo ""

# Make API call
response=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"query\": \"$QUERY\",
    \"limit\": 5
  }")

# Check if response is valid JSON
if echo "$response" | jq empty 2>/dev/null; then
    echo "📥 Response received:"
    echo "$response" | jq '.'
    
    # Extract count
    count=$(echo "$response" | jq -r '.count // 0')
    echo ""
    echo "✅ Found $count products"
    
    # Show product names
    if [ "$count" -gt 0 ]; then
        echo ""
        echo "📋 Product Names:"
        echo "$response" | jq -r '.results[] | "  - \(.name) (\(.category))"'
    fi
else
    echo "❌ Error: Invalid response"
    echo "$response"
fi

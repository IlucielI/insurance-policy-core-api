#!/bin/bash
# Test Chat API with Real LLM Integration

BASE_URL="http://localhost:8080/api/v1"
SESSION_ID="test-session-$(date +%s)"

echo "=== Testing Chat API with Real LLM ==="
echo "Session ID: $SESSION_ID"
echo ""

# Test 1: Simple greeting
echo "Test 1: Greeting in Indonesian"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"message\": \"Halo, saya ingin tahu tentang produk asuransi kesehatan\"
  }" | jq '.'

echo -e "\n\n"

# Test 2: Product inquiry
echo "Test 2: Premium calculation question"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"message\": \"Bagaimana cara menghitung premi asuransi jiwa?\"
  }" | jq '.'

echo -e "\n\n"

# Test 3: Claim process
echo "Test 3: Claim inquiry"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"$SESSION_ID\",
    \"message\": \"Apa saja dokumen yang diperlukan untuk klaim?\"
  }" | jq '.'

echo -e "\n\nDone! Check responses for AI-generated content (not hardcoded mock)."

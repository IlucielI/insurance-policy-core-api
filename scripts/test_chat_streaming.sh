#!/bin/bash
# Test script for LLM chat streaming integration

set -e

API_URL="${API_URL:-http://localhost:8080}"
LLM_URL="${LLM_URL:-http://100.103.220.104:20128/v1}"
SESSION_ID="test_session_$(date +%s)"

echo "========================================"
echo "LLM Chat Streaming Test"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Check LLM API availability
echo "Test 1: Checking LLM API connectivity..."
if curl -s -f "${LLM_URL}/models" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ LLM API is reachable${NC}"
else
    echo -e "${RED}✗ LLM API is NOT reachable at ${LLM_URL}${NC}"
    echo "  Make sure OmniRoute is running on 100.103.220.104:20128"
fi
echo ""

# Test 2: Test direct LLM streaming
echo "Test 2: Testing direct LLM streaming..."
echo "Request: 'Halo' to LLM API..."
curl -N -X POST "${LLM_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-sonnet-4.5",
    "messages": [
      {"role": "system", "content": "Anda adalah asisten asuransi. Jawab singkat."},
      {"role": "user", "content": "Halo"}
    ],
    "stream": true
  }' 2>/dev/null | head -20
echo -e "\n${YELLOW}(output truncated)${NC}"
echo ""

# Test 3: Test backend chat endpoint (streaming)
echo "Test 3: Testing backend /api/v1/chat (streaming)..."
echo "Session ID: ${SESSION_ID}"
echo "Message: 'Produk asuransi apa saja yang tersedia?'"
echo ""
echo "Response:"
curl -N -X POST "${API_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"${SESSION_ID}\",
    \"message\": \"Produk asuransi apa saja yang tersedia?\",
    \"stream\": true
  }" 2>/dev/null || echo -e "${RED}✗ Backend not reachable. Start with: cd /home/bayu/Project/insurance-policy-core-api && ./bin/api${NC}"
echo ""

# Test 4: Test non-streaming endpoint
echo "Test 4: Testing backend /api/v1/chat (non-streaming)..."
echo "Message: 'Bagaimana cara klaim?'"
echo ""
RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"${SESSION_ID}\",
    \"message\": \"Bagaimana cara klaim?\",
    \"stream\": false
  }" 2>/dev/null)

if [ $? -eq 0 ]; then
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo -e "${RED}✗ Request failed${NC}"
fi
echo ""

# Test 5: Test conversation history
echo "Test 5: Testing conversation with context..."
echo "Follow-up message: 'Berapa lama prosesnya?'"
echo ""
curl -N -X POST "${API_URL}/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": \"${SESSION_ID}\",
    \"message\": \"Berapa lama prosesnya?\",
    \"stream\": true
  }" 2>/dev/null || echo -e "${RED}✗ Backend not reachable${NC}"
echo ""

echo "========================================"
echo "Test Summary"
echo "========================================"
echo -e "${YELLOW}Manual frontend test:${NC}"
echo "1. cd /home/bayu/insurance-policy-app"
echo "2. npm run dev"
echo "3. Open http://localhost:3000"
echo "4. Click floating chat button (bottom-right)"
echo "5. Send message: 'Halo, produk apa saja yang ada?'"
echo "6. Verify streaming response appears token-by-token"
echo ""
echo -e "${GREEN}Integration complete!${NC}"
echo "See LLM_CHAT_INTEGRATION.md for full documentation"

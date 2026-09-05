#!/bin/bash
# Quick test for LLM chat streaming

echo "🧪 Testing LLM API..."
curl -s -X POST http://100.103.220.104:20128/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-sonnet-4.5",
    "messages": [{"role": "user", "content": "Halo, sebutkan 3 jenis asuransi"}],
    "stream": false
  }' | jq -r '.choices[0].message.content' 2>/dev/null || echo "LLM API tidak merespons atau jq tidak tersedia"

echo ""
echo "✅ Jika muncul response di atas, LLM API siap digunakan"
echo ""
echo "📝 Next steps:"
echo "1. Build backend: cd /home/bayu/Project/insurance-policy-core-api && go build -o bin/api cmd/api/main.go"
echo "2. Run backend: ./bin/api"
echo "3. Test frontend: cd /home/bayu/insurance-policy-app && npm run dev"
echo "4. Open http://localhost:3000 dan klik floating chat button"

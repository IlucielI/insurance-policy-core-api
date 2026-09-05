# LLM Integration - OmniRoute Client

## Overview
Chat Assistant now uses **real OmniRoute LLM API** (OpenAI-compatible) instead of MockLLMService.

## Architecture

```
POST /api/v1/chat
    ↓
ChatHandler.SendMessage()
    ↓
ChatUsecase.SendMessage()
    ↓ (RAG: embedding search)
    ↓ (Build context + history)
    ↓
OmniRouteClient.GenerateCompletion()
    ↓
HTTP POST http://100.103.220.104:20128/v1/chat/completions
    ↓
Real AI Response (Claude Sonnet 4.5)
```

## Implementation

### 1. OmniRouteClient (`internal/infrastructure/llm/omniroute_client.go`)
- Implements `LLMService` interface
- Methods:
  - `GenerateCompletion(ctx, messages) -> (string, error)` - Chat completion
  - `GenerateEmbedding(ctx, text) -> ([]float32, error)` - Text embeddings
- OpenAI-compatible API format
- 60s timeout per request
- Full error handling + status code checks

### 2. Main.go Integration (`cmd/api/main.go`)
```go
// Initialize LLM service with OmniRoute
llmBaseURL := os.Getenv("LLM_BASE_URL")
if llmBaseURL == "" {
    llmBaseURL = "http://100.103.220.104:20128/v1"
}

llmModel := os.Getenv("LLM_MODEL")
if llmModel == "" {
    llmModel = "claude-sonnet-4.5"
}

llmService := llm.NewOmniRouteClient(llmBaseURL, llmModel)
chatUsecase := usecase.NewChatUsecase(chatRepo, llmService)
```

### 3. Chat Flow (usecase/chat_usecase.go)
1. User sends message via POST `/api/v1/chat`
2. **RAG Pipeline**:
   - Generate embedding for user query
   - Search top 3 similar product embeddings (pgvector)
   - Build context from retrieved documents
3. **LLM Completion**:
   - System prompt: Professional insurance assistant (Indonesian)
   - Context: Product info from RAG
   - History: Last 5 messages
   - User message
4. **Response**: Real AI-generated reply + context sources

## Environment Variables

```bash
# LLM API Configuration
LLM_BASE_URL=http://100.103.220.104:20128/v1    # OmniRoute endpoint
LLM_MODEL=claude-sonnet-4.5                      # Model name

# Database (for embeddings)
DATABASE_URL=postgres://postgres:***@localhost:5432/insurance_policy?sslmode=disable
```

## Testing

### Manual Test (requires running server)
```bash
# Start server
go run cmd/api/main.go

# Run test script
./test_chat_llm.sh
```

### Expected Response Format
```json
{
  "reply": "Asuransi kesehatan kami menyediakan perlindungan...",
  "sources": [
    {
      "product_id": "prod_123",
      "chunk_type": "description",
      "text": "Asuransi Kesehatan Premium..."
    }
  ]
}
```

## API Endpoint

**POST** `/api/v1/chat`

**Request:**
```json
{
  "session_id": "user-session-123",
  "message": "Bagaimana cara menghitung premi asuransi jiwa?"
}
```

**Response:**
```json
{
  "reply": "Premi asuransi jiwa dihitung berdasarkan beberapa faktor...",
  "sources": [
    {
      "product_id": "prod_life_001",
      "chunk_type": "premium_calculation",
      "text": "Formula perhitungan premi..."
    }
  ]
}
```

## System Prompt (Indonesian)
```
Anda adalah asisten AI untuk sistem asuransi. Gunakan konteks berikut untuk menjawab pertanyaan customer:

[RAG Context from Product Embeddings]

Jawab dengan ramah, jelas, dan fokus pada informasi produk asuransi. 
Jika pertanyaan di luar konteks, arahkan ke customer service.
```

## Removed Files
- `internal/service/mock_llm_service.go` - No longer used (replaced by OmniRouteClient)

## Key Changes
1. **Real LLM responses** - No more hardcoded fallback responses
2. **RAG-enabled** - Retrieves relevant product info via embeddings
3. **Conversation history** - Maintains context across messages
4. **Indonesian language** - System prompt enforces Indonesian responses
5. **Error handling** - Graceful fallback if LLM fails

## Security Notes
- No API key required (internal OmniRoute endpoint)
- 60s timeout prevents hanging requests
- Context length managed by limiting history to 5 messages
- User input sanitized through JSON parsing

## Monitoring
Check logs for:
```
🤖 LLM Config: http://100.103.220.104:20128/v1 (model: claude-sonnet-4.5)
```

Errors logged as:
- "failed to send request" - Network issue
- "API error (status 500)" - OmniRoute endpoint error
- "no choices in response" - Malformed API response

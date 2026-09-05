# Semantic Search Integration Summary

## Overview
Integrated LLM embeddings for semantic search on insurance products using bge-m3 model via OmniRoute.

## Components Implemented

### 1. Infrastructure Layer
- **File**: `internal/infrastructure/llm/embeddings_client.go`
- **Purpose**: Client for generating embeddings via OmniRoute API
- **Key Method**: `GenerateEmbedding(ctx, text) -> []float32`
- **Configuration**: 
  - `LLM_BASE_URL`: http://100.103.220.104:20128/v1
  - `EMBEDDINGS_MODEL`: bge-m3
  - Vector dimension: 1024

### 2. Repository Layer
- **File**: `internal/repository/product_repository.go`
- **New Methods**:
  - `SaveEmbedding()`: Store embeddings in product_embeddings table
  - `SemanticSearch()`: Vector similarity search using pgvector
  - `GetAllProducts()`: Fetch all products for embedding generation

### 3. Use Case Layer
- **File**: `internal/usecase/product_usecase.go`
- **New Methods**:
  - `SemanticSearchProducts()`: Semantic search flow
  - `GenerateProductEmbeddings()`: Generate embeddings for all products

### 4. API Handler Layer
- **File**: `internal/delivery/http/product_handler.go`
- **New Endpoint**: `POST /api/v1/products/search`
- **Request Body**:
  ```json
  {
    "query": "asuransi untuk keluarga",
    "limit": 5
  }
  ```
- **Response**:
  ```json
  {
    "query": "asuransi untuk keluarga",
    "results": [...],
    "count": 5
  }
  ```

### 5. Configuration
- **File**: `config/config.go`
- **New Fields**:
  - `LLMBaseURL`: Base URL for embeddings API
  - `EmbeddingsModel`: Model name (bge-m3)

### 6. Main Application
- **File**: `cmd/api/main.go`
- **Changes**:
  - Initialize embeddings client
  - Wire embeddings service to product usecase
  - Register semantic search route

### 7. Database Migration
- **File**: `migrations/005_add_product_embeddings_unique_constraint.sql`
- **Purpose**: Add unique constraint on (product_id, chunk_type)

### 8. Scripts
- **File**: `scripts/generate_embeddings.go`
- **Purpose**: Go script to generate embeddings for existing products
- **File**: `scripts/generate_embeddings.py`
- **Purpose**: Python script (alternative) to generate embeddings
- **File**: `scripts/test_semantic_search.py`
- **Purpose**: Test semantic search endpoint

## Database Schema

The `product_embeddings` table structure:
```sql
CREATE TABLE product_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    chunk_type VARCHAR(50) NOT NULL,
    chunk_text TEXT NOT NULL,
    embedding vector(1024),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_product_chunk_type UNIQUE (product_id, chunk_type)
);

CREATE INDEX idx_product_embeddings_product_id ON product_embeddings(product_id);
CREATE INDEX idx_product_embeddings_embedding ON product_embeddings USING ivfflat (embedding vector_cosine_ops);
```

## How to Use

### 1. Generate Embeddings for Existing Products
```bash
# Using Go
cd /home/bayu/Project/insurance-policy-core-api
go run scripts/generate_embeddings.go

# OR using Python (if DATABASE_URL is accessible)
python3 scripts/generate_embeddings.py
```

### 2. Test Semantic Search
```bash
# Start the API server
./insurance-api

# Test the endpoint
curl -X POST http://localhost:8080/api/v1/products/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "asuransi untuk keluarga",
    "limit": 5
  }'
```

### 3. Example Search Queries
- "asuransi untuk keluarga" - Family insurance
- "asuransi kesehatan" - Health insurance
- "asuransi jiwa" - Life insurance
- "asuransi kendaraan" - Vehicle insurance
- "proteksi kecelakaan" - Accident protection

## Architecture

```
User Request
    ↓
[POST /products/search] Handler
    ↓
ProductUsecase.SemanticSearchProducts()
    ↓
EmbeddingsClient.GenerateEmbedding(query)
    ↓
ProductRepository.SemanticSearch(embedding)
    ↓
PostgreSQL pgvector similarity search
    ↓
Return top N products
```

## Vector Similarity Search

Uses pgvector's cosine similarity operator:
- Operator: `<=>` (cosine distance)
- Similarity: `1 - (embedding <=> query_embedding)`
- Index: IVFFlat for fast approximate nearest neighbor search
- Default limit: 5 products

## Search Text Generation

For each product:
```
"Kategori: {category}. {name}. {description}"
```

Example:
```
"Kategori: life. Asuransi Jiwa Sejahtera. Perlindungan jiwa untuk masa depan keluarga Anda dengan manfaat asuransi jiwa dan investasi."
```

## Environment Variables

Add to `.env`:
```env
LLM_BASE_URL=http://100.103.220.104:20128/v1
EMBEDDINGS_MODEL=bge-m3
```

## Testing Requirements

1. **Database must have pgvector extension enabled**
2. **Products must exist in database**
3. **Embeddings must be generated for products**
4. **OmniRoute API must be accessible**

## Semantic Search vs Keyword Search

### Keyword Search
- Exact match on text fields
- `WHERE name ILIKE '%keyword%' OR description ILIKE '%keyword%'`
- Fast but less accurate
- Misses semantic relationships

### Semantic Search (This Implementation)
- Vector similarity on meaning
- Understands context and relationships
- More accurate for natural language queries
- Finds products by intent, not just words
- Example: "asuransi untuk keluarga" matches products about family protection even if they don't contain those exact words

## Performance Considerations

- **Embedding generation**: ~100-500ms per product (one-time)
- **Search query**: ~50-200ms (including embedding generation + vector search)
- **Index**: IVFFlat index speeds up searches significantly
- **Cache**: Consider caching frequent query embeddings

## Next Steps

1. Run migration to add unique constraint
2. Generate embeddings for all existing products
3. Test semantic search with various queries
4. Compare accuracy with keyword search
5. Monitor performance and optimize if needed
6. Consider adding:
   - Query embedding cache
   - Hybrid search (semantic + keyword)
   - Re-ranking based on business rules
   - User feedback loop for relevance tuning

## Files Modified

1. `config/config.go` - Added LLM_BASE_URL and EmbeddingsModel
2. `internal/infrastructure/llm/embeddings_client.go` - NEW
3. `internal/repository/product_repository.go` - Added semantic search methods
4. `internal/usecase/product_usecase.go` - Added semantic search logic
5. `internal/delivery/http/product_handler.go` - Added search endpoint
6. `cmd/api/main.go` - Wired up embeddings service
7. `migrations/005_add_product_embeddings_unique_constraint.sql` - NEW
8. `scripts/generate_embeddings.go` - NEW
9. `scripts/generate_embeddings.py` - NEW
10. `.env` - Added LLM_BASE_URL and EMBEDDINGS_MODEL

## Verification Checklist

- [x] Embeddings client created
- [x] Repository methods for semantic search added
- [x] Use case methods implemented
- [x] API endpoint handler created
- [x] Configuration updated
- [x] Main.go wired up
- [x] Migration script created
- [x] Embedding generation scripts created
- [ ] Migration executed (requires database access)
- [ ] Embeddings generated for existing products
- [ ] Semantic search tested with query "asuransi untuk keluarga"
- [ ] Results verified to be more accurate than keyword search

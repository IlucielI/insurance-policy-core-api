# Insurance Policy Core API

Backend API for Insurance Policy Management System built with Go, Fiber, and PostgreSQL.

## Features

- **Product Management**: CRUD operations for insurance products
- **Application Workflow**: Draft → Submitted → Under Review → Approved/Rejected
- **Premium Calculator**: Age-based premium calculation with coverage factors
- **AI RAG Chatbot**: Semantic search over products using pgvector + LLM integration
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: OpenAPI 3.0 specification included

## Tech Stack

- **Framework**: Go 1.26 + Fiber v2
- **Database**: PostgreSQL 16 + pgvector
- **Cache**: Redis
- **Architecture**: Clean Architecture + DDD
- **AI/ML**: pgvector embeddings + LLM API integration

## Quick Start

### Prerequisites

- Go 1.26+
- PostgreSQL 16+
- Redis (optional)

### Local Development

1. **Clone repository**
   ```bash
   git clone https://github.com/IlucielI/insurance-policy-core-api.git
   cd insurance-policy-core-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup database**
   ```bash
   # Create database
   createdb insurance_policy
   
   # Run migrations
   psql -U postgres -d insurance_policy < migrations/001_init_schema.sql
   psql -U postgres -d insurance_policy < migrations/003_seed_correct.sql
   ```

4. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

   Required environment variables:
   ```env
   DATABASE_URL=postgres://postgres:password@localhost:5432/insurance_policy?sslmode=disable
   PORT=8080
   JWT_SECRET=your-secret-key
   REDIS_URL=redis://localhost:6379/0
   LLM_API_URL=http://your-llm-endpoint:20128/v1
   ```

5. **Run server**
   ```bash
   go run cmd/api/main.go
   ```

   Server will start on `http://localhost:8080`

### Docker Deployment

```bash
# Build image
docker build -t insurance-api .

# Run container
docker run -d \
  --name insurance-api \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  insurance-api
```

## API Documentation

### OpenAPI Specification

Full API documentation available in [docs/openapi.yaml](docs/openapi.yaml)

**Interactive API docs:**
- Swagger UI: Coming soon
- ReDoc: Coming soon

### Quick Reference

**Base URL**: `http://localhost:8080/api/v1`

**Health Check**
```bash
curl http://localhost:8080/health
```

**List Products**
```bash
curl http://localhost:8080/api/v1/products
```

**Get Product by Slug**
```bash
curl http://localhost:8080/api/v1/products/asuransi-jiwa-premium
```

**Calculate Premium**
```bash
curl -X POST http://localhost:8080/api/v1/products/asuransi-jiwa-premium/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "age": 35,
    "sum_assured": 500000000,
    "payment_term": 120
  }'
```

**Chat with AI**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Produk asuransi apa saja yang tersedia?",
    "session_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

## Project Structure

```
insurance-policy-core-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Domain entities
│   ├── usecase/                    # Business logic
│   ├── repository/                 # Data access layer
│   ├── delivery/
│   │   └── http/                   # HTTP handlers (Fiber)
│   ├── infrastructure/
│   │   ├── database/               # PostgreSQL + pgvector
│   │   ├── cache/                  # Redis
│   │   └── llm/                    # LLM client
│   └── middleware/                 # Auth, CORS, logging
├── migrations/                     # SQL migrations
├── docs/
│   └── openapi.yaml               # OpenAPI 3.0 spec
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Database Schema

### Products
- Stores insurance products with pricing rules
- Slugs for SEO-friendly URLs
- Age factors for premium calculation
- Category-based filtering (life, health, vehicle)

### Applications
- Policy application workflow
- Status: draft → submitted → under_review → approved/rejected
- JSONB fields for flexible applicant data
- Underwriter notes and rejection reasons

### Product Embeddings (pgvector)
- Vector embeddings for semantic search
- Enables AI-powered product recommendations
- RAG (Retrieval Augmented Generation) for chatbot

## AI Features

### RAG Chatbot

The chatbot uses **Retrieval Augmented Generation (RAG)** with pgvector:

1. **Semantic Search**: User query → embedding → pgvector similarity search
2. **Context Retrieval**: Top-k relevant products + chat history
3. **LLM Generation**: Context + query → LLM → natural language response

**Tech Stack:**
- pgvector extension for embeddings storage
- Cosine similarity for semantic search
- LLM API for response generation
- Redis for session management

**Example Flow:**
```
User: "Asuransi untuk usia 40 tahun dengan budget 500rb/bulan?"
  ↓
Embedding → pgvector search
  ↓
Context: [Product A, Product B]
  ↓
LLM prompt: "Based on these products... recommend..."
  ↓
Response: "Saya merekomendasikan Asuransi Jiwa Premium..."
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/usecase/...
```

## Production Deployment

**Live API**: https://insurance-app-api.bayuanugerah.my.id

Deployed on Oracle Cloud ARM64 VM with:
- Docker + NGINX reverse proxy
- Let's Encrypt SSL certificates
- PostgreSQL 16 + pgvector
- Auto-restart on failure

## Development Workflow

1. Create feature branch: `git checkout -b feature/new-feature`
2. Make changes with atomic commits (Angular convention)
3. Push and create PR: `git push origin feature/new-feature`
4. Code review (automated + manual)
5. Merge to main
6. Auto-deploy to production

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes (atomic commits, Angular convention)
4. Push to branch
5. Create Pull Request

## License

MIT License - see LICENSE file for details

## Contact

**Bayu Anugerah**
- Email: bayu.anugerah99@gmail.com
- GitHub: [@IlucielI](https://github.com/IlucielI)
- Portfolio: https://bayuanugerah.my.id

## Acknowledgments

Built as part of technical assessment for PT Logika Sarana Teknologi (Logique).

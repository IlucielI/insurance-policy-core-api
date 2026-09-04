# Insurance Policy Core API

Backend API untuk Insurance Policy Application System, dibangun menggunakan Go + Fiber dengan Clean Architecture.

## 🏗️ Tech Stack

- **Go 1.26+** (Fiber v2 web framework)
- **PostgreSQL 16** dengan pgvector extension
- **Redis** untuk session storage
- **Clean Architecture + DDD** (Domain-Driven Design)

## 📁 Project Structure

```
cmd/api/              # Application entry point
internal/
  ├── domain/         # Business entities (User, Product, Application)
  ├── usecase/        # Business logic layer
  ├── repository/     # Data access layer (pgx driver)
  └── delivery/http/  # HTTP handlers (Fiber routes)
config/               # Environment configuration
migrations/           # Database migrations (pgvector setup)
```

## 🚀 Quick Start

### Prerequisites

- Go 1.26 or higher
- PostgreSQL 16 dengan pgvector extension
- Redis server

### 1. Clone Repository

```bash
git clone https://github.com/IlucielI/insurance-policy-core-api.git
cd insurance-policy-core-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment Variables

Copy `.env.example` ke `.env`:

```bash
cp .env.example .env
```

Edit `.env` file:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/insurance_policy?sslmode=disable
REDIS_URL=redis://localhost:6379/0
SESSION_SECRET=your-secret-key-here
LLM_ENDPOINT=http://localhost:11434/v1  # Ollama atau OpenAI-compatible endpoint
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

### 4. Setup Database

**Install pgvector extension:**

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

**Run migrations:**

```bash
psql $DATABASE_URL -f migrations/001_init_schema.sql
```

Atau manual:

```sql
-- Copy isi migrations/001_init_schema.sql dan execute
```

### 5. Run Application

**Development mode:**

```bash
go run cmd/api/main.go
```

**Build & run:**

```bash
go build -o insurance-api cmd/api/main.go
./insurance-api
```

Server akan jalan di `http://localhost:8080`

## 📡 API Endpoints

### Health Check
```bash
GET /health
```

### Products
```bash
GET /api/v1/products                    # List all products
GET /api/v1/products/:id                # Get product detail
POST /api/v1/products/:id/calculate-premium  # Calculate premium
```

### Applications
```bash
POST /api/v1/applications               # Submit new application
GET /api/v1/applications/:id            # Get application detail
```

### Chat (RAG Chatbot)
```bash
POST /api/v1/chat                       # Send message to AI assistant
GET /api/v1/chat/:sessionId/history     # Get chat history
```

**Example request:**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "demo-session",
    "message": "Apa saja produk asuransi yang tersedia?"
  }'
```

### Admin Endpoints
```bash
PUT /api/v1/admin/applications/:id/status   # Update application status
GET /api/v1/admin/applications              # List all applications
POST /api/v1/admin/products                 # Create product (admin)
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...
```

## 🐳 Docker

**Build image:**

```bash
docker build -t insurance-api .
```

**Run container:**

```bash
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/db \
  -e REDIS_URL=redis://host:6379/0 \
  --name insurance-api \
  insurance-api
```

## 🏛️ Architecture

**Clean Architecture layers:**

1. **Domain Layer** (`internal/domain/`): Business entities
   - Pure Go structs, no external dependencies
   - Entities: User, Product, Application, Payment, ChatSession, Message

2. **Use Case Layer** (`internal/usecase/`): Business logic
   - Orchestrates data flow between repositories
   - Implements business rules (premium calculation, RAG chatbot)

3. **Repository Layer** (`internal/repository/`): Data access
   - PostgreSQL queries via pgx driver
   - pgvector similarity search

4. **Delivery Layer** (`internal/delivery/http/`): HTTP handlers
   - Fiber routes
   - Request validation & response formatting

## 🤖 AI Features

### RAG Chatbot (Retrieval-Augmented Generation)

**How it works:**

1. User message → embedding (via LLM endpoint)
2. pgvector searches similar past messages (cosine similarity)
3. Top 3 similar messages retrieved as context
4. Context + user message → LLM
5. LLM response → stored with embedding

**Fallback mode:**

Jika LLM endpoint down, chatbot return static helpful responses:

```go
// internal/usecase/chat_usecase.go
fallbackResponses := map[string]string{
    "products": "Kami menyediakan 3 jenis asuransi: Jiwa, Kesehatan, dan Kendaraan.",
    "premium": "Premi dihitung berdasarkan usia, uang pertanggungan, dan jenis produk.",
    // ...
}
```

## 🔧 Configuration

**Environment variables:**

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/db` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `SESSION_SECRET` | Session encryption key | `your-secret-key` |
| `LLM_ENDPOINT` | OpenAI-compatible API endpoint | `http://localhost:11434/v1` |
| `CORS_ORIGINS` | Allowed origins (comma-separated) | `http://localhost:3000` |

## 📊 Database Schema

**Key tables:**

- `users`: User accounts (email, password_hash, role)
- `products`: Insurance products (name, category, base_premium, embedding)
- `applications`: Policy applications (applicant data, status, premium_amount)
- `chat_sessions`: Chat sessions per user
- `messages`: Chat messages with embeddings (vector similarity search)

**pgvector indexes:**

```sql
CREATE INDEX ON products USING ivfflat (embedding vector_cosine_ops);
CREATE INDEX ON messages USING ivfflat (embedding vector_cosine_ops);
```

## 🚨 Troubleshooting

**Problem: CORS errors**

Solution: Check `CORS_ORIGINS` env var matches frontend URL

**Problem: Database connection failed**

Solution: Verify PostgreSQL running and credentials correct

**Problem: pgvector extension not found**

Solution: Install pgvector:

```bash
# Ubuntu/Debian
sudo apt install postgresql-16-pgvector

# Or compile from source
git clone https://github.com/pgvector/pgvector.git
cd pgvector
make
sudo make install
```

**Problem: LLM endpoint timeout**

Solution: App works with fallback responses. LLM optional for demo.

## 📝 Development Workflow

1. Create feature branch: `git checkout -b feat/feature-name`
2. Make changes, test locally
3. Commit (atomic): `git commit -m "feat(scope): description"`
4. Push: `git push -u origin feat/feature-name`
5. Create PR, wait for code review
6. Merge after approval

## 📄 License

MIT License - Bayu Anugerah

## 🔗 Related Repositories

- Frontend App: https://github.com/IlucielI/insurance-policy-app
- Admin CMS: https://github.com/IlucielI/insurance-policy-cms

## 📧 Contact

**Bayu Anugerah**  
Email: bayu.anugerah99@gmail.com  
GitHub: https://github.com/IlucielI

# Insurance Policy Core API

Backend API untuk Insurance Policy System (Technical Test PT Logika Sarana Teknologi).

## Tech Stack

- **Go 1.23+**
- **Fiber v2** (web framework)
- **PostgreSQL + pgvector** (database + vector search)
- **Redis** (session store)
- **golang-migrate** (migrations)

## Architecture

Clean Architecture + DDD:

```
cmd/api/          # Entry point
internal/
  domain/         # Entities (User, Product, Application)
  usecase/        # Business logic
  repository/     # Data access
  delivery/http/  # HTTP handlers
  middleware/     # Auth, logging
config/           # Config loaders
migrations/       # SQL migrations
```

## Setup

```bash
# Copy env
cp .env.example .env

# Edit .env dengan credentials real

# Run migrations
psql $DATABASE_URL -f migrations/001_init_schema.sql

# Install deps
go mod download

# Run
go run cmd/api/main.go
```

## Endpoints

### Auth
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login (session-based)
- `POST /api/v1/auth/logout` - Logout

### Products (Public)
- `GET /api/v1/products` - List products (filter: category, active)
- `GET /api/v1/products/:id` - Get product detail

### Applications (Customer)
- `POST /api/v1/applications` - Create application
- `GET /api/v1/applications/:id` - Get application detail

### Chat (RAG)
- `POST /api/v1/chat` - Send message to AI assistant
  - Body: `{ "session_id": "uuid", "message": "..." }`
  - Response: `{ "reply": "...", "sources": [...] }`

### Admin (Protected)
- `GET /api/v1/admin/applications` - List all applications
- `PATCH /api/v1/admin/applications/:id/status` - Update status (approve/reject)
- `POST /api/v1/admin/products` - Create product
- `PUT /api/v1/admin/products/:id` - Update product

## Database Schema

See `migrations/001_init_schema.sql`:

- **users** - Customer + admin accounts
- **products** - Insurance products (life, health, vehicle)
- **applications** - Policy applications (draft → submitted → under review → approved/rejected)
- **payments** - Midtrans payment records
- **chat_sessions** - RAG chatbot sessions
- **chat_messages** - Chat history
- **product_embeddings** - Vector embeddings for RAG (pgvector)
- **activity_logs** - Audit trail

## Development

```bash
# Run with hot reload (air)
air

# Run tests
go test ./...

# Build
go build -o bin/api cmd/api/main.go
```

## Deployment

Production: Oracle VM 16gb-bayu (100.70.163.113)

```bash
# Build binary
GOOS=linux GOARCH=amd64 go build -o insurance-api cmd/api/main.go

# Deploy via SCP
scp insurance-api 16gb-bayu:/home/ubuntu/apps/

# Systemd service (TBD)
```

## Environment Variables

| Var | Example | Description |
|-----|---------|-------------|
| `PORT` | `8080` | HTTP port |
| `DATABASE_URL` | `postgres://...` | PostgreSQL connection string |
| `REDIS_URL` | `redis://...` | Redis connection string |
| `SESSION_SECRET` | `random-string` | Session encryption key |
| `LLM_ENDPOINT` | `http://...` | LLM API endpoint (OmniRoute) |
| `MIDTRANS_SERVER_KEY` | `SB-...` | Midtrans server key |
| `RECAPTCHA_SECRET` | `6Le...` | Google reCAPTCHA secret |
| `CORS_ORIGINS` | `https://...` | Allowed CORS origins (comma-separated) |

## License

Proprietary - Technical Test Project

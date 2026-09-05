#!/bin/bash
# Quick Setup Script for MinIO Integration

echo "=========================================="
echo "MinIO Integration - Quick Setup"
echo "=========================================="
echo ""

# Step 1: Start Docker services
echo "Step 1: Starting Docker services (PostgreSQL, Redis, MinIO)..."
docker-compose up -d postgres redis minio

# Wait for services to be healthy
echo "Waiting for services to be ready..."
sleep 10

# Check MinIO
if curl -f -s http://localhost:9000/minio/health/live > /dev/null; then
    echo "✅ MinIO is running on port 9000"
    echo "   Console: http://localhost:9001 (minioadmin/minioadmin)"
else
    echo "❌ MinIO failed to start"
    exit 1
fi

# Check PostgreSQL
if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "✅ PostgreSQL is running on port 5432"
else
    echo "❌ PostgreSQL failed to start"
    exit 1
fi

# Check Redis
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is running on port 6379"
else
    echo "❌ Redis failed to start"
    exit 1
fi

echo ""
echo "Step 2: Environment configuration..."
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ Created .env file"
    echo "⚠️  Please update .env with your actual credentials if needed"
else
    echo "✅ .env file already exists"
fi

echo ""
echo "Step 3: Go dependencies..."
if command -v go &> /dev/null; then
    echo "Running go mod tidy..."
    go mod tidy
    go mod download
    echo "✅ Go dependencies installed"
else
    echo "⚠️  Go not found. Please install Go and run: go mod tidy"
fi

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Start the API:"
echo "   go run cmd/api/main.go"
echo ""
echo "2. Test the integration:"
echo "   ./test-minio.sh"
echo ""
echo "3. Access MinIO Console:"
echo "   http://localhost:9001"
echo "   Username: minioadmin"
echo "   Password: minioadmin"
echo ""
echo "4. API Endpoints:"
echo "   POST   http://localhost:8080/api/v1/documents/upload"
echo "   GET    http://localhost:8080/api/v1/documents"
echo "   GET    http://localhost:8080/api/v1/documents/:id/download"
echo "   DELETE http://localhost:8080/api/v1/documents/:id"
echo ""
echo "Documentation: docs/MINIO_INTEGRATION.md"
echo ""

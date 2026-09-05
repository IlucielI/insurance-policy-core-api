#!/bin/bash
# Deploy insurance API with Redis caching to 16gb-bayu VM

set -e

VM_HOST="16gb-bayu-oracle"
REDIS_HOST="150.230.61.73"  # 6gb-bayu VM
APP_DIR="/home/ubuntu/insurance-api"

echo "=== Deploying Insurance API with Redis Cache ==="

# Test Redis connection first
echo "1. Testing Redis connection..."
redis-cli -h $REDIS_HOST ping || { echo "❌ Redis not reachable"; exit 1; }
echo "✅ Redis connected"

# Build and push Docker image
echo "2. Building Docker image..."
docker build -t nbsbayu/insurance-api:redis-cache .
docker push nbsbayu/insurance-api:redis-cache

# Deploy to VM
echo "3. Deploying to VM..."
ssh $VM_HOST << 'REMOTE_SCRIPT'
    cd /home/ubuntu/insurance-api
    
    # Pull new image
    docker pull nbsbayu/insurance-api:redis-cache
    
    # Stop old container
    docker stop insurance-api 2>/dev/null || true
    docker rm insurance-api 2>/dev/null || true
    
    # Start new container with Redis
    docker run -d \
        --name insurance-api \
        --restart unless-stopped \
        -p 8080:8080 \
        -e REDIS_URL=redis://150.230.61.73:6379 \
        -e DATABASE_URL="${DATABASE_URL}" \
        -e LLM_BASE_URL=http://100.103.220.104:20128/v1 \
        nbsbayu/insurance-api:redis-cache
    
    # Wait for health check
    sleep 5
    docker logs insurance-api --tail 20
REMOTE_SCRIPT

echo "✅ Deployment complete!"
echo "Test at: http://161.33.39.119:8080/api/v1/products"

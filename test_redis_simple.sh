#!/bin/bash
# Simple Redis cache test via docker exec

REDIS_CONTAINER="redis-redis-1"
REDIS_VM="6gb-bayu-oracle"

echo "=== Testing Redis Cache ==="

# Test connection
echo "1. Testing Redis..."
ssh $REDIS_VM "docker exec $REDIS_CONTAINER redis-cli ping"

# Show current keys
echo "2. Current cache keys:"
ssh $REDIS_VM "docker exec $REDIS_CONTAINER redis-cli KEYS '*'"

# Show stats
echo "3. Cache stats:"
ssh $REDIS_VM "docker exec $REDIS_CONTAINER redis-cli INFO stats | grep keyspace"

echo "✅ Redis ready for caching!"

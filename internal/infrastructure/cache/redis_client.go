package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	
	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	val, err := r.client.Get(ctxTimeout, key).Result()
	if err == redis.Nil {
		return "", nil // Cache miss
	}
	return val, err
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	return r.client.Set(ctxTimeout, key, value, expiration).Err()
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	return r.client.Del(ctxTimeout, keys...).Err()
}

// Ping checks Redis connectivity
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// PoolStats returns Redis connection pool statistics
func (r *RedisClient) PoolStats() *redis.PoolStats {
	return r.client.PoolStats()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

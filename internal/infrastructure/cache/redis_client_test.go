package cache

import (
	"context"
	"testing"
	"time"
)

// Manual test - requires Redis running on localhost:6379
func TestRedisCache(t *testing.T) {
	t.Skip("Manual test - run with Redis container")
	
	client, err := NewRedisClient("redis://localhost:6379/0")
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	
	// Test Set and Get
	key := "test:key"
	value := "test value"
	
	err = client.Set(ctx, key, value, 10*time.Second)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	
	result, err := client.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	
	if result != value {
		t.Errorf("Expected %s, got %s", value, result)
	}
	
	// Test Delete
	err = client.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	
	// Verify deletion
	result, err = client.Get(ctx, key)
	if err != nil {
		t.Errorf("Get after delete failed: %v", err)
	}
	
	if result != "" {
		t.Errorf("Expected empty string after delete, got %s", result)
	}
	
	t.Log("✅ Redis cache test passed")
}

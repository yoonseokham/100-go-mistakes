package redis

import (
	"testing"
)

// Testing with init() is problematic
// The init() function runs before tests, we cannot control it
func TestStore(t *testing.T) {
	// This test depends on init() having run successfully
	err := Store("test", "value")
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	got, err := Get("test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != "value" {
		t.Errorf("Get() = %v, want %v", got, "value")
	}
}

// Testing with explicit initialization is much better
func TestNewCache(t *testing.T) {
	// We control when initialization happens
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	// Test with fresh, isolated state
	err = cache.Store("test", "value")
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	got, err := cache.Get("test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if got != "value" {
		t.Errorf("Get() = %v, want %v", got, "value")
	}
}

func TestCache_MultipleInstances(t *testing.T) {
	// With explicit initialization, we can create multiple independent instances
	cache1, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	cache2, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	// Each cache is independent
	cache1.Store("key", "value1")
	cache2.Store("key", "value2")

	val1, _ := cache1.Get("key")
	val2, _ := cache2.Get("key")

	if val1 != "value1" {
		t.Errorf("cache1 Get() = %v, want value1", val1)
	}

	if val2 != "value2" {
		t.Errorf("cache2 Get() = %v, want value2", val2)
	}
}

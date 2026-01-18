package redis

import (
	"errors"
	"fmt"
)

// Global variable set by init function
var cache map[string]string

// Bad: init function with side effects
// Problem: Runs automatically on package import, cannot control or test
func init() {
	fmt.Println("Initializing redis cache...")
	cache = make(map[string]string)
	// Imagine this connects to actual Redis
	// If connection fails, we can only panic or ignore the error
}

// Store saves a key-value pair
// Problem: Relies on init() having run successfully
func Store(key, value string) error {
	if cache == nil {
		return errors.New("cache not initialized")
	}
	cache[key] = value
	fmt.Printf("Stored: %s = %s\n", key, value)
	return nil
}

// Get retrieves a value by key
func Get(key string) (string, error) {
	if cache == nil {
		return "", errors.New("cache not initialized")
	}
	value, exists := cache[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

// Better approach: Explicit initialization with error handling
type Cache struct {
	data map[string]string
}

func NewCache() (*Cache, error) {
	// Can return error if initialization fails
	fmt.Println("Creating new cache instance...")
	return &Cache{
		data: make(map[string]string),
	}, nil
}

func (c *Cache) Store(key, value string) error {
	c.data[key] = value
	fmt.Printf("Stored in cache: %s = %s\n", key, value)
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	value, exists := c.data[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

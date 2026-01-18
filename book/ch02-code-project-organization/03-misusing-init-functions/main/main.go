package main

import (
	"fmt"
	"log"

	"github.com/yourusername/go-100-mistakes/book/ch02-code-project-organization/03-misusing-init-functions/redis"
)

// Multiple init functions in same file execute in order
func init() {
	fmt.Println("init 1")
}

func init() {
	fmt.Println("init 2")
}

func main() {
	fmt.Println("Starting main...")

	// Bad approach: Using global redis functions that rely on init()
	fmt.Println("\n--- Bad approach: Using init() ---")
	err := redis.Store("foo", "bar")
	if err != nil {
		log.Printf("Error storing: %v", err)
	}

	value, err := redis.Get("foo")
	if err != nil {
		log.Printf("Error getting: %v", err)
	} else {
		fmt.Printf("Retrieved: %s\n", value)
	}

	// Good approach: Explicit initialization
	fmt.Println("\n--- Good approach: Explicit initialization ---")
	cache, err := redis.NewCache()
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}

	if err := cache.Store("hello", "world"); err != nil {
		log.Fatalf("Failed to store: %v", err)
	}

	value, err = cache.Get("hello")
	if err != nil {
		log.Fatalf("Failed to get: %v", err)
	}
	fmt.Printf("Retrieved from cache: %s\n", value)
}

package misuinginit

import (
	"database/sql"
	"errors"
	"fmt"
)

// Bad Example 1: init() with side effects and error handling issues
// Problem: Cannot return errors, hard to test, order-dependent
var db *sql.DB

// This demonstrates why init() is problematic:
// Uncommenting this will cause panic on package import!
// Real code would look like this, but it breaks on import:
//
// func init() {
//     var err error
//     db, err = sql.Open("postgres", "connection-string")
//     if err != nil {
//         panic(err)  // Only option: panic or ignore error
//     }
// }

// Bad Example 2: init() doing too much work
var globalCache map[string]string

func init() {
	// Problem: Heavy initialization in init
	// - Increases startup time
	// - Cannot be controlled or mocked in tests
	// - No way to inject dependencies
	globalCache = make(map[string]string)
	loadCacheFromFile() // Expensive operation
}

func loadCacheFromFile() {
	// Simulate expensive file I/O
	fmt.Println("Loading cache from file...")
}

// Bad Example 3: Order-dependent initialization
var (
	config Config
	client *Client
)

func init() {
	// Problem: Relies on another init being called first
	// Order is not guaranteed across files
	config = loadConfig()
}

func init() {
	// This might run before the config init!
	// Leads to bugs that are hard to track
	client = NewClient(config)
}

// Good Example 1: Explicit initialization with error handling
type Database struct {
	conn *sql.DB
}

// Constructor that can return errors
func NewDatabase(connectionString string) (*Database, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Database{conn: conn}, nil
}

// Good Example 2: Lazy initialization
type Cache struct {
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Load() error {
	// Load when needed, not at init time
	// Can return error, can be tested
	return nil
}

// Good Example 3: Dependency injection
type Client struct {
	config Config
}

type Config struct {
	Host string
	Port int
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}

// When init() IS appropriate:
// - Registering with registries (database drivers, encoding types)
// - Setting up truly constant data
// - No error handling needed

var supportedFormats map[string]bool

func init() {
	// This is OK: simple, deterministic, no errors
	supportedFormats = map[string]bool{
		"json": true,
		"xml":  true,
		"yaml": true,
	}
}

// Example of proper usage pattern
type Service struct {
	db     *Database
	cache  *Cache
	client *Client
}

// Explicit initialization with error handling
func NewService(dbConnStr string, config Config) (*Service, error) {
	db, err := NewDatabase(dbConnStr)
	if err != nil {
		return nil, err
	}

	cache := NewCache()
	if err := cache.Load(); err != nil {
		return nil, err
	}

	client := NewClient(config)

	return &Service{
		db:     db,
		cache:  cache,
		client: client,
	}, nil
}

// Helper to demonstrate why init is problematic
func GetDB() (*sql.DB, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}
	return db, nil
}

func GetCache() map[string]string {
	return globalCache
}

func GetClient() *Client {
	return client
}

func loadConfig() Config {
	return Config{Host: "localhost", Port: 8080}
}

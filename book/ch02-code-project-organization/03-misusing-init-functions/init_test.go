package misuinginit

import (
	"testing"
)

// Testing code that uses init() is difficult
func TestGetDB(t *testing.T) {
	// We cannot control what init() did
	// We cannot mock the database
	// We cannot test different scenarios easily
	_, err := GetDB()
	// This test depends on init() succeeding
	if err != nil {
		t.Skip("Skipping test because init() may have failed")
	}
}

// Testing explicit initialization is easy
func TestNewDatabase(t *testing.T) {
	// We can test with different connection strings
	// We can control when initialization happens
	// We can properly handle and test errors

	tests := []struct {
		name    string
		connStr string
		wantErr bool
	}{
		{"empty connection", "", true},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatabase(tt.connStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCache(t *testing.T) {
	// Easy to test: no global state, no init() dependency
	cache := NewCache()
	if cache == nil {
		t.Fatal("expected cache to be non-nil")
	}

	if cache.data == nil {
		t.Error("expected cache.data to be initialized")
	}
}

func TestNewClient(t *testing.T) {
	config := Config{
		Host: "testhost",
		Port: 9090,
	}

	client := NewClient(config)
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	if client.config.Host != "testhost" {
		t.Errorf("expected host testhost, got %s", client.config.Host)
	}
}

func TestNewService(t *testing.T) {
	// With explicit initialization, we can:
	// 1. Test with different configurations
	// 2. Use mock databases
	// 3. Control the initialization order
	// 4. Handle errors properly

	config := Config{
		Host: "localhost",
		Port: 8080,
	}

	// This will fail because we're using real sql.Open
	// In real tests, you'd use dependency injection or interfaces
	_, err := NewService("invalid-connection", config)
	if err == nil {
		t.Error("expected error with invalid connection string")
	}
}

func TestSupportedFormats(t *testing.T) {
	// This is OK to use init() for: simple, deterministic data
	tests := []struct {
		format    string
		supported bool
	}{
		{"json", true},
		{"xml", true},
		{"yaml", true},
		{"csv", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := supportedFormats[tt.format]
			if got != tt.supported {
				t.Errorf("supportedFormats[%s] = %v, want %v", tt.format, got, tt.supported)
			}
		})
	}
}

// Benchmark showing init() impact on startup time
func BenchmarkInitOverhead(b *testing.B) {
	// init() runs before all tests, adding to startup time
	// With explicit initialization, you only pay when you need it
	for i := 0; i < b.N; i++ {
		_ = NewCache()
	}
}

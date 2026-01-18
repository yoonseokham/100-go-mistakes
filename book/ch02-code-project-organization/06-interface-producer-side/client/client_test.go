package client

import (
	"testing"

	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/06-interface-producer-side/store"
)

// Mock implementation - only implements what Client needs
type mockCustomersGetter struct {
	customers []store.Customer
}

func (m *mockCustomersGetter) GetAllCustomers() ([]store.Customer, error) {
	return m.customers, nil
}

// Easy to test - only mock 1 method!
func TestClient_ListAllCustomers(t *testing.T) {
	mock := &mockCustomersGetter{
		customers: []store.Customer{
			{ID: "1", Name: "Alice"},
			{ID: "2", Name: "Bob"},
		},
	}

	client := NewClient(mock)
	customers, err := client.ListAllCustomers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}
}

// Test with real Store
func TestClient_WithRealStore(t *testing.T) {
	customerStore := store.NewStore()
	customerStore.StoreCustomer(store.Customer{ID: "1", Name: "Alice"})
	customerStore.StoreCustomer(store.Customer{ID: "2", Name: "Bob"})

	client := NewClient(customerStore)
	customers, err := client.ListAllCustomers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}
}

// Mock for Reporter - different interface!
type mockBalanceReader struct {
	customers []store.Customer
}

func (m *mockBalanceReader) GetCustomersWithNegativeBalance() ([]store.Customer, error) {
	return m.customers, nil
}

func TestReporter_GenerateNegativeBalanceReport(t *testing.T) {
	mock := &mockBalanceReader{
		customers: []store.Customer{
			{ID: "1", Name: "Alice", Balance: -100},
		},
	}

	reporter := NewReporter(mock)
	customers, err := reporter.GenerateNegativeBalanceReport()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(customers) != 1 {
		t.Errorf("expected 1 customer, got %d", len(customers))
	}
}

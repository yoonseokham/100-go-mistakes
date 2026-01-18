package interfacepollution

import (
	"database/sql"
	"fmt"
)

// Decoupling Example: Using Interfaces for Testing

type Customer struct {
	ID   string
	Name string
}

// ============================================
// BAD EXAMPLE: Tight Coupling
// ============================================

// Bad: Service depends on concrete implementation
type BadCustomerService struct {
	store *MySQLStore // Tightly coupled to MySQL
}

type MySQLStore struct {
	db *sql.DB
}

func (m *MySQLStore) StoreCustomer(customer Customer) error {
	if m.db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := m.db.Exec("INSERT INTO customers (id, name) VALUES (?, ?)", customer.ID, customer.Name)
	return err
}

func NewBadCustomerService(store *MySQLStore) *BadCustomerService {
	return &BadCustomerService{store: store}
}

func (cs *BadCustomerService) CreateNewCustomer(id, name string) error {
	customer := Customer{ID: id, Name: name}
	return cs.store.StoreCustomer(customer)
}

// Problems with BadCustomerService:
// 1. Cannot easily test without real MySQL database
// 2. Cannot swap implementation (e.g., PostgreSQL, MongoDB)
// 3. Tightly coupled to MySQLStore
// 4. Hard to mock for unit tests

// ============================================
// GOOD EXAMPLE: Decoupled with Interface
// ============================================

// Step 1: Define interface in consumer package
type CustomerStorer interface {
	StoreCustomer(customer Customer) error
}

// Step 2: Service depends on interface
type GoodCustomerService struct {
	storer CustomerStorer // Depends on abstraction
}

func NewGoodCustomerService(storer CustomerStorer) *GoodCustomerService {
	return &GoodCustomerService{storer: storer}
}

func (cs *GoodCustomerService) CreateNewCustomer(id, name string) error {
	customer := Customer{ID: id, Name: name}
	return cs.storer.StoreCustomer(customer)
}

// Step 3: Concrete implementations

// MySQLStore already implements CustomerStorer (has StoreCustomer method)

// PostgreSQL implementation
type PostgresStore struct {
	db *sql.DB
}

func (p *PostgresStore) StoreCustomer(customer Customer) error {
	if p.db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := p.db.Exec("INSERT INTO customers (id, name) VALUES ($1, $2)", customer.ID, customer.Name)
	return err
}

// Test implementation
type MockStore struct {
	StoreFn func(Customer) error
}

func (m *MockStore) StoreCustomer(customer Customer) error {
	if m.StoreFn != nil {
		return m.StoreFn(customer)
	}
	return nil
}

// Benefits of GoodCustomerService:
// 1. Easy to test with MockStore
// 2. Can swap implementations (MySQL, Postgres, etc.)
// 3. Decoupled from concrete implementation
// 4. Follows Dependency Inversion Principle

// Usage examples:
// Production with MySQL:  NewGoodCustomerService(&MySQLStore{db})
// Production with Postgres: NewGoodCustomerService(&PostgresStore{db})
// Testing: NewGoodCustomerService(&MockStore{...})

package interfacepollution

// Decoupling Example: Using Interfaces for Testing

// Step 1: Define interface in consumer package
type CustomerStorer interface {
	StoreCustomer(customer Customer) error
}

type Customer struct {
	ID   string
	Name string
}

// Step 2: Service depends on interface
type CustomerService struct {
	storer CustomerStorer
}

func NewCustomerService(storer CustomerStorer) *CustomerService {
	return &CustomerService{storer: storer}
}

func (cs *CustomerService) CreateNewCustomer(id, name string) error {
	customer := Customer{ID: id, Name: name}
	return cs.storer.StoreCustomer(customer)
}

// Step 3: Concrete implementations

// Real implementation
type MySQLStore struct {
	// connection details
}

func (m *MySQLStore) StoreCustomer(customer Customer) error {
	// Store in MySQL
	return nil
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

// Usage:
// Production: NewCustomerService(&MySQLStore{})
// Testing: NewCustomerService(&MockStore{...})

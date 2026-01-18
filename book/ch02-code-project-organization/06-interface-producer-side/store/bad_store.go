package store

// Bad Example: Producer defines large interface

// ❌ Producer-side interface forces abstraction on all clients
type CustomerStorage interface {
	StoreCustomer(customer Customer) error
	GetCustomer(id string) (Customer, error)
	UpdateCustomer(customer Customer) error
	GetAllCustomers() ([]Customer, error)
	GetCustomersWithoutContract() ([]Customer, error)
	GetCustomersWithNegativeBalance() ([]Customer, error)
}

// Problems:
// 1. Too large (6 methods) - weak abstraction
// 2. Clients that only need 1-2 methods are forced to depend on all 6
// 3. Hard to implement - requires all 6 methods
// 4. Low reusability
// 5. Producer decides abstraction level, not client

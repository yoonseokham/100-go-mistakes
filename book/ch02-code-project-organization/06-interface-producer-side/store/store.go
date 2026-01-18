package store

// Producer provides concrete type, NOT interface

type Customer struct {
	ID      string
	Name    string
	Balance int
	HasContract bool
}

// Store is the concrete implementation
type Store struct {
	customers map[string]Customer
}

func NewStore() *Store {
	return &Store{
		customers: make(map[string]Customer),
	}
}

func (s *Store) StoreCustomer(customer Customer) error {
	s.customers[customer.ID] = customer
	return nil
}

func (s *Store) GetCustomer(id string) (Customer, error) {
	customer, exists := s.customers[id]
	if !exists {
		return Customer{}, nil
	}
	return customer, nil
}

func (s *Store) UpdateCustomer(customer Customer) error {
	s.customers[customer.ID] = customer
	return nil
}

func (s *Store) GetAllCustomers() ([]Customer, error) {
	customers := make([]Customer, 0, len(s.customers))
	for _, customer := range s.customers {
		customers = append(customers, customer)
	}
	return customers, nil
}

func (s *Store) GetCustomersWithoutContract() ([]Customer, error) {
	var customers []Customer
	for _, customer := range s.customers {
		if !customer.HasContract {
			customers = append(customers, customer)
		}
	}
	return customers, nil
}

func (s *Store) GetCustomersWithNegativeBalance() ([]Customer, error) {
	var customers []Customer
	for _, customer := range s.customers {
		if customer.Balance < 0 {
			customers = append(customers, customer)
		}
	}
	return customers, nil
}

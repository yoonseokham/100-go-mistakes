package service

import "errors"

// Domain model — belongs to service layer
type Customer struct {
	ID   string
	Name string
}

var ErrNotFound = errors.New("customer not found")

// CustomerRepo is defined HERE (consumer side) — Chapter 6/7 principle
// Service knows the interface, not the concrete DB implementation
type CustomerRepo interface {
	Save(customer Customer) (Customer, error)
	FindByID(id string) (Customer, error)
}

// CustomerService contains ONLY business logic
// It does NOT know about HTTP or which DB is used
type CustomerService struct {
	repo CustomerRepo
}

func NewCustomerService(repo CustomerRepo) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(name string) (Customer, error) {
	// Business rules live here
	if name == "" {
		return Customer{}, errors.New("name is required")
	}
	if len(name) > 50 {
		return Customer{}, errors.New("name must be 50 characters or less")
	}

	return s.repo.Save(Customer{Name: name})
}

func (s *CustomerService) GetCustomer(id string) (Customer, error) {
	if id == "" {
		return Customer{}, errors.New("id is required")
	}
	return s.repo.FindByID(id)
}

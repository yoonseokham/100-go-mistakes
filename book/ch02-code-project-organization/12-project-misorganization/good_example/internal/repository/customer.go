package repository

import (
	"database/sql"
	"errors"

	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/12-project-misorganization/good_example/internal/service"
)

// CustomerRepository knows only about the DB
// It does NOT know about HTTP or business rules

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Save(customer service.Customer) (service.Customer, error) {
	result, err := r.db.Exec("INSERT INTO customers (name) VALUES (?)", customer.Name)
	if err != nil {
		return service.Customer{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return service.Customer{}, err
	}
	customer.ID = string(rune(id))
	return customer, nil
}

func (r *CustomerRepository) FindByID(id string) (service.Customer, error) {
	var customer service.Customer
	err := r.db.QueryRow("SELECT id, name FROM customers WHERE id = ?", id).
		Scan(&customer.ID, &customer.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.Customer{}, service.ErrNotFound
		}
		return service.Customer{}, err
	}
	return customer, nil
}

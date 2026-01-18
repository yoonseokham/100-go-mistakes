package interfacepollution

import (
	"database/sql"
	"fmt"
)

// Good Example: Start with Concrete Type
// Define interface only when needed (discovered, not created)

// Step 1: Concrete implementation
type ConcreteUserService struct {
	db *sql.DB
}

func NewConcreteUserService(db *sql.DB) *ConcreteUserService {
	return &ConcreteUserService{db: db}
}

func (s *ConcreteUserService) CreateUser(name string) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := s.db.Exec("INSERT INTO users (name) VALUES (?)", name)
	return err
}

func (s *ConcreteUserService) GetUser(id int) (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}
	var name string
	err := s.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	return name, err
}

func (s *ConcreteUserService) DeleteUser(id int) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// Step 2: When testing becomes necessary, define interface in consumer
// (see good_example_consumer.go)

// Benefits:
// 1. Simple and direct
// 2. Easy to understand
// 3. Interface added only when actually needed
// 4. No premature abstraction

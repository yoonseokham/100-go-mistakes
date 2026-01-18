package interfacepollution

import "database/sql"

// Bad Example: Interface Pollution
// Creating interface without immediate need

type UserService interface {
	CreateUser(name string) error
	GetUser(id int) (string, error)
	DeleteUser(id int) error
}

// Only one implementation exists
type UserServiceImpl struct {
	db *sql.DB
}

func NewUserServiceImpl(db *sql.DB) UserService {
	return &UserServiceImpl{db: db}
}

func (s *UserServiceImpl) CreateUser(name string) error {
	_, err := s.db.Exec("INSERT INTO users (name) VALUES (?)", name)
	return err
}

func (s *UserServiceImpl) GetUser(id int) (string, error) {
	var name string
	err := s.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	return name, err
}

func (s *UserServiceImpl) DeleteUser(id int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// Problems:
// 1. Unnecessary abstraction - only one implementation
// 2. Added complexity with no benefit
// 3. Harder to navigate code (interface -> impl -> actual code)
// 4. Created "just in case" rather than for actual need

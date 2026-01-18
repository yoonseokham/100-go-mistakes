package client

import "github.com/yourusername/go-100-mistakes/book/ch02-code-project-organization/06-interface-producer-side/store"

// Good Example: Consumer defines small, focused interface

// ✅ Client defines only what it needs
type customersGetter interface {
	GetAllCustomers() ([]store.Customer, error)
}

type Client struct {
	getter customersGetter
}

func NewClient(getter customersGetter) *Client {
	return &Client{getter: getter}
}

func (c *Client) ListAllCustomers() ([]store.Customer, error) {
	return c.getter.GetAllCustomers()
}

// Benefits:
// 1. Small interface (1 method) - strong abstraction
// 2. Only depends on what it actually uses
// 3. Easy to mock for testing
// 4. Client controls abstraction level

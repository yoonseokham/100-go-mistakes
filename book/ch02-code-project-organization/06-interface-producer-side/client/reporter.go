package client

import "github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/06-interface-producer-side/store"

// Different consumer, different interface

// ✅ Reporter defines its own interface with different needs
type customerBalanceReader interface {
	GetCustomersWithNegativeBalance() ([]store.Customer, error)
}

type Reporter struct {
	reader customerBalanceReader
}

func NewReporter(reader customerBalanceReader) *Reporter {
	return &Reporter{reader: reader}
}

func (r *Reporter) GenerateNegativeBalanceReport() ([]store.Customer, error) {
	return r.reader.GetCustomersWithNegativeBalance()
}

// Same store.Store satisfies both:
// - customersGetter (used by Client)
// - customerBalanceReader (used by Reporter)
// Each client defines what it needs!

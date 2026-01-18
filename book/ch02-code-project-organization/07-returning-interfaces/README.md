# Returning Interfaces

## TL;DR

To prevent being restricted in terms of flexibility, a function shouldn't return interfaces but concrete implementations in most cases. Conversely, a function should accept interfaces whenever possible.

## Core Principle

> **Return concrete types, accept interfaces.**

This means:
- **Return values**: Concrete types (structs, concrete implementations)
- **Parameters**: Interfaces (abstractions)
- **Abstractions**: Should be discovered by clients, not forced by producers

## Why Return Concrete Types?

### Problem 1: Restricts Flexibility

When a function returns an interface, all clients are forced to depend on that specific abstraction. If a client needs additional methods not in the interface, they're out of luck.

### Problem 2: Increases Package Dependencies

Returning interfaces creates tight coupling between the producer package and all consumer packages through the shared interface definition.

### Problem 3: Forces Same Abstraction on Everyone

Different clients may need different levels of abstraction, but returning an interface forces everyone to use the same one.

## The Problem: Returning Interfaces

### Bad Example: Function Returns Interface

```go
// store/store.go (producer)
package store

// ❌ Function returns interface
type CustomerStorage interface {
    StoreCustomer(customer Customer) error
    GetCustomer(id string) (Customer, error)
    UpdateCustomer(customer Customer) error
    GetAllCustomers() ([]Customer, error)
}

type InMemoryStore struct {
    customers map[string]Customer
}

func (s *InMemoryStore) StoreCustomer(customer Customer) error { ... }
func (s *InMemoryStore) GetCustomer(id string) (Customer, error) { ... }
func (s *InMemoryStore) UpdateCustomer(customer Customer) error { ... }
func (s *InMemoryStore) GetAllCustomers() ([]Customer, error) { ... }

type Customer struct {
    ID   string
    Name string
}

// ❌ Returns interface instead of concrete type
func NewInMemoryStore() CustomerStorage {
    return &InMemoryStore{
        customers: make(map[string]Customer),
    }
}
```

**Problems:**

1. **Clients are restricted** - Can only use the 4 methods defined in CustomerStorage
2. **No access to concrete methods** - If InMemoryStore has additional methods, clients can't use them
3. **Hard to extend** - Adding new functionality requires changing the interface, breaking all clients
4. **Package coupling** - All clients must import and depend on the store package's interface

### Client's Limitation

```go
// client/client.go
package client

import "myapp/store"

type Client struct {
    storage store.CustomerStorage  // ❌ Stuck with this interface
}

func NewClient() *Client {
    return &Client{
        storage: store.NewInMemoryStore(),  // Returns interface
    }
}

// What if InMemoryStore has a useful Debug() method?
// We can't access it because the interface doesn't include it!
```

## The Solution: Return Concrete Types

### Good Example: Function Returns Concrete Type

```go
// store/store.go (producer)
package store

// ✅ Concrete type, not interface
type InMemoryStore struct {
    customers map[string]Customer
}

func (s *InMemoryStore) StoreCustomer(customer Customer) error { ... }
func (s *InMemoryStore) GetCustomer(id string) (Customer, error) { ... }
func (s *InMemoryStore) UpdateCustomer(customer Customer) error { ... }
func (s *InMemoryStore) GetAllCustomers() ([]Customer, error) { ... }
func (s *InMemoryStore) Debug() string { ... }  // Additional method

type Customer struct {
    ID   string
    Name string
}

// ✅ Returns concrete type
func NewInMemoryStore() *InMemoryStore {
    return &InMemoryStore{
        customers: make(map[string]Customer),
    }
}
```

### Client Defines Its Own Interface

```go
// client/client.go
package client

import "myapp/store"

// ✅ Client defines only what it needs
type customerGetter interface {
    GetCustomer(id string) (store.Customer, error)
}

type Client struct {
    getter customerGetter  // Uses small, focused interface
}

// ✅ Accepts interface, not concrete type
func NewClient(getter customerGetter) *Client {
    return &Client{getter: getter}
}

func (c *Client) GetCustomerDetails(id string) (store.Customer, error) {
    return c.getter.GetCustomer(id)
}
```

### Usage Example

```go
// main.go
package main

func main() {
    // ✅ Get concrete type from constructor
    concreteStore := store.NewInMemoryStore()

    // ✅ Pass to client - automatically satisfies interface
    client := client.NewClient(concreteStore)

    // Can also access additional methods if needed
    debugInfo := concreteStore.Debug()
}
```

**Benefits:**

1. **Maximum flexibility** - Clients can access all methods of the concrete type
2. **Client-specific abstractions** - Each client defines the interface it needs
3. **No forced coupling** - Clients don't depend on producer-defined interfaces
4. **Easy to extend** - Adding methods to the concrete type doesn't break anything

## Accept Interfaces, Return Structs

This is a fundamental Go design principle.

### The Pattern

```go
// ❌ BAD: Returns interface, accepts concrete type
func NewService(db *sql.DB) ServiceInterface { ... }

// ✅ GOOD: Returns concrete type, accepts interface
func NewService(db Database) *Service { ... }
```

### Complete Example

```go
// service/service.go
package service

// Client can define their own interface for database
type Database interface {
    Query(query string) ([]Row, error)
}

// ✅ Concrete type
type UserService struct {
    db Database  // ✅ Accepts interface
}

// ✅ Returns concrete type, accepts interface
func NewUserService(db Database) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (User, error) { ... }
func (s *UserService) CreateUser(name string) error { ... }

type User struct {
    ID   int
    Name string
}
```

```go
// main.go
package main

import (
    "database/sql"
    "myapp/service"
)

func main() {
    // sql.DB satisfies the Database interface
    db, _ := sql.Open("postgres", "connection-string")

    // Get concrete type back
    userService := service.NewUserService(db)

    // Full access to all UserService methods
    userService.GetUser(1)
    userService.CreateUser("John")
}
```

## When It's OK to Return Interfaces

In particular contexts, when we **know** (not foresee) that an abstraction will be helpful for clients, we can consider returning an interface.

### Example: Standard Library Functions

```go
// os.Open returns io.ReadCloser interface
func Open(name string) (io.ReadCloser, error)

// io.Pipe returns interfaces
func Pipe() (*PipeReader, *PipeWriter)
```

**Why this is acceptable:**

1. **Known to be universally useful** - These abstractions are proven to be needed
2. **Stdlib convention** - Consistent with Go standard library patterns
3. **Composable** - Small interfaces that can be combined
4. **Don't restrict** - The abstractions are minimal and widely applicable

### Guidelines for Returning Interfaces

If you do return an interface:

1. **Know, don't foresee** - Be certain it's needed, not just anticipated
2. **Keep it minimal** - Fewer methods = stronger abstraction
3. **Prove the abstraction** - Have clear evidence it benefits clients
4. **Follow stdlib patterns** - Be consistent with Go conventions

## Real-World Example

### Bad: E-commerce Service Returns Interface

```go
// ❌ BAD
package payment

type PaymentProcessor interface {
    ProcessPayment(amount float64) error
    RefundPayment(transactionID string) error
    GetTransactionHistory() ([]Transaction, error)
}

// Returns interface - locks clients into this abstraction
func NewStripeProcessor(apiKey string) PaymentProcessor {
    return &stripeProcessor{apiKey: apiKey}
}
```

**Problems:**
- What if Stripe adds a new feature like "hold payment"?
- Clients can't access it without changing the interface
- All clients forced to update when interface changes

### Good: Returns Concrete Type

```go
// ✅ GOOD
package payment

type StripeProcessor struct {
    apiKey string
}

func (p *StripeProcessor) ProcessPayment(amount float64) error { ... }
func (p *StripeProcessor) RefundPayment(transactionID string) error { ... }
func (p *StripeProcessor) GetTransactionHistory() ([]Transaction, error) { ... }
func (p *StripeProcessor) HoldPayment(amount float64, duration time.Duration) error { ... }

// Returns concrete type
func NewStripeProcessor(apiKey string) *StripeProcessor {
    return &StripeProcessor{apiKey: apiKey}
}
```

```go
// Client defines what it needs
package checkout

type paymentProcessor interface {
    ProcessPayment(amount float64) error
}

type CheckoutService struct {
    processor paymentProcessor
}

// Accepts interface
func NewCheckoutService(processor paymentProcessor) *CheckoutService {
    return &CheckoutService{processor: processor}
}
```

## Comparison

| Returning Interfaces | Returning Concrete Types |
|---------------------|-------------------------|
| Forces abstraction on clients | Clients choose their abstraction |
| Restricts access to methods | Full access to all methods |
| Hard to extend | Easy to add new methods |
| Tight package coupling | Loose coupling |
| Producer decides abstraction | Consumer decides abstraction |

## Key Principles

1. **Default: Return concrete types** - Give clients maximum flexibility
2. **Accept interfaces in parameters** - Allow clients to pass any compatible type
3. **Let clients define interfaces** - They know their needs best
4. **Producer stays concrete** - Provide implementations, not abstractions
5. **Exception: Proven abstractions** - stdlib patterns and known universal needs

## Summary

**In most cases:**
- Functions should return concrete types (structs, implementations)
- Functions should accept interfaces as parameters
- Clients define their own interfaces based on their specific needs

**Exception:**
- Return interfaces when the abstraction is proven to be universally useful
- Follow standard library conventions (like io.ReadCloser)

**Bottom line:** Returning concrete types gives clients the flexibility to define their own abstractions while still accessing all functionality. Otherwise, it can make our design more complex due to package dependencies and can restrict flexibility because all clients would have to rely on the same abstraction.

**Remember:** If we know (not foresee) that an abstraction will be helpful for clients, we can consider returning an interface. Otherwise, we shouldn't force abstractions; they should be discovered by clients.

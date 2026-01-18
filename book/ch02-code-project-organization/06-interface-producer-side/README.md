# Interface on the Producer Side

## TL;DR

Keeping interfaces on the client side avoids unnecessary abstractions.

## Implicit Interface Satisfaction

Interfaces are satisfied implicitly in Go, which tends to be a gamechanger compared to languages with explicit implementation. This means:

- No `implements` keyword needed
- Types automatically satisfy interfaces if they have the required methods
- Flexible and decoupled design

## Core Principle

> **Abstractions should be discovered, not created.**

This means:
- ❌ It's NOT up to the producer to force a given abstraction for all clients
- ✅ It's up to the client to decide whether it needs some form of abstraction
- ✅ The client determines the best abstraction level for its needs

## The Problem: Producer-Side Interface

### Bad Example: Large Interface on Producer Side

```go
// store/store.go (producer)
package store

// ❌ Producer forces large interface on all clients
type CustomerStorage interface {
    StoreCustomer(customer Customer) error
    GetCustomer(id string) (Customer, error)
    UpdateCustomer(customer Customer) error
    GetAllCustomers() ([]Customer, error)
    GetCustomersWithoutContract() ([]Customer, error)
    GetCustomersWithNegativeBalance() ([]Customer, error)
}

type Customer struct{}
```

**Problems:**
1. **Too large** - 6 methods = weak abstraction
2. **Forces clients to depend on entire interface** even if they only need 1-2 methods
3. **Low reusability** - Hard to satisfy large interfaces
4. **Imposed abstraction** - Clients cannot choose their own abstraction level

### Client's Perspective

```go
// client/client.go
package client

import "myapp/store"

// Client only needs GetAllCustomers
type Client struct {
    storage store.CustomerStorage  // ❌ Forced to depend on 6 methods
}

func (c *Client) ListCustomers() ([]store.Customer, error) {
    return c.storage.GetAllCustomers()  // Only uses 1 method!
}
```

**Client is forced to:**
- Import entire `CustomerStorage` interface
- Depend on 5 methods it doesn't use
- Mock 6 methods for testing (when only 1 is needed)

## The Solution: Consumer-Side Interface

### Good Example: Interface on Client Side

```go
// store/store.go (producer)
package store

// ✅ Producer provides concrete type only
type Store struct {
    // implementation details
}

func (s *Store) StoreCustomer(customer Customer) error { ... }
func (s *Store) GetCustomer(id string) (Customer, error) { ... }
func (s *Store) UpdateCustomer(customer Customer) error { ... }
func (s *Store) GetAllCustomers() ([]Customer, error) { ... }
func (s *Store) GetCustomersWithoutContract() ([]Customer, error) { ... }
func (s *Store) GetCustomersWithNegativeBalance() ([]Customer, error) { ... }

type Customer struct {
    ID   string
    Name string
}
```

```go
// client/client.go (consumer)
package client

import "myapp/store"

// ✅ Client defines only what it needs
type customersGetter interface {
    GetAllCustomers() ([]store.Customer, error)
}

type Client struct {
    getter customersGetter  // Small, focused interface
}

func NewClient(getter customersGetter) *Client {
    return &Client{getter: getter}
}

func (c *Client) ListCustomers() ([]store.Customer, error) {
    return c.getter.GetAllCustomers()
}
```

**Benefits:**
1. **Small interface** - 1 method = strong abstraction
2. **Client-specific needs** - Each client defines what it needs
3. **Easy to mock** - Only mock what's actually used
4. **High reusability** - Small interfaces are easier to satisfy

### Another Client Example

```go
// reporter/reporter.go (different consumer)
package reporter

import "myapp/store"

// Different client, different needs, different interface
type customerReader interface {
    GetCustomer(id string) (store.Customer, error)
    GetCustomersWithNegativeBalance() ([]store.Customer, error)
}

type Reporter struct {
    reader customerReader  // Only needs 2 methods
}
```

**Same `Store` satisfies multiple different interfaces!**

## When Producer-Side Interface is Appropriate

In particular contexts, when we **know** (not foresee) that an abstraction will be helpful for consumers, we may want to have it on the producer side.

### Example: io.Reader and io.Writer

```go
// io package (producer)
package io

// ✅ Known to be useful for all consumers
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Why this is OK:**
1. **Known to be universally useful** (not just anticipated)
2. **Minimal** - Only 1 method each
3. **Highly reusable** - Simple to implement
4. **Easily composable** - Can combine into ReadWriter

### Guidelines for Producer-Side Interfaces

If you do place an interface on the producer side:

1. **Keep it minimal** - Fewer methods = stronger abstraction
2. **Ensure universal utility** - Actually useful for most/all clients
3. **Make it composable** - Small interfaces can be combined
4. **Don't foresee, know** - Be certain it's needed, not just anticipated

## Comparison

| Producer-Side Interface | Consumer-Side Interface |
|------------------------|-------------------------|
| Forces abstraction on all clients | Clients choose their abstraction |
| Often too large for specific needs | Tailored to exact needs |
| Lower reusability | Higher reusability |
| io.Reader (universally useful) | customersGetter (specific need) |

## Key Principles

1. **Default: Consumer-side** - Interfaces should live on the client side in most cases
2. **Small is powerful** - "The bigger the interface, the weaker the abstraction"
3. **Client decides** - Let clients determine their abstraction needs
4. **Producer stays concrete** - Provide concrete types, let clients abstract
5. **Exception: Universal abstractions** - Like io.Reader, when truly universal

## Summary

**In most cases:**
- Producer provides concrete types
- Consumer defines interfaces for its specific needs
- Multiple consumers can define different interfaces for the same producer

**Exception:**
- Producer-side interfaces like io.Reader when universally useful
- Keep them minimal for maximum reusability

**Bottom line:** Keeping interfaces on the client side avoids unnecessary abstractions and gives clients the flexibility to define exactly what they need.

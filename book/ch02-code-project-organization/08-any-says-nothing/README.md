# Any Says Nothing

## TL;DR

Only use `any` if you need to accept or return any possible type, such as `json.Marshal`. Otherwise, `any` doesn't provide meaningful information and can lead to compile-time issues by allowing a caller to call methods with any data type.

## Understanding `any`

`any` is a type alias for `interface{}` that can hold any type in Go.

```go
package main

func main() {
    var i any

    i = 42
    i = "foo"
    i = struct {
        s string
    }{
        s: "bar",
    }
    i = f

    _ = i
}

func f() {}
```

## The Problem: Overusing `any`

### Bad Example

```go
package store

type Customer struct {
    // Some fields
}

type Contract struct {
    // Some fields
}

type Store struct{}

func (s *Store) Get(id string) (any, error) {
    return nil, nil
}

func (s *Store) Set(id string, v any) error {
    return nil
}
```

**Problems:**
- No compile-time type checking
- Caller must do type assertions
- Can lead to runtime panics
- Unclear what types are expected

### Good Example

```go
package store

func (s *Store) GetContract(id string) (Contract, error) {
    return Contract{}, nil
}

func (s *Store) SetContract(id string, contract Contract) error {
    return nil
}

func (s *Store) GetCustomer(id string) (Customer, error) {
    return Customer{}, nil
}

func (s *Store) SetCustomer(id string, customer Customer) error {
    return nil
}
```

**Benefits:**
- Type-safe at compile time
- No type assertions needed
- Clear intent
- Better IDE support

## When to Use `any`

Use `any` when you genuinely need to accept or return any possible type:

- **JSON marshaling**: `json.Marshal(v any)`
- **Formatting**: `fmt.Println(a ...any)`
- Other cases where any type must be handled

## Summary

The `any` type can be helpful if there is a genuine need for accepting or returning any possible type (for instance, when it comes to marshaling or formatting). In general, we should avoid overgeneralizing the code we write at all costs. Perhaps a little bit of duplicated code might occasionally be better if it improves other aspects such as code expressiveness.

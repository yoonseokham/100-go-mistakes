# Any Says Nothing

## TL;DR

Only use `any` if you need to accept or return any possible type, such as `json.Marshal`. Otherwise, `any` doesn't provide meaningful information and can lead to compile-time issues by allowing a caller to call methods with any data type.

## What is `any`?

`any` is a type alias for `interface{}` introduced in Go 1.18. It represents any type in Go.

```go
// These are equivalent
var x any
var y interface{}
```

Since `any` can hold any type, it provides maximum flexibility but at the cost of type safety.

## Understanding `any`

### How `any` Works

```go
package main

func main() {
    var i any

    i = 42                    // int
    i = "foo"                // string
    i = struct {             // struct
        s string
    }{
        s: "bar",
    }
    i = f                     // function

    _ = i
}

func f() {}
```

**Key point:** `any` can hold **any value** - integers, strings, structs, functions, anything.

## The Problem: Overusing `any`

### Bad Example: Generic Store with `any`

```go
package store

type Customer struct {
    ID   string
    Name string
}

type Contract struct {
    ID     string
    Terms  string
}

type Store struct{}

// ❌ BAD: Using any loses all type information
func (s *Store) Get(id string) (any, error) {
    return nil, nil
}

func (s *Store) Set(id string, v any) error {
    return nil
}
```

### Problems with This Approach

**1. No Type Safety**

```go
store := &Store{}

// Compiler allows anything!
store.Set("customer123", "wrong type")     // Compiles fine
store.Set("customer123", 42)               // Compiles fine
store.Set("customer123", []byte{1, 2, 3})  // Compiles fine
store.Set("customer123", Customer{})       // Also compiles fine

// No way to know what's the correct type!
```

**2. Type Assertions Required**

```go
// Caller must do type assertion
result, _ := store.Get("customer123")

// Risky - can panic at runtime!
customer := result.(Customer)

// Safer but verbose
customer, ok := result.(Customer)
if !ok {
    // Handle wrong type
}
```

**3. No Compile-Time Checks**

```go
// These mistakes won't be caught until runtime
store.Set("contract456", Customer{})  // Wrong type, but compiles!

result, _ := store.Get("contract456")
contract := result.(Contract)  // Runtime panic!
```

**4. No IDE Support**

- No autocomplete
- No refactoring support
- No type hints
- Harder to understand code

**5. Unclear Intent**

```go
// What does this function do? What type does it return?
func (s *Store) Get(id string) (any, error)

// Compare with:
func (s *Store) GetCustomer(id string) (Customer, error)
// Immediately clear!
```

## The Solution: Use Specific Types

### Good Example: Type-Specific Methods

```go
package store

type Customer struct {
    ID   string
    Name string
}

type Contract struct {
    ID     string
    Terms  string
}

type Store struct{}

// ✅ GOOD: Specific types for each entity
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

### Benefits

**1. Type Safety**

```go
store := &Store{}

// Compiler catches mistakes!
store.SetCustomer("123", Customer{ID: "123", Name: "John"})  // ✅ OK
store.SetCustomer("123", Contract{})                          // ❌ Compile error!
store.SetCustomer("123", "wrong")                             // ❌ Compile error!
```

**2. No Type Assertions Needed**

```go
// Direct usage, no casting
customer, err := store.GetCustomer("123")
if err != nil {
    return err
}
// customer is already Customer type!
fmt.Println(customer.Name)
```

**3. Clear Intent**

```go
// Immediately obvious what this does
func (s *Store) GetCustomer(id string) (Customer, error)
func (s *Store) GetContract(id string) (Contract, error)
```

**4. IDE Support**

- Full autocomplete
- Refactoring works correctly
- Type information in tooltips
- Better code navigation

## When to Use `any`

### Legitimate Use Cases

Use `any` when you **genuinely need** to accept or return **any possible type**.

#### 1. JSON Marshaling/Unmarshaling

```go
// ✅ GOOD: json.Marshal needs to handle any type
func Marshal(v any) ([]byte, error)
func Unmarshal(data []byte, v any) error

// Usage
customer := Customer{ID: "123", Name: "John"}
data, _ := json.Marshal(customer)  // Works with any type

contract := Contract{ID: "456"}
data2, _ := json.Marshal(contract)  // Also works
```

#### 2. Formatting and Logging

```go
// ✅ GOOD: fmt functions need to format any type
func Println(a ...any) (n int, err error)
func Sprintf(format string, a ...any) string

// Usage
fmt.Println("Customer:", customer)
fmt.Println("Count:", 42)
fmt.Println("Active:", true)
```

#### 3. Generic Data Structures (with caution)

```go
// ✅ OK: Cache that stores any value
type Cache struct {
    data map[string]any
}

func (c *Cache) Set(key string, value any) {
    c.data[key] = value
}

func (c *Cache) Get(key string) (any, bool) {
    val, ok := c.data[key]
    return val, ok
}
```

#### 4. Configuration/Plugin Systems

```go
// ✅ OK: Configuration values can be various types
type Config struct {
    settings map[string]any
}

// But validate types when retrieving
func (c *Config) GetString(key string) (string, error) {
    val, ok := c.settings[key]
    if !ok {
        return "", fmt.Errorf("key not found")
    }
    str, ok := val.(string)
    if !ok {
        return "", fmt.Errorf("value is not a string")
    }
    return str, nil
}
```

## Code Duplication vs Type Safety

### The Trade-off

Some developers avoid duplication at all costs:

```go
// ❌ Avoiding duplication with any
func (s *Store) Get(id string) (any, error)
func (s *Store) Set(id string, v any) error

// One implementation handles everything
```

But this sacrifices type safety:

```go
// ✅ A little duplication is better than losing type safety
func (s *Store) GetCustomer(id string) (Customer, error)
func (s *Store) SetCustomer(id string, customer Customer) error
func (s *Store) GetContract(id string) (Contract, error)
func (s *Store) SetContract(id string, contract Contract) error

// Yes, more code, but much safer and clearer
```

## Key Principles

1. **`any` says nothing** - It provides no information about what type to expect
2. **Prefer specific types** - Even if it means some duplication
3. **Type safety matters** - Compile-time errors are better than runtime panics
4. **Use `any` sparingly** - Only when truly necessary

## Comparison

| Using `any` | Using Specific Types |
|------------|---------------------|
| No compile-time type checking | Full compile-time type checking |
| Requires type assertions | Direct usage, no casting |
| Runtime panics possible | Type errors caught at compile time |
| No IDE support | Full IDE support |
| Unclear intent | Clear and explicit |
| Less code | More code but safer |

## Summary

**The Problem:** `any` provides no type information and removes compile-time type safety.

**The Solution:** Use specific types whenever possible, even if it means writing more code.

**The Exception:** Use `any` only when you genuinely need to accept or return any possible type (marshaling, formatting, etc.).

**The Wisdom:** A little bit of duplicated code might occasionally be better if it improves other aspects such as code expressiveness and type safety.

**Bottom line:** Don't overgeneralize your code at all costs. Specific types make your code safer, clearer, and easier to maintain.

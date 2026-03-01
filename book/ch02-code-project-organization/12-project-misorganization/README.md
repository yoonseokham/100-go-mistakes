# Project Misorganization (#12)

## TL;DR

Avoid premature packaging. Name packages after what they provide. Minimize exports. Use `internal/` to hide implementation details. Pick a structure (by context or by layer) and stay consistent.

---

## Project Structure: Two Schools of Thought

### By Context (Domain-driven)

Group by business concept:

```
shop/
  customer/
    handler.go
    service.go
    repository.go
  order/
    handler.go
    service.go
    repository.go
```

"Where is the customer code?" → Everything in `customer/`.

### By Layer (Hexagonal Architecture)

Group by technical responsibility:

```
shop/
  handler/
    customer.go
    order.go
  service/
    customer.go
    order.go
  repository/
    customer.go
    order.go
```

"Where are all DB queries?" → Everything in `repository/`.

**Both are valid. The key is consistency — never mix the two.**

---

## The Three Layers

### Handler
- Receives HTTP requests, parses them, writes responses
- Calls Service
- Contains NO business logic

```go
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    json.NewDecoder(r.Body).Decode(&req)   // parse
    customer, err := h.svc.CreateCustomer(req.Name)  // delegate
    json.NewEncoder(w).Encode(customer)    // respond
}
```

### Service
- Contains ALL business rules
- Does NOT know whether caller is HTTP or gRPC
- Does NOT know whether DB is MySQL or Postgres
- Depends only on a Repository **interface**

```go
func (s *CustomerService) CreateCustomer(name string) (Customer, error) {
    if name == "" {
        return Customer{}, errors.New("name is required")  // business rule
    }
    return s.repo.Save(Customer{Name: name})  // doesn't know which DB
}
```

### Repository
- Executes DB queries only
- Contains NO business logic

```go
func (r *CustomerRepository) Save(c Customer) (Customer, error) {
    _, err := r.db.Exec("INSERT INTO customers (name) VALUES (?)", c.Name)
    return c, err
}
```

---

## Dependency Flow

Dependencies flow in one direction only:

```
Handler → Service → Repository → DB
```

Each layer knows the layer below only through an **interface** (not a concrete type):

```go
// Handler knows Service via interface
type customerService interface {
    CreateCustomer(name string) (Customer, error)
}

// Service knows Repository via interface
type CustomerRepo interface {
    Save(customer Customer) (Customer, error)
}
```

Interfaces are defined on the **consumer side** (Chapters 6 & 7).

### Wiring happens only in main

```go
func main() {
    db   := connectDB()
    repo := repository.NewCustomerRepository(db)  // inject DB
    svc  := service.NewCustomerService(repo)      // inject repo
    h    := handler.NewCustomerHandler(svc)       // inject svc
}
```

---

## Why This Matters

| Problem (Bad) | Solution (Good) |
|---|---|
| Can't test without real DB | Inject MockRepo into Service |
| Swapping DB requires touching handler | Only swap Repository |
| Business rules scattered everywhere | All rules in Service |
| `init()` makes testing hard | Explicit constructors |

---

## Package Best Practices

### Avoid premature packaging

Don't force perfect structure up front. Let it evolve naturally.

### Avoid nano packages

```
❌ Too granular              ✅ Balanced
  parse/parse.go               encoding/
  validate/validate.go           encode.go
  format/format.go               decode.go
                                 validate.go
```

### Name packages after what they provide, not what they contain

```
❌ utils, helpers, common, misc
✅ http, auth, encoding, storage
```

Short, concise, expressive, single lowercase word.

### Minimize exports

```go
// Default: unexported
type internalHelper struct { ... }

// Export only when necessary
type PublicAPI struct { ... }
```

When unsure → don't export. You can always export later.

### Use `internal/`

```
project/
  internal/       ← cannot be imported by other modules
    service/
    repository/
  cmd/
    server/
      main.go
```

Code in `internal/` can be freely refactored without breaking external users.

---

## Recommended Server Structure

```
project/
  go.mod
  internal/
    handler/
      customer.go
    service/
      customer.go
    repository/
      customer.go
  cmd/
    server/
      main.go
```

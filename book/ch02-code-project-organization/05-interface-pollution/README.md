# Interface Pollution

## TL;DR

Abstractions should be discovered, not created. To prevent unnecessary complexity, create an interface when you need it and not when you foresee needing it, or if you can at least prove the abstraction to be a valid one.

## What is Interface Pollution?

Interface pollution is about overwhelming our code with unnecessary abstractions, making it harder to understand. It's a common mistake made by developers coming from languages with different habits (Java, C#).

## Understanding Go Interfaces

### Implicit Satisfaction

Go interfaces are satisfied **implicitly**. There is no explicit keyword like `implements` to mark that an object implements an interface.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// *os.File satisfies io.Reader automatically
// No need to declare "implements Reader"
file, _ := os.Open("file.txt")
var r io.Reader = file  // Just works!
```

### The Power of io.Reader and io.Writer

The `io` package provides fundamental abstractions:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Why are these powerful?**

```go
// Generic function using interfaces
func copySourceToDest(source io.Reader, dest io.Writer) error {
    buf := make([]byte, 1024)
    for {
        n, err := source.Read(buf)
        if n > 0 {
            dest.Write(buf[:n])
        }
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
    }
}

// Works with files, networks, memory, anything!
file1, _ := os.Open("source.txt")
file2, _ := os.Create("dest.txt")
copySourceToDest(file1, file2)

// Works with strings and buffers
source := strings.NewReader("data")
dest := &bytes.Buffer{}
copySourceToDest(source, dest)
```

### Interface Granularity

> **"The bigger the interface, the weaker the abstraction."** - Rob Pike

- Small interfaces are more reusable
- io.Reader and io.Writer are powerful because they cannot get any simpler
- Fine-grained interfaces can be combined into higher-level abstractions

```go
type ReadWriter interface {
    Reader
    Writer
}
```

## When to Use Interfaces

### 1. Common Behavior

Use interfaces when multiple types implement a common behavior.

**Example: sort.Interface**

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

This interface factors out the common behavior to sort any index-based collection.

### 2. Decoupling

Use interfaces to decouple code from a concrete implementation, enabling the Liskov Substitution Principle.

**Bad: Tight Coupling**

```go
type CustomerService struct {
    store mysql.Store  // Depends on concrete implementation
}

func (cs CustomerService) CreateNewCustomer(id string) error {
    customer := Customer{id: id}
    return cs.store.StoreCustomer(customer)
}
```

**Good: Decoupled with Interface**

```go
type customerStorer interface {
    StoreCustomer(Customer) error
}

type CustomerService struct {
    storer customerStorer  // Depends on abstraction
}

func (cs CustomerService) CreateNewCustomer(id string) error {
    customer := Customer{id: id}
    return cs.storer.StoreCustomer(customer)
}
```

**Benefits:**
- Integration tests with real implementation
- Unit tests with mocks
- Easy to swap implementations

### 3. Restricting Behavior

Use interfaces to restrict a type to a specific behavior.

```go
type IntConfig struct {
    // ...
}

func (c *IntConfig) Get() int {
    // Retrieve configuration
}

func (c *IntConfig) Set(value int) {
    // Update configuration
}

// Restrict to read-only
type intConfigGetter interface {
    Get() int
}

type Foo struct {
    threshold intConfigGetter  // Read-only
}

func (f Foo) Bar() {
    threshold := f.threshold.Get()  // Can only read
    // f.threshold.Set(100)  // Compile error!
}
```

## Interface Pollution Problems

### The Problem

It's common to see interfaces being overused in Go projects, especially by developers from C# or Java backgrounds.

**Bad Example: Unnecessary Interface**

```go
// Interface pollution
type UserService interface {
    CreateUser(name string) error
    GetUser(id int) (string, error)
}

type UserServiceImpl struct {
    db *sql.DB
}

// Only one implementation exists!
```

**Why is this bad?**
- Adds useless level of indirection
- Creates worthless abstraction
- Makes code more difficult to read and reason about
- No immediate reason for the abstraction

### The Main Caveat

> **Abstractions should be discovered, not created.**

- Don't start creating abstractions if there is no immediate reason
- Don't design with interfaces but wait for a concrete need
- Create an interface when you **need** it, not when you **foresee** that you could need it

### Rob Pike's Wisdom

> **"Don't design with interfaces, discover them."**

## Good Example: Discovered Abstraction

### Step 1: Start with Concrete Type

```go
// userservice/userservice.go
package userservice

type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) CreateUser(name string) error {
    _, err := s.db.Exec("INSERT INTO users (name) VALUES (?)", name)
    return err
}

func (s *UserService) GetUser(id int) (string, error) {
    var name string
    err := s.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
    return name, err
}
```

### Step 2: Discover Interface When Needed

```go
// handler/handler.go
package handler

// Define interface in consumer package
type UserCreator interface {
    CreateUser(name string) error
}

type UserHandler struct {
    creator UserCreator  // Depend on interface
}

func NewUserHandler(creator UserCreator) *UserHandler {
    return &UserHandler{creator: creator}
}
```

### Step 3: Automatically Satisfied

```go
// main.go
package main

func main() {
    db, _ := sql.Open("mysql", "connection-string")

    // Create concrete type
    userSvc := userservice.NewUserService(db)

    // Pass to handler - automatically satisfies interface
    handler := handler.NewUserHandler(userSvc)
}
```

## Key Principles

1. **Concrete types first** - Implement what you need now
2. **Interfaces when needed** - Add abstraction only when there's a clear benefit
3. **Consumer defines interface** - Define interfaces where they're used
4. **Small interfaces** - Keep them minimal and focused
5. **Question abstractions** - If it's unclear how an interface makes code better, consider removing it

## When NOT to Use Interfaces

❌ Don't create interfaces when:
- You have only one implementation
- "Just in case" we need it later
- Following habits from other languages
- The benefit is unclear

## When to Use Interfaces

✅ Do create interfaces when:
- Multiple implementations exist or are clearly needed
- Testing requires mocks
- Decoupling from implementation details
- Restricting behavior to specific operations
- Factoring out common behavior

## Summary

**The Problem:** Overwhelming code with unnecessary abstractions

**The Solution:** Discover abstractions when needed, don't create them preemptively

**The Guideline:** If it's unclear how an interface makes the code better, probably don't add it

Keep Go code simple and pragmatic. Let abstractions emerge naturally from actual needs rather than anticipated ones.

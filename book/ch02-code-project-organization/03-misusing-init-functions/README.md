# Misusing Init Functions

## TL;DR
When initializing variables, remember that init functions have limited error handling and make state handling and testing more complex. In most cases, initializations should be handled as specific functions.

## What is an Init Function?

An init function is used to initialize the state of an application. It:
- Takes no arguments and returns no result (`func()` signature)
- Executes automatically when a package is initialized
- Runs before the `main()` function

### Execution Order

When a package is initialized:
1. All constant and variable declarations are evaluated
2. All `init()` functions execute (in declaration order)
3. Then `main()` runs (if it's the main package)

If multiple `init()` functions exist:
- They execute in the order they appear in the source file
- Across multiple files, they execute in file name alphabetical order

## Example

### Execution Order Demo

**redis/redis.go:**
```go
func init() {
    fmt.Println("Initializing redis cache...")
    cache = make(map[string]string)
}
```

**main/main.go:**
```go
func init() {
    fmt.Println("init 1")
}

func init() {
    fmt.Println("init 2")
}

func main() {
    fmt.Println("Starting main...")
    redis.Store("foo", "bar")
}
```

**Output:**
```
Initializing redis cache...
init 1
init 2
Starting main...
Stored: foo = bar
```

## Problems with Init Functions

### 1. Limited Error Management

Init functions cannot return errors. Your only options are:
- Panic (stops the entire application)
- Ignore the error (leaves invalid state)

```go
func init() {
    var err error
    db, err = sql.Open("postgres", "connection-string")
    if err != nil {
        panic(err)  // Only option: panic or ignore
    }
}
```

### 2. Testing Complications

- Init functions run automatically before tests
- Cannot mock or control initialization
- Must set up external dependencies even for unit tests that don't need them
- Cannot create isolated test instances

```go
// Bad: Testing depends on init() side effects
func TestStore(t *testing.T) {
    // Relies on init() having run successfully
    redis.Store("key", "value")
}
```

### 3. Global State Requirement

If initialization requires state, you must use global variables:

```go
var cache map[string]string  // Global variable

func init() {
    cache = make(map[string]string)  // Sets global state
}
```

This creates:
- Unexpected behavior across codebase
- Difficult unit testing
- Hard-to-reason-about code

## Solution: Explicit Initialization

### Good Example: Constructor Pattern

```go
type Cache struct {
    data map[string]string
}

func NewCache() (*Cache, error) {
    // Can return error if initialization fails
    return &Cache{
        data: make(map[string]string),
    }, nil
}

func (c *Cache) Store(key, value string) error {
    c.data[key] = value
    return nil
}
```

**Benefits:**
- Proper error handling
- Easy to test with multiple instances
- No global state
- Explicit initialization timing
- Can inject dependencies

### Testing is Easy

```go
func TestCache(t *testing.T) {
    // Create fresh, isolated instance
    cache, err := NewCache()
    if err != nil {
        t.Fatalf("NewCache failed: %v", err)
    }

    // Test with controlled state
    cache.Store("key", "value")
}
```

## When Init Functions ARE Appropriate

Use `init()` only for:

1. **Static configuration** - No errors possible
2. **Registry registration** - e.g., database drivers
3. **Constant data setup** - Deterministic, no side effects

```go
var supportedFormats = map[string]bool{}

func init() {
    // OK: Simple, deterministic, no errors
    supportedFormats["json"] = true
    supportedFormats["xml"] = true
}
```

## Key Principles

- ❌ Don't use `init()` for anything that can fail
- ❌ Don't use `init()` for expensive operations
- ❌ Don't use `init()` when you need dependency injection
- ✅ Use explicit constructors that return errors
- ✅ Use dependency injection for testability
- ✅ Make initialization explicit and controllable

## Run Example

```bash
# Run the demo
bazel run //book/ch02-code-project-organization/03-misusing-init-functions/main

# Run tests
bazel test //book/ch02-code-project-organization/03-misusing-init-functions/redis:redis_test
```

## Conclusion

**Be cautious with init functions.** They can be helpful in some situations, such as defining static configuration. Otherwise, and in most cases, you should handle initializations through ad hoc functions.

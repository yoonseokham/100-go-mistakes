# Misusing init Functions

## Go Package Initialization Order

Understanding the initialization order is crucial:

```
package A imports B, C

Initialization order:
1. B's global variables
2. B's init() functions
3. C's global variables
4. C's init() functions
5. A's global variables
6. A's init() functions
7. A's main() function (if main package)
```

Within a package:
- Global variables are initialized in declaration order (respecting dependencies)
- `init()` functions run after all variables are initialized
- Multiple `init()` functions in the same package run in source file name order
- Multiple `init()` in the same file run in declaration order

## Problems with init()

### Problem 1: Cannot Return Errors

```go
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", "connection-string")
    if err != nil {
        // Only options: panic or ignore
        panic(err)  // Makes testing hard
    }
}
```

**Issues:**
- No way to gracefully handle errors
- Must panic or leave invalid state
- Cannot retry or provide fallback

### Problem 2: Hard to Test

```go
var globalCache map[string]string

func init() {
    globalCache = make(map[string]string)
    loadCacheFromFile()  // Runs on import!
}
```

**Issues:**
- Runs automatically on package import
- Cannot mock or control in tests
- Cannot inject dependencies
- Increases test startup time

### Problem 3: Order Dependencies

```go
var config Config
var client *Client

func init() {
    config = loadConfig()
}

func init() {
    // Might run before config is loaded!
    client = NewClient(config)
}
```

**Issues:**
- Order across files is not guaranteed
- Leads to subtle, hard-to-debug issues
- Fragile when refactoring

## Solutions: Explicit Initialization

### Solution 1: Constructor with Error Handling

```go
type Database struct {
    conn *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
    conn, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    return &Database{conn: conn}, nil
}
```

**Benefits:**
- Can return errors properly
- Easy to test with different inputs
- Explicit about when initialization happens

### Solution 2: Lazy Initialization

```go
type Cache struct {
    data map[string]string
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
    }
}

func (c *Cache) Load() error {
    // Load when needed, not at package init
    return nil
}
```

**Benefits:**
- Delays expensive operations
- Only initializes what's actually used
- Can handle errors at call site

### Solution 3: Dependency Injection

```go
type Service struct {
    db     *Database
    cache  *Cache
    client *Client
}

func NewService(dbConnStr string, config Config) (*Service, error) {
    db, err := NewDatabase(dbConnStr)
    if err != nil {
        return nil, err
    }

    cache := NewCache()
    client := NewClient(config)

    return &Service{
        db:     db,
        cache:  cache,
        client: client,
    }, nil
}
```

**Benefits:**
- Clear dependency chain
- Easy to test with mocks
- Explicit initialization order
- Proper error handling

## When init() IS Appropriate

Use `init()` only for:

1. **Registering with registries:**
```go
func init() {
    sql.Register("mydriver", &MyDriver{})
}
```

2. **Setting up constant, deterministic data:**
```go
var supportedFormats map[string]bool

func init() {
    supportedFormats = map[string]bool{
        "json": true,
        "xml":  true,
    }
}
```

3. **No error handling needed**
4. **No external dependencies or I/O**
5. **Truly one-time setup**

## Key Principles

- ❌ Don't use `init()` for anything that can fail
- ❌ Don't use `init()` for expensive operations
- ❌ Don't use `init()` when you need dependency injection
- ✅ Use explicit constructors that return errors
- ✅ Use dependency injection for testability
- ✅ Initialize only when needed (lazy initialization)

## Run Tests

```bash
bazel test //book/ch02-code-project-organization/03-misusing-init-functions:misuinginit_test
```

# Not Using the Functional Options Pattern (#11)

## TL;DR

To handle options conveniently and in an API-friendly manner, use the functional options pattern.

---

## The Problem

When a constructor has many optional parameters, there are two naive approaches — both with drawbacks.

### Bad V1: Parameter Explosion

```go
func NewServer(addr string, port int, timeout time.Duration, maxConn int) (*http.Server, error)
```

- Adding a new option changes the signature → breaks all callers
- Callers must pass every argument even if they don't care about it

### Bad V2: Config Struct

```go
func NewServer(addr string, cfg Config) (*http.Server, error)
```

- Passing `Config{}` is awkward and unclear
- Can't distinguish "port not set" from "port explicitly set to 0"
- No clean place to validate each field

---

## The Solution: Functional Options Pattern

### Three building blocks

**1. Unexported `options` struct — holds all config**

```go
type options struct {
    port    *int          // pointer: nil = not set, 0 = random port
    timeout time.Duration
    maxConn int
}
```

**2. `Option` function type — mutates options**

```go
type Option func(opts *options) error
```

**3. `WithXxx` functions — return an Option, validate inline**

```go
func WithPort(port int) Option {
    return func(opts *options) error {
        if port < 0 {
            return errors.New("port should be positive")
        }
        opts.port = &port
        return nil
    }
}
```

**4. Constructor applies options in order**

```go
func NewServer(addr string, opts ...Option) (*http.Server, error) {
    var o options
    for _, opt := range opts {
        if err := opt(&o); err != nil {
            return nil, err
        }
    }
    // apply defaults ...
}
```

---

## Nil Pointer Trick

Use `*int` instead of `int` to distinguish three states:

| Value | Meaning |
|-------|---------|
| `nil` | Option not provided → use default |
| `&0`  | Explicitly set to 0 → use random port |
| `&n`  | Explicitly set to n → use n |

---

## Usage

```go
NewServer("localhost")                                          // all defaults
NewServer("localhost", WithPort(8080))                         // custom port
NewServer("localhost", WithPort(8080), WithTimeout(30*time.Second)) // multiple options
NewServer("localhost", WithPort(-1))                           // returns error
```

---

## Comparison

| | Parameter Explosion | Config Struct | Functional Options |
|---|---|---|---|
| Adding new option | Breaks callers | Non-breaking | Non-breaking |
| Omitting options | Must pass zero values | Pass `Config{}` | Just omit |
| Validation | In constructor | In constructor | In each `WithXxx` |
| nil vs zero distinction | ❌ | ❌ | ✅ via pointer |
| Go idiomatic | ❌ | △ | ✅ |

---

## Key Principles

- **Required values** → regular parameters (`addr string`)
- **Optional values** → `Option` functions (`WithPort`, `WithTimeout`)
- **Default vs not provided** → use pointer (`*int`) and check for `nil`
- **Validation** → inside each `WithXxx`, not in the constructor


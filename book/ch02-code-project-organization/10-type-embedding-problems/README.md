# Not Being Aware of the Possible Problems with Type Embedding

## TL;DR

Using type embedding can help avoid boilerplate code; however, ensure that doing so doesn't lead to visibility issues where some fields should have remained hidden.

## What is Type Embedding?

In Go, a struct field is called **embedded** if it's declared without a name:

```go
type Foo struct {
    Bar // Embedded field (no name)
}

type Bar struct {
    Baz int
}
```

In the `Foo` struct, the `Bar` type is declared without an associated name; hence, it's an embedded field.

### Field and Method Promotion

Embedding **promotes** the fields and methods of an embedded type. Because `Bar` contains a `Baz` field, this field is promoted to `Foo`:

```go
foo := Foo{}
foo.Baz = 42  // Baz is promoted and accessible from Foo
```

## When to Use Type Embedding

Type embedding is **rarely a necessity**. Whatever the use case, we can probably solve it without type embedding. Type embedding is mainly used for **convenience**: in most cases, to promote behaviors.

### Good Use Case: Promoting Behavior

Consider a `Logger` that wraps `io.WriteCloser`:

```go
type Logger struct {
    writeCloser io.WriteCloser
}

func (l Logger) Write(p []byte) (int, error) {
    return l.writeCloser.Write(p)
}

func (l Logger) Close() error {
    return l.writeCloser.Close()
}

func main() {
    l := Logger{writeCloser: os.Stdout}
    _, _ = l.Write([]byte("foo"))
    _ = l.Close()
}
```

This requires forwarding methods. With type embedding:

```go
type Logger struct {
    io.WriteCloser  // Embedded - methods promoted
}

func main() {
    l := Logger{WriteCloser: os.Stdout}
    _, _ = l.Write([]byte("foo"))  // Write() promoted
    _ = l.Close()                   // Close() promoted
}
```

The forwarding methods are no longer needed, reducing boilerplate.

## When NOT to Use Type Embedding

### 1. Syntactic Sugar Only

Don't embed a type solely to simplify field access:

```go
// ❌ BAD: Only for syntactic convenience
type Foo struct {
    Bar
}

type Bar struct {
    Baz int
}

// Just to write foo.Baz instead of foo.Bar.Baz
```

If this is the only rationale, **don't embed the inner type** and use a field instead.

### 2. Exposing Internal Implementation

Don't embed types that promote data or behavior that should remain hidden.

#### Bad Example: Embedding sync.Mutex

```go
type InMem struct {
    sync.Mutex  // ❌ Embedded - Lock/Unlock become public!
    m map[string]int
}

func New() *InMem {
    return &InMem{m: make(map[string]int)}
}

func (i *InMem) Get(key string) (int, bool) {
    i.Lock()
    v, contains := i.m[key]
    i.Unlock()
    return v, contains
}
```

**Problem:**

```go
inMem := New()

// ❌ Clients can call Lock/Unlock directly!
inMem.Lock()
// ...
inMem.Unlock()

// This breaks encapsulation - internal locking behavior is exposed
```

The locking behavior should remain **private** to the struct, but embedding `sync.Mutex` promotes `Lock()` and `Unlock()` methods, making them public.

#### Good Example: Use a Named Field

```go
type InMem struct {
    mu sync.Mutex  // ✅ Named field - not promoted
    m  map[string]int
}

func New() *InMem {
    return &InMem{m: make(map[string]int)}
}

func (i *InMem) Get(key string) (int, bool) {
    i.mu.Lock()    // Internal use only
    v, contains := i.m[key]
    i.mu.Unlock()
    return v, contains
}
```

Now `Lock()` and `Unlock()` are not accessible from outside.

## Key Constraints

If we decide to use type embedding, keep two main constraints in mind:

1. **Don't use it solely as syntactic sugar** to simplify accessing a field (such as `Foo.Baz()` instead of `Foo.Bar.Baz()`). If this is the only rationale, use a field instead.

2. **Don't promote data or behavior we want to hide** from the outside. For example, don't embed types that allow clients to access locking behavior that should remain private to the struct.

## Basic Example

```go
type Foo struct {
    Bar  // Embedded
}

type Bar struct {
    Baz int
}

func main() {
    foo := Foo{}
    foo.Baz = 42  // Bar.Baz is promoted
}
```

## Summary

Using type embedding consciously by keeping these constraints in mind can help avoid boilerplate code with additional forwarding methods. However, let's make sure we don't do it solely for cosmetics and not promote elements that should remain hidden.

**Remember:**
- Type embedding is for **convenience** and **promoting behaviors**
- It's **not mandatory** - you can always solve problems without it
- Be careful not to expose internal implementation details
- Use named fields when you need encapsulation

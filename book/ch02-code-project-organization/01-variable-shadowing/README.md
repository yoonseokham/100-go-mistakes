# Variable Shadowing

## Problem

Using `:=` inside if/else blocks creates new variables that shadow outer scope variables, causing the outer variable to remain uninitialized.

```go
var client *http.Client
if tracing {
    client, err := createTracingClient()  // Creates NEW client in this scope
    // ...
}
// client is still nil here!
```

## Solutions

### Solution 1: Use `=` instead of `:=`

Assign to existing variables instead of creating new ones:

```go
var client *http.Client
var err error

if tracing {
    client, err = createTracingClient()  // Assigns to outer scope
    // ...
}
```

### Solution 2: Use temporary variable

Use a temporary variable inside the block and assign after:

```go
var client *http.Client

if tracing {
    c, err := createTracingClient()  // Temporary variable c
    if err != nil {
        return nil, err
    }
    client = c  // Assign to outer scope
}
```

## Run Tests

```bash
bazel test //book/ch02-code-project-organization/01-variable-shadowing:variableshadowing_test
```

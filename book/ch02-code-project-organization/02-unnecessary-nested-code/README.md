# Unnecessary Nested Code

## Problem

Deep nesting makes code hard to read and maintain. It increases cognitive load and makes the happy path unclear.

```go
// Bad: Deep nesting
if data != nil {
    if val, ok := data["user"]; ok {
        if user, ok := val.(map[string]interface{}); ok {
            if name, ok := user["name"]; ok {
                // ... more nesting
            }
        }
    }
}
```

## Solutions

### 1. Use Early Returns

Return early for error cases, keeping the happy path aligned to the left:

```go
// Good: Early returns, happy path on left
if data == nil {
    return "", errors.New("data is nil")
}

val, ok := data["user"]
if !ok {
    return "", errors.New("user field not found")
}

// Happy path continues...
return nameStr, nil
```

### 2. Use Continue in Loops

For loops, use `continue` to skip invalid items early:

```go
for _, user := range users {
    if user == nil {
        continue
    }

    if !isActive(user) {
        continue
    }

    // Process valid user
}
```

### 3. Eliminate Unnecessary Else

When using early returns, `else` blocks become unnecessary:

```go
// Bad: Unnecessary else
if value > 0 {
    return "positive"
} else {
    return "negative"
}

// Good: No else needed
if value > 0 {
    return "positive"
}
return "negative"
```

## Key Principles

1. **Happy path on the left**: Main logic should be left-aligned
2. **Error cases in if blocks**: Handle errors early with returns
3. **Minimize nesting depth**: Keep code flat and readable
4. **Avoid else after return**: Early returns make else unnecessary

## Run Tests

```bash
bazel test //book/ch02-code-project-organization/02-unnecessary-nested-code:unnecessarynested_test
```

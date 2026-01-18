# Unnecessary Nested Code

## Problem

Deep nesting with unnecessary `else` blocks makes code hard to read and maintain. Each level of nesting increases cognitive load.

### Bad Example: Deep Nesting with Else

```go
func join1(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    } else {
        if s2 == "" {
            return "", errors.New("s2 is empty")
        } else {
            concat, err := concatenate(s1, s2)
            if err != nil {
                return "", err
            } else {
                if len(concat) > max {
                    return concat[:max], nil
                } else {
                    return concat, nil
                }
            }
        }
    }
}
```

**Problems:**
- Deep nesting (4 levels)
- Unnecessary `else` blocks after `return`
- Happy path buried deep inside
- Hard to follow logic flow

## Solution: Early Returns and Flat Structure

### Good Example: Early Returns, No Else

```go
func join2(s1, s2 string, max int) (string, error) {
    if s1 == "" {
        return "", errors.New("s1 is empty")
    }
    if s2 == "" {
        return "", errors.New("s2 is empty")
    }
    concat, err := concatenate(s1, s2)
    if err != nil {
        return "", err
    }
    if len(concat) > max {
        return concat[:max], nil
    }
    return concat, nil
}
```

**Benefits:**
- Flat structure (no nesting)
- No unnecessary `else` blocks
- Happy path clearly visible at the end
- Easy to read and understand

## Key Principles

1. **Use early returns**: Return immediately when error conditions are met
2. **Eliminate else after return**: When you return in an `if` block, `else` is never needed
3. **Happy path on the left**: Main logic should be left-aligned and at the bottom
4. **Guard clauses first**: Validate inputs at the top with early returns

## Additional Patterns

### Use Continue in Loops

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

### Avoid Deep Nesting in Validation

```go
// Bad: nested validation
if valid {
    if authorized {
        if hasPermission {
            // do something
        }
    }
}

// Good: early returns
if !valid {
    return errors.New("invalid")
}
if !authorized {
    return errors.New("unauthorized")
}
if !hasPermission {
    return errors.New("no permission")
}
// do something
```

## Run Tests

```bash
bazel test //book/ch02-code-project-organization/02-unnecessary-nested-code:unnecessarynested_test
```

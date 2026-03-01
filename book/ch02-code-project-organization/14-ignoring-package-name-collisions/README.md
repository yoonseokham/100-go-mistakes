# Ignoring Package Name Collisions (#14)

## TL;DR

To avoid naming collisions between variables and packages, use unique names for each one. If that's not feasible, use an import alias — but only as a last resort.

---

## The Problem

When a variable name matches a package name, the package becomes inaccessible:

```go
import "redis"

func process() {
    redis := redis.NewClient()  // ❌ redis is now a variable, not the package
    redis.Ping()                // refers to the variable

    // Can no longer call redis.NewClient() below
    // "redis" now means the variable, not the package
}
```

---

## Solutions (In Order of Preference)

### 1. Rename the Variable (Best)

The simplest and most readable fix:

```go
redisClient := redis.NewClient()   // clear and specific
```

### 2. Use a More Descriptive Name

```go
customerCache := redis.NewClient()  // even better — says what it's for
```

### 3. Import Alias (Last Resort)

```go
import redisapi "redis"

redis := redisapi.NewClient()
```

**Why last resort?**
- Readers see `redisapi` and wonder "what is this?"
- Inconsistency across files — some use `redis`, others `redisapi`
- Adds cognitive load for no real benefit

---

## When Import Alias Is OK

When two packages genuinely have the same name — no other choice:

```go
import (
    "crypto/rand"
    mrand "math/rand"  // unavoidable — both are "rand"
)
```

---

## Common Collision Examples

```go
// ❌ Collides with context package
context := req.Context()
// ✅
ctx := req.Context()

// ❌ Collides with http package
http := doRequest()
// ✅
resp := doRequest()

// ❌ Collides with json package
json := marshal(data)
// ✅
data := marshal(data)

// ❌ Collides with sql package
sql := buildQuery()
// ✅
query := buildQuery()
```

---

## Key Principle

Most collisions are avoidable by choosing better variable names. Import aliases should only be used when two imported packages have the same name and there's no alternative.

| Priority | Approach | When to Use |
|----------|----------|-------------|
| 1st | Rename variable | Always try this first |
| 2nd | More descriptive name | When context helps clarity |
| 3rd | Import alias | Only when two packages share the same name |

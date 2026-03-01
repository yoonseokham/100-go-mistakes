# Creating Utility Packages (#13)

## TL;DR

Naming is a critical piece of application design. Creating packages such as `common`, `util`, and `shared` doesn't bring much value for the reader. Refactor such packages into meaningful and specific package names.

---

## The Problem

### Bad: Generic Package Names

```
project/
  utils/
    string.go     ← string helpers
    time.go       ← time helpers
    convert.go    ← conversion helpers
  common/
    errors.go
    constants.go
  shared/
    types.go
```

**Usage:**
```go
utils.NewStringSet("a", "b")     // what does "utils" tell you?
utils.FormatTime(t)              // utils does everything?
common.ErrNotFound               // common = junk drawer
```

**Problems:**
- Package name carries no meaning — "utils" could contain anything
- Unrelated code gets dumped into the same package
- Becomes a junk drawer that grows endlessly
- Callers can't understand the package's purpose at a glance

---

## The Solution

### Good: Name After What It Provides

```go
package stringset

type Set map[string]struct{}

func New(values ...string) Set { ... }
func (s Set) Sort() []string { ... }
```

**Usage:**
```go
stringset.New("a", "b", "c")   // immediately clear
```

### More Examples

```
❌ utils.FormatTime()      → ✅ timeutil.Format()
❌ utils.ParseJSON()       → ✅ jsonparser.Parse()
❌ common.ErrNotFound      → ✅ customerr.ErrNotFound
❌ helpers.Validate()      → ✅ validation.Run()
```

---

## The Standard Library Does This

Go's stdlib follows this principle:

| Package | What it provides |
|---------|-----------------|
| `strings` | String manipulation |
| `strconv` | String conversion |
| `net/http` | HTTP client/server |
| `encoding/json` | JSON encoding/decoding |
| `net/http/httputil` | HTTP utilities (not just "utils") |

Even `httputil` is scoped — it's HTTP utilities, not generic utilities.

---

## Key Principle

Name packages after **what they provide**, not **what they contain**:

```
❌ What they contain       ✅ What they provide
  utils                      stringset
  helpers                    httputil
  common                     auth
  shared                     encoding
  base                       storage
  misc                       validation
```

Package names should be **short, concise, expressive, and a single lowercase word**.

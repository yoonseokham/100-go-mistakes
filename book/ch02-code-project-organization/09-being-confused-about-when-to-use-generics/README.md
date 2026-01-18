# Being Confused About When to Use Generics

## TL;DR

Relying on generics and type parameters can prevent writing boilerplate code to factor out elements or behaviors. However, do not use type parameters prematurely, but only when you see a concrete need for them. Otherwise, they introduce unnecessary abstractions and complexity.

## The Problem Before Generics

Consider a function that extracts all keys from a `map[string]int`:

```go
func getKeys(m map[string]int) []string {
    var keys []string
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

What if we want to use this for `map[int]string`? Before generics, we had limited options:
- Code generation
- Reflection
- Code duplication

### Using `any` (interface{})

```go
func getKeys(m any) ([]any, error) {
    switch t := m.(type) {
    default:
        return nil, fmt.Errorf("unknown type: %T", t)
    case map[string]int:
        var keys []any
        for k := range t {
            keys = append(keys, k)
        }
        return keys, nil
    case map[int]string:
        // Copy the extraction logic
    }
}
```

**Problems:**
- Increases boilerplate code (duplicating the range loop for each case)
- Type checking happens at runtime instead of compile-time
- Returns `[]any`, losing type safety
- Requires error handling for unknown types

## Solution: Type Parameters

Thanks to generics, we can refactor using type parameters:

```go
func getKeys[K comparable, V any](m map[K]V) []K {
    var keys []K
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

**Type parameters:**
- `K comparable` - Map keys must be comparable (can use `==` or `!=`)
- `V any` - Values can be any type
- Instantiation happens at compile time, maintaining type safety

**Usage:**
```go
m := map[string]int{"one": 1, "two": 2, "three": 3}
keys := getKeys(m)  // Go infers the type arguments
```

## Constraints

A constraint is an interface type that can contain:
- A set of behaviors (methods)
- Arbitrary types

### Custom Constraint Example

```go
type customConstraint interface {
   ~int | ~string
}

func getKeys[K customConstraint, V any](m map[K]V) []K {
   // Same implementation
}
```

This restricts the key type to either `int` or `string` types.

### Understanding `~int` vs `int`

- `int` - Restricts to exactly that type
- `~int` - Restricts to all types whose underlying type is `int`

**Example:**

```go
type customConstraint interface {
   ~int
   String() string
}
```

This constraint requires:
1. Underlying type is `int`
2. Implements `String() string` method

```go
type customInt int

func (i customInt) String() string {
   return strconv.Itoa(int(i))
}
```

`customInt` satisfies this constraint because its underlying type is `int` and it implements `String()`.

**Note:** Using `int` instead of `~int` in the constraint would fail because the base `int` type doesn't implement `String()`.

## Generics with Data Structures

```go
type Node[T any] struct {
    Val  T
    next *Node[T]
}

func (n *Node[T]) Add(next *Node[T]) {
    n.next = next
}
```

**Note:** Type parameters cannot be used on methods, only on functions.

```go
type Foo struct {}

func (Foo) bar[T any](t T) {}  // ❌ Compilation error

// Error: methods cannot have type parameters
```

## Common Uses

Use generics when:

1. **Data structures** - Binary tree, linked list, heap, etc.
2. **Functions working with slices, maps, channels of any type**

```go
func merge[T any](ch1, ch2 <-chan T) <-chan T {
    // ...
}
```

3. **Factoring out behaviors**

```go
type sliceFn[T any] struct {
   s       []T
   compare func(T, T) bool
}

func (s sliceFn[T]) Len() int           { return len(s.s) }
func (s sliceFn[T]) Less(i, j int) bool { return s.compare(s.s[i], s.s[j]) }
func (s sliceFn[T]) Swap(i, j int)      { s.s[i], s.s[j] = s.s[j], s.s[i] }
```

## When NOT to Use Generics

Don't use generics when:

1. **Just calling a method of the type argument**

```go
func foo[T io.Writer](w T) {
   b := getBytes()
   _, _ = w.Write(b)
}
```

In this case, just use the interface directly:

```go
func foo(w io.Writer) {
   b := getBytes()
   _, _ = w.Write(b)
}
```

2. **When it makes code more complex**

Generics are never mandatory. Go developers lived without them for over a decade. If generics don't make code clearer, reconsider using them.

## Summary

We can find similarities with when not to use interfaces. Indeed, generics introduce a form of abstraction, and unnecessary abstractions introduce complexity.

**Key principles:**
- Don't pollute code with needless abstractions
- Focus on solving concrete problems
- Don't use type parameters prematurely
- Wait until you are about to write boilerplate code before considering generics

In general, we should avoid overgeneralizing the code we write at all costs. Perhaps a little bit of duplicated code might occasionally be better if it improves other aspects such as code expressiveness.

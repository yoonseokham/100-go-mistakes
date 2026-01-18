# Overusing Getters and Setters

## TL;DR

Forcing the use of getters and setters isn't idiomatic in Go. Being pragmatic and finding the right balance between efficiency and blindly following certain idioms should be the way to go.

## What is Data Encapsulation?

Data encapsulation refers to hiding the values or state of an object. Getters and setters are means to enable encapsulation by providing exported methods on top of unexported object fields.

### Example from Other Languages

In languages like Java or C#, this pattern is common:

```java
// Java style
class Person {
    private String name;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
```

## Getters and Setters in Go

### No Automatic Support

In Go, there is **no automatic support** for getters and setters as we see in some languages.

### Not Mandatory or Idiomatic

It is considered **neither mandatory nor idiomatic** to use getters and setters to access struct fields.

### Bad: Unnecessary Getters/Setters

```go
type Person struct {
    name string
}

// Unnecessary: just wrapping field access
func (p *Person) GetName() string {
    return p.name
}

func (p *Person) SetName(name string) {
    p.name = name
}
```

### Good: Direct Field Access

```go
type Person struct {
    Name string  // Exported field, direct access
}

// Usage
person := Person{Name: "John"}
fmt.Println(person.Name)  // Direct access
person.Name = "Jane"      // Direct modification
```

## Key Principles

### Don't Overwhelm Code

We **shouldn't overwhelm our code** with getters and setters on structs if they don't bring any value.

### Be Pragmatic

We should be pragmatic and strive to find the **right balance** between:
- Efficiency
- Following idioms that are sometimes considered indisputable in other programming paradigms

### Go's Design Philosophy

Remember that Go is a **unique language** designed for many characteristics, including **simplicity**.

## When to Use Getters and Setters

However, if we find a **need** for getters and setters or, as mentioned, **foresee a future need** while guaranteeing **forward compatibility**, there's nothing wrong with using them.

### Legitimate Use Cases

Use getters and setters when they **provide actual value**:

- Adding validation logic
- Transforming data on access
- Computing derived values
- Ensuring forward compatibility for future changes
- Adding behavior beyond simple field access

## Summary

| Approach | When to Use |
|----------|-------------|
| Direct field access | Default choice - simple and idiomatic |
| Getters/Setters | Only when they provide actual value |

**Bottom line:** Don't blindly follow patterns from other languages. Be pragmatic and choose the approach that best fits Go's philosophy of simplicity while meeting your actual needs.

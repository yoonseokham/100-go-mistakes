# Go 100 Mistakes Examples

This project contains example code from "100 Go Mistakes and How to Avoid Them" book, built and tested with Bazel.

## Installation

Install Bazel:
```bash
# macOS
brew install bazelisk
```

## Usage

### Build
```bash
# Build all targets
bazel build //...

# Build specific chapter
bazel build //book/ch02-code-project-organization/01-variable-shadowing:variableshadowing
```

### Test
```bash
# Run all tests
bazel test //...

# Run specific test
bazel test //book/ch02-code-project-organization/01-variable-shadowing:variableshadowing_test
```

### Gazelle (Auto-generate BUILD files)
```bash
# Auto-generate/update BUILD.bazel files
bazel run //:gazelle

# Update external dependencies from go.mod
bazel run //:gazelle-update-repos
```

## Project Structure

```
.
├── WORKSPACE              # Bazel workspace configuration
├── BUILD.bazel           # Root BUILD file
├── go.mod                # Go module definition
├── .bazelrc              # Bazel configuration
└── book/
    └── ch02-code-project-organization/
        └── 01-variable-shadowing/
            ├── BUILD.bazel      # Build rules
            ├── shadowing.go     # Example code
            ├── shadowing_test.go # Tests
            └── README.md        # Chapter notes
```

## Chapters

### Chapter 2: Code and Project Organization
- [01 - Variable Shadowing](book/ch02-code-project-organization/01-variable-shadowing/) - Avoiding unintended variable shadowing
- [02 - Unnecessary Nested Code](book/ch02-code-project-organization/02-unnecessary-nested-code/) - Using early returns and keeping happy path left-aligned
- [03 - Misusing Init Functions](book/ch02-code-project-organization/03-misusing-init-functions/) - Understanding init() problems and using explicit initialization
- [04 - Overusing Getters and Setters](book/ch02-code-project-organization/04-overusing-getters-setters/) - Being pragmatic about encapsulation in Go
- [05 - Interface Pollution](book/ch02-code-project-organization/05-interface-pollution/) - Discovering abstractions instead of creating them
- [06 - Interface on the Producer Side](book/ch02-code-project-organization/06-interface-producer-side/) - Keeping interfaces on the client side
- [07 - Returning Interfaces](book/ch02-code-project-organization/07-returning-interfaces/) - Returning concrete types, accepting interfaces
- [08 - Any Says Nothing](book/ch02-code-project-organization/08-any-says-nothing/) - Using specific types instead of any for type safety
- [09 - Being Confused About When to Use Generics](book/ch02-code-project-organization/09-being-confused-about-when-to-use-generics/) - Using generics appropriately to reduce boilerplate without adding complexity
- [10 - Type Embedding Problems](book/ch02-code-project-organization/10-type-embedding-problems/) - Understanding visibility issues with type embedding

## Bazel Build Rules

- `go_library`: Go package/library
- `go_binary`: Executable binary
- `go_test`: Go tests

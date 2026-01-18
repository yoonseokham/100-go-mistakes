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

## Bazel Build Rules

- `go_library`: Go package/library
- `go_binary`: Executable binary
- `go_test`: Go tests

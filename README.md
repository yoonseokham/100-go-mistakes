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

# Build specific example
bazel build //examples/mistake01:example
```

### Run
```bash
# Run binary
bazel run //examples/mistake01:example
```

### Test
```bash
# Run all tests
bazel test //...

# Run specific test
bazel test //examples/mistake01:mistake01_test
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
└── examples/
    └── mistake01/
        ├── BUILD.bazel   # Example-specific build rules
        ├── lib.go        # Library code
        ├── lib_test.go   # Test code
        └── main.go       # Executable example
```

## Bazel Build Rules

- `go_library`: Go package/library
- `go_binary`: Executable binary
- `go_test`: Go tests

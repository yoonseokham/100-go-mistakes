package main

import (
	"io"
	"os"
	"sync"
)

// Basic embedding example
type Foo struct {
	Bar
}

type Bar struct {
	Baz int
}

func basicExample() {
	foo := Foo{}
	foo.Baz = 42 // Bar.Baz is promoted to Foo
}

// Bad example: Embedding sync.Mutex exposes Lock/Unlock
type InMemBad struct {
	sync.Mutex // ❌ Embedded - Lock/Unlock become public
	m          map[string]int
}

func NewInMemBad() *InMemBad {
	return &InMemBad{m: make(map[string]int)}
}

func (i *InMemBad) Get(key string) (int, bool) {
	i.Lock()
	v, contains := i.m[key]
	i.Unlock()
	return v, contains
}

// Good example: Using named field for sync.Mutex
type InMemGood struct {
	mu sync.Mutex // ✅ Named field - not promoted
	m  map[string]int
}

func NewInMemGood() *InMemGood {
	return &InMemGood{m: make(map[string]int)}
}

func (i *InMemGood) Get(key string) (int, bool) {
	i.mu.Lock()
	v, contains := i.m[key]
	i.mu.Unlock()
	return v, contains
}

// Bad example: Forwarding methods (boilerplate)
type LoggerBad struct {
	writeCloser io.WriteCloser
}

func (l LoggerBad) Write(p []byte) (int, error) {
	return l.writeCloser.Write(p)
}

func (l LoggerBad) Close() error {
	return l.writeCloser.Close()
}

// Good example: Embedding to avoid forwarding methods
type LoggerGood struct {
	io.WriteCloser // ✅ Embedded - Write/Close promoted
}

func main() {
	// Basic example
	basicExample()

	// Bad InMem - exposes Lock/Unlock
	inMemBad := NewInMemBad()
	inMemBad.Lock()   // ❌ This shouldn't be accessible!
	inMemBad.Unlock() // ❌ This shouldn't be accessible!

	// Good InMem - Lock/Unlock not accessible
	inMemGood := NewInMemGood()
	// inMemGood.Lock() // ✅ Compilation error - mu is private

	// Bad Logger - requires forwarding methods
	loggerBad := LoggerBad{writeCloser: os.Stdout}
	_, _ = loggerBad.Write([]byte("foo"))
	_ = loggerBad.Close()

	// Good Logger - methods promoted, no forwarding needed
	loggerGood := LoggerGood{WriteCloser: os.Stdout}
	_, _ = loggerGood.Write([]byte("bar")) // Write() promoted
	_ = loggerGood.Close()                  // Close() promoted
}

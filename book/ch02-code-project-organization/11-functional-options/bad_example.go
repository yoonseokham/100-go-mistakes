package functionaloptions

import (
	"errors"
	"net/http"
	"time"
)

// Bad Example 1: Parameter explosion
// Adding more options means changing the function signature

func NewServerV1(addr string, port int, timeout time.Duration, maxConn int) (*http.Server, error) {
	if port < 0 {
		return nil, errors.New("port should be positive")
	}
	_ = timeout
	_ = maxConn
	return &http.Server{Addr: addr}, nil
}

// Problems:
// 1. Adding a new option (e.g. TLS) changes the signature → breaks all callers
// 2. Callers must pass all arguments even if they don't care
// NewServerV1("localhost", 8080, 30*time.Second, 100)

// Bad Example 2: Config struct — empty struct is awkward
type Config struct {
	Port    int
	Timeout time.Duration
	MaxConn int
}

func NewServerV2(addr string, cfg Config) (*http.Server, error) {
	if cfg.Port < 0 {
		return nil, errors.New("port should be positive")
	}
	return &http.Server{Addr: addr}, nil
}

// Problems:
// 1. Can't distinguish "port not set" vs "port = 0"
// 2. Passing empty Config{} is awkward
// 3. No way to set defaults cleanly
// NewServerV2("localhost", Config{})           // what does this mean?
// NewServerV2("localhost", Config{Port: 0})    // random port? or default?

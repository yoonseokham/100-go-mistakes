package functionaloptions

import (
	"errors"
	"math/rand"
	"net/http"
	"time"
)

const defaultHTTPPort = 8080

// Step 1: unexported options struct holds all config
type options struct {
	port    *int          // pointer: nil = not set, 0 = random port
	timeout time.Duration
	maxConn int
}

// Step 2: Option is a function type that mutates options
type Option func(opts *options) error

// Step 3: WithXxx functions return an Option
// Validation happens here — not in NewServer

func WithPort(port int) Option {
	return func(opts *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		opts.port = &port
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) error {
		if timeout <= 0 {
			return errors.New("timeout should be positive")
		}
		opts.timeout = timeout
		return nil
	}
}

func WithMaxConn(maxConn int) Option {
	return func(opts *options) error {
		if maxConn <= 0 {
			return errors.New("maxConn should be positive")
		}
		opts.maxConn = maxConn
		return nil
	}
}

// Step 4: NewServer applies all options in order
func NewServer(addr string, opts ...Option) (*http.Server, error) {
	// start with zero value — all fields unset
	var o options

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	// Apply defaults after options
	var port int
	if o.port == nil {
		port = defaultHTTPPort // not set → use default
	} else if *o.port == 0 {
		port = rand.Intn(65535) // explicitly 0 → random port
	} else {
		port = *o.port
	}

	_ = port
	return &http.Server{Addr: addr}, nil
}

// Benefits:
// 1. NewServer signature never changes when adding new options
// 2. Callers only pass what they need
// 3. Validation is co-located with each option
// 4. nil vs 0 distinction thanks to *int pointer
//
// Usage:
// NewServer("localhost")                                     // all defaults
// NewServer("localhost", WithPort(8080))                    // custom port
// NewServer("localhost", WithPort(8080), WithTimeout(30*time.Second)) // multiple
// NewServer("localhost", WithPort(-1))                      // returns error

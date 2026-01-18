package variableshadowing

import "net/http"

// Problem: Variable shadowing causes client to be nil after if/else block
func BadExample(tracing bool) (*http.Client, error) {
	var client *http.Client

	if tracing {
		// := creates NEW variables in this scope, shadowing outer client
		client, err := createTracingClient()
		_ = client // prevent unused variable error
		if err != nil {
			return nil, err
		}
	} else {
		// := creates NEW variables in this scope, shadowing outer client
		client, err := createDefaultClient()
		_ = client // prevent unused variable error
		if err != nil {
			return nil, err
		}
	}

	// client is still nil here! The inner clients were shadowed.
	return client, nil
}

// Solution 1: Use = instead of := to assign to existing variables
func GoodExample1(tracing bool) (*http.Client, error) {
	var client *http.Client
	var err error

	if tracing {
		// Use = to assign to outer scope variables
		client, err = createTracingClient()
		if err != nil {
			return nil, err
		}
	} else {
		// Use = to assign to outer scope variables
		client, err = createDefaultClient()
		if err != nil {
			return nil, err
		}
	}

	// client is properly assigned now
	return client, nil
}

// Solution 2: Use temporary variable and assign after
func GoodExample2(tracing bool) (*http.Client, error) {
	var client *http.Client

	if tracing {
		// Use temporary variable c
		c, err := createTracingClient()
		if err != nil {
			return nil, err
		}
		client = c
	} else {
		// Use temporary variable c
		c, err := createDefaultClient()
		if err != nil {
			return nil, err
		}
		client = c
	}

	return client, nil
}

// Helper functions to simulate client creation
func createTracingClient() (*http.Client, error) {
	return &http.Client{}, nil
}

func createDefaultClient() (*http.Client, error) {
	return &http.Client{}, nil
}

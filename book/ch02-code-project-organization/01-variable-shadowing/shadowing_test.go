package variableshadowing

import "testing"

func TestBadExample(t *testing.T) {
	// This demonstrates the bug: client will be nil
	client, err := BadExample(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client != nil {
		t.Error("expected client to be nil due to shadowing, but it wasn't")
	}
}

func TestGoodExample1(t *testing.T) {
	tests := []struct {
		name    string
		tracing bool
	}{
		{"with tracing", true},
		{"without tracing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := GoodExample1(tt.tracing)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if client == nil {
				t.Error("expected client to be non-nil, but got nil")
			}
		})
	}
}

func TestGoodExample2(t *testing.T) {
	tests := []struct {
		name    string
		tracing bool
	}{
		{"with tracing", true},
		{"without tracing", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := GoodExample2(tt.tracing)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if client == nil {
				t.Error("expected client to be non-nil, but got nil")
			}
		})
	}
}

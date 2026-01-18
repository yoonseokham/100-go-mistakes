package unnecessarynested

import (
	"testing"
)

func TestBadExample(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "valid data",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
				},
			},
			want:    "John",
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing user field",
			data:    map[string]interface{}{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BadExample(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("BadExample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BadExample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoodExample(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "valid data",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
				},
			},
			want:    "John",
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			want:    "",
			wantErr: true,
		},
		{
			name:    "missing user field",
			data:    map[string]interface{}{},
			want:    "",
			wantErr: true,
		},
		{
			name: "empty name",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "",
				},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoodExample(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoodExample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoodExample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBadExampleWithLoop(t *testing.T) {
	users := []map[string]interface{}{
		{"name": "Alice", "active": true},
		{"name": "Bob", "active": false},
		{"name": "Charlie", "active": true},
		nil,
		{"name": "Dave"},
	}

	result := BadExampleWithLoop(users)
	expected := []string{"Alice", "Charlie"}

	if len(result) != len(expected) {
		t.Errorf("BadExampleWithLoop() length = %d, want %d", len(result), len(expected))
	}

	for i, name := range expected {
		if i >= len(result) || result[i] != name {
			t.Errorf("BadExampleWithLoop()[%d] = %v, want %v", i, result[i], name)
		}
	}
}

func TestGoodExampleWithLoop(t *testing.T) {
	users := []map[string]interface{}{
		{"name": "Alice", "active": true},
		{"name": "Bob", "active": false},
		{"name": "Charlie", "active": true},
		nil,
		{"name": "Dave"},
	}

	result := GoodExampleWithLoop(users)
	expected := []string{"Alice", "Charlie"}

	if len(result) != len(expected) {
		t.Errorf("GoodExampleWithLoop() length = %d, want %d", len(result), len(expected))
	}

	for i, name := range expected {
		if i >= len(result) || result[i] != name {
			t.Errorf("GoodExampleWithLoop()[%d] = %v, want %v", i, result[i], name)
		}
	}
}

func TestBadExampleElse(t *testing.T) {
	tests := []struct {
		value int
		want  string
	}{
		{5, "positive"},
		{-5, "negative"},
		{0, "zero"},
	}

	for _, tt := range tests {
		got := BadExampleElse(tt.value)
		if got != tt.want {
			t.Errorf("BadExampleElse(%d) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestGoodExampleElse(t *testing.T) {
	tests := []struct {
		value int
		want  string
	}{
		{5, "positive"},
		{-5, "negative"},
		{0, "zero"},
	}

	for _, tt := range tests {
		got := GoodExampleElse(tt.value)
		if got != tt.want {
			t.Errorf("GoodExampleElse(%d) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		amount  float64
		wantErr bool
	}{
		{"valid order", "ORD123", 100.50, false},
		{"empty orderID", "", 100.50, true},
		{"zero amount", "ORD123", 0, true},
		{"negative amount", "ORD123", -10, true},
		{"exceeds limit", "ORD123", 15000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessOrder(tt.orderID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

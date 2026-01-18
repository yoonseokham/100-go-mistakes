package interfacepollution

// Restricting Behavior Example: Read-only Configuration

// Full configuration with read and write
type IntConfig struct {
	value int
}

func (c *IntConfig) Get() int {
	return c.value
}

func (c *IntConfig) Set(value int) {
	c.value = value
}

// Interface that restricts to read-only
type IntConfigGetter interface {
	Get() int
	// Intentionally excludes Set
}

// Service that should only read configuration
type ThresholdService struct {
	threshold IntConfigGetter // Read-only access
}

func NewThresholdService(threshold IntConfigGetter) *ThresholdService {
	return &ThresholdService{threshold: threshold}
}

func (s *ThresholdService) IsAboveThreshold(value int) bool {
	threshold := s.threshold.Get() // Can read
	// s.threshold.Set(100)  // Compile error! Method not available
	return value > threshold
}

// Usage example
func ExampleRestrictingBehavior() {
	config := &IntConfig{}
	config.Set(100) // External code can set

	service := NewThresholdService(config) // Passed as IntConfigGetter
	// service can only Get(), not Set()
	service.IsAboveThreshold(150)
}

// Benefits:
// 1. Semantically enforces read-only access
// 2. Compile-time safety
// 3. Clear intent in code
// 4. Prevents accidental modification

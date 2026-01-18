package unnecessarynested

import (
	"errors"
	"fmt"
)

// BadExample: Deep nesting makes code hard to read
func BadExample(data map[string]interface{}) (string, error) {
	if data != nil {
		if val, ok := data["user"]; ok {
			if user, ok := val.(map[string]interface{}); ok {
				if name, ok := user["name"]; ok {
					if nameStr, ok := name.(string); ok {
						if nameStr != "" {
							return nameStr, nil
						} else {
							return "", errors.New("name is empty")
						}
					} else {
						return "", errors.New("name is not a string")
					}
				} else {
					return "", errors.New("name field not found")
				}
			} else {
				return "", errors.New("user is not a map")
			}
		} else {
			return "", errors.New("user field not found")
		}
	}
	return "", errors.New("data is nil")
}

// GoodExample: Use early returns to avoid nesting
// Keep the happy path aligned to the left, error cases in if blocks
func GoodExample(data map[string]interface{}) (string, error) {
	if data == nil {
		return "", errors.New("data is nil")
	}

	val, ok := data["user"]
	if !ok {
		return "", errors.New("user field not found")
	}

	user, ok := val.(map[string]interface{})
	if !ok {
		return "", errors.New("user is not a map")
	}

	name, ok := user["name"]
	if !ok {
		return "", errors.New("name field not found")
	}

	nameStr, ok := name.(string)
	if !ok {
		return "", errors.New("name is not a string")
	}

	if nameStr == "" {
		return "", errors.New("name is empty")
	}

	// Happy path: aligned to the left, easy to read
	return nameStr, nil
}

// BadExampleWithLoop: Nested loops and conditions
func BadExampleWithLoop(users []map[string]interface{}) []string {
	var activeUsers []string

	if users != nil {
		for _, user := range users {
			if user != nil {
				if active, ok := user["active"]; ok {
					if isActive, ok := active.(bool); ok {
						if isActive {
							if name, ok := user["name"]; ok {
								if nameStr, ok := name.(string); ok {
									activeUsers = append(activeUsers, nameStr)
								}
							}
						}
					}
				}
			}
		}
	}

	return activeUsers
}

// GoodExampleWithLoop: Use continue for early loop iteration exit
func GoodExampleWithLoop(users []map[string]interface{}) []string {
	var activeUsers []string

	if users == nil {
		return activeUsers
	}

	for _, user := range users {
		if user == nil {
			continue
		}

		active, ok := user["active"]
		if !ok {
			continue
		}

		isActive, ok := active.(bool)
		if !ok || !isActive {
			continue
		}

		name, ok := user["name"]
		if !ok {
			continue
		}

		nameStr, ok := name.(string)
		if !ok {
			continue
		}

		// Happy path: clean and readable
		activeUsers = append(activeUsers, nameStr)
	}

	return activeUsers
}

// BadExampleElse: Unnecessary else blocks
func BadExampleElse(value int) string {
	if value > 0 {
		return "positive"
	} else {
		if value < 0 {
			return "negative"
		} else {
			return "zero"
		}
	}
}

// GoodExampleElse: Eliminate else blocks with early returns
func GoodExampleElse(value int) string {
	if value > 0 {
		return "positive"
	}

	if value < 0 {
		return "negative"
	}

	return "zero"
}

// Example: Processing with validation
func ProcessOrder(orderID string, amount float64) error {
	// Bad: nested validation
	badProcess := func() error {
		if orderID != "" {
			if amount > 0 {
				if amount <= 10000 {
					fmt.Printf("Processing order %s for $%.2f\n", orderID, amount)
					return nil
				} else {
					return errors.New("amount exceeds limit")
				}
			} else {
				return errors.New("amount must be positive")
			}
		}
		return errors.New("orderID is required")
	}
	_ = badProcess

	// Good: early returns for validation
	if orderID == "" {
		return errors.New("orderID is required")
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if amount > 10000 {
		return errors.New("amount exceeds limit")
	}

	// Happy path: business logic
	fmt.Printf("Processing order %s for $%.2f\n", orderID, amount)
	return nil
}

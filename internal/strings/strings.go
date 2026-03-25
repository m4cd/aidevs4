package strings

import (
	"encoding/json"
	"fmt"
)

// Unmarshalls JSON stored from string
func UnmarshalJSON[T any](s string) (*T, error) {

	var result *T

	unmarshalError := json.Unmarshal([]byte(s), &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from string: %w", unmarshalError)
	}

	return result, nil
}

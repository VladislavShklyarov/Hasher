package utils

import (
	"encoding/json"
	"fmt"
)

type FlexString string

func (fs *FlexString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("failed to parse string: %w", err)
		}
		*fs = FlexString(s)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("failed to parse number: %w", err)
	}
	*fs = FlexString(num.String())
	return nil
}

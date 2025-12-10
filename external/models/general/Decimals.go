package general

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Decimals string

func (d *Decimals) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*d = Decimals(s)
		return nil
	}

	// If that fails, try as int
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*d = Decimals(fmt.Sprintf("%d", i))
		return nil
	}

	return fmt.Errorf("decimals must be string or int")
}

func (d Decimals) Int() (int, error) {
	return strconv.Atoi(string(d))
}

// Helper method to get the value as a string
func (d Decimals) String() string {
	return string(d)
}

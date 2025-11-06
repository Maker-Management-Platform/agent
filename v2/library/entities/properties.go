package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Properties map[string]any

func (n *Properties) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON string:", src))
	}
	a := json.Unmarshal([]byte(str), &n)
	return a
}

func (n Properties) Value() (driver.Value, error) {
	val, err := json.Marshal(n)
	return string(val), err
}

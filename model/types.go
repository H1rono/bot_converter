package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type StringSlice []string

func (s *StringSlice) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("failed to unmarshal JSON value: %v", src))
	}
	return json.Unmarshal(bytes, &s)
}

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

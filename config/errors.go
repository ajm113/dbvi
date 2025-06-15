package config

import (
	"errors"
	"fmt"
)

type ConnectionValidationError struct {
	Field string
	Desc  string
}

func (e *ConnectionValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Desc, e.Field)
}

var (
	ErrConfigNotFound = errors.New("config not found")
)

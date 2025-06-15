package config

import (
	"slices"
)

var ConnectionTypes = []string{
	"postgres",
	"mysql",
	"redis",
}

type Connection struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
	ReadOnly bool   `yaml:"read_only"`
}

func validateConnection(c Connection) error {
	if c.Name == "" {
		return &ConnectionValidationError{Field: "name", Desc: "missing field"}
	}

	if !slices.Contains(ConnectionTypes, c.Type) {
		return &ConnectionValidationError{Field: "type", Desc: "invalid: " + c.Type}
	}

	if c.Host == "" {
		return &ConnectionValidationError{Field: "host", Desc: "missing field"}
	}

	// TODO: Do more validation here based on the type provided.
	return nil
}

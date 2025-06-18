package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Connections          []Connection
	UseConnection        string
	HasUnmaskedPasswords bool
}

type config struct {
	Connections   []Connection `yaml:"connections"`
	UseConnection string       `yaml:"use_connection"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// check for any YAML syntax errors.
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	// we want strict decoding to make yaml mistakes clear
	// so the user doesn't have to guess why things aren't working.
	cfg := &config{}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // THIS causes unknown fields to raise errors

	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	hasUnmaskedPasswords := false
	// Do in-depth validation of the config.
	for i, c := range cfg.Connections {
		err := validateConnection(c)

		if err != nil && !errors.Is(err, ErrPasswordUnmasked) {
			return nil, fmt.Errorf("error at connections[%d]: %w", i, err)
		}

		if errors.Is(err, ErrPasswordUnmasked) {
			hasUnmaskedPasswords = true
		}
	}

	return &Config{
		Connections:          cfg.Connections,
		UseConnection:        cfg.UseConnection,
		HasUnmaskedPasswords: hasUnmaskedPasswords,
	}, nil
}

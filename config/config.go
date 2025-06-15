package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
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
		fmt.Println("YAML syntax error:", err)
		return nil, err
	}

	// we want strict decoding to make yaml mistakes clear
	// so the user doesn't have to guess why things aren't working.
	cfg := &Config{}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true) // THIS causes unknown fields to raise errors

	if err := dec.Decode(&cfg); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return nil, err
	}

	// Do in-depth validation of the config.
	for i, c := range cfg.Connections {
		err := validateConnection(c)

		if err != nil {
			return nil, fmt.Errorf("error at connections[%d]: %w", i, err)
		}
	}

	return cfg, nil
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMaskPasswords(t *testing.T) {
	inputFile := filepath.Join("testdata", "valid_connection.yaml")
	outputFile := filepath.Join("testdata", "valid_connection_masked.yaml")

	err := MaskPasswords(inputFile, outputFile)
	if err != nil {
		t.Fatalf("failed masking yaml file %v", err)
	}

	config, err := Load(outputFile)
	if err != nil {
		t.Fatalf("failed loading masked yaml file %v", err)
	}

	for i, c := range config.Connections {
		if c.Password != MaskStr {
			t.Errorf("password for [%d].%s isn't masked", i, c.Name)
		}
	}

	if !t.Failed() {
		os.Remove(outputFile)
	}
}

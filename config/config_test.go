package config

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		file string
		c    *Config
		err  error
	}{
		{
			file: "blank.yaml",
			c:    nil,
			err:  io.EOF,
		},
		{
			file: "malformed_schema.yaml",
			c:    nil,
			err:  io.EOF,
		},
		{
			file: "valid_connection.yaml",
			c:    &Config{},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test config: %s", tt.file), func(t *testing.T) {
			c, err := Load(filepath.Join("testdata", tt.file))

			if tt.c == nil && c != nil {
				t.Errorf("expected no config, but got %v", c)
				return
			}

			if tt.c != nil {
				if !reflect.DeepEqual(c, tt.c) {
					t.Errorf("got %+v, want %+v", c, tt.c)
				}
			}

			if !errors.Is(err, tt.err) {
				t.Errorf("got %+v, want %+v", err, tt.err)
			}
		})
	}

}

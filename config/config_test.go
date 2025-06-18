package config

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		file          string
		c             *Config
		err           error
		wantTypeError bool
	}{
		{
			file: "blank.yaml",
			c:    nil,
			err:  io.EOF,
		},
		{
			file:          "malformed_schema.yaml",
			c:             nil,
			wantTypeError: true,
		},
		{
			file: "valid_connection.yaml",
			c: &Config{
				Connections: []Connection{
					{
						Name:     "Test",
						Type:     "postgres",
						Host:     "localhost",
						Port:     5432,
						Database: "postgres",
						ReadOnly: false,
						Username: "test",
						Password: "test",
					},
				},
				UseConnection:        "test",
				HasUnmaskedPasswords: true,
			},
			err: nil,
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

			var typeErr *yaml.TypeError
			gotTypeError := errors.As(err, &typeErr)

			if gotTypeError != tt.wantTypeError {
				t.Errorf("expected TypeError=%v, got %v (err=%v)", tt.wantTypeError, gotTypeError, err)
			}

			if !errors.Is(err, tt.err) && !tt.wantTypeError {
				t.Errorf("got %T, want %T", err, tt.err)
			}
		})
	}

}

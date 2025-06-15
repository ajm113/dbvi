package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateConnection(t *testing.T) {
	tests := []struct {
		c    Connection
		want error
	}{
		{
			c: Connection{
				Name: "Postgres",
				Type: "postgres",
				Host: "localhost",
				Port: 5432,
			},
			want: nil,
		},
		{
			c: Connection{
				Name: "MySQL",
				Type: "mysql",
				Host: "localhost",
				Port: 3306,
			},
			want: nil,
		},
		{
			c: Connection{
				Name: "Redis",
				Type: "redis",
				Host: "localhost",
				Port: 6379,
			},
			want: nil,
		},
		{
			c: Connection{
				Name: "Mongo - Invalid Type",
				Type: "mongo",
				Host: "localhost",
			},
			want: &ConnectionValidationError{Field: "type", Desc: "invalid: mongo"},
		},
		{
			c: Connection{
				Name: "",
				Type: "postgres",
				Host: "localhost",
			},
			want: &ConnectionValidationError{Field: "name", Desc: "missing field"},
		},
		{
			c: Connection{
				Name: "No Host",
				Type: "postgres",
				Host: "",
			},
			want: &ConnectionValidationError{Field: "host", Desc: "missing field"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test connection: %s:%s", tt.c.Name, tt.c.Type), func(t *testing.T) {
			err := validateConnection(tt.c)

			if tt.want == nil && err != nil {
				t.Errorf("expected no error, but got %s", err)
				return
			}

			if tt.want != nil {
				var errTyped *ConnectionValidationError
				if !errors.As(err, &errTyped) {
					t.Errorf("got \"%s\", want \"%s\"", err, tt.want)
					return
				}

				if !reflect.DeepEqual(errTyped, tt.want) {
					t.Errorf("got \"%s\", want \"%s\"", err, tt.want)
				}
			}
		})
	}
}

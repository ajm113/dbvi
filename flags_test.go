package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFlags(t *testing.T) {
	tests := []struct {
		args      []string
		wantArgs  *cliArguments
		wantError bool
	}{
		{
			args:      []string{"test"},
			wantArgs:  &cliArguments{},
			wantError: false,
		},
		{
			args:      []string{"test", "-h"},
			wantArgs:  &cliArguments{Commands: []cliCommand{{Name: "help"}}},
			wantError: false,
		},
		{
			args:      []string{"test", "--help"},
			wantArgs:  &cliArguments{Commands: []cliCommand{{Name: "help"}}},
			wantError: false,
		},
		{
			args:      []string{"test", "--help", "file.sql"},
			wantArgs:  &cliArguments{Commands: []cliCommand{{Name: "help"}}, Files: []string{"file.sql"}},
			wantError: false,
		},
		{
			args:      []string{"test", "--help", "file.sql"},
			wantArgs:  &cliArguments{Commands: []cliCommand{{Name: "help"}}, Files: []string{"file.sql"}},
			wantError: false,
		},
		{
			args:      []string{"test", "--", "file.sql", "file_no_extension"},
			wantArgs:  &cliArguments{Files: []string{"file.sql", "file_no_extension"}},
			wantError: false,
		},
		{
			args:      []string{"test", "-F", "invalid_argument_option"},
			wantArgs:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test flags: %v", tt.args), func(t *testing.T) {
			cliArgs, err := parseFlags(tt.args)

			if tt.wantArgs == nil && cliArgs != nil {
				t.Errorf("expected no cliArgs, but got %v", cliArgs)
				return
			}

			if tt.wantArgs != nil {
				if !reflect.DeepEqual(cliArgs, tt.wantArgs) {
					t.Errorf("got %+v, want %+v", cliArgs, tt.wantArgs)
				}
			}

			if err == nil && tt.wantError {
				t.Errorf("got %T, but wasn't expecting an error", err)
			}
		})
	}

}

package config

import (
	"testing"
)

// Test version validation
func TestVersionValidation(t *testing.T) {
	tests := []struct {
		version int
		valid   bool
	}{
		{1, true},
		{0, false},
		{2, false},
		{-1, false},
	}

	for _, tt := range tests {
		config := &Config{
			Version:  tt.version,
			Commands: []Command{{Name: "test", Script: "echo test"}},
		}

		err := validateConfig(config)

		if tt.valid && err != nil {
			t.Errorf("Version %d: expected valid, got error: %v", tt.version, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Version %d: expected invalid, got no error", tt.version)
		}
	}
}

// Test variable name validation
func TestVariableNameValidation(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"valid_name", true},
		{"validName", true},
		{"valid-name", true},
		{"", false},
		{"invalid name", false},
		{"123invalid", false},
		{"_invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variable := Variable{Name: tt.name, Value: "test"}
			err := validateVariable(variable)

			if tt.valid && err != nil {
				t.Errorf("Expected '%s' to be valid, got error: %v", tt.name, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected '%s' to be invalid, got no error", tt.name)
			}
		})
	}
}

// Test option type validation
func TestOptionTypeValidation(t *testing.T) {
	validTestTypes := []string{"bool", "str", "int", "float"}
	invalidTestTypes := []string{"string", "boolean", "number", "invalid", ""}

	for _, optType := range validTestTypes {
		t.Run("Valid: "+optType, func(t *testing.T) {
			opt := Option{Name: "test", Type: optType}
			err := validateOption(opt)
			if err != nil {
				t.Errorf("Expected type '%s' to be valid, got error: %v", optType, err)
			}
		})
	}

	for _, optType := range invalidTestTypes {
		t.Run("Invalid: "+optType, func(t *testing.T) {
			opt := Option{Name: "test", Type: optType}
			err := validateOption(opt)
			if err == nil {
				t.Errorf("Expected type '%s' to be invalid, got no error", optType)
			}
		})
	}
}

// Test shorthand validation
func TestShorthandValidation(t *testing.T) {
	tests := []struct {
		shorthand string
		valid     bool
	}{
		{"v", true},
		{"x", true},
		{"A", true},
		{"", true}, // empty is ok (optional)
		{"vv", false},
		{"1", false},
		{"-", false},
		{"h", false}, // reserved
		{"i", false}, // reserved
	}

	for _, tt := range tests {
		t.Run("Shorthand: "+tt.shorthand, func(t *testing.T) {
			opt := Option{Name: "test", Type: "bool", Shorthand: tt.shorthand}
			err := validateOption(opt)

			if tt.valid && err != nil {
				t.Errorf("Expected shorthand '%s' to be valid, got error: %v", tt.shorthand, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected shorthand '%s' to be invalid, got no error", tt.shorthand)
			}
		})
	}
}

// Test command name validation
func TestCommandNameValidation(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"deploy", true},
		{"build-app", true},
		{"test_unit", true},
		{"", false},
		{"invalid name", false},
		{"invalid@name", false},
		{"123invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command{Name: tt.name, Script: "echo test"}
			err := validateCommand(cmd)

			if tt.valid && err != nil {
				t.Errorf("Expected name '%s' to be valid, got error: %v", tt.name, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected name '%s' to be invalid, got no error", tt.name)
			}
		})
	}
}

// Test duplicate detection
func TestDuplicateCommandNames(t *testing.T) {
	config := &Config{
		Version: 1,
		Commands: []Command{
			{Name: "deploy", Script: "echo 1"},
			{Name: "deploy", Script: "echo 2"},
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Expected error for duplicate command names")
	}
}

func TestDuplicateAliases(t *testing.T) {
	config := &Config{
		Version: 1,
		Commands: []Command{
			{Name: "deploy", Aliases: []string{"d"}, Script: "echo 1"},
			{Name: "delete", Aliases: []string{"d"}, Script: "echo 2"},
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Expected error for duplicate aliases")
	}
}

func TestDuplicateShorthands(t *testing.T) {
	cmd := Command{
		Name: "test",
		Options: []Option{
			{Name: "verbose", Shorthand: "v", Type: "bool"},
			{Name: "version", Shorthand: "v", Type: "bool"},
		},
		Script: "echo test",
	}

	err := validateCommand(cmd)
	if err == nil {
		t.Error("Expected error for duplicate shorthands")
	}
}

package config

import (
	"strings"
	"testing"
)

// Test valid configurations
func TestLoadValidConfigs(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"Complete config", "testdata/valid/complete.yaml"},
		{"Minimal config", "testdata/valid/minimal.yaml"},
		{"All features", "testdata/valid/with_all_features.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Load(tt.filename)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if config == nil {
				t.Error("Expected config to be loaded, got nil")
			}
		})
	}
}

// Test invalid configurations
func TestLoadInvalidConfigs(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectError string
	}{
		{
			name:        "Missing version",
			filename:    "testdata/invalid/missing_version.yaml",
			expectError: "version",
		},
		{
			name:        "Wrong version",
			filename:    "testdata/invalid/wrong_version.yaml",
			expectError: "unsupported config version",
		},
		{
			name:        "Missing commands",
			filename:    "testdata/invalid/missing_commands.yaml",
			expectError: "command",
		},
		{
			name:        "Invalid option type",
			filename:    "testdata/invalid/invalid_option_type.yaml",
			expectError: "type",
		},
		{
			name:        "Invalid shorthand",
			filename:    "testdata/invalid/invalid_shorthand.yaml",
			expectError: "shorthand",
		},
		{
			name:        "Empty file",
			filename:    "testdata/invalid/empty_file.yaml",
			expectError: "empty",
		},
		{
			name:        "Malformed YAML",
			filename:    "testdata/invalid/malformed.yaml",
			expectError: "yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Load(tt.filename)

			if err == nil {
				t.Errorf("Expected error containing '%s', got no error", tt.expectError)
			}

			if config != nil {
				t.Error("Expected nil config on error, got valid config")
			}

			if tt.expectError != "" && err != nil {
				if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.expectError)) {
					t.Errorf("Expected error containing '%s', got: %v", tt.expectError, err)
				}
			}
		})
	}
}

// Test file not found
func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// Test that VarMap is properly built
func TestVarMapBuilding(t *testing.T) {
	config, err := Load("testdata/valid/complete.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.VarMap == nil {
		t.Error("Expected VarMap to be initialized")
	}

	// Check expected variables from complete.yaml
	expectedVars := map[string]string{
		"app_name": "testapp",
		"registry": "docker.io",
	}

	for name, expectedValue := range expectedVars {
		actualValue, exists := config.VarMap[name]
		if !exists {
			t.Errorf("Expected VarMap to contain '%s'", name)
			continue
		}
		if actualValue != expectedValue {
			t.Errorf("Expected VarMap[%s] = %s, got: %v", name, expectedValue, actualValue)
		}
	}
}

// Test empty commands list
func TestEmptyCommandsList(t *testing.T) {
	config := &Config{
		Version:  1,
		Commands: []Command{},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Expected error for empty commands list")
	}
}

// Test missing script
func TestMissingScript(t *testing.T) {
	cmd := Command{
		Name:   "test",
		Script: "",
	}

	err := validateCommand(cmd)
	if err == nil {
		t.Error("Expected error for missing script")
	}
}

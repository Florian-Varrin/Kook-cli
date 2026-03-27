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
		{"With namespace", "testdata/valid/with_namespace.yaml"},
		{"With includes", "testdata/valid/with_includes/Kookfile"},
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

// Test that includes are loaded correctly
func TestIncludesLoading(t *testing.T) {
	cfg, err := Load("testdata/valid/with_includes/Kookfile")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.IncludedConfigs) != 1 {
		t.Fatalf("Expected 1 included config, got %d", len(cfg.IncludedConfigs))
	}

	included := cfg.IncludedConfigs[0]
	if included.Namespace != "cms" {
		t.Errorf("Expected included namespace 'cms', got '%s'", included.Namespace)
	}
	if !included.IsIncluded {
		t.Error("Expected IsIncluded to be true")
	}
	if included.Dir == "" {
		t.Error("Expected Dir to be populated for included config")
	}
	if len(included.Commands) != 1 || included.Commands[0].Name != "deploy" {
		t.Errorf("Expected included command 'deploy', got %+v", included.Commands)
	}
}

// Test that includes of includes are not processed (one level only)
func TestIncludesAreNotRecursive(t *testing.T) {
	cfg, err := Load("testdata/valid/with_includes/Kookfile")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	for _, inc := range cfg.IncludedConfigs {
		if len(inc.IncludedConfigs) != 0 {
			t.Error("Expected included configs to not have their own includes processed")
		}
	}
}

// Test invalid includes
func TestInvalidIncludes(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectError string
	}{
		{
			name:        "Included Kookfile missing namespace",
			filename:    "testdata/invalid/includes/root_missing_namespace.yaml",
			expectError: "must define a namespace",
		},
		{
			name:        "Duplicate namespace between root and include",
			filename:    "testdata/invalid/includes/root_duplicate_namespace.yaml",
			expectError: "already used",
		},
		{
			name:        "Non-existent include path",
			filename:    "testdata/invalid/includes/root_nonexistent.yaml",
			expectError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.filename)
			if err == nil {
				t.Errorf("Expected error containing '%s', got no error", tt.expectError)
				return
			}
			if !strings.Contains(err.Error(), tt.expectError) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectError, err)
			}
		})
	}
}

// Test that Namespace is loaded and Dir is populated
func TestNamespaceAndDirLoading(t *testing.T) {
	cfg, err := Load("testdata/valid/with_namespace.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Namespace != "cms" {
		t.Errorf("Expected namespace 'cms', got '%s'", cfg.Namespace)
	}

	if cfg.Dir == "" {
		t.Error("Expected Dir to be populated")
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

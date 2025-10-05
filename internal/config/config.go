package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = "Kookfile"

// FindAndLoad searches for a Kookfile in the current directory and parent directories
func FindAndLoad() (*Config, error) {
	// Look for Kookfile in current directory first
	if _, err := os.Stat(ConfigFileName); err == nil {
		return Load(ConfigFileName)
	}

	// Search in parent directories
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		configPath := filepath.Join(dir, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return Load(configPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return nil, fmt.Errorf("no %s found in current directory or parent directories", ConfigFileName)
}

// Load reads and parses a Kookfile with validation
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate the config
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// Build variable map for template access
	config.VarMap = make(map[string]interface{})
	for _, v := range config.Variables {
		config.VarMap[v.Name] = v.Value
	}

	return &config, nil
}

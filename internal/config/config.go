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

// Load reads and parses a Kookfile
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Validate version
	if config.Version != 1 {
		return nil, fmt.Errorf("unsupported config version: %d (expected 1)", config.Version)
	}

	// Build variable map for template access
	config.VarMap = make(map[string]interface{})
	for _, v := range config.Variables {
		config.VarMap[v.Name] = v.Value
	}

	return &config, nil
}

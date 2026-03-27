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

// Load reads and parses a Kookfile, including any declared includes (one level deep).
func Load(filename string) (*Config, error) {
	cfg, err := loadFile(filename)
	if err != nil {
		return nil, err
	}

	namespaces := map[string]bool{}
	if cfg.Namespace != "" {
		namespaces[cfg.Namespace] = true
	}

	for _, includePath := range cfg.Includes {
		kookfilePath, err := resolveIncludePath(cfg.Dir, includePath)
		if err != nil {
			return nil, fmt.Errorf("include '%s': %w", includePath, err)
		}

		included, err := loadFile(kookfilePath)
		if err != nil {
			return nil, fmt.Errorf("include '%s': %w", includePath, err)
		}

		if included.Namespace == "" {
			return nil, fmt.Errorf("include '%s': included Kookfile must define a namespace", includePath)
		}
		if namespaces[included.Namespace] {
			return nil, fmt.Errorf("include '%s': namespace '%s' is already used", includePath, included.Namespace)
		}
		namespaces[included.Namespace] = true

		included.IsIncluded = true
		cfg.IncludedConfigs = append(cfg.IncludedConfigs, included)
	}

	return cfg, nil
}

// loadFile parses and validates a single Kookfile without processing its includes.
func loadFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}
	cfg.Dir = filepath.Dir(absPath)

	cfg.VarMap = make(map[string]interface{})
	for _, v := range cfg.Variables {
		cfg.VarMap[v.Name] = v.Value
	}

	return &cfg, nil
}

// resolveIncludePath resolves an include path relative to baseDir.
// If the path points to a directory, it looks for a Kookfile inside it.
func resolveIncludePath(baseDir, includePath string) (string, error) {
	if !filepath.IsAbs(includePath) {
		includePath = filepath.Join(baseDir, includePath)
	}

	info, err := os.Stat(includePath)
	if err != nil {
		return "", fmt.Errorf("path not found: %s", includePath)
	}

	if info.IsDir() {
		kookfilePath := filepath.Join(includePath, ConfigFileName)
		if _, err := os.Stat(kookfilePath); err != nil {
			return "", fmt.Errorf("no Kookfile found in directory '%s'", includePath)
		}
		return kookfilePath, nil
	}

	return includePath, nil
}

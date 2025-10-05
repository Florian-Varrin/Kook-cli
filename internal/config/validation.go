package config

import (
	"fmt"
	"regexp"
)

var (
	validNamePattern      = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	validShorthandPattern = regexp.MustCompile(`^[a-zA-Z]$`)
	validVarPattern       = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	reservedShorthands    = map[string]bool{"h": true, "i": true}
	validTypes            = map[string]bool{
		"bool":  true,
		"str":   true,
		"int":   true,
		"float": true,
	}
)

// validateConfig validates the entire configuration
func validateConfig(config *Config) error {
	// Validate version
	if config.Version != 1 {
		return fmt.Errorf("unsupported config version: %d (expected 1)", config.Version)
	}

	// Must have at least one command
	if len(config.Commands) == 0 {
		return fmt.Errorf("config must have at least one command")
	}

	// Validate variables
	for i, v := range config.Variables {
		if err := validateVariable(v); err != nil {
			return fmt.Errorf("variable %d: %w", i, err)
		}
	}

	// Validate commands
	commandNames := make(map[string]bool)
	for i, cmd := range config.Commands {
		if err := validateCommand(cmd); err != nil {
			return fmt.Errorf("command %d (%s): %w", i, cmd.Name, err)
		}

		// Check for duplicate command names
		if commandNames[cmd.Name] {
			return fmt.Errorf("duplicate command name: %s", cmd.Name)
		}
		commandNames[cmd.Name] = true

		// Check for duplicate aliases
		for _, alias := range cmd.Aliases {
			if commandNames[alias] {
				return fmt.Errorf("duplicate command name/alias: %s", alias)
			}
			commandNames[alias] = true
		}
	}

	return nil
}

// validateVariable validates a single variable
func validateVariable(v Variable) error {
	if v.Name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if !validNamePattern.MatchString(v.Name) {
		return fmt.Errorf("invalid variable name '%s': must start with letter and contain only letters, numbers, hyphens, and underscores", v.Name)
	}

	return nil
}

// validateCommand validates a single command
func validateCommand(cmd Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if !validNamePattern.MatchString(cmd.Name) {
		return fmt.Errorf("invalid command name '%s': must start with letter and contain only letters, numbers, hyphens, and underscores", cmd.Name)
	}

	if cmd.Script == "" {
		return fmt.Errorf("command script cannot be empty")
	}

	// Validate aliases
	for _, alias := range cmd.Aliases {
		if !validNamePattern.MatchString(alias) {
			return fmt.Errorf("invalid alias '%s': must start with letter and contain only letters, numbers, hyphens, and underscores", alias)
		}
	}

	// Validate options
	optionNames := make(map[string]bool)
	shorthands := make(map[string]bool)

	for i, opt := range cmd.Options {
		if err := validateOption(opt); err != nil {
			return fmt.Errorf("option %d (%s): %w", i, opt.Name, err)
		}

		// Check for duplicate option names
		if optionNames[opt.Name] {
			return fmt.Errorf("duplicate option name: %s", opt.Name)
		}
		optionNames[opt.Name] = true

		// Check for duplicate shorthands
		if opt.Shorthand != "" {
			if shorthands[opt.Shorthand] {
				return fmt.Errorf("duplicate shorthand: %s", opt.Shorthand)
			}
			shorthands[opt.Shorthand] = true
		}
	}

	return nil
}

// validateOption validates a single option
func validateOption(opt Option) error {
	if opt.Name == "" {
		return fmt.Errorf("option name cannot be empty")
	}

	if !validNamePattern.MatchString(opt.Name) {
		return fmt.Errorf("invalid option name '%s': must start with letter and contain only letters, numbers, hyphens, and underscores", opt.Name)
	}

	// Validate type
	if !validTypes[opt.Type] {
		return fmt.Errorf("invalid option type '%s': must be bool, str, int, or float", opt.Type)
	}

	// Validate shorthand
	if opt.Shorthand != "" {
		if !validShorthandPattern.MatchString(opt.Shorthand) {
			return fmt.Errorf("invalid shorthand '%s': must be a single letter", opt.Shorthand)
		}

		if reservedShorthands[opt.Shorthand] {
			return fmt.Errorf("shorthand '%s' is reserved (used by -h/--help or -i/--interactive)", opt.Shorthand)
		}
	}

	// Validate var name if provided
	if opt.Var != "" {
		if !validVarPattern.MatchString(opt.Var) {
			return fmt.Errorf("invalid var name '%s': must start with letter or underscore and contain only letters, numbers, and underscores", opt.Var)
		}
	}

	return nil
}

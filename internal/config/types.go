package config

import (
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version         int                    `yaml:"version"`
	Namespace       string                 `yaml:"namespace,omitempty"`
	Includes        []Include              `yaml:"includes,omitempty"`
	Variables       []Variable             `yaml:"variables"`
	Commands        []Command              `yaml:"commands"`
	VarMap          map[string]interface{} `yaml:"-"`
	Dir             string                 `yaml:"-"`
	IncludedConfigs []*Config              `yaml:"-"`
	IsIncluded      bool                   `yaml:"-"`
	Hidden          bool                   `yaml:"-"`
	Only            []string               `yaml:"-"`
}

// IsCommandHidden reports whether a command is left out of help and completion.
func (c *Config) IsCommandHidden(name string) bool {
	if c.Hidden {
		return true
	}
	if len(c.Only) == 0 {
		return false
	}
	for _, pattern := range c.Only {
		if matched, _ := path.Match(pattern, name); matched {
			return false
		}
	}
	return true
}

func (c *Config) HasVisibleCommands() bool {
	for _, cmd := range c.Commands {
		if !c.IsCommandHidden(cmd.Name) {
			return true
		}
	}
	return false
}

// Include accepts either a plain path or a {path, hidden, only} mapping.
type Include struct {
	Path   string   `yaml:"path"`
	Hidden bool     `yaml:"hidden"`
	Only   []string `yaml:"only"`
}

func (i *Include) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		return node.Decode(&i.Path)
	case yaml.MappingNode:
		type plain Include
		return node.Decode((*plain)(i))
	default:
		return fmt.Errorf("line %d: include must be a path or a {path, hidden, only} mapping", node.Line)
	}
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Command struct {
	Name        string   `yaml:"name"`
	Aliases     []string `yaml:"aliases"`
	Description string   `yaml:"description,omitempty"`
	Help        string   `yaml:"help,omitempty"`
	Options     []Option `yaml:"options"`
	Script      string   `yaml:"script"`
	Silent      bool     `yaml:"silent,omitempty"`
}

type Option struct {
	Name        string `yaml:"name"`
	Shorthand   string `yaml:"shorthand,omitempty"`
	Description string `yaml:"description,omitempty"`
	Var         string `yaml:"var,omitempty"`
	Type        string `yaml:"type"`
	Mandatory   bool   `yaml:"mandatory,omitempty"`
}

func (o Option) GetVarName() string {
	if o.Var != "" {
		return o.Var
	}
	return strings.ReplaceAll(o.Name, "-", "_")
}

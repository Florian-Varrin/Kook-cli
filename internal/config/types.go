package config

import "strings"

type Config struct {
	Version   int                    `yaml:"version"`
	Variables []Variable             `yaml:"variables"`
	Commands  []Command              `yaml:"commands"`
	VarMap    map[string]interface{} `yaml:"-"`
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

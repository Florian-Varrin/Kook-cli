package cli

import (
	"testing"

	"kook/internal/config"
)

func TestBuildCommandHiddenInclude(t *testing.T) {
	cmd := config.Command{Name: "up", Script: "echo up"}

	tests := []struct {
		name   string
		hidden bool
	}{
		{"visible include", false},
		{"hidden include", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Namespace: "docker", IsIncluded: true, Hidden: tt.hidden}
			cobraCmd := buildCommand(cfg, cmd, cfg.Namespace)

			if cobraCmd.Use != "docker:up" {
				t.Errorf("Expected use 'docker:up', got '%s'", cobraCmd.Use)
			}
			if cobraCmd.Hidden != tt.hidden {
				t.Errorf("Expected Hidden=%v, got %v", tt.hidden, cobraCmd.Hidden)
			}
		})
	}
}

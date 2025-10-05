package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"kook/internal/config"

	"github.com/spf13/cobra"
)

// Execute runs a command with the given configuration and cobra command
func Execute(cfg *config.Config, cmd config.Command, cobraCmd *cobra.Command) error {
	// Build template context with variables and options
	ctx := make(map[string]interface{})

	// Add all variables
	for k, v := range cfg.VarMap {
		ctx[k] = v
	}

	// Add all option values using their var names
	for _, opt := range cmd.Options {
		val, err := getOptionValue(cobraCmd, opt)
		if err != nil {
			return fmt.Errorf("failed to get option '%s': %w", opt.Name, err)
		}
		ctx[opt.GetVarName()] = val
	}

	// Parse and execute template
	tmpl, err := template.New(cmd.Name).Parse(cmd.Script)
	if err != nil {
		return fmt.Errorf("failed to parse script template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return fmt.Errorf("failed to execute script template: %w", err)
	}

	scriptCmd := buf.String()

	// Print execution message unless silent mode is enabled
	if !cmd.Silent {
		fmt.Printf("Executing: %s\n", scriptCmd)
	}

	// Execute the command using bash
	bashCmd := exec.Command("bash", "-c", scriptCmd)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr
	bashCmd.Stdin = os.Stdin

	return bashCmd.Run()
}

func getOptionValue(cobraCmd *cobra.Command, opt config.Option) (interface{}, error) {
	switch opt.Type {
	case "bool":
		return cobraCmd.Flags().GetBool(opt.Name)
	case "str":
		return cobraCmd.Flags().GetString(opt.Name)
	case "int":
		return cobraCmd.Flags().GetInt(opt.Name)
	case "float":
		return cobraCmd.Flags().GetFloat64(opt.Name)
	default:
		return nil, fmt.Errorf("unknown option type: %s", opt.Type)
	}
}

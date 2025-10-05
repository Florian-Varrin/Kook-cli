package cli

import (
	"fmt"
	"os"

	"kook/internal/config"
	"kook/internal/executor"

	"github.com/spf13/cobra"
)

// Execute is the main entry point for the CLI
func Execute() error {
	rootCmd := buildRootCommand()

	// Try to load config and add dynamic commands
	cfg, err := config.FindAndLoad()
	if err != nil {
		// If no config found, still allow completion command to work
		if len(os.Args) > 1 && os.Args[1] == "completion" {
			return rootCmd.Execute()
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Add all commands from config
	for _, cmd := range cfg.Commands {
		rootCmd.AddCommand(buildCommand(cfg, cmd))
	}

	return rootCmd.Execute()
}

func buildRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kook",
		Short: "A simple CLI tool configured via Kookfile",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Dynamic completion: load current directory's Kookfile
			cfg, err := config.FindAndLoad()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completions []string
			for _, c := range cfg.Commands {
				completions = append(completions, c.Name)
				completions = append(completions, c.Aliases...)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		},
	}

	// Add completion command
	rootCmd.AddCommand(buildCompletionCommand())

	return rootCmd
}

func buildCommand(cfg *config.Config, cmd config.Command) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:     cmd.Name,
		Aliases: cmd.Aliases,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return executor.Execute(cfg, cmd, cobraCmd)
		},
	}

	// Add flags for each option
	for _, opt := range cmd.Options {
		addFlag(cobraCmd, opt)

		// Mark as required if mandatory
		if opt.Mandatory {
			cobraCmd.MarkFlagRequired(opt.Name)
		}
	}

	return cobraCmd
}

func addFlag(cobraCmd *cobra.Command, opt config.Option) {
	switch opt.Type {
	case "bool":
		cobraCmd.Flags().Bool(opt.Name, false, "")
	case "str":
		cobraCmd.Flags().String(opt.Name, "", "")
	case "int":
		cobraCmd.Flags().Int(opt.Name, 0, "")
	case "float":
		cobraCmd.Flags().Float64(opt.Name, 0.0, "")
	default:
		fmt.Fprintf(os.Stderr, "Warning: unknown option type '%s' for option '%s'\n", opt.Type, opt.Name)
	}
}

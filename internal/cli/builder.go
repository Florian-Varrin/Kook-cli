package cli

import (
	"fmt"
	"os"
	"strconv"

	"kook/internal/config"
	"kook/internal/executor"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Execute is the main entry point for the CLI
func Execute(version string) error {
	rootCmd := buildRootCommand(version)

	// Try to load config and add dynamic commands
	cfg, err := config.FindAndLoad()
	if err != nil {
		// If no config found, still allow completion and version commands to work
		if len(os.Args) > 1 && (os.Args[1] == "completion" || os.Args[1] == "version") {
			return rootCmd.Execute()
		}
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Add root commands
	for _, cmd := range cfg.Commands {
		rootCmd.AddCommand(buildCommand(cfg, cmd, cfg.Namespace))
	}

	// Add commands from included configs
	for _, includedCfg := range cfg.IncludedConfigs {
		for _, cmd := range includedCfg.Commands {
			rootCmd.AddCommand(buildCommand(includedCfg, cmd, includedCfg.Namespace))
		}
	}

	return rootCmd.Execute()
}

func buildRootCommand(version string) *cobra.Command {
	long := `Kook is a task runner that reads commands from a Kookfile.

Each project can have its own Kookfile with custom commands,
options, and variables. Commands support Go templates for
dynamic script generation.`

	cfg, _ := config.FindAndLoad()
	if cfg != nil && cfg.Namespace != "" {
		long += fmt.Sprintf("\n\nNamespace: %s\nCommands are listed as \"%s:<command>\" but also available without the namespace prefix.", cfg.Namespace, cfg.Namespace)
	}
	if cfg != nil {
		for _, inc := range cfg.IncludedConfigs {
			long += fmt.Sprintf("\n\nNamespace: %s (included)\nCommands available as \"%s:<command>\".", inc.Namespace, inc.Namespace)
		}
	}

	rootCmd := &cobra.Command{
		Use:     "kook",
		Short:   "A simple CLI tool configured via Kookfile",
		Long:    long,
		Version: version,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			cfg, err := config.FindAndLoad()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			var completions []string

			addEntry := func(name, description string) {
				if description != "" {
					completions = append(completions, fmt.Sprintf("%s\t%s", name, description))
				} else {
					completions = append(completions, name)
				}
			}

			addCommandCompletions := func(commands []config.Command, namespace string, isIncluded bool) {
				for _, c := range commands {
					if namespace != "" {
						addEntry(fmt.Sprintf("%s:%s", namespace, c.Name), c.Description)
						for _, alias := range c.Aliases {
							addEntry(fmt.Sprintf("%s:%s", namespace, alias), c.Description)
						}
					}
					// Bare names only for root commands
					if !isIncluded {
						addEntry(c.Name, c.Description)
						for _, alias := range c.Aliases {
							addEntry(alias, c.Description)
						}
					}
				}
			}

			addCommandCompletions(cfg.Commands, cfg.Namespace, false)
			for _, includedCfg := range cfg.IncludedConfigs {
				addCommandCompletions(includedCfg.Commands, includedCfg.Namespace, true)
			}

			return completions, cobra.ShellCompDirectiveNoFileComp
		},
	}

	rootCmd.AddCommand(buildCompletionCommand())

	return rootCmd
}

func buildCommand(cfg *config.Config, cmd config.Command, namespace string) *cobra.Command {
	var use string
	var aliases []string

	if cfg.IsIncluded {
		// Included commands: only accessible via namespace prefix.
		use = fmt.Sprintf("%s:%s", namespace, cmd.Name)
		for _, alias := range cmd.Aliases {
			aliases = append(aliases, fmt.Sprintf("%s:%s", namespace, alias))
		}
	} else if namespace != "" {
		// Root commands with namespace: primary name is namespaced, bare names are aliases.
		use = fmt.Sprintf("%s:%s", namespace, cmd.Name)
		aliases = append(aliases, cmd.Name)
		for _, alias := range cmd.Aliases {
			aliases = append(aliases, alias)
			aliases = append(aliases, fmt.Sprintf("%s:%s", namespace, alias))
		}
	} else {
		// Root commands without namespace.
		use = cmd.Name
		aliases = make([]string, len(cmd.Aliases))
		copy(aliases, cmd.Aliases)
	}

	cobraCmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   cmd.Description,
		Long:    cmd.Help,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			interactive, _ := cobraCmd.Flags().GetBool("interactive")

			if interactive {
				if err := promptForOptions(cobraCmd, cmd); err != nil {
					return fmt.Errorf("interactive prompt failed: %w", err)
				}

				// Validate mandatory fields after interactive input
				for _, opt := range cmd.Options {
					if opt.Mandatory {
						var isEmpty bool
						switch opt.Type {
						case "bool":
							isEmpty = false
						case "str":
							val, _ := cobraCmd.Flags().GetString(opt.Name)
							isEmpty = val == ""
						case "int":
							isEmpty = !cobraCmd.Flags().Changed(opt.Name)
						case "float":
							isEmpty = !cobraCmd.Flags().Changed(opt.Name)
						}

						if isEmpty {
							return fmt.Errorf("required option '%s' not provided", opt.Name)
						}
					}
				}
			}

			return executor.Execute(cfg, cmd, cobraCmd)
		},
	}

	cobraCmd.Flags().BoolP("interactive", "i", false, "Use interactive mode to select options")

	for _, opt := range cmd.Options {
		addFlag(cobraCmd, opt)
		// Don't use MarkFlagRequired - we'll validate manually
	}

	// Custom flag validation that checks if we're in interactive mode
	cobraCmd.PreRunE = func(cobraCmd *cobra.Command, args []string) error {
		interactive, _ := cobraCmd.Flags().GetBool("interactive")
		if !interactive {
			// Only validate required flags if NOT in interactive mode
			for _, opt := range cmd.Options {
				if opt.Mandatory && !cobraCmd.Flags().Changed(opt.Name) {
					return fmt.Errorf("required flag(s) \"%s\" not set", opt.Name)
				}
			}
		}
		return nil
	}

	return cobraCmd
}

func promptForOptions(cobraCmd *cobra.Command, cmd config.Command) error {
	for _, opt := range cmd.Options {
		// Skip if flag was already provided via command line
		if cobraCmd.Flags().Changed(opt.Name) {
			continue
		}

		var prompt survey.Prompt
		message := opt.Name
		if opt.Description != "" {
			message = opt.Description
		}

		switch opt.Type {
		case "bool":
			prompt = &survey.Select{
				Message: message,
				Options: []string{"Yes", "No"},
				Default: "No",
			}
			var answer string
			if err := survey.AskOne(prompt, &answer); err != nil {
				return err
			}
			value := answer == "Yes"
			cobraCmd.Flags().Set(opt.Name, strconv.FormatBool(value))

		case "str":
			prompt = &survey.Input{
				Message: message,
			}
			var answer string
			if err := survey.AskOne(prompt, &answer, survey.WithValidator(func(ans interface{}) error {
				if opt.Mandatory && ans.(string) == "" {
					return fmt.Errorf("this field is required")
				}
				return nil
			})); err != nil {
				return err
			}
			cobraCmd.Flags().Set(opt.Name, answer)

		case "int":
			prompt = &survey.Input{
				Message: message,
			}
			var answer string
			if err := survey.AskOne(prompt, &answer, survey.WithValidator(func(ans interface{}) error {
				str := ans.(string)
				if opt.Mandatory && str == "" {
					return fmt.Errorf("this field is required")
				}
				if str != "" {
					if _, err := strconv.Atoi(str); err != nil {
						return fmt.Errorf("must be a valid integer")
					}
				}
				return nil
			})); err != nil {
				return err
			}
			if answer != "" {
				cobraCmd.Flags().Set(opt.Name, answer)
			}

		case "float":
			prompt = &survey.Input{
				Message: message,
			}
			var answer string
			if err := survey.AskOne(prompt, &answer, survey.WithValidator(func(ans interface{}) error {
				str := ans.(string)
				if opt.Mandatory && str == "" {
					return fmt.Errorf("this field is required")
				}
				if str != "" {
					if _, err := strconv.ParseFloat(str, 64); err != nil {
						return fmt.Errorf("must be a valid number")
					}
				}
				return nil
			})); err != nil {
				return err
			}
			if answer != "" {
				cobraCmd.Flags().Set(opt.Name, answer)
			}
		}
	}

	return nil
}

func addFlag(cobraCmd *cobra.Command, opt config.Option) {
	usage := opt.Description

	switch opt.Type {
	case "bool":
		if opt.Shorthand != "" {
			cobraCmd.Flags().BoolP(opt.Name, opt.Shorthand, false, usage)
		} else {
			cobraCmd.Flags().Bool(opt.Name, false, usage)
		}
	case "str":
		if opt.Shorthand != "" {
			cobraCmd.Flags().StringP(opt.Name, opt.Shorthand, "", usage)
		} else {
			cobraCmd.Flags().String(opt.Name, "", usage)
		}
	case "int":
		if opt.Shorthand != "" {
			cobraCmd.Flags().IntP(opt.Name, opt.Shorthand, 0, usage)
		} else {
			cobraCmd.Flags().Int(opt.Name, 0, usage)
		}
	case "float":
		if opt.Shorthand != "" {
			cobraCmd.Flags().Float64P(opt.Name, opt.Shorthand, 0.0, usage)
		} else {
			cobraCmd.Flags().Float64(opt.Name, 0.0, usage)
		}
	default:
		fmt.Fprintf(os.Stderr, "Warning: unknown option type '%s' for option '%s'\n", opt.Type, opt.Name)
	}
}

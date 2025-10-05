# Kook

A simple, powerful CLI task runner configured via `Kookfile`. Think of it as `Just` with templating superpowers - define commands once, run them anywhere in your projects.

## Why Kook?

- **Project-specific commands**: Each project can have its own `Kookfile` with custom commands
- **Template-driven**: Use Go templates to create dynamic commands based on options and variables
- **Type-safe options**: Define boolean, string, integer, and float options with validation
- **Interactive mode**: Use `--interactive` flag to get prompted for options with a user-friendly interface
- **Auto-completion**: Smart shell completion that adapts to each project's `Kookfile`
- **Zero config**: Just drop a `Kookfile` in your project and you're ready to go

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/kook.git
cd kook

# Build and install
go build -o kook
sudo mv kook /usr/local/bin/

# Or use go install
go install
```

## Quick Start

1. **Create a `Kookfile` in your project:**

```yaml
version: 1

variables:
  - name: container
    value: my-app

commands:
  - name: start
    description: Start the application
    aliases:
      - up
    script: |
      docker compose up -d
      echo "Started {{ .container }}"

  - name: logs
    description: Show container logs
    options:
      - name: follow
        description: Follow log output
        type: bool
    script: |
      docker logs {{ .container }} {{- if .follow }} -f{{- end }}
```

2. **Run your commands:**

```bash
kook start
kook logs --follow
```

That's it! 🎉

## Interactive Mode

Every command automatically supports `--interactive` (or `-i`) mode, which prompts you for options instead of requiring them on the command line.

### Usage

```bash
# Regular mode with flags
kook deploy --environment production --tag v1.2.3 --dry-run

# Interactive mode - prompts for all options
kook deploy --interactive

# Shorthand
kook deploy -i

# Mix both - interactive only prompts for missing options
kook deploy --environment production -i
```

### Example Interaction

```bash
$ kook deploy -i
? Target environment (staging, production): production
? Docker image tag to deploy: v1.2.3
? Show what would happen without executing
  ❯ Yes
    No

Executing: kubectl set image deployment/app app=v1.2.3 -n production
```

### Benefits

- ✅ **Easier to use** - No need to remember flag names
- ✅ **Guided experience** - Option descriptions help users understand what to enter
- ✅ **Mandatory validation** - Won't proceed without required fields
- ✅ **Type validation** - Ensures correct input types (int, float, etc.)
- ✅ **Mix with flags** - Combine CLI flags with interactive prompts

## Kookfile Structure

### Basic Format

```yaml
version: 1              # Required: config version (only "1" supported)

variables:              # Optional: global variables
  - name: var_name
    value: var_value

commands:               # Required: your custom commands
  - name: command-name
    # ... command definition
```

### Command Definition

```yaml
commands:
  - name: deploy                    # Required: command name
    description: Deploy to env      # Optional: short one-line description
    help: |                         # Optional: long multi-line help text
      Detailed explanation of what this command does.
      
      Can include multiple paragraphs and examples.
    aliases:                        # Optional: command shortcuts
      - d
    silent: false                   # Optional: hide "Executing..." output (default: false)
    options:                        # Optional: command options/flags
      - name: environment           # Required: option name (use hyphens for CLI)
        description: Target env     # Optional: option description/help text
        var: env                    # Optional: variable name in template (default: auto-convert hyphens to underscores)
        type: str                   # Required: bool, str, int, or float
        mandatory: true             # Optional: make option required (default: false)
    script: |                       # Required: command script (supports Go templates)
      kubectl apply -f deploy.yaml --namespace {{ .env }}
```

### Variables

Variables are globally accessible in all command scripts:

```yaml
variables:
  - name: app_name
    value: myapp
  - name: docker_registry
    value: registry.example.com

commands:
  - name: build
    description: Build Docker image
    script: |
      docker build -t {{ .docker_registry }}/{{ .app_name }}:latest .
```

### Options

Options define command-line flags:

#### Option Types

- **`bool`**: Boolean flag (true when present, false otherwise)
- **`str`**: String value
- **`int`**: Integer value
- **`float`**: Float value

#### Option Properties

```yaml
options:
  - name: dry-run              # CLI flag: --dry-run
    description: Preview only  # Optional: option description
    var: dryRun                # Template variable: .dryRun (optional, defaults to dry_run)
    type: bool                 # Type: bool, str, int, float
    mandatory: true            # Make it required (optional, default: false)
```

**Important**: CLI flags use hyphens (`--dry-run`), but template variables use underscores or custom names:
- `--dry-run` → `.dry_run` (automatic)
- `--dry-run` with `var: dryRun` → `.dryRun` (explicit)

### Templates

Kook uses [Go templates](https://pkg.go.dev/text/template) in scripts:

#### Accessing Variables and Options

```yaml
script: |
  # Variables
  echo {{ .my_variable }}
  
  # Options
  echo {{ .option_name }}
```

#### Conditionals

```yaml
script: |
  {{- if .verbose }}
  echo "Verbose mode enabled"
  {{- end }}
  
  {{- if eq .environment "prod" }}
  echo "Production deployment"
  {{- else }}
  echo "Non-production deployment"
  {{- end }}
```

#### Loops

```yaml
# Not directly supported via Kookfile, but you can use bash loops in your script
script: |
  for service in api worker scheduler; do
    docker restart {{ .container }}-$service
  done
```

#### Template Tips

- Use `{{- ... }}` to trim whitespace before
- Use `... -}}` to trim whitespace after
- Use `{{ if .flag }}...{{ end }}` for conditionals
- Use `{{ .var }}` to access variables/options

## Examples

### Development Environment

```yaml
version: 1

variables:
  - name: compose_file
    value: docker-compose.yml

commands:
  - name: start
    description: Start development environment
    aliases:
      - up
    options:
      - name: detach
        description: Run in detached mode
        type: bool
    script: |
      docker compose -f {{ .compose_file }} up {{- if .detach }} -d{{- end }}

  - name: stop
    description: Stop development environment
    aliases:
      - down
    script: |
      docker compose -f {{ .compose_file }} down

  - name: logs
    description: Show service logs
    options:
      - name: service
        description: Service name
        type: str
        mandatory: true
      - name: follow
        description: Follow log output
        type: bool
    script: |
      docker compose -f {{ .compose_file }} logs {{ .service }} {{- if .follow }} -f{{- end }}

  - name: restart
    description: Restart a service
    options:
      - name: service
        description: Service name
        type: str
        mandatory: true
    script: |
      docker compose -f {{ .compose_file }} restart {{ .service }}
```

Usage:
```bash
kook start --detach
kook logs --service api --follow
kook restart --service worker
kook stop

# Or use interactive mode
kook logs -i
kook restart -i
```
### Cache Management

```yaml
version: 1

variables:
  - name: container
    value: app_php

commands:
  - name: clear-cache
    description: Clear application cache
    aliases:
      - cc
    options:
      - name: full
        description: Perform full cache clear including warmup
        type: bool
    script: |
      docker exec {{ .container }} bash -c 'rm -rf var/cache/* && chmod a+rw -R var/cache{{- if .full }} && php bin/console cache:clear{{- end }}'
      echo "Cache cleared{{- if .full }} (full){{- end }}"

  - name: warmup-cache
    description: Warm up application cache
    aliases:
      - warm
    script: |
      docker exec {{ .container }} php bin/console cache:warmup
```

Usage:
```bash
kook cc
kook cc --full
kook warm

# Interactive mode
kook cc -i
```

## Shell Completion

Kook provides dynamic shell completion that adapts to each project's `Kookfile`. Completions automatically update based on the commands available in your current directory.

### How It Works

1. When you press `Tab`, your shell calls `kook __complete`
2. Kook reads the `Kookfile` in your current directory (or parent directories)
3. Kook returns available commands and options
4. Your shell displays the completions

This means each project gets its own custom completions!

### Setup Completion

#### Bash

**One-time setup:**
```bash
# Generate completion script
kook completion bash > ~/.kook-completion.bash

# Add to your ~/.bashrc
echo 'source ~/.kook-completion.bash' >> ~/.bashrc

# Reload
source ~/.bashrc
```

#### Zsh

**One-time setup:**
```bash
# Generate completion script
kook completion zsh > "${fpath[1]}/_kook"

# Reload completions
compinit
```

Or add to `~/.zshrc`:
```bash
# Add this to ~/.zshrc
autoload -U compinit; compinit
```

#### Fish

**One-time setup:**
```bash
# Generate completion script
kook completion fish > ~/.config/fish/completions/kook.fish

# Reload (or restart Fish)
source ~/.config/fish/completions/kook.fish
```

#### PowerShell

**One-time setup:**
```powershell
# Generate and run completion script
kook completion powershell | Out-String | Invoke-Expression

# For persistent completions, add to your PowerShell profile
kook completion powershell > kook-completion.ps1
# Then add this to your profile: . /path/to/kook-completion.ps1
```

### Using Completion

Once set up, completion works automatically:

```bash
# Complete command names
kook <Tab>
# Shows: clear-cache, cc, deploy, logs, start, up, ...

# Complete with aliases
kook c<Tab>
# Shows: clear-cache, cc

# Complete options
kook deploy --<Tab>
# Shows: --environment, --tag, --dry-run, --help

# Complete option values (for mandatory options)
kook deploy --environment <Tab>
# (Your shell may provide additional completion based on history)
```

### Per-Project Completions

The beauty of Kook is that completions change based on your location:

```bash
# In project-a/
cd ~/projects/project-a
kook <Tab>
# Shows commands from project-a/Kookfile

# In project-b/
cd ~/projects/project-b
kook <Tab>
# Shows commands from project-b/Kookfile (different commands!)

# In a subdirectory
cd ~/projects/project-a/src/components
kook <Tab>
# Still finds ~/projects/project-a/Kookfile and shows its commands
```

## Configuration Search

Kook searches for `Kookfile` in the following order:

1. Current directory
2. Parent directories (recursively up to root)

This means you can run `kook` commands from anywhere within your project tree!

```bash
# All of these work if Kookfile is at ~/projects/myapp/Kookfile
cd ~/projects/myapp && kook start
cd ~/projects/myapp/src && kook start
cd ~/projects/myapp/src/components && kook start
```

## Development

```bash
# Run without building
go run main.go <command>

# Build
go build -o kook

# Run tests
go test ./...

# Install locally
go install
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
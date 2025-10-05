# Kook

A simple, powerful CLI task runner configured via `Kookfile`. Think of it as `Just` with templating superpowers - define commands once, run them anywhere in your projects.

## Why Kook?

- **Project-specific commands**: Each project can have its own `Kookfile` with custom commands
- **Template-driven**: Use Go templates to create dynamic commands based on options and variables
- **Type-safe options**: Define boolean, string, integer, and float options with validation
- **Interactive mode**: Use `--interactive` flag to get prompted for options with a user-friendly interface
- **Auto-completion**: Smart shell completion that adapts to each project's `Kookfile`
- **IDE support**: JSON Schema for auto-completion and validation in VS Code, JetBrains IDEs, and more
- **Zero config**: Just drop a `Kookfile` in your project and you're ready to go

## Installation

### Stable Releases

Download the latest stable release from the [releases page](https://github.com/Florian-Varrin/Kook-cli/releases/latest).

### Pre-release Versions

Want to test upcoming features? Pre-release versions (alpha, beta, RC) are available on the [releases page](https://github.com/Florian-Varrin/Kook-cli/releases).

⚠️ **Note:** Pre-releases may contain bugs or incomplete features. Use stable releases for production environments.

### Linux

#### Debian/Ubuntu (.deb)

```bash
# Download the latest .deb package
wget https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_linux_amd64.deb

# Install
sudo dpkg -i kook_<version>_linux_amd64.deb

# Or use apt (handles dependencies better)
sudo apt install ./kook_<version>_linux_amd64.deb
```

For ARM64:
```bash
wget https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_linux_arm64.deb
sudo apt install ./kook_<version>_linux_arm64.deb
```

#### RedHat/Fedora/CentOS (.rpm)

```bash
# Download the latest .rpm package
wget https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_linux_amd64.rpm

# Install
sudo rpm -i kook_<version>_linux_amd64.rpm

# Or with dnf
sudo dnf install kook_<version>_linux_amd64.rpm
```

#### Alpine (.apk)

```bash
# Download the latest .apk package
wget https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_linux_amd64.apk

# Install
sudo apk add --allow-untrusted kook_<version>_linux_amd64.apk
```

#### Manual Installation (Any Linux)

```bash
# Download and extract
curl -L https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_linux_amd64.tar.gz | tar xz

# Move to PATH
sudo mv kook /usr/local/bin/

# Verify installation
kook --version
```

### macOS

```bash
# Download and extract
curl -L https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_darwin_amd64.tar.gz | tar xz

# Move to PATH
sudo mv kook /usr/local/bin/

# For Apple Silicon (M1/M2/M3)
curl -L https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_darwin_arm64.tar.gz | tar xz
sudo mv kook /usr/local/bin/
```

### Windows

1. Download the latest Windows release from [releases page](https://github.com/Florian-Varrin/Kook-cli/releases)
2. Extract the `.zip` file
3. Move `kook.exe` to a directory in your PATH

Or use PowerShell:
```powershell
# Download (replace <version> with actual version number)
Invoke-WebRequest -Uri "https://github.com/Florian-Varrin/Kook-cli/releases/latest/download/kook_<version>_windows_amd64.zip" -OutFile "kook.zip"

# Extract
Expand-Archive -Path kook.zip -DestinationPath .

# Move to a directory in your PATH (example)
Move-Item kook.exe C:\Windows\System32\
```

### From Source

```bash
# Clone the repository
git clone https://github.com/Florian-Varrin/Kook-cli.git
cd Kook-cli

# Build and install
go build -o kook
sudo mv kook /usr/local/bin/

# Or use go install
go install
```

### Verify Installation

After installation, verify it works:

```bash
kook --version
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

## IDE Support

Kook provides JSON Schema for auto-completion and validation in your IDE.

### VS Code

1. Install the [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
2. Add this comment to the top of your `Kookfile`:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Florian-Varrin/Kook-cli/master/kookfile-schema.json

version: 1
commands:
  - name: example
    # IDE now provides auto-completion, validation, and inline documentation!
    description: Example command
    script: echo "Hello"
```

**What you get:**
- ✅ Auto-completion for all fields (type `.` to see suggestions)
- ✅ Inline documentation on hover
- ✅ Real-time validation and error detection
- ✅ Dropdown suggestions for enums (e.g., `type: bool|str|int|float`)

### JetBrains IDEs (IntelliJ, WebStorm, PyCharm, etc.)

**First time setup:**

1. Open your `Kookfile`
2. IntelliJ will likely not recognize it as YAML initially
3. Right-click the file → **Associate with File Type** → **YAML**
4. Alternatively, go to **Settings** → **Editor** → **File Types** → **YAML** → Add pattern `Kookfile`

**Option 1: Inline schema (recommended)**

Add this comment to the top of your `Kookfile`:
```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/Florian-Varrin/Kook-cli/master/kookfile-schema.json
```

**Option 2: Configure globally**

1. Go to **Settings** → **Languages & Frameworks** → **Schemas and DTDs** → **JSON Schema Mappings**
2. Click **+** to add new mapping
3. Name: `Kookfile`
4. Schema URL: `https://raw.githubusercontent.com/Florian-Varrin/Kook-cli/master/kookfile-schema.json`
5. Add file path pattern: `Kookfile`

### Vim/Neovim

**With coc.nvim:**

Add to your `coc-settings.json`:
```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/Florian-Varrin/Kook-cli/master/kookfile-schema.json": "Kookfile"
  }
}
```

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
        shorthand: e                # Optional: single letter shorthand (e.g., 'e' for -e)
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
    shorthand: d               # Optional: short flag -d
    description: Preview only  # Optional: option description
    var: dryRun                # Template variable: .dryRun (optional, defaults to dry_run)
    type: bool                 # Type: bool, str, int, float
    mandatory: true            # Make it required (optional, default: false)
```

**Important**:
- CLI flags use hyphens (`--dry-run`), but template variables use underscores or custom names:
    - `--dry-run` → `.dry_run` (automatic)
    - `--dry-run` with `var: dryRun` → `.dryRun` (explicit)
- Shorthand must be a single letter (e.g., `d`, `v`, `e`)
- Reserved shorthands: `-h` (help), `-i` (interactive)

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
        shorthand: d
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
        shorthand: s
        description: Service name
        type: str
        mandatory: true
      - name: follow
        shorthand: f
        description: Follow log output
        type: bool
    script: |
      docker compose -f {{ .compose_file }} logs {{ .service }} {{- if .follow }} -f{{- end }}

  - name: restart
    description: Restart a service
    options:
      - name: service
        shorthand: s
        description: Service name
        type: str
        mandatory: true
    script: |
      docker compose -f {{ .compose_file }} restart {{ .service }}
```

Usage:
```bash
kook start --detach
kook start -d              # Using shorthand

kook logs --service api --follow
kook logs -s api -f        # Using shorthands

kook restart --service worker
kook restart -s worker     # Using shorthand

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
        shorthand: f
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
kook cc -f         # Using shorthand

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

## Releases

Kook follows a structured release process using GitFlow and automated releases via GitHub Actions.

### Release Types

We use [Semantic Versioning](https://semver.org/) with the following release types:

- **Stable Releases**: `v1.0.0`, `v1.1.0`, `v2.0.0`
    - Production-ready
    - Fully tested
    - Marked as "Latest" on GitHub

- **Release Candidates (RC)**: `v1.0.0-rc.1`, `v1.0.0-rc.2`
    - Feature-complete
    - Final testing before stable release
    - May have minor bugs

- **Beta Releases**: `v1.0.0-beta.1`, `v1.0.0-beta.2`
    - Major features implemented
    - Still under active testing
    - May have known issues

- **Alpha Releases**: `v1.0.0-alpha.1`, `v1.0.0-alpha.2`
    - Early preview
    - Experimental features
    - Not recommended for production

### How Releases Work

#### 1. Automated Release Process

When a version tag is pushed to GitHub, our CI/CD pipeline automatically:

1. ✅ Runs all tests
2. ✅ Builds binaries for all platforms (Linux, macOS, Windows)
3. ✅ Creates package formats (.deb, .rpm, .apk)
4. ✅ Generates checksums
5. ✅ Creates a GitHub Release
6. ✅ Uploads all artifacts
7. ✅ Generates changelog from commits

#### 2. Pre-release Detection

The system automatically detects pre-releases:
- Tags containing `alpha`, `beta`, or `rc` are marked as pre-releases
- Pre-releases don't appear as "Latest" release
- Users can opt-in to download pre-releases from the releases page

#### 3. Release Workflow (for Contributors)

**For stable releases:**
```bash
# Using GitFlow
git flow release start 1.1.0
# Make final changes, update docs...
git flow release finish 1.1.0
git push origin master develop --tags
```

**For release candidates:**
```bash
# Create RC before finishing release
git flow release start 1.1.0
# ... make changes ...
git tag -a v1.1.0-rc.1 -m "Release candidate 1"
git push origin v1.1.0-rc.1
# Test the RC...
# When ready, finish release normally
git flow release finish 1.1.0
git push origin master develop --tags
```

**For experimental features:**
```bash
# On a feature branch
git tag -a v1.2.0-alpha.1 -m "Early alpha of new feature"
git push origin v1.2.0-alpha.1
```

### Available Platforms

Each release includes binaries for:

**Linux:**
- amd64 (64-bit Intel/AMD)
- arm64 (64-bit ARM)
- Formats: .tar.gz, .deb, .rpm, .apk

**macOS:**
- amd64 (Intel Macs)
- arm64 (Apple Silicon M1/M2/M3)
- Format: .tar.gz

**Windows:**
- amd64 (64-bit)
- arm64 (ARM64)
- Format: .zip

### Choosing the Right Release

- **Production use**: Use the latest stable release (no suffix)
- **Testing new features**: Use the latest RC (release candidate)
- **Early preview**: Use beta or alpha releases (expect bugs!)

### Downloading Releases

Visit the [releases page](https://github.com/Florian-Varrin/Kook-cli/releases) to:
- View all releases and their changelogs
- Download binaries for your platform
- See checksums for verification
- Access pre-release versions

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
# DevGen CLI - Development Generation Tool with Charm UI

A powerful command-line interface for development artifact generation, project management, and workflow orchestration, built with Go and the Charm ecosystem for beautiful terminal user interfaces.

## 🎯 Overview

DevGen CLI provides an interactive terminal experience for:
- **Playbook Execution**: Run development workflows with real-time progress tracking
- **Project Management**: Initialize, configure, and manage development projects
- **Template System**: Install and create reusable project templates
- **Development Server**: Start and monitor development servers with live reload
- **Configuration Management**: Interactive configuration editing and validation

## 🌟 Features

### 🚀 Interactive Playbook Execution
- Real-time progress tracking with animated spinners
- Multi-branch parallel execution support
- Conditional step execution with dependency management
- Comprehensive logging and error handling
- Pause/resume functionality
- Multiple view modes (overview, branches, logs, progress)

### 🎨 Beautiful Terminal UI
- **Cyber Theme**: High-contrast neon colors for maximum visibility
- **Pastel Theme**: Soft colors for comfortable extended use
- Responsive layouts that adapt to terminal size
- Keyboard navigation with vim-style bindings
- Help system with context-aware shortcuts
- Progress bars and status indicators

### 📦 Template Management
- Browse available templates with filtering
- Interactive template installation with progress tracking
- Create custom templates with guided setup
- Version management and automatic updates
- Local and remote template repositories

### 🏗️ Project Lifecycle
- Interactive project initialization
- Real-time project status dashboard
- Artifact generation (APIs, components, configs)
- Health monitoring and metrics
- Integration with DevQ.ai standards

### ⚙️ Advanced Configuration
- Interactive configuration editor
- YAML validation and syntax highlighting
- Environment-specific profiles
- Hot reload configuration changes
- Backup and restore settings

## 🛠️ Installation

### Prerequisites
- Go 1.21 or later
- Terminal with true color support (recommended)

### Build from Source
```bash
git clone https://github.com/devq-ai/devgen-cli.git
cd devgen-cli
go mod tidy
go build -o devgen .
```

### Install Binary
```bash
# Install to /usr/local/bin
sudo cp devgen /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)
```

## 🚀 Quick Start

### Initialize Configuration
```bash
# Create default configuration
devgen config init

# Edit configuration interactively
devgen config edit
```

### Run a Playbook
```bash
# List available playbooks
devgen playbook list

# Run a specific playbook
devgen playbook run my-workflow.yaml

# Create a new playbook
devgen playbook create
```

### Project Management
```bash
# Initialize a new project
devgen project init my-app

# Show project status
devgen project status

# Generate artifacts
devgen project generate api
```

### Template Operations
```bash
# List available templates
devgen template list

# Install a template
devgen template install fastapi-basic

# Create a custom template
devgen template create
```

### Development Server
```bash
# Start development server with monitoring
devgen server start

# Check server status
devgen server status

# Stop server
devgen server stop
```

## 📋 Command Reference

### Global Flags
```bash
-c, --config string      Config file path (default: ~/.devgen/config.yaml)
-v, --verbose           Enable verbose logging
    --log-level string  Log level (debug, info, warn, error)
-o, --output string     Output directory for generated files
-i, --interactive       Enable interactive mode
```

### Playbook Commands
```bash
devgen playbook run <file>        # Execute playbook with interactive UI
devgen playbook validate <file>   # Validate playbook configuration
devgen playbook create            # Create new playbook interactively
devgen playbook list              # List available playbooks
```

### Template Commands
```bash
devgen template list              # Show available templates
devgen template install <name>    # Install template with progress
devgen template create            # Create new template
```

### Project Commands
```bash
devgen project init [name]        # Initialize new project
devgen project status             # Show project status dashboard
devgen project generate <type>    # Generate project artifacts
```

### Server Commands
```bash
devgen server start               # Start development server
devgen server stop                # Stop development server
devgen server status              # Show server status
```

### Configuration Commands
```bash
devgen config init                # Initialize default configuration
devgen config edit                # Edit configuration interactively
devgen config show                # Display current configuration
```

## 🎨 UI Themes

### Cyber Theme (Default)
Perfect for high-energy development sessions with maximum contrast:
- **Primary**: Electric Cyan (#00ffff)
- **Secondary**: Neon Pink (#ff0080)
- **Success**: Matrix Green (#00ff00)
- **Warning**: Laser Yellow (#ffff00)
- **Background**: Void Black (#000000)

### Pastel Theme
Comfortable for extended use with gentle colors:
- **Primary**: Sky Blue (#b3e5fc)
- **Secondary**: Blush Pink (#ffb3ba)
- **Success**: Mint Green (#a8e6a3)
- **Warning**: Cream Yellow (#fff9c4)
- **Background**: Midnight Black (#000000)

### Theme Configuration
```yaml
ui:
  theme:
    name: "cyber"  # or "pastel"
    dark_mode: true
    colors:
      primary: "#00ffff"
      secondary: "#ff0080"
      # ... custom colors
```

## ⌨️ Keyboard Shortcuts

### Global Navigation
- `q` / `Ctrl+C`: Quit application
- `?`: Toggle help panel
- `Tab`: Switch between views
- `↑`/`↓` or `k`/`j`: Navigate up/down
- `←`/`→` or `h`/`l`: Navigate left/right

### Playbook Execution
- `e`: Execute/start playbook
- `p`: Pause/resume execution
- `r`: Reset playbook to initial state
- `d`: Show detailed step information
- `Space`: Toggle step selection
- `Enter`: Select/confirm action

### List Navigation
- `/`: Filter/search items
- `Enter`: Select item
- `Esc`: Clear filter/cancel

## 📁 Configuration Structure

### Main Configuration (`~/.devgen/config.yaml`)
```yaml
version: "1.0.0"
devgen:
  default_output_dir: "./output"
  default_template: "fastapi-basic"
  auto_save: true

templates:
  repository: "https://github.com/devq-ai/templates.git"
  local_path: "~/.devgen/templates"
  auto_update: true

projects:
  workspace: "~/projects"
  default_structure:
    directories: ["src", "tests", "docs"]

servers:
  default:
    host: "localhost"
    port: 8080
    auto_restart: true

logging:
  level: "info"
  format: "json"
  colors: true

ui:
  theme:
    name: "cyber"
    dark_mode: true
```

### Playbook Structure
```yaml
name: "Development Workflow"
version: "1.0.0"
description: "Full-stack development pipeline"
author: "DevQ.ai Team"

variables:
  project_name: "my-app"
  environment: "development"

branches:
  - name: "backend-setup"
    description: "Initialize backend services"
    parallel: false
    steps:
      - name: "setup-database"
        agent: "database-manager"
        action: "create and configure database"
        condition: "start"
        timeout: "5m"

      - name: "start-api"
        agent: "api-server"
        action: "start FastAPI server"
        condition: "database-ready"
        depends: ["setup-database"]

  - name: "frontend-setup"
    description: "Initialize frontend application"
    parallel: true
    steps:
      - name: "install-deps"
        agent: "package-manager"
        action: "install Node.js dependencies"
        condition: "start"

      - name: "start-dev-server"
        agent: "dev-server"
        action: "start Next.js development server"
        condition: "deps-installed"
```

## 🔧 Development

### Project Structure
```
cli_frontend/
├── main.go              # CLI entry point and command structure
├── ui.go                # Main UI components and playbook execution
├── engine.go            # Execution engine with step management
├── config.go            # Configuration management and YAML handling
├── components/
│   └── ui_components.go # Additional UI components (lists, forms, etc.)
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── README.md            # This file
├── style_guide.md       # UI design system documentation
├── charm_libraries.md   # Charm ecosystem component reference
└── cli_testing_strategies.md # Testing approaches and examples
```

### Key Dependencies
- **Bubble Tea**: TUI framework with Elm architecture
- **Lip Gloss**: CSS-like styling for terminal layouts
- **Bubbles**: Pre-built UI components (lists, inputs, tables)
- **Glamour**: Markdown rendering with syntax highlighting
- **Huh**: Terminal forms and prompts
- **Cobra**: CLI command structure and parsing
- **YAML v3**: Configuration file parsing

### Building
```bash
# Development build
go build -o devgen .

# Production build with optimizations
go build -ldflags="-s -w" -o devgen .

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o devgen-linux .
GOOS=darwin GOARCH=amd64 go build -o devgen-macos .
GOOS=windows GOARCH=amd64 go build -o devgen.exe .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestPlaybookExecution

# Benchmark tests
go test -bench=.
```

## 🎯 Advanced Usage

### Custom Agents
Create custom agents for specific tasks:

```yaml
# In playbook
steps:
  - name: "custom-deployment"
    agent: "kubernetes-deployer"
    action: "deploy to production cluster"
    parameters:
      namespace: "production"
      replicas: 3
      image: "my-app:latest"
```

### Environment Variables
Configure runtime behavior:

```bash
export DEVGEN_CONFIG_DIR="./custom-config"
export DEVGEN_LOG_LEVEL="debug"
export DEVGEN_THEME="pastel"
export DEVGEN_AUTO_UPDATE="false"
```

### Integration with CI/CD
Use in automated pipelines:

```bash
# Non-interactive mode
devgen playbook run --non-interactive deployment.yaml

# JSON output for parsing
devgen project status --output json

# Exit codes for pipeline control
devgen playbook validate *.yaml || exit 1
```

### Custom Templates
Create reusable project templates:

```yaml
# template.yaml
name: "My Custom Template"
description: "Custom project template"
files:
  - src: "templates/main.go.tmpl"
    dest: "main.go"
  - src: "templates/config.yaml.tmpl"
    dest: "config.yaml"
variables:
  - name: "project_name"
    description: "Name of the project"
    required: true
  - name: "port"
    description: "Server port"
    default: "8080"
```

## 🐛 Troubleshooting

### Common Issues

**Configuration not found**
```bash
# Initialize default configuration
devgen config init

# Verify configuration location
devgen config show
```

**Playbook validation errors**
```bash
# Validate playbook syntax
devgen playbook validate my-workflow.yaml

# Check for missing dependencies
devgen playbook run --dry-run my-workflow.yaml
```

**UI rendering issues**
```bash
# Check terminal capabilities
echo $TERM
tput colors

# Force basic rendering
TERM=xterm-256color devgen playbook run workflow.yaml
```

**Permission errors**
```bash
# Check file permissions
ls -la ~/.devgen/

# Fix permissions
chmod 755 ~/.devgen/
chmod 644 ~/.devgen/config.yaml
```

### Debug Mode
Enable verbose logging for troubleshooting:

```bash
# Enable debug logging
devgen --log-level debug playbook run workflow.yaml

# Save logs to file
devgen --verbose playbook run workflow.yaml 2>&1 | tee debug.log
```

## 📚 Documentation

- [Style Guide](style_guide.md) - UI design system and theming
- [Charm Libraries](charm_libraries.md) - Component library reference
- [Testing Strategies](cli_testing_strategies.md) - Testing approaches and examples
- [API Reference](https://pkg.go.dev/github.com/devq-ai/devgen-cli) - Go package documentation

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/devgen-cli.git`
3. Install dependencies: `go mod tidy`
4. Create a feature branch: `git checkout -b feature-name`
5. Make your changes and add tests
6. Run tests: `go test ./...`
7. Submit a pull request

### Code Style
- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comments for public APIs
- Include tests for new functionality
- Update documentation for user-facing changes

### UI Guidelines
- Follow the established design system
- Test with both cyber and pastel themes
- Ensure keyboard navigation works properly
- Test with different terminal sizes
- Validate accessibility features

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) - For the amazing terminal UI framework
- [DevQ.ai](https://devq.ai/) - For the development standards and workflow patterns
- [Go team](https://golang.org/) - For the excellent programming language
- Contributors and community members

## 🔗 Links

- [GitHub Repository](https://github.com/devq-ai/devgen-cli)
- [Issue Tracker](https://github.com/devq-ai/devgen-cli/issues)
- [Discussions](https://github.com/devq-ai/devgen-cli/discussions)
- [DevQ.ai Documentation](https://docs.devq.ai/)
- [Charm Documentation](https://charm.sh/docs/)

# DevGen CLI - Complete User Guide

A comprehensive guide to using the DevGen CLI tool for development workflow automation, project management, and beautiful terminal interfaces.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Features](#core-features)
- [Command Reference](#command-reference)
- [Playbook System](#playbook-system)
- [UI Navigation](#ui-navigation)
- [Configuration](#configuration)
- [Templates](#templates)
- [Project Management](#project-management)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)
- [FAQ](#faq)

---

## Overview

DevGen CLI is a powerful command-line interface built with Go and the Charm ecosystem that provides:

- **🎨 Beautiful Terminal UI** - Modern, responsive interfaces with cyber and pastel themes
- **📋 Playbook Execution** - Automated workflow orchestration with real-time monitoring
- **🏗️ Project Management** - Complete project lifecycle management
- **📦 Template System** - Reusable project templates and scaffolding
- **⚙️ Configuration Management** - Interactive configuration editing and validation
- **🔄 Development Server** - Built-in development server with hot reload

### Key Benefits

- **Visual Workflow Management** - See your development processes in action
- **Automated Task Orchestration** - Execute complex multi-step workflows
- **Real-time Progress Tracking** - Monitor execution with live updates
- **Interactive Configuration** - Edit settings with guided forms
- **Cross-platform Compatibility** - Works on macOS, Linux, and Windows

---

## Installation

### Prerequisites

- **Go 1.21+** - For building from source
- **Terminal with true color support** - For optimal visual experience
- **Git** - For template management and version control

### Method 1: Pre-built Binary (Recommended)

```bash
# Download the latest release
curl -L https://github.com/devq-ai/devgen-cli/releases/latest/download/devgen-$(uname -s)-$(uname -m) -o devgen

# Make executable
chmod +x devgen

# Move to PATH
sudo mv devgen /usr/local/bin/
```

### Method 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/devq-ai/devgen-cli.git
cd devgen-cli

# Run automated setup
chmod +x setup.sh
./setup.sh

# Build using Makefile
make build

# Install binary
make install
```

### Method 3: Go Install

```bash
go install github.com/devq-ai/devgen-cli/cmd/devgen@latest
```

### Verification

```bash
# Check installation
devgen --version

# View help
devgen --help
```

---

## Quick Start

### 1. Initialize Configuration

```bash
# Create default configuration
devgen config init

# Edit configuration interactively
devgen config edit
```

### 2. Run Your First Playbook

```bash
# Use the example playbook
devgen playbook run example-playbook.yaml

# Or create a new one
devgen playbook create
```

### 3. Initialize a Project

```bash
# Initialize new project
devgen project init my-awesome-app

# Check project status
devgen project status
```

### 4. Explore Templates

```bash
# List available templates
devgen template list

# Install a template
devgen template install fastapi-basic
```

---

## Core Features

### 🎨 Beautiful Terminal UI

DevGen CLI features two distinctive themes:

#### Cyber Theme (High Energy)
- **Electric colors** - Neon cyan, matrix green, laser yellow
- **High contrast** - Maximum visibility for intense development sessions
- **Retro-futuristic** - Perfect for late-night coding sessions

#### Pastel Theme (Comfortable)
- **Soft colors** - Sky blue, mint green, blush pink
- **Eye-friendly** - Comfortable for extended use
- **Modern aesthetic** - Clean and professional appearance

### 📋 Interactive Playbook Execution

- **Real-time progress tracking** with animated progress bars
- **Parallel execution** support for independent tasks
- **Dependency management** with automatic ordering
- **Error handling** with retry mechanisms
- **Live logging** with searchable output
- **Pause/resume** functionality for long-running workflows

### 🏗️ Project Management

- **Interactive initialization** with guided setup
- **Status dashboards** showing project health
- **Artifact generation** for APIs, components, and configurations
- **Environment management** with profile support
- **Health monitoring** with real-time metrics

### 📦 Template System

- **Template discovery** with filtering and search
- **Interactive installation** with progress tracking
- **Custom template creation** with guided setup
- **Version management** with automatic updates
- **Local and remote** template repositories

---

## Command Reference

### Global Options

All commands support these global options:

```bash
-c, --config string      Configuration file path (default: ~/.devgen/config.yaml)
-v, --verbose           Enable verbose logging
    --log-level string  Log level: debug, info, warn, error (default: info)
-o, --output string     Output directory for generated files (default: current directory)
-i, --interactive       Force interactive mode
    --theme string      UI theme: cyber, pastel (default: cyber)
    --no-color          Disable colored output
```

### Command Structure

```bash
devgen [global-options] <command> [command-options] [arguments]
```

---

## Playbook System

Playbooks are YAML files that define automated workflows with steps, dependencies, and conditions.

### Basic Playbook Structure

```yaml
name: "My Development Workflow"
version: "1.0.0"
description: "Automated development pipeline"
author: "Your Name"

variables:
  project_name: "my-app"
  environment: "development"

branches:
  - name: "setup"
    description: "Initial setup tasks"
    parallel: false
    steps:
      - name: "install-dependencies"
        agent: "package-manager"
        action: "install project dependencies"
        condition: "start"
        timeout: "5m"
```

### Playbook Commands

#### Run Playbook
```bash
devgen playbook run workflow.yaml
```

**Features:**
- Interactive execution with real-time progress
- Multiple view modes (overview, branches, logs, progress)
- Keyboard navigation with vim-style bindings
- Pause/resume functionality
- Error handling with retry options

#### Validate Playbook
```bash
devgen playbook validate workflow.yaml
```

**Validation checks:**
- YAML syntax validation
- Required field verification
- Dependency cycle detection
- Agent availability confirmation
- Variable substitution validation

#### Create Playbook
```bash
devgen playbook create
```

**Interactive creation process:**
1. Basic information (name, description, author)
2. Variable definitions
3. Branch structure setup
4. Step configuration
5. Dependency mapping
6. Validation and preview

#### List Playbooks
```bash
devgen playbook list
```

**Features:**
- Searchable list with filtering
- Detailed playbook information
- Execution history
- Quick actions (run, edit, delete)

### Playbook Components

#### Variables
Define reusable values throughout your playbook:

```yaml
variables:
  project_name: "myapp"
  api_port: "8000"
  database_url: "postgresql://localhost:5432/${project_name}"
```

#### Branches
Group related steps that can run in parallel or sequence:

```yaml
branches:
  - name: "backend"
    parallel: false
    timeout: "10m"
    steps: [...]

  - name: "frontend"
    parallel: true
    prerequisites: ["backend"]
    steps: [...]
```

#### Steps
Individual tasks within a branch:

```yaml
steps:
  - name: "setup-database"
    agent: "database-manager"
    action: "create and configure database"
    condition: "start"
    timeout: "5m"
    retries: 3
    parameters:
      database_name: "${project_name}_db"
      port: 5432
    environment:
      POSTGRES_PASSWORD: "secure_password"
    artifacts:
      - name: "database-logs"
        path: "./logs/postgres.log"
        type: "log"
```

#### Conditions and Dependencies
Control step execution flow:

```yaml
steps:
  - name: "build-app"
    condition: "dependencies-installed"
    depends: ["install-deps", "setup-config"]

  - name: "run-tests"
    condition: "build-complete"
    depends: ["build-app"]
```

---

## UI Navigation

### Keyboard Shortcuts

#### Global Navigation
- `q` / `Ctrl+C` - Quit application
- `?` - Toggle help panel
- `Tab` - Switch between views
- `↑`/`↓` or `k`/`j` - Navigate up/down
- `←`/`→` or `h`/`l` - Navigate left/right

#### Playbook Execution
- `e` - Execute/start playbook
- `p` - Pause/resume execution
- `r` - Reset playbook to initial state
- `s` - Stop execution
- `d` - Show detailed step information
- `Space` - Toggle step selection
- `Enter` - Select/confirm action

#### List Navigation
- `/` - Filter/search items
- `Enter` - Select item
- `Esc` - Clear filter/cancel
- `u` - Update/refresh list

#### View Modes
- `1` - Overview mode
- `2` - Branches mode
- `3` - Logs mode
- `4` - Progress mode

### View Modes

#### Overview Mode
- High-level playbook status
- Branch execution progress
- Overall timing information
- Quick navigation to details

#### Branches Mode
- Detailed branch information
- Step-by-step progress
- Dependency visualization
- Error highlighting

#### Logs Mode
- Real-time log streaming
- Searchable log output
- Log level filtering
- Copy/export functionality

#### Progress Mode
- Visual progress indicators
- Performance metrics
- Resource usage monitoring
- Execution timeline

---

## Configuration

### Configuration File Structure

DevGen uses a YAML configuration file located at `~/.devgen/config.yaml`:

```yaml
version: "1.0.0"

# Core settings
devgen:
  default_output_dir: "./output"
  default_template: "fastapi-basic"
  auto_save: true
  log_level: "info"

# Template configuration
templates:
  repository: "https://github.com/devq-ai/templates.git"
  local_path: "~/.devgen/templates"
  auto_update: true
  update_interval: "24h"

# Project settings
projects:
  workspace: "~/projects"
  default_structure:
    directories: ["src", "tests", "docs", "config"]
    files: ["README.md", ".gitignore", "requirements.txt"]

# Server configuration
servers:
  default:
    host: "localhost"
    port: 8080
    auto_restart: true
    hot_reload: true

# UI settings
ui:
  theme:
    name: "cyber"  # or "pastel"
    dark_mode: true
    animations: true
    sound_effects: false

  # Color customization
  colors:
    primary: "#00ffff"
    secondary: "#ff0080"
    success: "#00ff00"
    warning: "#ffff00"
    error: "#ff0040"

# Logging configuration
logging:
  level: "info"
  format: "json"
  colors: true
  file: "./logs/devgen.log"
  max_size: "10MB"
  max_backups: 5

# Agent configuration
agents:
  timeout: "30s"
  retries: 3
  log_level: "info"
```

### Configuration Commands

#### Initialize Configuration
```bash
devgen config init
```

Creates a default configuration file with sensible defaults.

#### Edit Configuration
```bash
devgen config edit
```

Opens an interactive editor with:
- Syntax highlighting
- Validation on save
- Help tooltips
- Auto-completion
- Error detection

#### Show Configuration
```bash
devgen config show
```

Displays current configuration with:
- Syntax highlighting
- Section navigation
- Search functionality
- Value validation status

### Environment Variables

Override configuration with environment variables:

```bash
# Core settings
export DEVGEN_CONFIG_DIR="./custom-config"
export DEVGEN_LOG_LEVEL="debug"
export DEVGEN_OUTPUT_DIR="./output"

# UI settings
export DEVGEN_THEME="pastel"
export DEVGEN_NO_COLOR="true"

# Template settings
export DEVGEN_TEMPLATES_REPO="https://github.com/myorg/templates.git"
export DEVGEN_AUTO_UPDATE="false"
```

---

## Templates

Templates provide reusable project structures and configurations.

### Template Structure

```
template-name/
├── template.yaml          # Template metadata
├── files/                 # Template files
│   ├── src/
│   │   └── main.py.tmpl
│   ├── tests/
│   │   └── test_main.py.tmpl
│   └── README.md.tmpl
├── hooks/                 # Lifecycle hooks
│   ├── pre-install.sh
│   └── post-install.sh
└── README.md              # Template documentation
```

### Template Commands

#### List Templates
```bash
devgen template list
```

**Features:**
- Filter by category, author, or tags
- Search by name or description
- View template details
- Check compatibility
- See installation status

#### Install Template
```bash
devgen template install template-name

# Interactive installation
devgen template install --interactive template-name

# Specify version
devgen template install template-name@v1.2.0
```

**Installation process:**
1. Template validation
2. Dependency checking
3. Variable collection
4. File processing
5. Hook execution
6. Verification

#### Create Template
```bash
devgen template create
```

**Interactive creation:**
1. Template metadata (name, description, author)
2. Variable definitions
3. File structure setup
4. Hook configuration
5. Documentation generation
6. Testing and validation

### Template Metadata

```yaml
# template.yaml
name: "FastAPI Basic"
description: "Basic FastAPI application with DevQ.ai standards"
version: "1.0.0"
author: "DevQ.ai Team"
license: "MIT"
homepage: "https://github.com/devq-ai/templates"

# Template requirements
requirements:
  devgen_version: ">=1.0.0"
  system:
    - "python>=3.8"
    - "node>=16.0.0"

# Template variables
variables:
  - name: "project_name"
    description: "Name of the project"
    type: "string"
    required: true
    validation: "^[a-zA-Z][a-zA-Z0-9_-]*$"

  - name: "author_name"
    description: "Author name"
    type: "string"
    default: "DevQ.ai Team"

  - name: "port"
    description: "API server port"
    type: "integer"
    default: 8000
    validation: "1024-65535"

# File processing
files:
  - src: "**/*.tmpl"
    dest: "."
    template: true
    exclude: ["*.test.tmpl"]

  - src: "static/**/*"
    dest: "static/"
    template: false

# Lifecycle hooks
hooks:
  pre_install:
    - "hooks/pre-install.sh"

  post_install:
    - "hooks/post-install.sh"
    - "hooks/verify-installation.sh"

# Categories and tags
categories: ["web", "api", "python"]
tags: ["fastapi", "async", "rest", "openapi"]
```

---

## Project Management

### Project Structure

DevGen projects follow a standardized structure:

```
my-project/
├── .devgen/
│   ├── config.yaml        # Project-specific configuration
│   ├── state.json         # Project state and metadata
│   └── templates/         # Local template cache
├── src/                   # Source code
├── tests/                 # Test files
├── docs/                  # Documentation
├── config/                # Configuration files
├── scripts/               # Build and deployment scripts
├── .gitignore
├── README.md
└── requirements.txt       # Dependencies
```

### Project Commands

#### Initialize Project
```bash
devgen project init [project-name]
```

**Interactive initialization:**
1. Project name and description
2. Template selection
3. Configuration options
4. Directory structure setup
5. Git repository initialization
6. Initial commit

#### Project Status
```bash
devgen project status
```

**Status dashboard includes:**
- Project information
- Health metrics
- Recent activity
- Dependencies status
- Build status
- Test results
- Deployment status

#### Generate Artifacts
```bash
devgen project generate <type>
```

**Available generators:**
- `api` - Generate API endpoints
- `model` - Generate data models
- `test` - Generate test files
- `config` - Generate configuration files
- `docs` - Generate documentation
- `deployment` - Generate deployment configs

**Example:**
```bash
# Generate API endpoint
devgen project generate api --name users --methods get,post,put,delete

# Generate model
devgen project generate model --name User --fields name:string,email:string,age:integer

# Generate tests
devgen project generate test --target src/api/users.py
```

### Project Configuration

```yaml
# .devgen/config.yaml
project:
  name: "my-awesome-app"
  description: "An awesome application"
  version: "1.0.0"
  author: "Your Name"
  license: "MIT"

# Build configuration
build:
  output_dir: "./build"
  source_dir: "./src"
  test_dir: "./tests"

# Development settings
development:
  auto_reload: true
  port: 8000
  debug: true

# Dependencies
dependencies:
  python: ">=3.8"
  node: ">=16.0.0"

# Custom generators
generators:
  api:
    template: "fastapi-endpoint"
    output_dir: "./src/api"

  model:
    template: "pydantic-model"
    output_dir: "./src/models"
```

---

## Advanced Usage

### Custom Agents

Create custom agents for specific tasks:

```yaml
# In playbook
agents:
  custom-deployer:
    type: "shell"
    command: "./scripts/deploy.sh"
    environment:
      - "DEPLOY_ENV=production"
      - "API_KEY=${API_KEY}"

  database-migrator:
    type: "python"
    module: "scripts.migrate"
    function: "run_migrations"
    parameters:
      database_url: "${DATABASE_URL}"
      migration_path: "./migrations"

# Use in steps
steps:
  - name: "deploy-application"
    agent: "custom-deployer"
    action: "deploy to production"
    parameters:
      version: "v1.2.0"
      environment: "production"
```

### Environment Profiles

Manage different environments:

```yaml
# config.yaml
environments:
  development:
    database_url: "postgresql://localhost:5432/myapp_dev"
    api_port: 8000
    debug: true

  staging:
    database_url: "postgresql://staging-db:5432/myapp"
    api_port: 80
    debug: false

  production:
    database_url: "${DATABASE_URL}"
    api_port: 80
    debug: false
    monitoring: true
```

```bash
# Use specific environment
devgen --env production playbook run deploy.yaml
```

### Hooks and Callbacks

Add custom logic at different stages:

```yaml
# In playbook
hooks:
  before_execution:
    - type: "shell"
      command: "echo 'Starting workflow'"

  after_step:
    - type: "notification"
      service: "slack"
      message: "Step {{.step_name}} completed"

  on_error:
    - type: "shell"
      command: "./scripts/cleanup.sh"

  after_execution:
    - type: "report"
      template: "./templates/report.md"
      output: "./reports/execution-{{.timestamp}}.md"
```

### Integration with CI/CD

Use DevGen in automated pipelines:

```yaml
# .github/workflows/devgen.yml
name: DevGen Workflow

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  devgen:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Setup DevGen
        run: |
          curl -L https://github.com/devq-ai/devgen-cli/releases/latest/download/devgen-Linux-x86_64 -o devgen
          chmod +x devgen
          sudo mv devgen /usr/local/bin/

      - name: Run DevGen Playbook
        run: |
          devgen playbook run --non-interactive ci-pipeline.yaml
        env:
          DEVGEN_LOG_LEVEL: info
          DEVGEN_OUTPUT_DIR: ./reports
```

---

## Troubleshooting

### Common Issues

#### Configuration Not Found
```bash
# Problem: Configuration file not found
# Solution: Initialize default configuration
devgen config init

# Verify configuration location
devgen config show
```

#### Playbook Validation Errors
```bash
# Problem: Invalid playbook syntax
# Solution: Validate and fix issues
devgen playbook validate my-workflow.yaml

# Check for detailed error information
devgen --log-level debug playbook validate my-workflow.yaml
```

#### UI Rendering Issues
```bash
# Problem: UI not displaying correctly
# Solution: Check terminal capabilities
echo $TERM
tput colors

# Force basic rendering
TERM=xterm-256color devgen playbook run workflow.yaml

# Disable colors if needed
devgen --no-color playbook run workflow.yaml
```

#### Template Installation Failures
```bash
# Problem: Template installation fails
# Solution: Check template repository access
devgen template list --verbose

# Clear template cache
rm -rf ~/.devgen/templates
devgen template list --update
```

#### Performance Issues
```bash
# Problem: Slow execution
# Solution: Enable performance profiling
devgen --log-level debug playbook run workflow.yaml 2>&1 | grep -i "timing\|performance"

# Check resource usage
top -p $(pgrep devgen)
```

### Debug Mode

Enable comprehensive debugging:

```bash
# Enable verbose logging
devgen --verbose --log-level debug playbook run workflow.yaml

# Save debug output
devgen --verbose playbook run workflow.yaml 2>&1 | tee debug.log

# Profile memory usage
devgen --profile-memory playbook run workflow.yaml
```

### Log Analysis

DevGen provides detailed logging:

```bash
# View recent logs
tail -f ~/.devgen/logs/devgen.log

# Search logs for errors
grep -i error ~/.devgen/logs/devgen.log

# Filter by log level
grep -i "level=error" ~/.devgen/logs/devgen.log
```

---

## Best Practices

### Playbook Design

#### Structure Organization
```yaml
# Good: Clear, logical structure
name: "Backend API Setup"
description: "Complete backend API initialization with database and authentication"

branches:
  - name: "infrastructure"
    description: "Set up core infrastructure"
    steps: [...]

  - name: "application"
    description: "Deploy application services"
    prerequisites: ["infrastructure"]
    steps: [...]
```

#### Error Handling
```yaml
# Good: Comprehensive error handling
steps:
  - name: "database-setup"
    agent: "database-manager"
    action: "create database"
    timeout: "5m"
    retries: 3

    on_error:
      - type: "cleanup"
        command: "docker-compose down"
      - type: "log"
        message: "Database setup failed - check logs"
      - type: "exit"
        code: 1
```

#### Resource Management
```yaml
# Good: Proper resource limits
steps:
  - name: "build-application"
    agent: "builder"
    action: "build application"
    timeout: "10m"

    resources:
      memory: "2GB"
      cpu: "2"
      disk: "1GB"
```

### Template Design

#### Variable Validation
```yaml
# Good: Comprehensive validation
variables:
  - name: "port"
    description: "Server port"
    type: "integer"
    required: true
    validation: "1024-65535"
    example: "8000"

  - name: "database_url"
    description: "Database connection URL"
    type: "string"
    required: true
    validation: "^postgresql://.*"
    example: "postgresql://localhost:5432/myapp"
```

#### File Organization
```
# Good: Logical file structure
template/
├── template.yaml
├── files/
│   ├── backend/
│   │   ├── src/
│   │   └── tests/
│   └── frontend/
│       ├── src/
│       └── tests/
└── docs/
    └── README.md
```

### Performance Optimization

#### Parallel Execution
```yaml
# Good: Utilize parallelism
branches:
  - name: "backend"
    parallel: true
    steps: [...]

  - name: "frontend"
    parallel: true
    steps: [...]

  - name: "integration"
    parallel: false
    prerequisites: ["backend", "frontend"]
    steps: [...]
```

#### Resource Caching
```yaml
# Good: Cache expensive operations
steps:
  - name: "install-dependencies"
    agent: "package-manager"
    action: "install packages"
    cache:
      key: "deps-{{.checksum}}"
      paths: ["node_modules", "venv"]
```

### Security Considerations

#### Sensitive Data
```yaml
# Good: Use environment variables for secrets
variables:
  database_password: "${DATABASE_PASSWORD}"
  api_key: "${API_KEY}"

# Never hardcode secrets
# Bad: database_password: "supersecret123"
```

#### File Permissions
```yaml
# Good: Set appropriate permissions
steps:
  - name: "create-config"
    agent: "file-manager"
    action: "create configuration file"
    parameters:
      permissions: "0600"  # Read/write for owner only
```

---

## FAQ

### General Questions

**Q: What makes DevGen CLI different from other automation tools?**
A: DevGen CLI combines powerful automation capabilities with a beautiful terminal UI, making complex workflows visual and interactive. It's designed specifically for development workflows with features like real-time progress tracking, dependency management, and integrated project management.

**Q: Can I use DevGen CLI in CI/CD pipelines?**
A: Yes! DevGen CLI supports non-interactive mode and provides structured output formats perfect for CI/CD integration. Use `--non-interactive` flag and `--output json` for machine-readable results.

**Q: How do I contribute templates to the community?**
A: Submit templates to the [DevGen Templates Repository](https://github.com/devq-ai/devgen-templates). Include comprehensive documentation, examples, and tests.

### Technical Questions

**Q: How do I handle secrets in playbooks?**
A: Use environment variables and avoid hardcoding sensitive data. DevGen CLI supports secure variable substitution:
```yaml
variables:
  api_key: "${API_KEY}"
  database_password: "${DB_PASSWORD}"
```

**Q: Can I extend DevGen CLI with custom agents?**
A: Yes! Create custom agents using shell scripts, Python modules, or any executable. Define them in your playbook:
```yaml
agents:
  my-custom-agent:
    type: "shell"
    command: "./scripts/my-agent.sh"
```

**Q: How do I debug failing playbooks?**
A: Use debug mode for detailed information:
```bash
devgen --log-level debug playbook run workflow.yaml
```

**Q: Can I run DevGen CLI on Windows?**
A: Yes! DevGen CLI is cross-platform. Use Windows Terminal or WSL for the best experience.

### Performance Questions

**Q: How do I optimize playbook execution time?**
A: Use parallel execution, caching, and efficient dependency management:
```yaml
branches:
  - name: "parallel-tasks"
    parallel: true
    steps: [...]
```

**Q: What's the resource overhead of DevGen CLI?**
A: DevGen CLI is lightweight with minimal memory footprint. The UI components are optimized for performance and use efficient rendering techniques.

### Customization Questions

**Q: How do I create custom themes?**
A: Modify the color configuration in your config file:
```yaml
ui:
  theme:
    name: "custom"
  colors:
    primary: "#your-color"
    secondary: "#your-color"
```

**Q: Can I disable the UI and use DevGen CLI in headless mode?**
A: Yes! Use the `--non-interactive` flag for headless operation:
```bash
devgen --non-interactive playbook run workflow.yaml
```

---

## Resources

### Documentation
- [Official Documentation](https://docs.devq.ai/devgen-cli)
- [API Reference](https://pkg.go.dev/github.com/devq-ai/devgen-cli)
- [Template Documentation](https://github.com/devq-ai/devgen-templates)

### Community
- [GitHub Repository](https://github.com/devq-ai/devgen-cli)
- [Issue Tracker](https://github.com/devq-ai/devgen-cli/issues)
- [Discussions](https://github.com/devq-ai/devgen-cli/discussions)
- [Discord Community](https://discord.gg/devq-ai)

### Learning Resources
- [Video Tutorials](https://www.youtube.com/devqai)
- [Blog Posts](https://blog.devq.ai/tags/devgen-cli)
- [Example Playbooks](https://github.com/devq-ai/devgen-examples)

### Support
- [Support Documentation](https://docs.devq.ai/support)
- [Community Forum](https://forum.devq.ai)
- [Professional Support](https://devq.ai/support)

---

## Changelog

### Version 1.0.0
- Initial release with core functionality
- Playbook execution with interactive UI
- Template management system
- Project management features
- Configuration management
- Cyber and Pastel themes

### Upcoming Features
- Plugin system for custom extensions
- Cloud synchronization for templates
- Advanced analytics and reporting
- Integration with popular development tools
- Mobile companion app

---

*DevGen CLI - Empowering developers with beautiful, interactive automation tools.*

**Made with ❤️ by the DevQ.ai Team**

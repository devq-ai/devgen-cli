#!/bin/bash

# DevGen CLI Setup Script
# Automated setup for the DevGen CLI development environment

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="devgen-cli"
GO_VERSION="1.21"
REQUIRED_TOOLS=("go" "git" "make")
OPTIONAL_TOOLS=("docker" "docker-compose")

# Functions
print_banner() {
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                     DevGen CLI Setup                        ║"
    echo "║              Development Generation Tool                     ║"
    echo "║                   with Charm UI                             ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_step() {
    echo -e "${BLUE}🔧 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${PURPLE}ℹ️  $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Get OS information
get_os() {
    case "$(uname -s)" in
        Darwin) echo "macos" ;;
        Linux) echo "linux" ;;
        CYGWIN*|MINGW*|MSYS*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

# Get architecture
get_arch() {
    case "$(uname -m)" in
        x86_64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        i386|i686) echo "386" ;;
        *) echo "unknown" ;;
    esac
}

# Check Go version
check_go_version() {
    if ! command_exists go; then
        return 1
    fi

    local version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    local major=$(echo $version | cut -d. -f1)
    local minor=$(echo $version | cut -d. -f2)

    if [ "$major" -gt 1 ] || ([ "$major" -eq 1 ] && [ "$minor" -ge 21 ]); then
        return 0
    else
        return 1
    fi
}

# Install Go if needed
install_go() {
    local os=$(get_os)
    local arch=$(get_arch)

    if [ "$os" = "unknown" ] || [ "$arch" = "unknown" ]; then
        print_error "Unsupported OS/architecture combination"
        exit 1
    fi

    print_step "Installing Go ${GO_VERSION}..."

    local go_tar="go${GO_VERSION}.${os}-${arch}.tar.gz"
    local download_url="https://golang.org/dl/${go_tar}"

    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download Go
    if command_exists wget; then
        wget -q "$download_url"
    elif command_exists curl; then
        curl -s -L -O "$download_url"
    else
        print_error "Need wget or curl to download Go"
        exit 1
    fi

    # Extract and install
    if [ "$os" = "macos" ] || [ "$os" = "linux" ]; then
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "$go_tar"

        # Add to PATH
        if ! grep -q "/usr/local/go/bin" ~/.bashrc 2>/dev/null; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        fi

        if ! grep -q "/usr/local/go/bin" ~/.zshrc 2>/dev/null; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
        fi

        export PATH=$PATH:/usr/local/go/bin
    else
        print_error "Automatic Go installation not supported on Windows"
        print_info "Please install Go manually from https://golang.org/dl/"
        exit 1
    fi

    cd - > /dev/null
    rm -rf "$temp_dir"

    print_success "Go ${GO_VERSION} installed successfully"
}

# Check and install required tools
check_tools() {
    print_step "Checking required tools..."

    local missing_tools=()

    for tool in "${REQUIRED_TOOLS[@]}"; do
        if command_exists "$tool"; then
            if [ "$tool" = "go" ]; then
                if check_go_version; then
                    print_success "$tool ($(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+')) is available"
                else
                    print_warning "$tool version is too old (need 1.21+)"
                    missing_tools+=("$tool")
                fi
            else
                print_success "$tool is available"
            fi
        else
            print_warning "$tool is not installed"
            missing_tools+=("$tool")
        fi
    done

    # Install missing tools
    for tool in "${missing_tools[@]}"; do
        case "$tool" in
            "go")
                install_go
                ;;
            "git")
                print_error "Git is required but not installed"
                print_info "Please install Git: https://git-scm.com/downloads"
                exit 1
                ;;
            "make")
                local os=$(get_os)
                if [ "$os" = "macos" ]; then
                    if command_exists brew; then
                        brew install make
                    else
                        print_error "Please install Homebrew and run: brew install make"
                        exit 1
                    fi
                elif [ "$os" = "linux" ]; then
                    if command_exists apt-get; then
                        sudo apt-get update && sudo apt-get install -y build-essential
                    elif command_exists yum; then
                        sudo yum groupinstall -y "Development Tools"
                    elif command_exists pacman; then
                        sudo pacman -S base-devel
                    else
                        print_error "Please install make using your package manager"
                        exit 1
                    fi
                fi
                ;;
        esac
    done

    # Check optional tools
    print_step "Checking optional tools..."
    for tool in "${OPTIONAL_TOOLS[@]}"; do
        if command_exists "$tool"; then
            print_success "$tool is available"
        else
            print_warning "$tool is not installed (optional)"
        fi
    done
}

# Initialize project
init_project() {
    print_step "Initializing project structure..."

    # Create project directory if it doesn't exist
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi

    # Create necessary directories
    mkdir -p build dist logs reports cache tmp profiles
    mkdir -p tests/{unit,integration,e2e}
    mkdir -p scripts/{dev,build,deploy}
    mkdir -p docs/{api,guides,examples}
    mkdir -p config/{dev,prod,test}

    print_success "Project directories created"
}

# Install Go dependencies
install_dependencies() {
    print_step "Installing Go dependencies..."

    # Download and tidy modules
    go mod download
    go mod tidy

    print_success "Go dependencies installed"
}

# Install development tools
install_dev_tools() {
    print_step "Installing development tools..."

    local dev_tools=(
        "github.com/cosmtrek/air@latest"                           # Hot reload
        "github.com/golangci/golangci-lint/cmd/golangci-lint@latest" # Linting
        "golang.org/x/tools/cmd/goimports@latest"                  # Import formatting
        "github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"   # Security scanning
        "golang.org/x/tools/cmd/godoc@latest"                      # Documentation
        "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"       # Markdown docs
        "github.com/git-chglog/git-chglog/cmd/git-chglog@latest"   # Changelog generation
    )

    for tool in "${dev_tools[@]}"; do
        local tool_name=$(echo "$tool" | cut -d'/' -f3 | cut -d'@' -f1)
        print_info "Installing $tool_name..."
        if go install "$tool" 2>/dev/null; then
            print_success "$tool_name installed"
        else
            print_warning "Failed to install $tool_name"
        fi
    done
}

# Setup configuration files
setup_config() {
    print_step "Setting up configuration files..."

    # Create .env file if it doesn't exist
    if [ ! -f ".env" ]; then
        cat > .env << 'EOF'
# DevGen CLI Environment Configuration

# Development settings
DEBUG=true
ENVIRONMENT=development

# Logging
DEVGEN_LOG_LEVEL=info
DEVGEN_LOG_FORMAT=text

# Paths
DEVGEN_CONFIG_DIR=./.devgen
DEVGEN_OUTPUT_DIR=./output
DEVGEN_CACHE_DIR=./.cache

# Server settings
DEVGEN_SERVER_HOST=localhost
DEVGEN_SERVER_PORT=8080

# Terminal settings
TERM=xterm-256color
COLORTERM=truecolor

# Development tools
DEVGEN_WATCH_MODE=true
DEVGEN_AUTO_RELOAD=true
EOF
        print_success ".env file created"
    else
        print_info ".env file already exists"
    fi

    # Create .gitignore if it doesn't exist
    if [ ! -f ".gitignore" ]; then
        cat > .gitignore << 'EOF'
# DevGen CLI - Git Ignore

# Build outputs
build/
dist/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test outputs
*.test
coverage.out
coverage.html
*.prof

# Temporary files
tmp/
temp/
*.tmp
*.temp

# Logs
logs/
*.log

# Cache
cache/
.cache/

# Configuration
.env.local
.env.production

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# DevGen specific
.devgen/
reports/
profiles/
EOF
        print_success ".gitignore created"
    else
        print_info ".gitignore already exists"
    fi

    # Create Air configuration for hot reload
    if [ ! -f ".air.toml" ]; then
        cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
args_bin = []
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ."
delay = 1000
exclude_dir = ["assets", "tmp", "vendor", "testdata", "build", "dist", "cache"]
exclude_file = []
exclude_regex = ["_test.go"]
exclude_unchanged = false
follow_symlink = false
full_bin = ""
include_dir = []
include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml"]
include_file = []
kill_delay = "0s"
log = "build-errors.log"
rerun = false
rerun_delay = 500
send_interrupt = false
stop_on_root = false

[color]
app = ""
build = "yellow"
main = "magenta"
runner = "green"
watcher = "cyan"

[log]
main_only = false
time = false

[misc]
clean_on_exit = false

[screen]
clear_on_rebuild = false
keep_scroll = true
EOF
        print_success ".air.toml created"
    else
        print_info ".air.toml already exists"
    fi
}

# Build the project
build_project() {
    print_step "Building project..."

    if make build 2>/dev/null; then
        print_success "Project built successfully"
    else
        print_warning "Build failed, trying direct go build..."
        if go build -o build/devgen .; then
            print_success "Project built with go build"
        else
            print_error "Build failed"
            return 1
        fi
    fi
}

# Run tests
run_tests() {
    print_step "Running tests..."

    if make test 2>/dev/null; then
        print_success "All tests passed"
    else
        print_warning "Make test failed, trying direct go test..."
        if go test ./...; then
            print_success "Tests passed with go test"
        else
            print_warning "Some tests failed (this is normal for a new setup)"
        fi
    fi
}

# Setup development environment
setup_dev_env() {
    print_step "Setting up development environment..."

    # Add GOPATH/bin to PATH if not already there
    local go_bin_path="$(go env GOPATH)/bin"
    if [[ ":$PATH:" != *":$go_bin_path:"* ]]; then
        if [ -f ~/.bashrc ]; then
            echo "export PATH=\$PATH:$go_bin_path" >> ~/.bashrc
        fi
        if [ -f ~/.zshrc ]; then
            echo "export PATH=\$PATH:$go_bin_path" >> ~/.zshrc
        fi
        export PATH=$PATH:$go_bin_path
        print_success "Added Go bin directory to PATH"
    fi

    # Create development scripts
    mkdir -p scripts

    # Create development start script
    cat > scripts/dev.sh << 'EOF'
#!/bin/bash
# Development script for DevGen CLI

echo "🚀 Starting DevGen CLI development environment..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Start with hot reload
if command -v air >/dev/null 2>&1; then
    echo "📡 Starting with hot reload (Air)..."
    air
else
    echo "🔨 Air not found, building and running..."
    go run .
fi
EOF

    chmod +x scripts/dev.sh

    print_success "Development environment configured"
}

# Create example files
create_examples() {
    print_step "Creating example files..."

    # Create example playbook if it doesn't exist
    if [ ! -f "example-playbook.yaml" ]; then
        print_info "example-playbook.yaml already exists"
    else
        print_success "Example playbook available"
    fi

    # Create sample config
    mkdir -p .devgen
    if [ ! -f ".devgen/config.yaml" ]; then
        cat > .devgen/config.yaml << 'EOF'
version: "1.0.0"
devgen:
  default_output_dir: "./output"
  default_template: "fastapi-basic"
  auto_save: true
  check_updates: true

ui:
  theme:
    name: "cyber"
    dark_mode: true

logging:
  level: "info"
  format: "text"
  colors: true

servers:
  default:
    host: "localhost"
    port: 8080
    auto_restart: true
EOF
        print_success "Sample configuration created"
    else
        print_info "Configuration already exists"
    fi
}

# Print usage instructions
print_usage() {
    echo
    print_info "Setup completed! Here's how to get started:"
    echo
    echo -e "${CYAN}Basic Usage:${NC}"
    echo "  make dev              # Start development server with hot reload"
    echo "  make build            # Build the application"
    echo "  make test             # Run tests"
    echo "  make run              # Build and run"
    echo
    echo -e "${CYAN}Development:${NC}"
    echo "  ./scripts/dev.sh      # Start development environment"
    echo "  make lint             # Run linter"
    echo "  make format           # Format code"
    echo
    echo -e "${CYAN}CLI Usage:${NC}"
    echo "  ./build/devgen --help                    # Show help"
    echo "  ./build/devgen playbook run example-playbook.yaml  # Run example"
    echo "  ./build/devgen config init              # Initialize config"
    echo "  ./build/devgen server start             # Start development server"
    echo
    echo -e "${CYAN}Testing:${NC}"
    echo "  make test-coverage    # Run tests with coverage"
    echo "  make test-race        # Run tests with race detection"
    echo "  make test-bench       # Run benchmarks"
    echo
    echo -e "${CYAN}Docker (optional):${NC}"
    echo "  make docker-build     # Build Docker image"
    echo "  make docker-run       # Run in container"
    echo
    echo -e "${GREEN}🎉 Happy coding with DevGen CLI!${NC}"
    echo
}

# Main execution
main() {
    print_banner

    # Check if we're in the right directory
    if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
        print_error "This script must be run from the DevGen CLI project root"
        print_info "Please navigate to the directory containing go.mod and main.go"
        exit 1
    fi

    print_info "Setting up DevGen CLI development environment..."
    echo

    # Run setup steps
    check_tools
    echo

    init_project
    echo

    install_dependencies
    echo

    install_dev_tools
    echo

    setup_config
    echo

    setup_dev_env
    echo

    create_examples
    echo

    build_project
    echo

    run_tests
    echo

    print_usage
}

# Handle script arguments
case "${1:-setup}" in
    "setup"|"")
        main
        ;;
    "tools")
        check_tools
        ;;
    "deps")
        install_dependencies
        ;;
    "dev-tools")
        install_dev_tools
        ;;
    "build")
        build_project
        ;;
    "test")
        run_tests
        ;;
    "clean")
        print_step "Cleaning build artifacts..."
        rm -rf build dist tmp cache logs reports profiles
        print_success "Clean completed"
        ;;
    "help"|"-h"|"--help")
        echo "DevGen CLI Setup Script"
        echo
        echo "Usage: $0 [command]"
        echo
        echo "Commands:"
        echo "  setup      Complete setup (default)"
        echo "  tools      Check and install required tools"
        echo "  deps       Install Go dependencies"
        echo "  dev-tools  Install development tools"
        echo "  build      Build the project"
        echo "  test       Run tests"
        echo "  clean      Clean build artifacts"
        echo "  help       Show this help"
        ;;
    *)
        print_error "Unknown command: $1"
        print_info "Run '$0 help' for usage information"
        exit 1
        ;;
esac

# DevGen CLI - Makefile for development and building
# This Makefile provides convenient commands for development, testing, and building

.PHONY: help build test clean install dev deps lint format check-format vet security update-deps cross-compile docker release

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME := devgen
MAIN_FILE := main.go
BUILD_DIR := build
DIST_DIR := dist
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# Build flags
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"
LDFLAGS_DEV := -ldflags "-X main.version=$(VERSION)-dev -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Platform targets for cross-compilation
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

help: ## Show this help message
	@echo "DevGen CLI - Development Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nDevelopment:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n%s:\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

dev: deps ## Start development mode with hot reload
	@echo "🚀 Starting development server..."
	@which air > /dev/null || (echo "Installing air for hot reload..." && go install github.com/cosmtrek/air@latest)
	@air -c .air.toml || go run $(LDFLAGS_DEV) $(MAIN_FILE)

run: build ## Build and run the application
	@echo "🏃 Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

demo: build ## Run demo with sample playbook
	@echo "🎭 Running demo with example playbook..."
	@./$(BUILD_DIR)/$(BINARY_NAME) playbook run example-playbook.yaml

##@ Dependencies

deps: ## Download and install dependencies
	@echo "📦 Installing dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

update-deps: ## Update all dependencies to latest versions
	@echo "🔄 Updating dependencies..."
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

check-deps: ## Check for outdated dependencies
	@echo "🔍 Checking for outdated dependencies..."
	@go list -u -m all | grep '\[.*\]'

##@ Building

build: deps ## Build the application for current platform
	@echo "🔨 Building $(BINARY_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-dev: deps ## Build development version with debug info
	@echo "🔨 Building development version..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS_DEV) -o $(BUILD_DIR)/$(BINARY_NAME)-dev $(MAIN_FILE)

build-static: deps ## Build static binary (no external dependencies)
	@echo "🔨 Building static binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-static $(MAIN_FILE)

cross-compile: deps ## Build for all platforms
	@echo "🌍 Cross-compiling for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then output_name=$$output_name.exe; fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$$output_name $(MAIN_FILE); \
	done
	@echo "✅ Cross-compilation complete. Binaries in $(DIST_DIR)/"

##@ Testing

test: ## Run all tests
	@echo "🧪 Running tests..."
	@$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

test-race: ## Run tests with race condition detection
	@echo "🧪 Running tests with race detection..."
	@$(GOTEST) -race -v ./...

test-bench: ## Run benchmark tests
	@echo "🏃‍♂️ Running benchmark tests..."
	@$(GOTEST) -bench=. -benchmem ./...

test-integration: build ## Run integration tests
	@echo "🔗 Running integration tests..."
	@./scripts/run-integration-tests.sh

##@ Code Quality

lint: ## Run linter
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

format: ## Format Go code
	@echo "💄 Formatting code..."
	@$(GOFMT) -s -w .
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@goimports -w .

check-format: ## Check if code is properly formatted
	@echo "🔍 Checking code formatting..."
	@test -z "$$($(GOFMT) -l .)" || (echo "❌ Code is not formatted. Run 'make format'" && exit 1)
	@echo "✅ Code is properly formatted"

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@$(GOVET) ./...

security: ## Run security checks
	@echo "🔒 Running security checks..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	@gosec ./...

check: lint vet security check-format ## Run all code quality checks

##@ Installation

install: build ## Install binary to system
	@echo "📥 Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"

install-dev: build-dev ## Install development version to system
	@echo "📥 Installing $(BINARY_NAME)-dev to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME)-dev /usr/local/bin/
	@echo "✅ $(BINARY_NAME)-dev installed successfully"

uninstall: ## Remove binary from system
	@echo "📤 Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)-dev
	@echo "✅ $(BINARY_NAME) uninstalled successfully"

##@ Docker

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t devgen:$(VERSION) -t devgen:latest .

docker-run: docker-build ## Build and run Docker container
	@echo "🐳 Running Docker container..."
	@docker run --rm -it -v $(PWD):/workspace devgen:latest

docker-shell: docker-build ## Open shell in Docker container
	@echo "🐳 Opening shell in Docker container..."
	@docker run --rm -it -v $(PWD):/workspace --entrypoint /bin/sh devgen:latest

##@ Release

release: clean cross-compile ## Create release packages
	@echo "📦 Creating release packages..."
	@mkdir -p $(DIST_DIR)/packages
	@for binary in $(DIST_DIR)/$(BINARY_NAME)-*; do \
		if [ -f "$$binary" ]; then \
			platform=$$(basename $$binary | sed 's/$(BINARY_NAME)-//'); \
			package_name=$(BINARY_NAME)-$(VERSION)-$$platform; \
			mkdir -p $(DIST_DIR)/packages/$$package_name; \
			cp $$binary $(DIST_DIR)/packages/$$package_name/; \
			cp README.md $(DIST_DIR)/packages/$$package_name/; \
			cp example-playbook.yaml $(DIST_DIR)/packages/$$package_name/; \
			tar -czf $(DIST_DIR)/packages/$$package_name.tar.gz -C $(DIST_DIR)/packages $$package_name; \
			rm -rf $(DIST_DIR)/packages/$$package_name; \
			echo "Created package: $$package_name.tar.gz"; \
		fi; \
	done
	@echo "✅ Release packages created in $(DIST_DIR)/packages/"

changelog: ## Generate changelog
	@echo "📝 Generating changelog..."
	@which git-chglog > /dev/null || (echo "Installing git-chglog..." && go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest)
	@git-chglog -o CHANGELOG.md

##@ Utilities

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build time: $(BUILD_TIME)"

info: ## Show build information
	@echo "🔍 Build Information:"
	@echo "  Binary name: $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit: $(COMMIT)"
	@echo "  Build time: $(BUILD_TIME)"
	@echo "  Go version: $$(go version)"
	@echo "  Build dir: $(BUILD_DIR)"
	@echo "  Dist dir: $(DIST_DIR)"

setup-dev: ## Setup development environment
	@echo "🛠️ Setting up development environment..."
	@$(GOMOD) download
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Creating .air.toml configuration..."
	@cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
args_bin = []
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ."
delay = 1000
exclude_dir = ["assets", "tmp", "vendor", "testdata", "build", "dist"]
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
	@echo "✅ Development environment setup complete"

watch: ## Watch files and run tests on changes
	@echo "👀 Watching for file changes..."
	@which fswatch > /dev/null || (echo "Please install fswatch: brew install fswatch (macOS) or apt-get install fswatch (Linux)" && exit 1)
	@fswatch -o . | xargs -n1 -I{} make test

##@ Documentation

docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	@echo "Starting godoc server on http://localhost:6060"
	@godoc -http=:6060

docs-md: ## Generate markdown documentation
	@echo "📝 Generating markdown documentation..."
	@which gomarkdoc > /dev/null || (echo "Installing gomarkdoc..." && go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest)
	@gomarkdoc ./... > API.md
	@echo "✅ API documentation generated: API.md"

##@ Benchmarks and Profiling

profile-cpu: build ## Run CPU profiling
	@echo "📊 Running CPU profiling..."
	@mkdir -p profiles
	@./$(BUILD_DIR)/$(BINARY_NAME) -cpuprofile=profiles/cpu.prof playbook run example-playbook.yaml
	@go tool pprof profiles/cpu.prof

profile-mem: build ## Run memory profiling
	@echo "📊 Running memory profiling..."
	@mkdir -p profiles
	@./$(BUILD_DIR)/$(BINARY_NAME) -memprofile=profiles/mem.prof playbook run example-playbook.yaml
	@go tool pprof profiles/mem.prof

stress-test: build ## Run stress tests
	@echo "💪 Running stress tests..."
	@for i in {1..10}; do \
		echo "Stress test iteration $$i"; \
		./$(BUILD_DIR)/$(BINARY_NAME) playbook validate example-playbook.yaml; \
	done

# Include custom targets if they exist
-include Makefile.local

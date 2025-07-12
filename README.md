# DevGen CLI - AI Development Platform

DevGen CLI is a comprehensive command-line interface for AI developers working with Model Context Protocol (MCP) servers, knowledge bases, and AI development workflows. Built with elegant terminal UI powered by Charm libraries.

## 🚀 Core Features

- **🎛️ Interactive Dashboard** - Beautiful terminal UI for real-time MCP server monitoring and management
- **🔌 MCP Server Management** - Toggle, monitor, and manage MCP servers with category-based organization
- **🌐 Registry Integration** - Centralized server discovery through HTTP-based MCP Registry
- **🔐 SSH Server** - Secure remote terminal access for public-facing deployments
- **📊 Health Monitoring** - Real-time server status and health across all registered servers
- **🎨 Modern Terminal UI** - Cyberpunk-inspired design with emoji indicators and clean layouts

## 📦 Installation

### Quick Install

```bash
# Clone the repository
git clone https://github.com/devq-ai/devgen-cli.git
cd devgen-cli

# Build and install
make build
make install-user

# Add to PATH (add to your shell config)
echo 'alias devgen="$HOME/.local/bin/devgen"' >> ~/.zshrc
source ~/.zshrc
```

### Prerequisites

- Go 1.21 or higher
- Access to MCP servers configuration file (`mcp_status.json`)

## 🎯 Quick Start

```bash
# Launch interactive dashboard
devgen dashboard
# or shorter alias
devgen d

# Show detailed help
devgen help

# Check version
devgen --version
```

## 🔧 Core Commands

DevGen CLI provides four main commands with intuitive aliases:

### 📊 Dashboard Command
```bash
devgen dashboard    # Launch interactive server dashboard
devgen dash         # Alias
devgen d            # Short alias
```

**Dashboard Features:**
- Real-time server status with emoji category indicators
- Interactive navigation (↑/↓ arrows or j/k keys)
- Toggle servers on/off with Enter/Space
- Single-column layout with text wrapping
- Category-based organization (🧠 knowledge, ⚡ development, 🌐 web, etc.)

**Dashboard Controls:**
- `↑/↓` or `j/k` - Navigate server list
- `Enter/Space` - Toggle selected server on/off
- `q` - Quit dashboard

### 🌐 Registry Command
```bash
devgen registry status    # Check MCP Registry status
devgen registry servers   # List all registered servers
devgen registry tools     # Show available tools
devgen registry start     # Start the registry server

# Aliases
devgen reg status
devgen r status
```

**Registry Features:**
- Centralized server discovery and management
- HTTP API for integration (default: http://127.0.0.1:31337)
- Real-time server registration and health monitoring
- Tool aggregation across all registered servers

### 🔐 SSH Command
```bash
devgen ssh                # Start SSH server on port 2222
devgen ssh --ssh-port 3000 --ssh-host 0.0.0.0

# Aliases
devgen server
devgen remote
```

**SSH Features:**
- Secure remote access to DevGen CLI
- Essential for public-facing web deployments
- Password authentication (demo/devq)
- Interactive terminal sessions
- Remote server management capabilities

**Connection:**
```bash
ssh -p 2222 demo@your-server.com
# Password: demo or devq
```

### 📖 Help Command
```bash
devgen help          # Show comprehensive help
devgen guide         # Alias
devgen docs          # Alias
```

## 🔌 Supported MCP Servers

DevGen manages 13+ MCP servers across multiple categories:

### 🧠 Knowledge Servers
- **context7-mcp** - Redis-backed contextual reasoning and document management
- **memory-mcp** - Memory management and persistence for AI workflows
- **sequential-thinking-mcp** - Step-by-step problem solving and reasoning chains

### ⚡ Development Servers
- **fastapi-mcp** - FastAPI project generation and management
- **pytest-mcp** - Python testing framework integration
- **pydantic-ai-mcp** - Pydantic AI agent management and orchestration

### 🌐 Web & Data Servers
- **crawl4ai-mcp** - Web crawling and content extraction
- **github-mcp** - GitHub repository operations and management
- **surrealdb-mcp** - Multi-model database operations

### 🔧 Framework Servers
- **fastmcp-mcp** - FastMCP framework status and management
- **registry-mcp** - MCP server discovery and registry management

### 💾 Database Servers
- **postgres-mcp** - PostgreSQL database operations
- **sqlite-mcp** - SQLite database management

### 🏗️ Infrastructure Servers
- **logfire-mcp** - Observability and logging platform integration

## 🎨 Design & UI

DevGen features a modern cyberpunk-inspired terminal interface:

**Color Palette:**
- Primary: Neon Pink (`#FF10F0`)
- Success: Neon Green (`#39FF14`) 
- Error: Neon Red (`#FF3131`)
- Info: Neon Cyan (`#00FFFF`)
- Text: Light Gray (`#E3E3E3`)

**UI Elements:**
- Category emoji indicators for visual organization
- Clean single-column layout with proper text wrapping
- Responsive design that adapts to terminal size
- Consistent styling across all commands

## ⚙️ Configuration

DevGen automatically searches for `mcp_status.json` configuration in:

1. Current directory (`./mcp_status.json`)
2. Parent directory (`../mcp_status.json`)
3. DevQAI machina directory (`/Users/dionedge/devqai/machina/mcp_status.json`)

**Custom Configuration:**
```bash
devgen --config /path/to/custom.json dashboard
```

**Global Flags:**
- `--config, -c FILE` - Configuration file path
- `--verbose, -v` - Enable verbose logging
- `--log-level LEVEL` - Set log level (debug, info, warn, error)
- `--ssh` - Start SSH server mode
- `--ssh-port PORT` - SSH server port (default: 2222)
- `--ssh-host HOST` - SSH server host (default: localhost)
- `--registry-url URL` - MCP registry URL (default: http://127.0.0.1:31337)
- `--use-registry` - Use MCP registry for server management

## 🚀 Planned Features

DevGen CLI includes comprehensive technical specifications for upcoming features:

### 🧠 Knowledge Base Management (`devgen kb`)
- Database statistics and analytics
- Knowledge base health monitoring
- Data import/export capabilities
- Vector search integration

### 🔍 RAG-Powered Search (`devgen search`)
- Semantic similarity search across knowledge bases
- Code pattern matching and discovery
- Multi-source search aggregation
- Knowledge graph exploration

### 🛡️ DeHallucinator (`devgen dehall`)
- AI hallucination detection and prevention
- Fact verification against knowledge bases
- Code accuracy validation
- Real-time verification during AI interactions

## 🛠️ Development

### Building from Source

```bash
# Install dependencies
go mod download

# Build for current platform
make build

# Install locally
make install-user

# Cross-compile for all platforms
make cross-compile

# Run tests
make test

# Development build with debug info
make build-dev
```

### Project Structure

```
devgen-cli/
├── src/                 # Go source files
│   ├── main.go         # Main CLI and command definitions
│   ├── dashboard.go    # Interactive dashboard implementation
│   └── registry.go     # Registry integration
├── docs/               # Documentation
│   └── TECHNICAL_SPECIFICATION.md
├── build/              # Build artifacts
├── Makefile           # Build and development tasks
└── README.md          # This file
```

### Make Targets

```bash
make build          # Build the CLI binary
make install-user   # Install to ~/.local/bin
make test          # Run tests
make clean         # Clean build artifacts
make help          # Show available targets
```

## 📊 Example Output

### Dashboard View
```
🚀 DevGen MCP Server Dashboard

🧠 context7-mcp                    active
   Redis-backed contextual reasoning and document management
   
⚡ fastapi-mcp                     active  
   FastAPI project generation and management
   
🌐 crawl4ai-mcp                    inactive
   Web crawling and content extraction

Press ↑/↓ to navigate, Enter to toggle, q to quit
```

### Registry Status
```
🌐 MCP Registry Status

✓ Registry is running at http://127.0.0.1:31337
✓ 13 servers registered
✓ 81+ tools available
✓ Last updated: 2025-07-12T10:30:15Z
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is part of the DevQ.ai ecosystem. See [LICENSE](LICENSE) for details.

## 🔗 Related Projects

- **[DevQ.ai Platform](https://devq.ai)** - AI-powered development tools
- **[MCP Registry](https://github.com/devq-ai/mcp-registry)** - Centralized MCP server discovery
- **[FastMCP Framework](https://github.com/devq-ai/fastmcp)** - Framework for building MCP servers

---

**Built with ❤️ by the DevQ.ai team using [Charm](https://charm.sh) libraries.**

*For detailed technical specifications and architecture, see [docs/TECHNICAL_SPECIFICATION.md](docs/TECHNICAL_SPECIFICATION.md)*
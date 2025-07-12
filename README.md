# DevGen CLI - MCP Server Management Tool

DevGen is a powerful command-line interface for managing MCP (Model Context Protocol) servers with beautiful terminal UI powered by Charm libraries.

## 🚀 Features

- **Interactive Dashboard** - Beautiful terminal UI for real-time server monitoring
- **Server Management** - Start, stop, and toggle MCP servers with simple commands
- **Health Monitoring** - Check server status and health across all registered servers
- **Registry Integration** - Works seamlessly with the Machina MCP Registry
- **Terminal UI** - Built with Charm libraries for an elegant terminal experience

## 📦 Installation

### From Source

```bash
# Clone and build
git clone https://github.com/devq-ai/devqai.git
cd devqai/devgen-cli
make build

# Install globally
make install
```

### Prerequisites

- Go 1.23 or higher
- Access to the Machina MCP Registry (`mcp_status.json`)

## 🔧 Usage

### Basic Commands

```bash
# List all MCP servers
devgen server list

# Show status of a specific server
devgen server status context7-mcp

# Start a server
devgen server start context7-mcp

# Stop a server
devgen server stop context7-mcp

# Toggle server on/off
devgen server toggle context7-mcp

# Check health of all servers
devgen server health
```

### Interactive Dashboard

Launch the interactive terminal dashboard:

```bash
devgen dashboard
```

**Dashboard Controls:**
- `↑/↓` - Navigate server list
- `Enter` - Toggle selected server on/off
- `r` - Refresh server status
- `q` - Quit dashboard

### Configuration

```bash
# Show current configuration
devgen config show

# Verify configuration file
devgen config init
```

## 🔌 MCP Server Registry

DevGen works with the Machina MCP Registry which includes:

### Knowledge Servers
- **context7-mcp** - Redis-backed contextual reasoning and document management
- **memory-mcp** - Memory management and persistence for AI workflows
- **sequential-thinking-mcp** - Step-by-step problem solving and reasoning chains

### Development Servers
- **fastapi-mcp** - FastAPI project generation and management
- **pytest-mcp** - Python testing framework integration
- **pydantic-ai-mcp** - Pydantic AI agent management and orchestration

### Web & Data Servers
- **crawl4ai-mcp** - Web crawling and content extraction
- **github-mcp** - GitHub repository operations and management
- **surrealdb-mcp** - Multi-model database operations

### Framework Servers
- **fastmcp-mcp** - FastMCP framework status and management
- **registry-mcp** - MCP server discovery and registry management

## 🎨 Terminal UI Design

DevGen uses a cyberpunk-inspired color scheme:

- **Primary**: Neon Pink (`#FF10F0`)
- **Success**: Neon Green (`#39FF14`)
- **Error**: Neon Red (`#FF3131`)
- **Warning**: Neon Yellow (`#E9FF32`)
- **Info**: Neon Cyan (`#00FFFF`)

## 📁 Configuration Files

DevGen automatically searches for `mcp_status.json` in these locations:

1. Current directory (`./mcp_status.json`)
2. Parent directory (`../mcp_status.json`)
3. Default Machina location (`/Users/dionedge/devqai/machina/mcp_status.json`)

You can specify a custom config file:

```bash
devgen --config /path/to/mcp_status.json server list
```

## 🔍 Command Reference

### Global Flags

- `--config, -c` - Config file path (default: `mcp_status.json`)
- `--verbose, -v` - Enable verbose logging
- `--log-level` - Set log level (debug, info, warn, error)
- `--interactive, -i` - Enable interactive mode
- `--help, -h` - Show help
- `--version` - Show version

### Server Commands

```bash
devgen server [command]

Available Commands:
  list      List all MCP servers
  status    Show server status [server-name]
  start     Start a server [server-name]
  stop      Stop a server [server-name]
  toggle    Toggle server on/off [server-name]
  health    Check health of all servers
```

### Dashboard Command

```bash
devgen dashboard

Launch interactive terminal dashboard for real-time server management
```

### Config Commands

```bash
devgen config [command]

Available Commands:
  show      Show current configuration
  init      Verify configuration file
```

## 🛠️ Development

### Building

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Build for development (with debug info)
make build-dev

# Cross-compile for all platforms
make cross-compile
```

### Testing

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Integration tests
make test-integration
```

### Code Quality

```bash
# Format code
make format

# Run linter
make lint

# Run security checks
make security

# Run all quality checks
make check
```

## 📊 Example Output

### Server List
```
🔌 MCP Server Registry

• context7-mcp [active]
  Redis-backed contextual reasoning and document management with vector search
  Endpoint: stdio://devqai/mcp/mcp-servers/context7-mcp/context7_mcp/server.py
  Tools: 15 available
  Category: knowledge
  Framework: FastMCP
  Last Health Check: 2025-07-10T13:25:07-05:00

• memory-mcp [production-ready]
  Memory management and persistence for AI workflows with search capabilities
  Endpoint: stdio://devqai/mcp/mcp-servers/memory-mcp/memory_mcp/server.py
  Tools: 8 available
  Category: knowledge
  Framework: FastMCP
  Last Health Check: Never
```

### Server Status
```
📊 Status: context7-mcp
Status: active
Version: 1.0.0
Description: Redis-backed contextual reasoning and document management with vector search
Endpoint: stdio://devqai/mcp/mcp-servers/context7-mcp/context7_mcp/server.py
Category: knowledge
Framework: FastMCP
Tools: 15 available
  Available tools: store_document, search_documents, get_context, clear_context, list_documents
  ... and 10 more tools
Last Health Check: 2025-07-10T13:25:07-05:00
Health Check Failures: 0
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is part of the DevQ.ai ecosystem and follows the same licensing terms.

## 🔗 Related Projects

- [Machina MCP Registry](../machina/README.md) - The main MCP server registry platform
- [FastMCP Framework](../machina/fastmcp/) - Framework for building MCP servers
- [DevQ.ai](https://devq.ai) - AI-powered development tools

---

Built with ❤️ by the DevQ.ai team using [Charm](https://charm.sh) libraries.\n\n---\n\n## 🌐 NEW: MCP Registry Integration\n\nThe CLI now supports integration with the HTTP-based MCP Registry system for centralized server management:\n\n### Registry Commands\n\n```bash\n# Check registry status\ndevgen registry status\n\n# List registered servers\ndevgen registry servers\n\n# List all available tools\ndevgen registry tools\n\n# Start the registry\ndevgen registry start\n```\n\n### Registry Features\n\n- **Centralized Discovery**: All 13 MCP servers accessible through a single HTTP endpoint\n- **Tool Aggregation**: Unified access to 81+ tools across all servers\n- **Health Monitoring**: Real-time status for all registered servers\n- **Auto-Registration**: Servers automatically register with the registry\n- **HTTP API**: RESTful endpoints for integration with other tools\n\n### Registry URLs\n\n- **Main Registry**: http://127.0.0.1:31337\n- **Servers Endpoint**: http://127.0.0.1:31337/servers\n- **Tools Endpoint**: http://127.0.0.1:31337/tools\n\n### Configuration\n\n```bash\n# Use registry for server management\ndevgen --use-registry dashboard\n\n# Custom registry URL\ndevgen --registry-url http://localhost:8080 registry status\n```

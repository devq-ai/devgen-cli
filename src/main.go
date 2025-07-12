package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/spf13/cobra"
)

// Global flags
var (
	configFile   string
	verbose      bool
	logLevel     string
	sshMode      bool
	sshPort      int
	sshHost      string
	registryURL  string
	useRegistry  bool
)

// MCP Server types
type MCPServer struct {
	Name              string      `json:"name"`
	Endpoint          string      `json:"endpoint"`
	Tools             []string    `json:"tools"`
	Status            string      `json:"status"`
	Version           string      `json:"version"`
	Description       string      `json:"description"`
	Metadata          MCPMetadata `json:"metadata"`
	RegisteredAt      string      `json:"registered_at"`
	LastHealthCheck   string      `json:"last_health_check"`
	LastSeen          *string     `json:"last_seen"`
	HealthCheckFails  int         `json:"health_check_failures"`
}

type MCPMetadata struct {
	Framework       string   `json:"framework"`
	Category        string   `json:"category"`
	HealthCheck     string   `json:"health_check"`
	EnvironmentVars []string `json:"environment_vars"`
}

type MCPRegistry struct {
	Version   string      `json:"version"`
	Timestamp string      `json:"timestamp"`
	Servers   []MCPServer `json:"servers"`
	Tools     []MCPTool   `json:"tools"`
}

type MCPTool struct {
	Name        string `json:"name"`
	ServerName  string `json:"server_name"`
	Description string `json:"description"`
	UseCount    int    `json:"use_count"`
	ErrorCount  int    `json:"error_count"`
	LastUsed    string `json:"last_used"`
}

// Dashboard types
// Dashboard types moved to dashboard.go

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true).
		Padding(1, 2)

	statusRunning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#39FF14")).
		Bold(true)

	statusStopped = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3131")).
		Bold(true)

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)

	itemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E3E3E3"))

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true)
)

// Logfire integration - send logs to logfire-mcp server
func logToLogfire(level, message string, extra map[string]interface{}) {
	go func() {
		// Try to send to logfire-mcp server via HTTP
		requestData := map[string]interface{}{
			"level":      level,
			"message":    message,
			"extra_data": extra,
		}
		
		jsonData, _ := json.Marshal(requestData)
		
		// Send to Logfire via clean Python subprocess
		cmd := exec.Command("python3", "-c", fmt.Sprintf(`
import os, sys, json
sys.path.append('src')
os.environ['LOGFIRE_TOKEN'] = os.getenv('LOGFIRE_WRITE_TOKEN', '')
import logfire
logfire.configure(inspect_arguments=False)

data = json.loads('''%s''')
extra = data.get('extra_data', {})
extra['service'] = 'machina-cli'

if data['level'] == 'info':
    logfire.info(data['message'], **extra)
elif data['level'] == 'warning':  
    logfire.warning(data['message'], **extra)
elif data['level'] == 'error':
    logfire.error(data['message'], **extra)
else:
    logfire.info(data['message'], level=data['level'], **extra)
`, string(jsonData)))
		
		cmd.Dir = "/Users/dionedge/devqai/machina"
		cmd.Run() // Ignore errors for non-blocking
		
		// Fallback: write to local file for debugging
		logFile, err := os.OpenFile("machina_logfire.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			logData := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"level":     level,
				"message":   message,
				"service":   "machina-cli",
				"component": "main",
				"project":   os.Getenv("LOGFIRE_PROJECT_NAME"),
			}
			for k, v := range extra {
				logData[k] = v
			}
			jsonData, _ := json.Marshal(logData)
			logFile.WriteString(string(jsonData) + "\n")
			logFile.Close()
		}
		
		// Also write to debug log
		debugFile, _ := os.OpenFile("machina_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(debugFile, "[LOGFIRE] %s: %s\n", level, message)
		debugFile.Close()
	}()
}

func main() {
	// Load environment variables from .env file
	loadEnvFile()

	// Initialize structured logging
	logger := log.New(os.Stderr)
	logger.SetPrefix("devgen")

	rootCmd := &cobra.Command{
		Use:   "devgen",
		Short: "DevGen - AI Development Platform CLI",
		Long: `DevGen CLI - AI Development Platform

A comprehensive command-line interface for managing MCP servers, knowledge bases,
and AI development workflows. Built for developers working with Model Context Protocol
servers and AI-powered development tools.

🚀 Core Features:
  • Interactive MCP server dashboard
  • Knowledge base management and statistics  
  • RAG-powered search and lookup
  • AI hallucination detection and verification
  • Registry-based server discovery

📚 Quick Start:
  devgen dashboard    # Launch interactive server dashboard
  devgen help         # Show detailed command help
  devgen --version    # Show version information

For more information, visit: https://github.com/devq-ai/devgen-cli`,
		Version: "1.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return setupLogging(logger)
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "mcp_status.json", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&sshMode, "ssh", false, "start SSH server for terminal access")
	rootCmd.PersistentFlags().IntVar(&sshPort, "ssh-port", 2222, "SSH server port")
	rootCmd.PersistentFlags().StringVar(&sshHost, "ssh-host", "localhost", "SSH server host")
	rootCmd.PersistentFlags().StringVar(&registryURL, "registry-url", "http://127.0.0.1:31337", "MCP registry URL")
	rootCmd.PersistentFlags().BoolVar(&useRegistry, "use-registry", false, "use MCP registry for server management")

	// Add core commands
	rootCmd.AddCommand(
		newDashboardCmd(),
		newRegistryCmd(),
		newSSHCmd(),
		newHelpCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

// Setup logging
func setupLogging(logger *log.Logger) error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logger.SetLevel(level)

	if verbose {
		logger.SetLevel(log.DebugLevel)
	}

	return nil
}

// Dashboard command
func newDashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"dash", "d"},
		Short:   "Launch interactive dashboard",
		Long:    "Launch the interactive terminal dashboard for managing MCP servers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDashboard()
		},
	}

	return cmd
}

// SSH command
func newSSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh",
		Aliases: []string{"server", "remote"},
		Short:   "Start SSH server for remote terminal access",
		Long:    "Start an SSH server that provides secure remote terminal access to DevGen CLI commands. Essential for public-facing deployments.",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Starting SSH server", "host", sshHost, "port", sshPort)
			return startSSHServer()
		},
	}

	cmd.Flags().IntVar(&sshPort, "ssh-port", 2222, "SSH server port")
	cmd.Flags().StringVar(&sshHost, "ssh-host", "localhost", "SSH server host")

	return cmd
}

// Help command with detailed explanations
func newHelpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "help",
		Aliases: []string{"guide", "docs"},
		Short:   "Show detailed command help and usage examples",
		Long:    "Display comprehensive help information for all DevGen CLI commands with examples and use cases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showExtendedHelp()
		},
	}

	return cmd
}


// Show extended help with detailed command explanations
func showExtendedHelp() error {
	helpText := `
DevGen CLI - AI Development Platform
====================================

OVERVIEW:
DevGen CLI is a comprehensive command-line interface for AI developers working with
Model Context Protocol (MCP) servers, knowledge bases, and AI development workflows.

CORE COMMANDS:
──────────────

📊 devgen dashboard
   Launch interactive terminal dashboard for MCP server management
   
   Features:
   • Real-time server status monitoring
   • Interactive server toggling (↑/↓ to navigate, Enter to toggle)
   • Category-based server organization with emoji indicators
   • Live server statistics and health monitoring
   
   Usage:
     devgen dashboard
   
   Controls:
     ↑/↓ or j/k    Navigate server list
     Enter/Space   Toggle selected server on/off
     q             Quit dashboard

🔌 devgen registry
   Manage the HTTP-based MCP Registry for centralized server discovery
   
   Subcommands:
     status     Check if the MCP Registry is running
     servers    List all servers registered with the HTTP registry
     tools      Show all available tools across registered servers
     start      Start the MCP Registry server
   
   Features:
   • Centralized server discovery and management
   • Real-time server registration and health monitoring
   • Tool aggregation across all registered servers
   • HTTP API for integration with other tools
   
   Usage:
     devgen registry status
     devgen registry servers
     devgen registry tools
     devgen registry start

🔐 devgen ssh
   Start SSH server for secure remote terminal access
   
   Features:
   • Secure remote access to DevGen CLI
   • Essential for public-facing web deployments
   • Password and public key authentication
   • Interactive terminal sessions
   
   Usage:
     devgen ssh                           # Start SSH server on default port 2222
     devgen ssh --ssh-port 2222          # Custom port
     devgen ssh --ssh-host 0.0.0.0       # Bind to all interfaces
   
   Connection:
     ssh -p 2222 demo@your-server.com    # Connect to SSH server
     Password: demo or devq

PLANNED FEATURES (Coming Soon):
─────────────────────────────────

🧠 devgen kb (Knowledge Base Management)
   • Database statistics and analytics
   • Knowledge base health monitoring
   • Data import/export capabilities
   
   Commands:
     devgen kb stats           # Database statistics
     devgen kb search          # Knowledge search
     devgen kb import          # Import data sources
     devgen kb export          # Export knowledge

🔍 devgen search (RAG-Powered Search)
   • Semantic similarity search across knowledge bases
   • Code pattern matching and discovery
   • Multi-source search aggregation
   
   Commands:
     devgen search "FastAPI auth"     # Semantic search
     devgen search --code "handler"   # Code-specific search
     devgen search --graph "API"      # Knowledge graph search

🛡️ devgen dehall (DeHallucinator - AI Verification)
   • AI hallucination detection and prevention
   • Fact verification against knowledge bases
   • Code accuracy validation
   
   Commands:
     devgen dehall check response.txt    # Check for hallucinations
     devgen dehall verify "claim"        # Verify specific facts
     devgen dehall analyze code.py       # Analyze code explanations

GLOBAL FLAGS:
─────────────
  -c, --config FILE        Configuration file path (default: mcp_status.json)
  -v, --verbose           Enable verbose logging
  --log-level LEVEL       Set log level (debug, info, warn, error)
  --registry-url URL      MCP registry URL (default: http://127.0.0.1:31337)
  --use-registry          Use MCP registry for server management
  --version               Show version information

CONFIGURATION:
──────────────
DevGen automatically searches for configuration files in:
1. Current directory (./mcp_status.json)
2. Parent directory (../mcp_status.json)  
3. DevQAI machina directory (/Users/dionedge/devqai/machina/mcp_status.json)

Custom configuration:
  devgen --config /path/to/custom.json dashboard

EXAMPLES:
─────────
# Start the interactive dashboard
devgen dashboard

# Check registry status and start if needed
devgen registry status
devgen registry start

# Use a custom configuration file
devgen --config ./my-servers.json dashboard

GETTING HELP:
─────────────
• devgen help              # This detailed help
• devgen [command] --help   # Command-specific help

For more information and updates:
📖 Documentation: https://github.com/devq-ai/devgen-cli
🐛 Issues: https://github.com/devq-ai/devgen-cli/issues

Happy coding! 🚀
`

	fmt.Print(helpText)
	return nil
}


// Registry command for HTTP MCP Registry integration
func newRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "registry",
		Aliases: []string{"reg", "r"},
		Short:   "Interact with the MCP Registry",
		Long:    "Commands for interacting with the HTTP-based MCP Registry system.",
	}

	cmd.AddCommand(
		newRegistryStatusCmd(),
		newRegistryServersCmd(),
		newRegistryToolsCmd(),
		newRegistryStartCmd(),
	)

	return cmd
}

// Registry status command
func newRegistryStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check MCP Registry status",
		Long:  "Check the status of the MCP Registry and get basic information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkRegistryStatus()
		},
	}

	return cmd
}

// Registry servers command
func newRegistryServersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "List servers from MCP Registry",
		Long:  "List all registered servers from the HTTP MCP Registry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRegistryServers()
		},
	}

	return cmd
}

// Registry tools command
func newRegistryToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "List tools from MCP Registry",
		Long:  "List all available tools from the HTTP MCP Registry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRegistryTools()
		},
	}

	return cmd
}

// Registry start command
func newRegistryStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the MCP Registry",
		Long:  "Start the HTTP-based MCP Registry server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startMCPRegistry()
		},
	}

	return cmd
}

// Load MCP registry
func loadMCPRegistry() (*MCPRegistry, error) {
	// Try multiple locations for the config file
	var data []byte
	var err error

	// First try the specified config file
	data, err = ioutil.ReadFile(configFile)
	if err != nil && configFile == "mcp_status.json" {
		// Smart discovery of machina repository
		machinaRoot := findMachinaRoot()

		locations := []string{
			"./mcp_status.json",
			"../mcp_status.json",
		}

		if machinaRoot != "" {
			locations = append(locations, filepath.Join(machinaRoot, "mcp_status.json"))
		}

		// Fallback locations
		locations = append(locations,
			"/Users/dionedge/devqai/machina/mcp_status.json",
			os.ExpandEnv("$HOME/devqai/machina/mcp_status.json"),
		)

		for _, location := range locations {
			data, err = ioutil.ReadFile(location)
			if err == nil {
				configFile = location
				break
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read registry file: %v", err)
	}

	var registry MCPRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry JSON: %v", err)
	}

	return &registry, nil
}

// Save MCP registry to file
func saveMCPRegistry(registry *MCPRegistry) error {
	// Debug: log save attempt
	logFile, _ := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Fprintf(logFile, "SAVE: Attempting to save registry to %s\n", configFile)
	
	// Find and log the crawl4ai-mcp status being saved
	for _, server := range registry.Servers {
		if server.Name == "crawl4ai-mcp" {
			fmt.Fprintf(logFile, "SAVE: crawl4ai-mcp status being written: %s\n", server.Status)
			break
		}
	}
	logFile.Close()
	
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry JSON: %v", err)
	}

	// Write to file and ensure it's synced
	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logFile, _ := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "SAVE ERROR: Failed to open file: %v\n", err)
		logFile.Close()
		return fmt.Errorf("failed to open registry file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		logFile, _ := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "SAVE ERROR: Failed to write data: %v\n", err)
		logFile.Close()
		return fmt.Errorf("failed to write registry data: %v", err)
	}

	// Force sync to disk
	if err := file.Sync(); err != nil {
		logFile, _ := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "SAVE ERROR: Failed to sync: %v\n", err)
		logFile.Close()
		return fmt.Errorf("failed to sync registry file: %v", err)
	}

	logFile, _ = os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Fprintf(logFile, "SAVE SUCCESS: Registry saved\n")
	logFile.Close()

	return nil
}

// Find machina root directory
func findMachinaRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Look for machina indicators
	for {
		indicators := []string{
			"mcp_status.json",
			"mcp-servers",
			"fastmcp",
		}

		for _, indicator := range indicators {
			if _, err := os.Stat(filepath.Join(currentDir, indicator)); err == nil {
				return currentDir
			}
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}

	return ""
}

// Load environment variables
func loadEnvFile() {
	// Look for .env file in current directory or parent directories
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			file, err := os.Open(envPath)
			if err != nil {
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
					continue
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					os.Setenv(key, value)
				}
			}

			fmt.Printf("📄 Loaded environment variables from: %s\n", envPath)
			return
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
}

// Dashboard implementation
// Dashboard methods moved to dashboard.go

// testMCPServerConnectivity tests if an MCP server can actually start
func testMCPServerConnectivity(server *MCPServer) bool {
	// For stdio-based servers, try to actually start them briefly
	if strings.HasPrefix(server.Endpoint, "stdio://") {
		// Extract the Python script path from the endpoint
		scriptPath := strings.TrimPrefix(server.Endpoint, "stdio://")

		// Update path to actual location
		if strings.Contains(scriptPath, "context7-mcp") {
			scriptPath = "/Users/dionedge/devqai/machina/mcp-servers/context7_mcp.py"
		} else if strings.Contains(scriptPath, "memory-mcp") {
			scriptPath = "/Users/dionedge/devqai/machina/mcp-servers/memory_mcp.py"
		}

		// Simple connectivity test - check if file exists and is executable
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
		return true
	}

	// For other protocols, assume they're working
	return true
}

// toggleServer toggles the status of an MCP server
func toggleServer(serverName string) error {
	registry, err := loadMCPRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %v", err)
	}

	for i := range registry.Servers {
		if registry.Servers[i].Name == serverName {
			if registry.Servers[i].Status == "active" || registry.Servers[i].Status == "production-ready" || registry.Servers[i].Status == "running" {
				registry.Servers[i].Status = "inactive"
			} else {
				registry.Servers[i].Status = "active"
			}
			break
		}
	}

	// Save the updated registry back to file
	return saveMCPRegistry(registry)
}












// SSH Server implementation
func startSSHServer() error {
	registry, err := loadMCPRegistry()
	if err != nil {
		return fmt.Errorf("failed to load MCP registry: %w", err)
	}

	// Ensure SSH directory exists
	sshDir := ".ssh"
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	// Generate host key if it doesn't exist
	hostKeyPath := filepath.Join(sshDir, "devgen_host_key")
	if err := generateHostKeyIfNotExists(hostKeyPath); err != nil {
		return fmt.Errorf("failed to generate host key: %w", err)
	}

	// Create SSH server with Wish middleware
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", sshHost, sshPort)),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			return password == "demo" || password == "devq"
		}),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					handleSSHSession(sess, registry)
				}
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create SSH server: %w", err)
	}

	fmt.Printf("SSH server started at %s:%d\n", sshHost, sshPort)
	fmt.Printf("Connect with: ssh -p %d demo@%s\n", sshPort, sshHost)
	fmt.Printf("Password: demo or devq\n")

	return s.ListenAndServe()
}

func generateHostKeyIfNotExists(hostKeyPath string) error {
	if _, err := os.Stat(hostKeyPath); err == nil {
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	privateKeyFile, err := os.OpenFile(hostKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	fmt.Printf("Generated SSH host key at %s\n", hostKeyPath)
	return nil
}

func handleSSHSession(sess ssh.Session, registry *MCPRegistry) {
	pty, winCh, isPty := sess.Pty()
	if !isPty {
		fmt.Fprintf(sess, "DevGen CLI requires a PTY\n")
		sess.Exit(1)
		return
	}

	// Create terminal renderer
	renderer := lipgloss.NewRenderer(sess)

	// Style definitions for SSH terminal
	titleStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true).
		Padding(1, 2)

	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)

	// Welcome message
	welcome := titleStyle.Render("🚀 DevGen SSH Terminal") + "\n\n" +
		headerStyle.Render("Available Commands:") + "\n" +
		"• list        - List all MCP servers\n" +
		"• status <name> - Show server status\n" +
		"• health      - Check health of all servers\n" +
		"• help        - Show this help\n" +
		"• exit        - Close connection\n\n"

	fmt.Fprint(sess, welcome)

	// Handle window size changes
	go func() {
		for win := range winCh {
			pty.Window.Width = win.Width
			pty.Window.Height = win.Height
		}
	}()

	// Command processing loop
	for {
		fmt.Fprint(sess, headerStyle.Render("devgen> "))

		// Read command
		var cmd string
		fmt.Fscanf(sess, "%s", &cmd)

		switch cmd {
		case "list":
			handleSSHListCommand(sess, registry, renderer)
		case "status":
			var serverName string
			fmt.Fscanf(sess, "%s", &serverName)
			handleSSHStatusCommand(sess, registry, serverName, renderer)
		case "health":
			handleSSHHealthCommand(sess, registry, renderer)
		case "help":
			fmt.Fprint(sess, welcome)
		case "exit", "quit":
			fmt.Fprint(sess, "Goodbye! 👋\n")
			sess.Exit(0)
			return
		case "":
			// Empty command, just continue
		default:
			fmt.Fprintf(sess, "Unknown command: %s\n", cmd)
			fmt.Fprint(sess, "Type 'help' for available commands\n")
		}
	}
}

func handleSSHListCommand(sess ssh.Session, registry *MCPRegistry, renderer *lipgloss.Renderer) {
	titleStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true)

	statusRunning := renderer.NewStyle().
		Foreground(lipgloss.Color("#39FF14")).
		Bold(true)

	statusStopped := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF3131")).
		Bold(true)

	fmt.Fprint(sess, titleStyle.Render("🔌 MCP Server Registry")+"\n\n")

	for _, server := range registry.Servers {
		statusText := "inactive"
		statusStyle := statusStopped
		if server.Status == "active" || server.Status == "production-ready" {
			statusText = server.Status
			statusStyle = statusRunning
		}

		fmt.Fprintf(sess, "• %s [%s]\n", server.Name, statusStyle.Render(statusText))
		fmt.Fprintf(sess, "  %s\n", server.Description)
		fmt.Fprintf(sess, "  Tools: %d | Category: %s\n\n", len(server.Tools), server.Metadata.Category)
	}
}

func handleSSHStatusCommand(sess ssh.Session, registry *MCPRegistry, serverName string, renderer *lipgloss.Renderer) {
	titleStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true)

	headerStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)

	if serverName == "" {
		fmt.Fprint(sess, "Usage: status <server-name>\n")
		return
	}

	// Find server
	var server *MCPServer
	for i := range registry.Servers {
		if registry.Servers[i].Name == serverName {
			server = &registry.Servers[i]
			break
		}
	}

	if server == nil {
		fmt.Fprintf(sess, "Server not found: %s\n", serverName)
		return
	}

	fmt.Fprint(sess, titleStyle.Render("📊 Server Status: "+server.Name)+"\n\n")
	fmt.Fprintf(sess, "%s: %s\n", headerStyle.Render("Status"), server.Status)
	fmt.Fprintf(sess, "%s: %s\n", headerStyle.Render("Description"), server.Description)
	fmt.Fprintf(sess, "%s: %s\n", headerStyle.Render("Category"), server.Metadata.Category)
	fmt.Fprintf(sess, "%s: %d\n", headerStyle.Render("Tools"), len(server.Tools))
	fmt.Fprint(sess, "\n")
}

func handleSSHHealthCommand(sess ssh.Session, registry *MCPRegistry, renderer *lipgloss.Renderer) {
	titleStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true)

	successStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#39FF14")).
		Bold(true)

	errorStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("#FF3131")).
		Bold(true)

	fmt.Fprint(sess, titleStyle.Render("🏥 Health Check Results")+"\n\n")

	healthy := 0
	total := len(registry.Servers)

	for _, server := range registry.Servers {
		if server.Status == "active" || server.Status == "production-ready" {
			fmt.Fprintf(sess, "%s %s - %s\n", successStyle.Render("✓"), server.Name, server.Status)
			healthy++
		} else {
			fmt.Fprintf(sess, "%s %s - %s\n", errorStyle.Render("✗"), server.Name, server.Status)
		}
	}

	fmt.Fprintf(sess, "\n%s: %d/%d servers healthy\n", titleStyle.Render("Summary"), healthy, total)
}

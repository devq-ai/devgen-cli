package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashboard-specific types and models
type dashboardModel struct {
	servers      []MCPServer
	spinner      spinner.Model
	loading      bool
	selected     int
	gridWidth    int
	gridHeight   int
	registry     *MCPRegistry
	dataLoadedAt time.Time
}

type serversLoadedMsg struct {
	registry *MCPRegistry
	loadedAt time.Time
}

type serverToggledMsg struct{}

// Dashboard styles
var (
	dashboardTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true).
		Padding(1, 2)

	dashboardHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)

	dashboardItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E3E3E3"))

	dashboardSelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF10F0")).
		Bold(true)

	dashboardStatusRunning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#39FF14")).
		Bold(true)

	dashboardStatusStopped = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF3131")).
		Bold(true)
)


// Initialize the dashboard
func (m dashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadServers(),
	)
}

// Debug function to log key events
func logKeyEvent(msg tea.KeyMsg) {
	logFile, _ := os.OpenFile("key_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	fmt.Fprintf(logFile, "Key: Type=%d, Alt=%t, String=%s, Runes=%v\n", 
		msg.Type, msg.Alt, msg.String(), msg.Runes)
}

// Update handles dashboard events  
func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Use msg.String() for more reliable key detection across terminals
		keyStr := msg.String()
		
		// Log key events to both file and Logfire
		logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "KEY EVENT: Type=%d, String='%s', Runes=%v\n", msg.Type, msg.String(), msg.Runes)
		fmt.Fprintf(logFile, "KEY STRING: '%s'\n", keyStr)
		logFile.Close()
		
		// Log to Logfire
		logToLogfire("info", "Dashboard key event", map[string]interface{}{
			"key_string": keyStr,
			"key_type":   msg.Type,
			"servers_count": len(m.servers),
			"selected_index": m.selected,
		})
		
		switch keyStr {
		case "ctrl+c":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "QUIT: Ctrl+C detected\n")
			logFile.Close()
			return m, tea.Quit
		case "enter":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "ENTER: servers=%d, selected=%d\n", len(m.servers), m.selected)
			logFile.Close()
			
			logToLogfire("info", "Dashboard enter key pressed", map[string]interface{}{
				"servers_count": len(m.servers),
				"selected_index": m.selected,
			})
			
			if len(m.servers) > 0 && m.selected < len(m.servers) {
				serverName := m.servers[m.selected].Name
				logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				fmt.Fprintf(logFile, "TOGGLE: Calling toggleServerCmd for %s\n", serverName)
				logFile.Close()
				
				logToLogfire("info", "Toggling server status", map[string]interface{}{
					"server_name": serverName,
					"current_status": m.servers[m.selected].Status,
				})
				
				return m, m.toggleServerCmd(serverName)
			}
			return m, nil
		case " ":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "SPACE: servers=%d, selected=%d\n", len(m.servers), m.selected)
			logFile.Close()
			if len(m.servers) > 0 && m.selected < len(m.servers) {
				logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				fmt.Fprintf(logFile, "TOGGLE: Calling toggleServerCmd for %s\n", m.servers[m.selected].Name)
				logFile.Close()
				return m, m.toggleServerCmd(m.servers[m.selected].Name)
			}
			return m, nil
		case "up", "k":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "UP: selected %d -> %d\n", m.selected, m.selected-1)
			logFile.Close()
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "down", "j":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "DOWN: selected %d -> %d\n", m.selected, m.selected+1)
			logFile.Close()
			if m.selected < len(m.servers)-1 {
				m.selected++
			}
			return m, nil
		case "left", "h":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "LEFT: selected %d -> %d\n", m.selected, m.selected-1)
			logFile.Close()
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case "right", "l":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "RIGHT: selected %d -> %d\n", m.selected, m.selected+1)
			logFile.Close()
			if m.selected < len(m.servers)-1 {
				m.selected++
			}
			return m, nil
		case "q":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "QUIT: q key detected\n")
			logFile.Close()
			return m, tea.Quit
		case "r":
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "REFRESH: r key detected, setting loading=true\n")
			logFile.Close()
			
			// Create new model with loading state and return it with the command
			newModel := m
			newModel.loading = true
			return newModel, newModel.loadServers()
		}
		return m, nil
		
	case serversLoadedMsg:
		// Log UI state update
		logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "MSG: serversLoadedMsg received\n")
		
		m.loading = false
		m.registry = msg.registry
		m.dataLoadedAt = msg.loadedAt
		if msg.registry != nil {
			m.servers = msg.registry.Servers
			fmt.Fprintf(logFile, "UI UPDATE: Set %d servers in model\n", len(m.servers))
			
			// Log the crawl4ai-mcp server status in the UI model
			for i, server := range m.servers {
				if server.Name == "crawl4ai-mcp" {
					fmt.Fprintf(logFile, "  UI MODEL TARGET: %s (status: %s) at index %d\n", server.Name, server.Status, i)
					break
				}
			}
			
			// Log currently selected server
			if m.selected < len(m.servers) {
				fmt.Fprintf(logFile, "  Selected server: %d = %s (status: %s)\n", m.selected, m.servers[m.selected].Name, m.servers[m.selected].Status)
			}
		} else {
			m.servers = []MCPServer{}
			fmt.Fprintf(logFile, "UI UPDATE: Set empty servers array\n")
		}
		if m.selected >= len(m.servers) {
			m.selected = len(m.servers) - 1
		}
		if m.selected < 0 {
			m.selected = 0
		}
		logFile.Close()
		return m, nil
	case serverToggledMsg:
		// Log that we received the toggle message
		logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "MSG: Received serverToggledMsg, triggering reload\n")
		logFile.Close()
		return m, m.loadServers()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// Render the dashboard view
func (m dashboardModel) View() string {
	if m.loading {
		return fmt.Sprintf("\n%s Loading servers...\n", m.spinner.View())
	}

	header := dashboardTitleStyle.Render("🔌 MCP Server Dashboard")
	footer := dashboardItemStyle.Render("Press 'enter/space' to toggle, 'q' to quit, arrow keys/hjkl to navigate")

	// Debug info with timestamp
	dataLoadedTime := "never"
	if !m.dataLoadedAt.IsZero() {
		dataLoadedTime = m.dataLoadedAt.Format("15:04:05")
	}
	debugInfo := fmt.Sprintf("Servers: %d | Data loaded at: %s", len(m.servers), dataLoadedTime)
	if len(m.servers) > 0 {
		selectedServer := "none"
		if m.selected < len(m.servers) {
			selectedServer = m.servers[m.selected].Name
		}
		debugInfo += fmt.Sprintf(" | Selected: %d (%s)", m.selected, selectedServer)
	}

	// Simple list - show ALL servers without pagination
	var serverList strings.Builder
	renderedCount := 0
	
	for i, server := range m.servers {
		serverLine := m.renderServerCard(server, i == m.selected)
		serverList.WriteString(serverLine)
		renderedCount++
		
		// Add single line spacing between servers
		if i < len(m.servers)-1 {
			serverList.WriteString("\n")
		}
	}
	debugInfo += fmt.Sprintf(" | Rendered: %d", renderedCount)

	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", header, debugInfo, serverList.String(), footer)
}


// Load servers from registry
func (m dashboardModel) loadServers() tea.Cmd {
	return func() tea.Msg {
		loadTime := time.Now()
		
		// Log reload attempt
		logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "LOAD: Starting loadServers from configFile=%s\n", configFile)
		logFile.Close()
		
		registry, err := loadMCPRegistry()
		if err != nil {
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "LOAD ERROR: %v\n", err)
			logFile.Close()
			
			// Create empty registry on error
			emptyRegistry := &MCPRegistry{
				Version:   "ERROR",
				Timestamp: fmt.Sprintf("Load error: %v", err),
				Servers:   []MCPServer{},
			}
			return serversLoadedMsg{registry: emptyRegistry, loadedAt: loadTime}
		}

		logFile, _ = os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "LOAD: Registry loaded successfully, %d servers\n", len(registry.Servers))
		logFile.Close()

		// NOTE: Removed connectivity test override that was causing all "active" servers 
		// to be displayed as "inactive". The stored status values are accurate.
		// If connectivity testing is needed, it should be implemented separately
		// without overriding the display status.

		// Log loaded servers for debugging
		logFile, _ = os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "LOAD: Returning serversLoadedMsg with %d servers\n", len(registry.Servers))
		
		// Find and log crawl4ai-mcp specifically
		for i, server := range registry.Servers {
			if server.Name == "crawl4ai-mcp" {
				fmt.Fprintf(logFile, "  LOAD TARGET: %s (status: %s) at index %d\n", server.Name, server.Status, i)
				break
			}
		}
		logFile.Close()

		return serversLoadedMsg{registry: registry, loadedAt: loadTime}
	}
}

// Toggle server command - implement inline like CLI to avoid context issues
func (m dashboardModel) toggleServerCmd(serverName string) tea.Cmd {
	return func() tea.Msg {
		// Log toggle attempt
		logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "TOGGLE CMD: Starting toggle for server '%s' using configFile=%s\n", serverName, configFile)
		logFile.Close()
		
		// Load registry fresh (same as CLI)
		registry, err := loadMCPRegistry()
		if err != nil {
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "TOGGLE CMD ERROR: Failed to load registry: %v\n", err)
			logFile.Close()
			return serverToggledMsg{} // Still trigger reload even on error
		}

		logFile, _ = os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(logFile, "TOGGLE CMD: Registry loaded, %d servers\n", len(registry.Servers))
		logFile.Close()

		// Toggle the server status (same logic as toggleServer function)
		found := false
		for i := range registry.Servers {
			if registry.Servers[i].Name == serverName {
				oldStatus := registry.Servers[i].Status
				if registry.Servers[i].Status == "active" || registry.Servers[i].Status == "production-ready" || registry.Servers[i].Status == "running" {
					registry.Servers[i].Status = "inactive"
				} else {
					registry.Servers[i].Status = "active"
				}
				logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				fmt.Fprintf(logFile, "TOGGLE CMD: Changed %s status from '%s' to '%s'\n", serverName, oldStatus, registry.Servers[i].Status)
				logFile.Close()
				found = true
				break
			}
		}

		if !found {
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "TOGGLE CMD ERROR: Server '%s' not found\n", serverName)
			logFile.Close()
		}

		// Save registry fresh (same as CLI)
		err = saveMCPRegistry(registry)
		if err != nil {
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "TOGGLE CMD ERROR: Failed to save registry: %v\n", err)
			logFile.Close()
		} else {
			logFile, _ := os.OpenFile("dashboard_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(logFile, "TOGGLE CMD: Registry saved successfully\n")
			logFile.Close()
		}
		
		// Add small delay to ensure file write completes before triggering reload
		time.Sleep(50 * time.Millisecond)
		
		return serverToggledMsg{}
	}
}

// Create and run the dashboard
func runDashboard() error {
	// Log dashboard startup to Logfire
	logToLogfire("info", "Dashboard starting up", map[string]interface{}{
		"config_file": configFile,
		"terminal": os.Getenv("TERM"),
	})
	
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF10F0"))

	// Create dashboard model
	m := dashboardModel{
		servers:      []MCPServer{},
		spinner:      s,
		loading:      true,
		selected:     0,
		gridWidth:    1,
		gridHeight:   13,
		registry:     nil,
		dataLoadedAt: time.Time{},
	}

	// Run the dashboard with Ghostty terminal optimizations
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	
	// Log dashboard shutdown
	logToLogfire("info", "Dashboard shutting down", map[string]interface{}{
		"error": fmt.Sprintf("%v", err),
	})
	
	return err
}
// Render a single server card with proper text wrapping
func (m dashboardModel) renderServerCard(server MCPServer, selected bool) string {
	// Determine status color
	var statusStyle lipgloss.Style
	if server.Status == "active" || server.Status == "production-ready" || server.Status == "running" {
		statusStyle = dashboardStatusRunning
	} else {
		statusStyle = dashboardStatusStopped
	}

	// Simple inline style for single column layout
	nameStyle := lipgloss.NewStyle().Bold(true)
	if selected {
		nameStyle = nameStyle.Foreground(lipgloss.Color("#FF10F0"))
	}
	
	// Wrap description to terminal width
	description := wrapText(server.Description, 80)
	
	// Category icon mapping
	categoryIcon := map[string]string{
		"knowledge":       "🧠",
		"development":     "⚡",
		"web":            "🌐",
		"framework":      "🔧",
		"database":       "💾",
		"infrastructure": "🏗️",
	}
	
	icon := categoryIcon[server.Metadata.Category]
	if icon == "" {
		icon = "📦"
	}
	
	// Build simple one-line format with wrapped description
	line1 := fmt.Sprintf("%s %s [%s • %d tools]", 
		icon,
		nameStyle.Render(server.Name),
		statusStyle.Render(server.Status),
		len(server.Tools))
	
	line2 := fmt.Sprintf("   %s", description)
	
	// Return simple formatted text
	if selected {
		return fmt.Sprintf("▶ %s\n%s", line1, line2)
	}
	return fmt.Sprintf("  %s\n%s", line1, line2)
}

// Helper function to wrap text
func wrapText(text string, width int) string {
	words := strings.Fields(text)
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		if currentLine.Len() + len(word) + 1 > width {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
		}
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	// Return first 2 lines
	if len(lines) > 2 {
		return lines[0] + "\n  " + lines[1] + "..."
	} else if len(lines) > 1 {
		return lines[0] + "\n  " + lines[1]
	}
	return text
}

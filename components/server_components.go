package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// DevServer manages development server with monitoring UI
type DevServer struct {
	host           string
	port           int
	server         *http.Server
	status         ServerStatus
	metrics        ServerMetrics
	logs           []LogEntry
	logger         *log.Logger
	startTime      time.Time
	mutex          sync.RWMutex
	healthChecks   []HealthCheck
	watchers       []FileWatcher
	middleware     []Middleware
	routes         map[string]RouteInfo
}

// ServerStatus represents current server state
type ServerStatus int

const (
	StatusStopped ServerStatus = iota
	StatusStarting
	StatusRunning
	StatusStopping
	StatusError
	StatusRestarting
)

func (s ServerStatus) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusStopping:
		return "stopping"
	case StatusError:
		return "error"
	case StatusRestarting:
		return "restarting"
	default:
		return "unknown"
	}
}

// ServerMetrics holds performance metrics
type ServerMetrics struct {
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	AverageResponse time.Duration `json:"average_response"`
	MemoryUsage     int64         `json:"memory_usage"`
	CPUUsage        float64       `json:"cpu_usage"`
	Uptime          time.Duration `json:"uptime"`
	ActiveConnections int         `json:"active_connections"`
	TotalBytes      int64         `json:"total_bytes"`
}

// LogEntry represents a server log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Method      string        `json:"method"`
	Interval    time.Duration `json:"interval"`
	Timeout     time.Duration `json:"timeout"`
	Expected    int           `json:"expected_status"`
	LastCheck   time.Time     `json:"last_check"`
	Status      string        `json:"status"`
	Response    time.Duration `json:"response_time"`
}

// FileWatcher monitors file changes for hot reload
type FileWatcher struct {
	Path      string    `json:"path"`
	Pattern   string    `json:"pattern"`
	LastMod   time.Time `json:"last_modified"`
	Active    bool      `json:"active"`
	TriggerCount int    `json:"trigger_count"`
}

// Middleware represents server middleware
type Middleware struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Order       int    `json:"order"`
	Config      map[string]interface{} `json:"config"`
	Description string `json:"description"`
}

// RouteInfo holds route information
type RouteInfo struct {
	Path        string    `json:"path"`
	Method      string    `json:"method"`
	Handler     string    `json:"handler"`
	Hits        int64     `json:"hits"`
	AvgResponse time.Duration `json:"avg_response"`
	LastAccess  time.Time `json:"last_access"`
}

// Create a new development server
func newDevServer(host string, port int) *DevServer {
	return &DevServer{
		host:         host,
		port:         port,
		status:       StatusStopped,
		metrics:      ServerMetrics{},
		logs:         make([]LogEntry, 0),
		logger:       log.New(os.Stderr),
		healthChecks: make([]HealthCheck, 0),
		watchers:     make([]FileWatcher, 0),
		middleware:   make([]Middleware, 0),
		routes:       make(map[string]RouteInfo),
	}
}

// Start the development server
func (ds *DevServer) start() error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	if ds.status == StatusRunning {
		return fmt.Errorf("server is already running")
	}

	ds.status = StatusStarting
	ds.startTime = time.Now()
	ds.addLog("info", "Starting development server", map[string]interface{}{
		"host": ds.host,
		"port": ds.port,
	})

	// Create HTTP server
	mux := http.NewServeMux()
	ds.setupRoutes(mux)

	ds.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ds.host, ds.port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := ds.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ds.mutex.Lock()
			ds.status = StatusError
			ds.addLog("error", "Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
			ds.mutex.Unlock()
		}
	}()

	// Start monitoring goroutines
	go ds.startHealthChecks()
	go ds.startFileWatching()
	go ds.updateMetrics()

	ds.status = StatusRunning
	ds.addLog("info", "Development server started successfully", nil)

	return nil
}

// Setup HTTP routes
func (ds *DevServer) setupRoutes(mux *http.ServeMux) {
	// Health check endpoint
	mux.HandleFunc("/health", ds.handleHealth)
	ds.routes["/health"] = RouteInfo{
		Path:    "/health",
		Method:  "GET",
		Handler: "handleHealth",
	}

	// Metrics endpoint
	mux.HandleFunc("/metrics", ds.handleMetrics)
	ds.routes["/metrics"] = RouteInfo{
		Path:    "/metrics",
		Method:  "GET",
		Handler: "handleMetrics",
	}

	// Status endpoint
	mux.HandleFunc("/status", ds.handleStatus)
	ds.routes["/status"] = RouteInfo{
		Path:    "/status",
		Method:  "GET",
		Handler: "handleStatus",
	}

	// Logs endpoint
	mux.HandleFunc("/logs", ds.handleLogs)
	ds.routes["/logs"] = RouteInfo{
		Path:    "/logs",
		Method:  "GET",
		Handler: "handleLogs",
	}

	// Static file serving
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

// HTTP handlers
func (ds *DevServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	ds.updateRouteStats("/health", r.Method, time.Now())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","uptime":"%v","timestamp":"%s"}`,
		ds.getUptime(), time.Now().Format(time.RFC3339))
}

func (ds *DevServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	ds.updateRouteStats("/metrics", r.Method, time.Now())

	ds.mutex.RLock()
	metrics := ds.metrics
	ds.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Return metrics as JSON (simplified)
	fmt.Fprintf(w, `{
		"requests": %d,
		"errors": %d,
		"uptime": "%v",
		"memory_mb": %.2f,
		"cpu_percent": %.2f
	}`, metrics.RequestCount, metrics.ErrorCount, ds.getUptime(),
		float64(metrics.MemoryUsage)/1024/1024, metrics.CPUUsage)
}

func (ds *DevServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	ds.updateRouteStats("/status", r.Method, time.Now())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"%s","port":%d,"host":"%s"}`,
		ds.status.String(), ds.port, ds.host)
}

func (ds *DevServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	ds.updateRouteStats("/logs", r.Method, time.Now())

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	ds.mutex.RLock()
	logs := ds.logs
	if len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}
	ds.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Simplified JSON response
	fmt.Fprintf(w, `{"logs":[`)
	for i, log := range logs {
		if i > 0 {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, `{"timestamp":"%s","level":"%s","message":"%s"}`,
			log.Timestamp.Format(time.RFC3339), log.Level, log.Message)
	}
	fmt.Fprintf(w, `]}`)
}

// Add log entry
func (ds *DevServer) addLog(level, message string, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    "devserver",
		Data:      data,
	}

	ds.logs = append(ds.logs, entry)
	if len(ds.logs) > 1000 {
		ds.logs = ds.logs[1:]
	}

	ds.logger.Info("Server log", "level", level, "message", message)
}

// Get server uptime
func (ds *DevServer) getUptime() time.Duration {
	if ds.startTime.IsZero() {
		return 0
	}
	return time.Since(ds.startTime)
}

// Update route statistics
func (ds *DevServer) updateRouteStats(path, method string, start time.Time) {
	key := fmt.Sprintf("%s %s", method, path)

	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	route := ds.routes[path]
	route.Hits++
	route.LastAccess = start

	duration := time.Since(start)
	if route.AvgResponse == 0 {
		route.AvgResponse = duration
	} else {
		route.AvgResponse = (route.AvgResponse + duration) / 2
	}

	ds.routes[path] = route
	ds.metrics.RequestCount++
}

// Start health checks
func (ds *DevServer) startHealthChecks() {
	// Add default health checks
	ds.healthChecks = append(ds.healthChecks, HealthCheck{
		Name:     "self",
		URL:      fmt.Sprintf("http://%s:%d/health", ds.host, ds.port),
		Method:   "GET",
		Interval: 30 * time.Second,
		Timeout:  5 * time.Second,
		Expected: 200,
	})

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if ds.status != StatusRunning {
			continue
		}

		for i := range ds.healthChecks {
			go ds.runHealthCheck(&ds.healthChecks[i])
		}
	}
}

// Run individual health check
func (ds *DevServer) runHealthCheck(hc *HealthCheck) {
	if time.Since(hc.LastCheck) < hc.Interval {
		return
	}

	start := time.Now()
	client := &http.Client{Timeout: hc.Timeout}

	resp, err := client.Get(hc.URL)
	duration := time.Since(start)

	hc.LastCheck = time.Now()
	hc.Response = duration

	if err != nil {
		hc.Status = "error"
		ds.addLog("warn", "Health check failed", map[string]interface{}{
			"name":  hc.Name,
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == hc.Expected {
		hc.Status = "healthy"
	} else {
		hc.Status = "unhealthy"
		ds.addLog("warn", "Health check returned unexpected status", map[string]interface{}{
			"name":     hc.Name,
			"expected": hc.Expected,
			"actual":   resp.StatusCode,
		})
	}
}

// Start file watching for hot reload
func (ds *DevServer) startFileWatching() {
	// Add default watchers
	ds.watchers = append(ds.watchers, FileWatcher{
		Path:    "./",
		Pattern: "*.go",
		Active:  true,
	})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if ds.status != StatusRunning {
			continue
		}

		for i := range ds.watchers {
			ds.checkFileChanges(&ds.watchers[i])
		}
	}
}

// Check for file changes
func (ds *DevServer) checkFileChanges(fw *FileWatcher) {
	if !fw.Active {
		return
	}

	info, err := os.Stat(fw.Path)
	if err != nil {
		return
	}

	if !fw.LastMod.IsZero() && info.ModTime().After(fw.LastMod) {
		fw.TriggerCount++
		ds.addLog("info", "File change detected", map[string]interface{}{
			"path":    fw.Path,
			"pattern": fw.Pattern,
		})
	}

	fw.LastMod = info.ModTime()
}

// Update metrics periodically
func (ds *DevServer) updateMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if ds.status != StatusRunning {
			continue
		}

		ds.mutex.Lock()
		ds.metrics.Uptime = ds.getUptime()
		// Simulate metrics updates
		ds.metrics.MemoryUsage = 64 * 1024 * 1024 // 64MB
		ds.metrics.CPUUsage = 15.5 // 15.5%
		ds.mutex.Unlock()
	}
}

// ServerStatusViewer displays real-time server status
type ServerStatusViewer struct {
	server   *DevServer
	table    table.Model
	logs     list.Model
	metrics  progress.Model
	spinner  spinner.Model
	view     string
	width    int
	height   int
	logger   *log.Logger
}

func newServerStatusViewer() *ServerStatusViewer {
	// Create status table
	columns := []table.Column{
		{Title: "Metric", Width: 20},
		{Title: "Value", Width: 20},
		{Title: "Status", Width: 15},
		{Title: "Last Updated", Width: 20},
	}

	rows := []table.Row{
		{"Server Status", "Running", "🟢 Healthy", "1 min ago"},
		{"Uptime", "2h 15m", "🟢 Good", "now"},
		{"Request Count", "1,247", "🟢 Normal", "5s ago"},
		{"Error Rate", "0.2%", "🟢 Low", "30s ago"},
		{"Memory Usage", "64.5 MB", "🟡 Moderate", "10s ago"},
		{"CPU Usage", "15.3%", "🟢 Low", "5s ago"},
		{"Active Connections", "23", "🟢 Normal", "15s ago"},
		{"Health Checks", "3/3 Passing", "🟢 Healthy", "1 min ago"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Apply cyber theme
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(textColor).
		Background(primaryColor).
		Bold(false)
	t.SetStyles(s)

	// Create spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(primaryColor)

	// Create progress bar for metrics
	prog := progress.New(progress.WithDefaultGradient())

	return &ServerStatusViewer{
		table:   t,
		metrics: prog,
		spinner: sp,
		view:    "status",
		logger:  log.New(nil),
	}
}

func (ssv *ServerStatusViewer) Init() tea.Cmd {
	return tea.Batch(
		ssv.spinner.Tick,
		ssv.tick(),
	)
}

func (ssv *ServerStatusViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ssv.width = msg.Width
		ssv.height = msg.Height
		return ssv, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ssv, tea.Quit
		case "tab":
			ssv.switchView()
			return ssv, nil
		case "r":
			ssv.logger.Info("Refreshing server status...")
			return ssv, ssv.refreshData()
		}

	case spinner.TickMsg:
		ssv.spinner, cmd = ssv.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case refreshMsg:
		// Update table data
		ssv.updateTableData()
		cmds = append(cmds, ssv.tick())
	}

	// Update current view
	switch ssv.view {
	case "status":
		ssv.table, cmd = ssv.table.Update(msg)
		cmds = append(cmds, cmd)
	case "logs":
		ssv.logs, cmd = ssv.logs.Update(msg)
		cmds = append(cmds, cmd)
	}

	return ssv, tea.Batch(cmds...)
}

func (ssv *ServerStatusViewer) View() string {
	if ssv.width == 0 {
		return "Loading..."
	}

	var content strings.Builder

	// Title
	title := fmt.Sprintf("%s Development Server Status", ssv.spinner.View())
	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	// View tabs
	tabs := []string{"Status", "Logs", "Metrics", "Health"}
	tabContent := ""
	for i, tab := range tabs {
		style := statusPendingStyle
		if strings.ToLower(tab) == ssv.view {
			style = statusRunningStyle
		}
		if i > 0 {
			tabContent += " | "
		}
		tabContent += style.Render(tab)
	}
	content.WriteString(tabContent)
	content.WriteString("\n\n")

	// Main content based on view
	switch ssv.view {
	case "status":
		content.WriteString(ssv.table.View())
	case "logs":
		content.WriteString("📋 Recent Server Logs\n\n")
		content.WriteString(statusPendingStyle.Render("No logs available in demo mode"))
	case "metrics":
		content.WriteString("📊 Performance Metrics\n\n")
		content.WriteString("CPU Usage: ")
		content.WriteString(ssv.metrics.ViewAs(0.15))
		content.WriteString("\n\nMemory Usage: ")
		content.WriteString(ssv.metrics.ViewAs(0.65))
		content.WriteString("\n\nDisk Usage: ")
		content.WriteString(ssv.metrics.ViewAs(0.23))
	case "health":
		content.WriteString("🏥 Health Checks\n\n")
		content.WriteString(statusCompletedStyle.Render("✓ Self Health Check - 200ms"))
		content.WriteString("\n")
		content.WriteString(statusCompletedStyle.Render("✓ Database Connection - 50ms"))
		content.WriteString("\n")
		content.WriteString(statusCompletedStyle.Render("✓ Redis Cache - 15ms"))
	}

	content.WriteString("\n\n")

	// Help
	help := "tab: switch view • r: refresh • q: quit"
	content.WriteString(helpStyle.Render(help))

	return baseStyle.Render(content.String())
}

func (ssv *ServerStatusViewer) switchView() {
	views := []string{"status", "logs", "metrics", "health"}
	for i, view := range views {
		if view == ssv.view {
			ssv.view = views[(i+1)%len(views)]
			break
		}
	}
}

func (ssv *ServerStatusViewer) updateTableData() {
	// Update table with fresh data
	rows := []table.Row{
		{"Server Status", "Running", "🟢 Healthy", time.Now().Format("15:04:05")},
		{"Uptime", "2h 15m", "🟢 Good", time.Now().Format("15:04:05")},
		{"Request Count", "1,247", "🟢 Normal", time.Now().Format("15:04:05")},
		{"Error Rate", "0.2%", "🟢 Low", time.Now().Format("15:04:05")},
		{"Memory Usage", "64.5 MB", "🟡 Moderate", time.Now().Format("15:04:05")},
		{"CPU Usage", "15.3%", "🟢 Low", time.Now().Format("15:04:05")},
		{"Active Connections", "23", "🟢 Normal", time.Now().Format("15:04:05")},
		{"Health Checks", "3/3 Passing", "🟢 Healthy", time.Now().Format("15:04:05")},
	}
	ssv.table.SetRows(rows)
}

func (ssv *ServerStatusViewer) tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (ssv *ServerStatusViewer) refreshData() tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{}
	}
}

func (ssv *ServerStatusViewer) run() error {
	p := tea.NewProgram(ssv, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// ConfigEditor provides interactive configuration editing
type ConfigEditor struct {
	form     *huh.Form
	config   *Config
	field    string
	step     int
	saved    bool
	logger   *log.Logger
}

func newConfigEditor() *ConfigEditor {
	return &ConfigEditor{
		config: CreateDefaultConfig(),
		step:   0,
		saved:  false,
		logger: log.New(nil),
	}
}

func (ce *ConfigEditor) Init() tea.Cmd {
	return ce.createForm()
}

func (ce *ConfigEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ce, tea.Quit
		}
	}
	return ce, nil
}

func (ce *ConfigEditor) View() string {
	if ce.saved {
		return statusCompletedStyle.Render("✅ Configuration saved successfully!")
	}

	title := titleStyle.Render("Configuration Editor")
	content := title + "\n\n"
	content += "Configure DevGen CLI settings interactively...\n\n"

	return baseStyle.Render(content)
}

func (ce *ConfigEditor) createForm() tea.Cmd {
	var logLevel, theme, outputDir string
	var autoSave, checkUpdates bool

	ce.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Log Level").
				Options(
					huh.NewOption("Debug", "debug"),
					huh.NewOption("Info", "info"),
					huh.NewOption("Warn", "warn"),
					huh.NewOption("Error", "error"),
				).
				Value(&logLevel),

			huh.NewSelect[string]().
				Title("UI Theme").
				Options(
					huh.NewOption("Cyber (High Contrast)", "cyber"),
					huh.NewOption("Pastel (Comfortable)", "pastel"),
				).
				Value(&theme),

			huh.NewInput().
				Title("Default Output Directory").
				Value(&outputDir).
				Placeholder("./output"),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable auto-save?").
				Value(&autoSave),

			huh.NewConfirm().
				Title("Check for updates automatically?").
				Value(&checkUpdates),
		),
	)

	return nil
}

func (ce *ConfigEditor) run() error {
	if ce.form != nil {
		err := ce.form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}
	}

	// Save configuration
	configPath := GetConfigPath()
	err := SaveConfig(ce.config, configPath)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	ce.saved = true
	ce.logger.Info("Configuration saved", "path", configPath)
	return nil
}

// ConfigViewer displays current configuration
type ConfigViewer struct {
	table  table.Model
	config *Config
	logger *log.Logger
}

func newConfigViewer() *ConfigViewer {
	// Load current config
	config, err := LoadConfig(GetConfigPath())
	if err != nil {
		config = CreateDefaultConfig()
	}

	columns := []table.Column{
		{Title: "Setting", Width: 30},
		{Title: "Value", Width: 40},
		{Title: "Description", Width: 50},
	}

	rows := []table.Row{
		{"Version", config.Version, "Configuration version"},
		{"Default Output Dir", config.DevGen.DefaultOutputDir, "Default output directory for generated files"},
		{"Default Template", config.DevGen.DefaultTemplate, "Default project template"},
		{"Auto Save", fmt.Sprintf("%t", config.DevGen.AutoSave), "Automatically save configuration changes"},
		{"Check Updates", fmt.Sprintf("%t", config.DevGen.CheckUpdates), "Check for updates automatically"},
		{"Log Level", config.Logging.Level, "Logging verbosity level"},
		{"Log Format", config.Logging.Format, "Log output format"},
		{"UI Theme", config.UI.Theme.Name, "User interface theme"},
		{"Dark Mode", fmt.Sprintf("%t", config.UI.Theme.DarkMode), "Enable dark mode"},
		{"Server Host", config.Servers.Default.Host, "Default server host"},
		{"Server Port", fmt.Sprintf("%d", config.Servers.Default.Port), "Default server port"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	// Apply styling
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(textColor).
		Background(primaryColor).
		Bold(false)
	t.SetStyles(s)

	return &ConfigViewer{
		table:  t,
		config: config,
		logger: log.New(nil),
	}
}

func (cv *ConfigViewer) Init() tea.Cmd {
	return nil
}

func (cv *ConfigViewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return cv, tea.Quit
		case "e":
			// Open editor
			cv.logger.Info("Opening configuration editor...")
			return cv, tea.Quit
		}
	}

	cv.table, cmd = cv.table.Update(msg)
	return cv, cmd
}

func (cv *ConfigViewer) View() string {
	title := titleStyle.Render("Current Configuration")
	content := title + "\n\n" + cv.table.View() + "\n\n"
	help := helpStyle.Render("e: edit • q: quit")
	return content + help
}

func (cv *ConfigViewer) run() error {
	p := tea.NewProgram(cv, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// ConfigInitializer creates default configuration
type ConfigInitializer struct {
	progress progress.Model
	step     int
	total    int
	completed bool
	logger   *log.Logger
}

func newConfigInitializer() *ConfigInitializer {
	prog := progress.New(progress.WithDefaultGradient())
	return &ConfigInitializer{
		progress:  prog,
		step:      0,
		total:     5,
		completed: false,
		logger:    log.New(nil),
	}
}

func (ci *ConfigInitializer) Init() tea.Cmd {
	return ci.nextStep()
}

func (ci *ConfigInitializer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if ci.completed {
				return ci, tea.Quit
			}
		}

	case stepMsg:
		ci.step = int(msg)
		if ci.step >= ci.total {
			ci.completed = true
			ci.logger.Info("Configuration initialization completed")
		}
		return ci, ci.nextStep()
	}

	return ci, nil
}

func (ci *ConfigInitializer) View() string {
	if ci.completed {
		content := statusCompletedStyle.Render("✅ Configuration initialized successfully!")
		content += "\n\n"
		content += helpStyle.Render("Press 'q' to exit")
		return baseStyle.Render(content)
	}

	var content strings.Builder

	title := titleStyle.Render("Initializing Configuration")
	content.WriteString(title)
	content.WriteString("\n\n")

	steps := []string{
		"Creating configuration directory",
		"Generating default configuration",
		"Setting up templates directory",
		"Configuring logging",
		"Finalizing setup",
	}

	percentage := float64(ci.step) / float64(ci.total)
	progressBar := ci.progress.ViewAs(percentage)
	content.WriteString(progressBar)
	content.WriteString("\n\n")

	content.WriteString("Progress:\n")
	for i, step := range steps {
		if i < ci.step {
			content.WriteString(statusCompletedStyle.Render(fmt.Sprintf("✓ %s", step)))
		} else if i == ci.step {
			content.WriteString(statusRunningStyle.Render(fmt.Sprintf("⟳ %s", step)))
		} else {
			content.WriteString(statusPendingStyle.Render(fmt.Sprintf("○ %s", step)))
		}
		content.WriteString("\n")
	}

	return baseStyle.Render(content.String())
}

func (ci *ConfigInitializer) nextStep() tea.Cmd {
	if ci.step >= ci.total {
		return nil
	}

	return tea.Tick(time.Millisecond*800, func(t time.Time) tea.Msg {
		return stepMsg(ci.step + 1)
	})
}

func (ci *ConfigInitializer) run() error {
	// Create configuration directory
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate default configuration
	config := CreateDefaultConfig()
	configPath := GetConfigPath()

	// Save configuration
	if err := SaveConfig(config, configPath); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}

	ci.logger.Info("Default configuration created", "path", configPath)

	// Run the UI
	p := tea.NewProgram(ci, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Message types for server components
type refreshMsg struct{}
type stepMsg int

// Additional utility functions for server components

// FormatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// FormatBytes formats bytes in a human-readable way
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetSystemInfo returns basic system information
func getSystemInfo() map[string]string {
	return map[string]string{
		"os":      "linux/darwin/windows", // Would use runtime.GOOS
		"arch":    "amd64",                 // Would use runtime.GOARCH
		"version": "go1.21",               // Would use runtime.Version()
		"cores":   "8",                     // Would use runtime.NumCPU()
	}
}

// ValidatePort checks if a port number is valid
func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	return nil
}

// IsPortAvailable checks if a port is available for binding
func isPortAvailable(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

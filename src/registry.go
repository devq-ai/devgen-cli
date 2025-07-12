package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// HTTP Registry Types
type HTTPRegistryServer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Port        int    `json:"port"`
	URL         string `json:"url"`
}

type HTTPRegistryTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Registry management functions
func checkRegistryStatus() error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	fmt.Printf("🔍 Checking MCP Registry Status\n")
	fmt.Printf("Registry URL: %s\n", registryURL)
	
	// Check servers endpoint
	resp, err := client.Get(registryURL + "/servers")
	if err != nil {
		fmt.Printf("❌ Registry not accessible: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Printf("❌ Registry returned status %d\n", resp.StatusCode)
		return fmt.Errorf("registry returned status %d", resp.StatusCode)
	}
	
	var servers []HTTPRegistryServer
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		fmt.Printf("❌ Failed to decode response: %v\n", err)
		return err
	}
	
	fmt.Printf("✅ Registry is active\n")
	fmt.Printf("📊 Registered servers: %d\n", len(servers))
	
	return nil
}

func listRegistryServers() error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(registryURL + "/servers")
	if err != nil {
		return fmt.Errorf("failed to connect to registry: %v", err)
	}
	defer resp.Body.Close()
	
	var servers []HTTPRegistryServer
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	
	fmt.Printf("🔌 MCP Registry Servers (%d total)\n\n", len(servers))
	
	for i, server := range servers {
		fmt.Printf("%d. %s\n", i+1, statusRunning.Render(server.Name))
		fmt.Printf("   📝 Description: %s\n", server.Description)
		fmt.Printf("   🌐 URL: %s:%d\n", server.URL, server.Port)
		fmt.Printf("\n")
	}
	
	return nil
}

func listRegistryTools() error {
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(registryURL + "/tools")
	if err != nil {
		return fmt.Errorf("failed to connect to registry: %v", err)
	}
	defer resp.Body.Close()
	
	var tools []HTTPRegistryTool
	if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	
	fmt.Printf("🛠️  MCP Registry Tools (%d total)\n\n", len(tools))
	
	// Group tools by server
	toolsByServer := make(map[string][]string)
	for _, tool := range tools {
		if strings.Contains(tool.Name, ".") {
			serverName := strings.Split(tool.Name, ".")[0]
			toolName := strings.Split(tool.Name, ".")[1]
			toolsByServer[serverName] = append(toolsByServer[serverName], toolName)
		} else {
			toolsByServer["Unknown"] = append(toolsByServer["Unknown"], tool.Name)
		}
	}
	
	for serverName, serverTools := range toolsByServer {
		fmt.Printf("📦 %s (%d tools):\n", headerStyle.Render(serverName), len(serverTools))
		for _, tool := range serverTools {
			fmt.Printf("   • %s\n", tool)
		}
		fmt.Printf("\n")
	}
	
	return nil
}

func startMCPRegistry() error {
	fmt.Printf("🚀 Starting MCP Registry...\n")
	
	// Check if already running
	client := &http.Client{Timeout: 2 * time.Second}
	if resp, err := client.Get(registryURL + "/servers"); err == nil {
		resp.Body.Close()
		fmt.Printf("✅ Registry already running at %s\n", registryURL)
		return nil
	}
	
	// Find and start the registry
	machinaRoot := findMachinaRoot()
	if machinaRoot == "" {
		return fmt.Errorf("could not find machina root directory")
	}
	
	registryPath := filepath.Join(machinaRoot, "start_registry_servers.py")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		registryPath = filepath.Join(machinaRoot, "mcp-registry", "start_registry_server.py")
		if _, err := os.Stat(registryPath); os.IsNotExist(err) {
			return fmt.Errorf("could not find registry start script")
		}
	}
	
	fmt.Printf("📂 Found registry script: %s\n", registryPath)
	
	// Start the registry in background
	cmd := exec.Command("python3", registryPath)
	cmd.Dir = filepath.Dir(registryPath)
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start registry: %v", err)
	}
	
	fmt.Printf("⏳ Waiting for registry to start...\n")
	time.Sleep(3 * time.Second)
	
	// Check if it started successfully
	if resp, err := client.Get(registryURL + "/servers"); err == nil {
		resp.Body.Close()
		fmt.Printf("✅ Registry started successfully at %s\n", registryURL)
		return nil
	} else {
		return fmt.Errorf("registry failed to start: %v", err)
	}
}
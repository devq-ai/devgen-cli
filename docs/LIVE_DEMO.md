# DevQ.ai Website Live Demo Integration

This guide shows how to set up live connections between the DevQ.ai website and your local Machina instance for real-time demos and showcases.

## 🎯 Overview

The DevGen CLI can now act as a web server, allowing the DevQ.ai website to connect directly to your local Machina instance for live demonstrations. This creates a seamless experience where website visitors can see real MCP servers in action.

```
DevQ.ai Website (https://devq.ai)
├── Interactive Demo UI
├── Real-time Server Status
└── Live Control Panel
    ↓ CORS-enabled API calls
    ↓ WebSocket connections
Local Machina Instance (Your Dev Machine)
├── DevGen Web Server (Port 8080)
├── FastAPI endpoints
├── Live MCP Registry
└── Real MCP Servers
```

## 🚀 Quick Start

### 1. Start the Web Server

```bash
# Basic web server (allows connections from devq.ai)
devgen web

# Custom port and origin
devgen web --web-port 8080 --allow-origin "https://devq.ai"

# Allow multiple origins for development
devgen web --allow-origin "*" --web-port 3000
```

### 2. Verify Connection

```bash
# Test local endpoints
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/mcp/servers
curl http://localhost:8080/api/v1/mcp/registry
```

### 3. Configure Port Forwarding (for external access)

```bash
# Option 1: SSH tunnel (secure)
ssh -R 8080:localhost:8080 your-server.com

# Option 2: ngrok (easy for demos)
ngrok http 8080

# Option 3: Expose directly (local network only)
devgen web --allow-origin "*"
```

## 📡 API Endpoints

The web server exposes these endpoints for the DevQ.ai website:

### Registry Management
```
GET  /api/v1/mcp/registry       # Get complete MCP registry
PUT  /api/v1/mcp/registry       # Update registry (admin only)
```

### Server Operations
```
GET  /api/v1/mcp/servers        # List all servers with health status
POST /api/v1/mcp/servers/{name} # Control individual servers
```

### Health Monitoring
```
GET  /api/v1/health            # Health check all servers
```

### Service Info
```
GET  /                         # API documentation and status
```

## 🎮 DevQ.ai Website Integration

### JavaScript Client Example

```javascript
// Connect to local DevGen instance
class DevGenClient {
    constructor(baseUrl = 'http://localhost:8080') {
        this.baseUrl = baseUrl;
    }

    async getServers() {
        const response = await fetch(`${this.baseUrl}/api/v1/mcp/servers`);
        return await response.json();
    }

    async toggleServer(serverName) {
        const response = await fetch(`${this.baseUrl}/api/v1/mcp/servers/${serverName}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ action: 'toggle' })
        });
        return await response.json();
    }

    async getHealth() {
        const response = await fetch(`${this.baseUrl}/api/v1/health`);
        return await response.json();
    }
}

// Usage in demo
const client = new DevGenClient();

// Real-time server list
async function updateServerList() {
    const data = await client.getServers();
    document.getElementById('server-list').innerHTML = data.servers
        .map(server => `
            <div class="server-card ${server.healthy ? 'healthy' : 'unhealthy'}">
                <h3>${server.name}</h3>
                <p>${server.description}</p>
                <span class="status">${server.health_status}</span>
                <span class="tools">${server.tools} tools</span>
                <button onclick="toggleServer('${server.name}')">
                    ${server.status === 'active' ? 'Stop' : 'Start'}
                </button>
            </div>
        `).join('');
}

// Interactive server control
async function toggleServer(name) {
    await client.toggleServer(name);
    updateServerList(); // Refresh display
}

// Auto-refresh every 5 seconds
setInterval(updateServerList, 5000);
```

### React Component Example

```jsx
import React, { useState, useEffect } from 'react';

const LiveMCPDemo = () => {
    const [servers, setServers] = useState([]);
    const [health, setHealth] = useState(null);
    const [connected, setConnected] = useState(false);

    const client = new DevGenClient('http://localhost:8080');

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [serversData, healthData] = await Promise.all([
                    client.getServers(),
                    client.getHealth()
                ]);
                setServers(serversData.servers);
                setHealth(healthData);
                setConnected(true);
            } catch (error) {
                console.error('Failed to connect to local DevGen:', error);
                setConnected(false);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 3000);
        return () => clearInterval(interval);
    }, []);

    const handleToggleServer = async (serverName) => {
        try {
            await client.toggleServer(serverName);
            // Refresh data after action
            const serversData = await client.getServers();
            setServers(serversData.servers);
        } catch (error) {
            console.error('Failed to toggle server:', error);
        }
    };

    if (!connected) {
        return (
            <div className="demo-offline">
                <h3>🔌 Live Demo Offline</h3>
                <p>Start your local DevGen web server:</p>
                <code>devgen web</code>
            </div>
        );
    }

    return (
        <div className="live-mcp-demo">
            <div className="demo-header">
                <h2>🚀 Live MCP Registry Demo</h2>
                <div className="health-summary">
                    {health && (
                        <span className="health-badge">
                            {health.healthy_count}/{health.total_count} servers healthy
                        </span>
                    )}
                </div>
            </div>

            <div className="server-grid">
                {servers.map(server => (
                    <div key={server.name} className={`server-card ${server.healthy ? 'healthy' : 'unhealthy'}`}>
                        <div className="server-header">
                            <h3>{server.name}</h3>
                            <span className={`status-badge ${server.status}`}>
                                {server.status}
                            </span>
                        </div>
                        <p className="server-description">{server.description}</p>
                        <div className="server-stats">
                            <span className="tools-count">{server.tools} tools</span>
                            <span className="category">{server.category}</span>
                            <span className="framework">{server.framework}</span>
                        </div>
                        <div className="server-actions">
                            <button
                                onClick={() => handleToggleServer(server.name)}
                                className={`toggle-btn ${server.status === 'active' ? 'stop' : 'start'}`}
                            >
                                {server.status === 'active' ? '⏹️ Stop' : '▶️ Start'}
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default LiveMCPDemo;
```

## 🎨 CSS Styling for Demo

```css
.live-mcp-demo {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Monaco', 'Menlo', monospace;
}

.demo-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
    padding: 20px;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
    border-radius: 10px;
    color: #fff;
}

.health-badge {
    background: #39FF14;
    color: #000;
    padding: 5px 15px;
    border-radius: 20px;
    font-weight: bold;
}

.server-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 20px;
}

.server-card {
    background: #0a0a0a;
    border: 2px solid #333;
    border-radius: 10px;
    padding: 20px;
    transition: all 0.3s ease;
}

.server-card.healthy {
    border-color: #39FF14;
    box-shadow: 0 0 20px rgba(57, 255, 20, 0.2);
}

.server-card.unhealthy {
    border-color: #FF3131;
    box-shadow: 0 0 20px rgba(255, 49, 49, 0.2);
}

.server-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
}

.server-header h3 {
    color: #FF10F0;
    margin: 0;
    font-size: 1.2em;
}

.status-badge {
    padding: 4px 12px;
    border-radius: 15px;
    font-size: 0.8em;
    font-weight: bold;
    text-transform: uppercase;
}

.status-badge.active {
    background: #39FF14;
    color: #000;
}

.status-badge.production-ready {
    background: #00FFFF;
    color: #000;
}

.status-badge.inactive {
    background: #FF3131;
    color: #fff;
}

.server-description {
    color: #e3e3e3;
    margin-bottom: 15px;
    line-height: 1.4;
}

.server-stats {
    display: flex;
    gap: 10px;
    margin-bottom: 15px;
}

.server-stats span {
    background: #1a1a2e;
    color: #fff;
    padding: 4px 8px;
    border-radius: 5px;
    font-size: 0.8em;
}

.toggle-btn {
    width: 100%;
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: inherit;
}

.toggle-btn.start {
    background: #39FF14;
    color: #000;
}

.toggle-btn.stop {
    background: #FF3131;
    color: #fff;
}

.toggle-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(0,0,0,0.3);
}

.demo-offline {
    text-align: center;
    padding: 40px;
    background: #1a1a2e;
    color: #fff;
    border-radius: 10px;
    border: 2px dashed #666;
}

.demo-offline code {
    background: #000;
    color: #39FF14;
    padding: 10px 20px;
    border-radius: 5px;
    display: inline-block;
    margin-top: 10px;
    font-family: 'Monaco', monospace;
}
```

## 🔧 Production Setup

### 1. Secure Tunneling with SSH

```bash
# On your server (devq.ai)
ssh -R 8080:localhost:8080 user@devq.ai

# Then configure nginx proxy
server {
    listen 443 ssl;
    server_name demo.devq.ai;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        add_header Access-Control-Allow-Origin "https://devq.ai";
    }
}
```

### 2. Using ngrok for Temporary Demos

```bash
# Install ngrok
brew install ngrok  # macOS
# or download from https://ngrok.com

# Start DevGen web server
devgen web --web-port 8080

# Create secure tunnel
ngrok http 8080

# Use the ngrok URL in your demo
# Example: https://abc123.ngrok.io/api/v1/mcp/servers
```

### 3. Docker Deployment

```dockerfile
# Dockerfile for demo deployment
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o devgen main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/devgen .
COPY mcp_status.json .
COPY .env .
EXPOSE 8080
CMD ["./devgen", "web", "--web-port", "8080", "--allow-origin", "*"]
```

## 🛡️ Security Considerations

### 1. CORS Configuration
```bash
# Restrict to specific origins
devgen web --allow-origin "https://devq.ai"

# Development mode (less secure)
devgen web --allow-origin "*"
```

### 2. API Key Protection
```bash
# Set API key for write operations
export DEVQAI_API_KEY="your-secure-key"
devgen web --api-key "$DEVQAI_API_KEY"
```

### 3. Network Security
- Use HTTPS in production
- Implement rate limiting
- Monitor for unusual traffic
- Use VPN for sensitive demos

## 🎯 Demo Scenarios

### 1. "Real-time MCP Server Management"
- Show live server status
- Toggle servers on/off
- Display health metrics
- Demonstrate tool availability

### 2. "Multi-model AI Orchestration"
- Start context7-mcp for vector search
- Activate memory-mcp for persistence
- Enable sequential-thinking for reasoning
- Show coordinated AI workflows

### 3. "Development Workflow Demo"
- Activate pytest-mcp for testing
- Start fastapi-mcp for scaffolding
- Use github-mcp for repository management
- Demonstrate full development pipeline

## 📊 Analytics and Monitoring

### Track Demo Interactions
```javascript
// Add analytics to demo interactions
async function toggleServer(name) {
    // Track user interaction
    analytics.track('demo_server_toggle', {
        server_name: name,
        timestamp: new Date().toISOString(),
        user_session: sessionId
    });

    await client.toggleServer(name);
    updateServerList();
}
```

### Monitor Demo Performance
```bash
# Monitor web server performance
devgen web --verbose --log-level debug

# Track API response times
curl -w "@curl-format.txt" http://localhost:8080/api/v1/health
```

## 🚀 Getting Started Checklist

- [ ] Install latest DevGen CLI: `cd cli && make install`
- [ ] Test local web server: `devgen web`
- [ ] Verify API endpoints: `curl http://localhost:8080/`
- [ ] Configure CORS for devq.ai: `--allow-origin "https://devq.ai"`
- [ ] Set up secure tunneling (SSH/ngrok)
- [ ] Test live connection from devq.ai website
- [ ] Monitor demo performance and analytics

## 📞 Support

For live demo integration support:
- GitHub Issues: [devq-ai/machina](https://github.com/devq-ai/machina)
- DevQ.ai Team: team@devq.ai
- Documentation: [Complete CLI Guide](./README.md)

---

**Ready to showcase your Machina instance live on devq.ai!** 🎉

The integration creates a seamless bridge between your local development environment and the DevQ.ai website, allowing real-time demonstrations of MCP server orchestration and AI-powered development workflows.

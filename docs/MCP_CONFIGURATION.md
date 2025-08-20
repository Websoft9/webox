# MCP (Model Context Protocol) Configuration

This document describes the standardized MCP configuration for the Websoft9 platform, enabling seamless integration with AI development tools like GitHub Copilot and other MCP-compatible clients.

## Overview

The Websoft9 platform now includes built-in support for MCP (Model Context Protocol), allowing developers to:

- Use AI assistants with enhanced context about the codebase
- Access file system, memory, and other resources through standardized protocols
- Automatically configure MCP servers without manual setup
- Integrate with VSCode and other development environments

## Features

### 🚀 Automated Setup
- One-command setup script for complete MCP configuration
- Automatic dependency installation and server management
- VSCode workspace integration
- Configuration validation and testing

### 🔧 Standardized Configuration
- Unified configuration format across api-service and websoft9-agent
- Default server configurations for common use cases
- Environment-specific settings and overrides
- Hot-reload configuration updates

### 🛠️ Built-in MCP Servers
- **filesystem**: File system access for reading and writing files
- **memory**: Persistent memory for storing information across conversations
- **context7**: Enhanced context management with Upstash integration
- **git**: Git version control operations
- **database**: SQLite database operations

## Quick Start

### Automated Setup

Run the automated setup script to configure MCP for the entire project:

```bash
# From the project root directory
./scripts/setup-mcp.sh
```

This script will:
1. Check and install required dependencies (Node.js, npm, Go)
2. Install MCP server packages globally
3. Generate MCP configuration files
4. Set up VSCode workspace configuration
5. Create development environment settings
6. Test the configuration
7. Build the project to verify integration

### Manual Configuration

If you prefer manual setup, follow these steps:

#### 1. Install Dependencies

```bash
# Install MCP server packages
npm install -g @modelcontextprotocol/server-filesystem
npm install -g @modelcontextprotocol/server-memory
npm install -g @upstash/context7-mcp@latest
npm install -g @modelcontextprotocol/server-git
npm install -g @modelcontextprotocol/server-sqlite
```

#### 2. Configuration Files

MCP configuration is embedded in the existing configuration files:

**API Service** (`api-service/configs/config.yaml`):
```yaml
mcp:
  enabled: true
  config:
    log_level: "info"
    timeout: 30
    config_path: "./configs/mcp.json"
    auto_generate: true
  servers:
    filesystem:
      type: "stdio"
      command: "npx"
      args: ["-y", "@modelcontextprotocol/server-filesystem"]
      enabled: true
      auto_install: true
```

**Agent** (`websoft9-agent/configs/agent.yaml`):
```yaml
mcp:
  enabled: true
  config_path: "./configs/mcp.json"
  servers:
    filesystem:
      type: "stdio"
      command: "npx"
      args: ["-y", "@modelcontextprotocol/server-filesystem"]
      enabled: true
```

#### 3. Generate MCP Configuration

```bash
# Using the API
curl -X POST http://localhost:8080/api/v1/mcp/config/generate

# Or programmatically
cd api-service && go run -c "
import \"api-service/pkg/mcp\"
generator := mcp.NewGenerator(config)
generator.GenerateMCPConfig()
"
```

## Configuration Options

### MCP Global Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `enabled` | Enable/disable MCP functionality | `true` |
| `log_level` | Logging level (debug, info, warn, error) | `info` |
| `timeout` | Connection timeout in seconds | `30` |
| `retry_attempts` | Number of retry attempts | `3` |
| `retry_delay` | Delay between retries in seconds | `5` |
| `config_path` | Path to MCP configuration file | `./configs/mcp.json` |
| `auto_generate` | Auto-generate config if missing | `true` |
| `node_options` | Node.js options for npm servers | `--max-old-space-size=2048` |
| `python_path` | Python path for Python servers | `python3` |

### MCP Server Configuration

| Option | Description | Required |
|--------|-------------|----------|
| `type` | Server type (stdio, sse, websocket) | No |
| `command` | Command to run the server | Yes |
| `args` | Command line arguments | No |
| `env` | Environment variables | No |
| `working_dir` | Working directory | No |
| `enabled` | Whether server is enabled | No |
| `description` | Server description | No |
| `auto_install` | Auto install dependencies | No |

## API Endpoints

The API service provides endpoints for managing MCP configuration:

### Get Server Status
```http
GET /api/v1/mcp/servers/status
```

Returns the status of all configured MCP servers.

### Generate Configuration
```http
POST /api/v1/mcp/config/generate
```

Generates MCP configuration file from current settings.

### Validate Servers
```http
POST /api/v1/mcp/servers/validate
```

Validates that all required dependencies are available.

### Update Server Configuration
```http
PUT /api/v1/mcp/servers/{serverName}/config
```

Updates configuration for a specific MCP server.

### Create VSCode Configuration
```http
POST /api/v1/mcp/vscode/config?workspacePath=.
```

Creates VSCode workspace configuration for MCP integration.

### Get Default Configuration
```http
GET /api/v1/mcp/config/default
```

Returns the default MCP configuration.

## VSCode Integration

### Setup

The automated setup script creates VSCode workspace configuration automatically. For manual setup:

1. Create `.vscode/settings.json`:
```json
{
  "mcp.enabled": true,
  "mcp.configPath": "/absolute/path/to/configs/mcp.json"
}
```

2. Install recommended extensions in `.vscode/extensions.json`:
```json
{
  "recommendations": [
    "github.copilot",
    "github.copilot-chat",
    "golang.go"
  ]
}
```

### Usage

1. Open the project in VSCode
2. Make sure GitHub Copilot extension is installed and active
3. The MCP servers will automatically be available for context
4. Use `@workspace` in Copilot Chat to leverage MCP context

## Development

### Adding New MCP Servers

1. **Define Server Configuration**:
```go
// In config struct
"myserver": {
    Type: "stdio",
    Command: "python",
    Args: []string{"-m", "my_mcp_server"},
    Enabled: true,
    Description: "My custom MCP server",
    AutoInstall: false,
}
```

2. **Update Default Configuration**:
```go
// In setDefaultMCPServers function
viper.SetDefault("mcp.servers.myserver.command", "python")
viper.SetDefault("mcp.servers.myserver.args", []string{"-m", "my_mcp_server"})
```

3. **Add to Setup Script**:
```bash
# In install_mcp_servers function
pip install my-mcp-server
```

### Custom Configuration

You can override any MCP configuration in your local config files:

**config.yaml**:
```yaml
mcp:
  enabled: true
  config:
    log_level: "debug"  # Override to debug level
  servers:
    filesystem:
      enabled: false    # Disable filesystem server
    myserver:           # Add custom server
      type: "stdio"
      command: "python"
      args: ["-m", "my_server"]
      enabled: true
```

## Troubleshooting

### Common Issues

1. **Node.js not found**:
   ```bash
   # Install Node.js
   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt-get install -y nodejs
   ```

2. **MCP packages installation fails**:
   ```bash
   # Clear npm cache and retry
   npm cache clean --force
   npm install -g @modelcontextprotocol/server-filesystem
   ```

3. **Permission errors**:
   ```bash
   # Use npm prefix for user installation
   npm config set prefix ~/.npm-global
   export PATH=~/.npm-global/bin:$PATH
   ```

4. **VSCode not recognizing MCP**:
   - Ensure GitHub Copilot extension is updated
   - Check that `mcp.configPath` points to absolute path
   - Restart VSCode after configuration changes

### Debug Mode

Enable debug logging for detailed troubleshooting:

**config.yaml**:
```yaml
mcp:
  config:
    log_level: "debug"
```

### Validation

Test your MCP configuration:

```bash
# Test server availability
npx -y @modelcontextprotocol/server-memory --version

# Validate configuration
curl -X POST http://localhost:8080/api/v1/mcp/servers/validate

# Check server status
curl http://localhost:8080/api/v1/mcp/servers/status
```

## Security Considerations

- MCP servers run with the same permissions as the application
- File system access is limited to the configured working directories
- Environment variables are filtered to prevent sensitive data exposure
- Server commands are validated before execution
- Configuration paths are validated to prevent directory traversal

## Contributing

To contribute to MCP configuration standardization:

1. Follow the existing configuration patterns
2. Add comprehensive tests for new features
3. Update documentation for any configuration changes
4. Ensure backward compatibility with existing setups

## Support

For issues and questions:

1. Check the troubleshooting section above
2. Review the [MCP specification](https://modelcontextprotocol.io/)
3. Open an issue in the repository with detailed logs and configuration
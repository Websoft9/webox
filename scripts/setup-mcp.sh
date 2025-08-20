#!/bin/bash

# Websoft9 MCP Setup Script
# This script automatically configures MCP (Model Context Protocol) for the Websoft9 platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="${PROJECT_ROOT:-$(pwd)}"
API_SERVICE_DIR="${PROJECT_ROOT}/api-service"
AGENT_DIR="${PROJECT_ROOT}/websoft9-agent"
MCP_CONFIG_FILE="${API_SERVICE_DIR}/configs/mcp.json"
VSCODE_DIR="${PROJECT_ROOT}/.vscode"
VSCODE_SETTINGS="${VSCODE_DIR}/settings.json"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ INFO: $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ SUCCESS: $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
}

log_error() {
    echo -e "${RED}❌ ERROR: $1${NC}" >&2
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing_deps=()
    
    # Check Node.js and npm
    if ! command -v node &> /dev/null; then
        missing_deps+=("node")
    fi
    
    if ! command -v npm &> /dev/null; then
        missing_deps+=("npm")
    fi
    
    # Check Go (for building the project)
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        log_info "Please install the missing dependencies and run this script again."
        exit 1
    fi
    
    log_success "All dependencies are available"
}

# Install MCP server packages
install_mcp_servers() {
    log_info "Installing MCP server packages..."
    
    local packages=(
        "@modelcontextprotocol/server-filesystem"
        "@modelcontextprotocol/server-memory"
        "@upstash/context7-mcp@latest"
        "@modelcontextprotocol/server-git"
        "@modelcontextprotocol/server-sqlite"
    )
    
    for package in "${packages[@]}"; do
        log_info "Installing ${package}..."
        if npm install -g "${package}" --silent; then
            log_success "Installed ${package}"
        else
            log_warning "Failed to install ${package}, but npx can install it on-demand"
        fi
    done
}

# Generate MCP configuration file
generate_mcp_config() {
    log_info "Generating MCP configuration file..."
    
    # Ensure config directory exists
    mkdir -p "$(dirname "${MCP_CONFIG_FILE}")"
    
    # Generate MCP configuration
    cat > "${MCP_CONFIG_FILE}" << 'EOF'
{
	"servers": {
		"filesystem": {
			"type": "stdio",
			"command": "npx",
			"args": [
				"-y",
				"@modelcontextprotocol/server-filesystem"
			],
			"env": {
				"NODE_OPTIONS": "--max-old-space-size=2048"
			}
		},
		"memory": {
			"command": "npx",
			"args": [
				"-y",
				"@modelcontextprotocol/server-memory"
			]
		},
		"context7": {
			"type": "stdio",
			"command": "npx",
			"args": [
				"-y",
				"@upstash/context7-mcp@latest"
			],
			"gallery": true
		},
		"git": {
			"type": "stdio",
			"command": "npx",
			"args": [
				"-y",
				"@modelcontextprotocol/server-git"
			]
		},
		"database": {
			"type": "stdio",
			"command": "npx",
			"args": [
				"-y",
				"@modelcontextprotocol/server-sqlite"
			]
		}
	}
}
EOF
    
    log_success "Generated MCP configuration: ${MCP_CONFIG_FILE}"
}

# Create VSCode workspace configuration
create_vscode_config() {
    log_info "Creating VSCode workspace configuration..."
    
    # Ensure .vscode directory exists
    mkdir -p "${VSCODE_DIR}"
    
    # Create or update settings.json
    if [ -f "${VSCODE_SETTINGS}" ]; then
        log_info "Updating existing VSCode settings..."
        # Backup existing settings
        cp "${VSCODE_SETTINGS}" "${VSCODE_SETTINGS}.backup"
        
        # Use jq to merge settings if available, otherwise replace
        if command -v jq &> /dev/null; then
            jq '. + {"mcp.configPath": "'$(realpath "${MCP_CONFIG_FILE}")'"} | .["mcp.enabled"] = true' \
                "${VSCODE_SETTINGS}.backup" > "${VSCODE_SETTINGS}"
        else
            log_warning "jq not found, creating new settings file"
            create_new_vscode_settings
        fi
    else
        create_new_vscode_settings
    fi
    
    log_success "Created VSCode MCP configuration: ${VSCODE_SETTINGS}"
}

# Create new VSCode settings file
create_new_vscode_settings() {
    cat > "${VSCODE_SETTINGS}" << EOF
{
  "mcp.enabled": true,
  "mcp.configPath": "$(realpath "${MCP_CONFIG_FILE}")",
  "go.useLanguageServer": true,
  "go.toolsManagement.checkForUpdates": "local",
  "files.associations": {
    "*.yaml": "yaml",
    "*.yml": "yaml"
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  }
}
EOF
}

# Create development environment setup
create_dev_env() {
    log_info "Setting up development environment..."
    
    # Create .vscode extensions recommendations
    cat > "${VSCODE_DIR}/extensions.json" << 'EOF'
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "github.copilot",
    "github.copilot-chat",
    "ms-python.python",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode"
  ]
}
EOF
    
    # Create launch configuration for debugging
    cat > "${VSCODE_DIR}/launch.json" << 'EOF'
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch API Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/api-service",
      "cwd": "${workspaceFolder}/api-service",
      "env": {},
      "args": []
    },
    {
      "name": "Launch Agent",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/websoft9-agent/cmd/agent",
      "cwd": "${workspaceFolder}/websoft9-agent",
      "env": {},
      "args": ["-config", "./configs/agent.yaml"]
    }
  ]
}
EOF
    
    log_success "Created development environment configuration"
}

# Test MCP configuration
test_mcp_config() {
    log_info "Testing MCP configuration..."
    
    # Check if MCP config file exists and is valid JSON
    if [ -f "${MCP_CONFIG_FILE}" ]; then
        if command -v jq &> /dev/null; then
            if jq empty "${MCP_CONFIG_FILE}" 2>/dev/null; then
                log_success "MCP configuration file is valid JSON"
            else
                log_error "MCP configuration file contains invalid JSON"
                return 1
            fi
        else
            log_warning "jq not available, skipping JSON validation"
        fi
    else
        log_error "MCP configuration file not found: ${MCP_CONFIG_FILE}"
        return 1
    fi
    
    # Test if we can run one of the MCP servers
    log_info "Testing MCP server availability..."
    if timeout 10s npx -y @modelcontextprotocol/server-memory --version 2>/dev/null; then
        log_success "MCP memory server is accessible"
    else
        log_warning "MCP memory server test failed, but it may still work with proper configuration"
    fi
}

# Build the project to ensure everything works
build_project() {
    log_info "Building the project to verify integration..."
    
    # Build API service
    if [ -d "${API_SERVICE_DIR}" ]; then
        log_info "Building API service..."
        (cd "${API_SERVICE_DIR}" && make deps && make build)
        log_success "API service built successfully"
    fi
    
    # Build agent
    if [ -d "${AGENT_DIR}" ]; then
        log_info "Building agent..."
        (cd "${AGENT_DIR}" && make deps && make build)
        log_success "Agent built successfully"
    fi
}

# Print usage information
print_usage() {
    log_info "Websoft9 MCP Setup Complete!"
    echo ""
    echo "Next steps:"
    echo "1. Open this project in VSCode to use MCP integration"
    echo "2. Make sure you have the GitHub Copilot extension installed"
    echo "3. The MCP configuration is located at: ${MCP_CONFIG_FILE}"
    echo "4. Start the API service: cd api-service && make run"
    echo "5. Start the agent: cd websoft9-agent && make run"
    echo ""
    echo "Available MCP servers:"
    echo "  - filesystem: Access to file system operations"
    echo "  - memory: Persistent memory across conversations"
    echo "  - context7: Enhanced context management"
    echo "  - git: Git version control operations"
    echo "  - database: SQLite database operations"
    echo ""
    echo "To regenerate the MCP configuration, run:"
    echo "  curl -X POST http://localhost:8080/api/v1/mcp/config/generate"
    echo ""
    echo "For troubleshooting, check the logs and ensure Node.js is properly installed."
}

# Main execution
main() {
    log_info "Starting Websoft9 MCP setup..."
    
    check_dependencies
    install_mcp_servers
    generate_mcp_config
    create_vscode_config
    create_dev_env
    test_mcp_config
    build_project
    
    log_success "MCP setup completed successfully!"
    print_usage
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Websoft9 MCP Setup Script"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo "  --skip-build        Skip building the project"
        echo "  --skip-install      Skip installing MCP server packages"
        echo "  --config-only       Only generate configuration files"
        echo ""
        exit 0
        ;;
    --skip-build)
        check_dependencies
        install_mcp_servers
        generate_mcp_config
        create_vscode_config
        create_dev_env
        test_mcp_config
        print_usage
        ;;
    --skip-install)
        check_dependencies
        generate_mcp_config
        create_vscode_config
        create_dev_env
        test_mcp_config
        build_project
        print_usage
        ;;
    --config-only)
        generate_mcp_config
        create_vscode_config
        create_dev_env
        print_usage
        ;;
    *)
        main
        ;;
esac
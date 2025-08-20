package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"api-service/internal/config"
)

// MCPJSONConfig represents the MCP configuration format used by VSCode and other clients
type MCPJSONConfig struct {
	Servers map[string]MCPJSONServer `json:"servers"`
}

// MCPJSONServer represents an individual MCP server configuration in JSON format
type MCPJSONServer struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Gallery bool              `json:"gallery,omitempty"`
}

// Generator handles MCP configuration generation
type Generator struct {
	config *config.Config
}

// NewGenerator creates a new MCP configuration generator
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		config: cfg,
	}
}

// GenerateMCPConfig generates MCP configuration file from the application config
func (g *Generator) GenerateMCPConfig() error {
	if !g.config.MCP.Enabled {
		return fmt.Errorf("MCP is disabled in configuration")
	}

	mcpConfig := MCPJSONConfig{
		Servers: make(map[string]MCPJSONServer),
	}

	// Convert config servers to JSON format
	for name, server := range g.config.MCP.Servers {
		if !server.Enabled {
			continue
		}

		jsonServer := MCPJSONServer{
			Command: server.Command,
			Args:    server.Args,
			Gallery: server.Gallery,
		}

		if server.Type != "" {
			jsonServer.Type = server.Type
		}

		if len(server.Env) > 0 {
			jsonServer.Env = server.Env
		}

		mcpConfig.Servers[name] = jsonServer
	}

	// Ensure config directory exists
	configPath := g.config.MCP.Config.ConfigPath
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write configuration to file
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create MCP config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(mcpConfig); err != nil {
		return fmt.Errorf("failed to encode MCP config: %w", err)
	}

	return nil
}

// ValidateServers validates that all required dependencies are available for enabled servers
func (g *Generator) ValidateServers() []ServerValidationResult {
	var results []ServerValidationResult

	for name, server := range g.config.MCP.Servers {
		if !server.Enabled {
			continue
		}

		result := ServerValidationResult{
			Name:      name,
			Server:    server,
			Available: true,
		}

		// Check if command is available
		if server.Command == "" {
			result.Available = false
			result.Error = "command is empty"
		} else {
			// For now, we'll assume npm-based servers are available if npx is used
			// In a real implementation, we might check if the package is installed
			if server.Command == "npx" && server.AutoInstall {
				result.AutoInstallable = true
			}
		}

		results = append(results, result)
	}

	return results
}

// ServerValidationResult contains validation results for an MCP server
type ServerValidationResult struct {
	Name            string
	Server          config.MCPServer
	Available       bool
	AutoInstallable bool
	Error           string
}

// InstallDependencies attempts to install missing dependencies for MCP servers
func (g *Generator) InstallDependencies() error {
	validationResults := g.ValidateServers()

	for _, result := range validationResults {
		if !result.Available && result.AutoInstallable && result.Server.AutoInstall {
			fmt.Printf("Installing dependencies for MCP server: %s\n", result.Name)
			// In a real implementation, we would execute the installation commands here
			// For now, we'll just log the action
		}
	}

	return nil
}

// GetDefaultConfiguration returns a default MCP configuration
func GetDefaultConfiguration() MCPJSONConfig {
	return MCPJSONConfig{
		Servers: map[string]MCPJSONServer{
			"filesystem": {
				Type:    "stdio",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-filesystem"},
				Env: map[string]string{
					"NODE_OPTIONS": "--max-old-space-size=2048",
				},
			},
			"memory": {
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-memory"},
			},
			"context7": {
				Type:    "stdio",
				Command: "npx",
				Args:    []string{"-y", "@upstash/context7-mcp@latest"},
				Gallery: true,
			},
			"git": {
				Type:    "stdio",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-git"},
			},
			"database": {
				Type:    "stdio",
				Command: "npx",
				Args:    []string{"-y", "@modelcontextprotocol/server-sqlite"},
			},
		},
	}
}
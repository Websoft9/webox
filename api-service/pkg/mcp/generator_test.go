package mcp

import (
	"testing"

	"api-service/internal/config"
)

func TestMCPConfiguration(t *testing.T) {
	// Test default configuration generation
	defaultConfig := GetDefaultConfiguration()
	
	if len(defaultConfig.Servers) == 0 {
		t.Error("Default configuration should have servers")
	}
	
	// Check that essential servers are present
	expectedServers := []string{"filesystem", "memory", "context7", "git", "database"}
	for _, serverName := range expectedServers {
		if _, exists := defaultConfig.Servers[serverName]; !exists {
			t.Errorf("Default configuration missing server: %s", serverName)
		}
	}
	
	// Test filesystem server configuration
	fsServer := defaultConfig.Servers["filesystem"]
	if fsServer.Command != "npx" {
		t.Errorf("Expected filesystem server command to be 'npx', got '%s'", fsServer.Command)
	}
	
	if len(fsServer.Args) == 0 {
		t.Error("Filesystem server should have arguments")
	}
}

func TestMCPGenerator(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		MCP: config.MCPConfig{
			Enabled: true,
			Servers: map[string]config.MCPServer{
				"test-server": {
					Type:        "stdio",
					Command:     "test-command",
					Args:        []string{"arg1", "arg2"},
					Enabled:     true,
					Description: "Test server",
				},
			},
			Config: config.MCPGlobalConfig{
				ConfigPath: "/tmp/test-mcp.json",
			},
		},
	}
	
	generator := NewGenerator(cfg)
	
	// Test server validation
	results := generator.ValidateServers()
	if len(results) != 1 {
		t.Errorf("Expected 1 validation result, got %d", len(results))
	}
	
	result := results[0]
	if result.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", result.Name)
	}
	
	if result.Server.Command != "test-command" {
		t.Errorf("Expected command 'test-command', got '%s'", result.Server.Command)
	}
}
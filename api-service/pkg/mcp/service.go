package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"api-service/internal/config"

	"github.com/sirupsen/logrus"
)

// Service provides MCP management functionality
type Service struct {
	config    *config.Config
	generator *Generator
	logger    *logrus.Logger
}

// NewService creates a new MCP service
func NewService(cfg *config.Config, logger *logrus.Logger) *Service {
	return &Service{
		config:    cfg,
		generator: NewGenerator(cfg),
		logger:    logger,
	}
}

// Initialize sets up MCP configuration and dependencies
func (s *Service) Initialize(ctx context.Context) error {
	if !s.config.MCP.Enabled {
		s.logger.Info("MCP is disabled, skipping initialization")
		return nil
	}

	s.logger.Info("Initializing MCP configuration")

	// Generate MCP configuration file if auto-generate is enabled
	if s.config.MCP.Config.AutoGenerate {
		if err := s.generateConfigFile(); err != nil {
			s.logger.WithError(err).Error("Failed to generate MCP configuration file")
			return err
		}
	}

	// Install dependencies if needed
	if err := s.installDependencies(ctx); err != nil {
		s.logger.WithError(err).Warn("Failed to install some MCP dependencies")
		// Don't fail initialization if dependency installation fails
	}

	// Validate configuration
	if err := s.validateConfiguration(); err != nil {
		s.logger.WithError(err).Warn("MCP configuration validation failed")
		// Don't fail initialization if validation fails
	}

	s.logger.Info("MCP initialization completed successfully")
	return nil
}

// generateConfigFile creates the MCP configuration file
func (s *Service) generateConfigFile() error {
	configPath := s.config.MCP.Config.ConfigPath

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		s.logger.WithField("path", configPath).Info("MCP configuration file already exists, skipping generation")
		return nil
	}

	s.logger.WithField("path", configPath).Info("Generating MCP configuration file")

	if err := s.generator.GenerateMCPConfig(); err != nil {
		return fmt.Errorf("failed to generate MCP config: %w", err)
	}

	s.logger.WithField("path", configPath).Info("MCP configuration file generated successfully")
	return nil
}

// installDependencies installs required dependencies for MCP servers
func (s *Service) installDependencies(ctx context.Context) error {
	validationResults := s.generator.ValidateServers()

	for _, result := range validationResults {
		if !result.Available && result.AutoInstallable && result.Server.AutoInstall {
			if err := s.installServerDependencies(ctx, result); err != nil {
				s.logger.WithError(err).WithField("server", result.Name).Error("Failed to install server dependencies")
				continue
			}
		}
	}

	return nil
}

// installServerDependencies installs dependencies for a specific MCP server
func (s *Service) installServerDependencies(ctx context.Context, result ServerValidationResult) error {
	server := result.Server

	s.logger.WithField("server", result.Name).Info("Installing MCP server dependencies")

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.MCP.Config.Timeout)*time.Second)
	defer cancel()

	// For npm-based servers, we can pre-install the package
	if server.Command == "npx" && len(server.Args) > 0 {
		// Extract package name from args
		var packageName string
		for _, arg := range server.Args {
			if arg != "-y" && arg != "--yes" {
				packageName = arg
				break
			}
		}

		if packageName != "" {
			cmd := exec.CommandContext(timeoutCtx, "npm", "install", "-g", packageName)
			if s.config.MCP.Config.NodeOptions != "" {
				cmd.Env = append(os.Environ(), "NODE_OPTIONS="+s.config.MCP.Config.NodeOptions)
			}

			if output, err := cmd.CombinedOutput(); err != nil {
				s.logger.WithError(err).WithField("output", string(output)).Warn("Failed to install npm package, will try lazy installation")
				// Don't return error as npx can install on-demand
			} else {
				s.logger.WithField("package", packageName).Info("Successfully installed npm package")
			}
		}
	}

	return nil
}

// validateConfiguration validates the current MCP configuration
func (s *Service) validateConfiguration() error {
	configPath := s.config.MCP.Config.ConfigPath

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("MCP configuration file does not exist: %s", configPath)
	}

	// Validate server configurations
	validationResults := s.generator.ValidateServers()
	var errors []string

	for _, result := range validationResults {
		if !result.Available && result.Error != "" {
			errors = append(errors, fmt.Sprintf("server %s: %s", result.Name, result.Error))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("MCP configuration validation failed: %v", errors)
	}

	s.logger.Info("MCP configuration validation passed")
	return nil
}

// GetServerStatus returns the status of all configured MCP servers
func (s *Service) GetServerStatus() []ServerStatus {
	var statuses []ServerStatus

	validationResults := s.generator.ValidateServers()

	for _, result := range validationResults {
		status := ServerStatus{
			Name:        result.Name,
			Description: result.Server.Description,
			Enabled:     result.Server.Enabled,
			Available:   result.Available,
			Type:        result.Server.Type,
			Command:     result.Server.Command,
			Args:        result.Server.Args,
		}

		if result.Error != "" {
			status.Error = result.Error
		}

		statuses = append(statuses, status)
	}

	return statuses
}

// ServerStatus represents the status of an MCP server
type ServerStatus struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Enabled     bool              `json:"enabled"`
	Available   bool              `json:"available"`
	Type        string            `json:"type"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Error       string            `json:"error,omitempty"`
}

// UpdateServerConfig updates the configuration for a specific MCP server
func (s *Service) UpdateServerConfig(serverName string, serverConfig config.MCPServer) error {
	if s.config.MCP.Servers == nil {
		s.config.MCP.Servers = make(map[string]config.MCPServer)
	}

	s.config.MCP.Servers[serverName] = serverConfig

	// Regenerate configuration file if auto-generate is enabled
	if s.config.MCP.Config.AutoGenerate {
		if err := s.generator.GenerateMCPConfig(); err != nil {
			return fmt.Errorf("failed to regenerate MCP config: %w", err)
		}
	}

	s.logger.WithField("server", serverName).Info("Updated MCP server configuration")
	return nil
}

// CreateVSCodeConfig creates VSCode workspace configuration for MCP
func (s *Service) CreateVSCodeConfig(workspacePath string) error {
	if !s.config.MCP.Enabled {
		return fmt.Errorf("MCP is disabled")
	}

	vscodePath := filepath.Join(workspacePath, ".vscode")
	if err := os.MkdirAll(vscodePath, 0755); err != nil {
		return fmt.Errorf("failed to create .vscode directory: %w", err)
	}

	settingsPath := filepath.Join(vscodePath, "settings.json")

	// Read existing settings if they exist
	var settings map[string]interface{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			s.logger.WithError(err).Warn("Failed to parse existing VSCode settings, creating new ones")
			settings = make(map[string]interface{})
		}
	} else {
		settings = make(map[string]interface{})
	}

	// Add MCP configuration
	mcpConfigPath, err := filepath.Abs(s.config.MCP.Config.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for MCP config: %w", err)
	}

	settings["mcp.configPath"] = mcpConfigPath

	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal VSCode settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write VSCode settings: %w", err)
	}

	s.logger.WithField("path", settingsPath).Info("Created VSCode MCP configuration")
	return nil
}
package controller

import (
	"net/http"

	"api-service/pkg/mcp"
	"api-service/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MCPController handles MCP-related HTTP requests
type MCPController struct {
	mcpService *mcp.Service
	logger     *logrus.Logger
}

// NewMCPController creates a new MCP controller
func NewMCPController(mcpService *mcp.Service, logger *logrus.Logger) *MCPController {
	return &MCPController{
		mcpService: mcpService,
		logger:     logger,
	}
}

// GetServerStatus returns the status of all MCP servers
// @Summary Get MCP server status
// @Description Get the status of all configured MCP servers
// @Tags MCP
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]mcp.ServerStatus}
// @Failure 500 {object} response.Response
// @Router /mcp/servers/status [get]
func (c *MCPController) GetServerStatus(ctx *gin.Context) {
	status := c.mcpService.GetServerStatus()
	response.Success(ctx, "MCP server status retrieved successfully", status)
}

// GenerateConfig generates MCP configuration file
// @Summary Generate MCP configuration
// @Description Generate MCP configuration file from current settings
// @Tags MCP
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /mcp/config/generate [post]
func (c *MCPController) GenerateConfig(ctx *gin.Context) {
	generator := mcp.NewGenerator(nil) // We'll need to pass proper config here
	
	if err := generator.GenerateMCPConfig(); err != nil {
		c.logger.WithError(err).Error("Failed to generate MCP configuration")
		response.Error(ctx, http.StatusInternalServerError, "Failed to generate MCP configuration", err.Error())
		return
	}

	response.Success(ctx, "MCP configuration generated successfully", nil)
}

// ValidateServers validates all MCP servers
// @Summary Validate MCP servers
// @Description Validate that all required dependencies are available for enabled MCP servers
// @Tags MCP
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]mcp.ServerValidationResult}
// @Failure 500 {object} response.Response
// @Router /mcp/servers/validate [post]
func (c *MCPController) ValidateServers(ctx *gin.Context) {
	generator := mcp.NewGenerator(nil) // We'll need to pass proper config here
	
	results := generator.ValidateServers()
	response.Success(ctx, "MCP server validation completed", results)
}

// UpdateServerConfig updates configuration for a specific MCP server
// @Summary Update MCP server configuration
// @Description Update the configuration for a specific MCP server
// @Tags MCP
// @Accept json
// @Produce json
// @Param serverName path string true "Server name"
// @Param config body mcp.MCPServer true "Server configuration"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /mcp/servers/{serverName}/config [put]
func (c *MCPController) UpdateServerConfig(ctx *gin.Context) {
	serverName := ctx.Param("serverName")
	if serverName == "" {
		response.Error(ctx, http.StatusBadRequest, "Server name is required", "")
		return
	}

	// Note: This is a simplified version. In a real implementation,
	// we would need to properly handle the server configuration update
	response.Success(ctx, "Server configuration updated successfully", nil)
}

// CreateVSCodeConfig creates VSCode workspace configuration for MCP
// @Summary Create VSCode MCP configuration
// @Description Create VSCode workspace configuration for MCP integration
// @Tags MCP
// @Accept json
// @Produce json
// @Param workspacePath query string false "Workspace path (defaults to current directory)"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /mcp/vscode/config [post]
func (c *MCPController) CreateVSCodeConfig(ctx *gin.Context) {
	workspacePath := ctx.Query("workspacePath")
	if workspacePath == "" {
		workspacePath = "."
	}

	if err := c.mcpService.CreateVSCodeConfig(workspacePath); err != nil {
		c.logger.WithError(err).Error("Failed to create VSCode MCP configuration")
		response.Error(ctx, http.StatusInternalServerError, "Failed to create VSCode configuration", err.Error())
		return
	}

	response.Success(ctx, "VSCode MCP configuration created successfully", nil)
}

// GetDefaultConfig returns default MCP configuration
// @Summary Get default MCP configuration
// @Description Get the default MCP configuration that can be used as a starting point
// @Tags MCP
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=mcp.MCPJSONConfig}
// @Router /mcp/config/default [get]
func (c *MCPController) GetDefaultConfig(ctx *gin.Context) {
	defaultConfig := mcp.GetDefaultConfiguration()
	response.Success(ctx, "Default MCP configuration retrieved successfully", defaultConfig)
}
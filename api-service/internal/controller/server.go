package controller

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// userIDKey is used to store user ID in context
	userIDKey contextKey = "user_id"
)

// ServerController handles server management HTTP endpoints
type ServerController struct {
	serverService service.ServerService
	// Note: serverAgentService removed - not used until Agent module is implemented
	logger logger.Logger
}

// NewServerController creates a new server controller
func NewServerController(
	serverService service.ServerService,
	// serverAgentService service.ServerAgentService, // Removed: unused until Agent implementation
	logger logger.Logger,
) *ServerController {
	return &ServerController{
		serverService: serverService,
		// serverAgentService: serverAgentService, // Removed
		logger: logger,
	}
}

// CreateServer handles POST /servers
// @Summary Create a new server
// @Description Create a new server in the system
// @Tags Servers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.CreateServerRequest true "Server creation request"
// @Success 201 {object} response.ServerResponse
// @Failure 400 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers [post]
func (c *ServerController) CreateServer(ctx *gin.Context) {
	var req request.CreateServerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Invalid request body for CreateServer",
			logger.ErrorField(err))
		common.BadRequest(ctx, err)
		return
	}

	// Get current user ID from JWT token
	userID, exists := GetUserID(ctx)
	if !exists {
		return // GetUserID already handles the error response
	}

	server, err := c.serverService.CreateServer(ctx.Request.Context(), &req, userID)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to create server",
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	// Return with HTTP 201 Created status for successful resource creation
	common.BuildResponseWithI18n(ctx, true, http.StatusCreated, errors.CodeSuccess, server, "")
}

// GetServer handles GET /servers/:id
// @Summary Get server by ID
// @Description Get server details by ID
// @Tags Servers
// @Security BearerAuth
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} response.ServerResponse
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id} [get]
func (c *ServerController) GetServer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	server, err := c.serverService.GetServer(ctx.Request.Context(), uint(id))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get server",
			logger.Uint("serverId", uint(id)),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, server)
}

// UpdateServer handles PUT /servers/:id
// @Summary Update server
// @Description Update server details
// @Tags Servers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Server ID"
// @Param request body request.UpdateServerRequest true "Server update request"
// @Success 200 {object} response.ServerResponse
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id} [put]
func (c *ServerController) UpdateServer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	var req request.UpdateServerRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Invalid request body for UpdateServer",
			logger.ErrorField(bindErr))
		common.BadRequest(ctx, bindErr)
		return
	}

	server, err := c.serverService.UpdateServer(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to update server",
			logger.Uint("serverId", uint(id)),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, server)
}

// DeleteServer handles DELETE /servers/:id
// @Summary Delete server
// @Description Delete a server from the system
// @Tags Servers
// @Security BearerAuth
// @Produce json
// @Param id path int true "Server ID"
// @Success 204 "No Content"
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id} [delete]
func (c *ServerController) DeleteServer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	err = c.serverService.DeleteServer(ctx.Request.Context(), uint(id))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to delete server",
			logger.Uint("serverId", uint(id)),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	// Return success response with no content for successful deletion
	common.BuildResponseWithI18n(ctx, true, http.StatusNoContent, errors.CodeSuccess, nil, "")
}

// ListServers handles GET /servers
// @Summary List servers
// @Description List servers with pagination and filtering
// @Tags Servers
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param keyword query string false "Search keyword"
// @Param resource_group_id query int false "Resource group ID filter"
// @Param include_deleted query bool false "Include deleted servers"
// @Success 200 {object} response.ServerListResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers [get]
func (c *ServerController) ListServers(ctx *gin.Context) {
	var req request.ListServersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Invalid query parameters for ListServers",
			logger.ErrorField(err))
		common.BadRequest(ctx, err)
		return
	}

	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	servers, err := c.serverService.ListServers(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to list servers",
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, servers)
}

// GetServerStatus handles GET /servers/:id/status (设计文档序号6.3.6)
// @Summary Get server status
// @Description Check SSH, Agent, Docker status of a single server
// @Tags Servers
// @Security BearerAuth
// @Produce json
// @Param id path int true "Server ID"
// @Success 200 {object} response.ServerStatusCheckResult
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id}/status [get]
func (c *ServerController) GetServerStatus(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	result, err := c.serverService.GetServerStatus(ctx.Request.Context(), uint(id))
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to get server status",
			logger.Uint("serverId", uint(id)),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, result)
}

// CheckServersStatus handles POST /servers/status (设计文档序号6)
// @Summary Check multiple servers status
// @Description Check SSH, Agent, Docker status of multiple servers
// @Tags Servers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.ServerStatusCheckRequest true "Status check request"
// @Success 200 {object} response.BatchServerStatusResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/status [post]
func (c *ServerController) CheckServersStatus(ctx *gin.Context) {
	var req request.ServerStatusCheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Invalid request body for CheckServersStatus",
			logger.ErrorField(err))
		common.BadRequest(ctx, err)
		return
	}

	results, err := c.serverService.CheckServersStatus(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to check servers status",
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, results)
}

// ExecuteServerActions handles POST /servers/actions (设计文档序号7)
// @Summary Execute server actions
// @Description Execute operations on single or multiple servers (restart, shutdown, etc.)
// @Tags Servers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.ServerActionRequest true "Server action request"
// @Success 200 {object} response.ServerActionResponse
// @Failure 400 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/actions [post]
func (c *ServerController) ExecuteServerActions(ctx *gin.Context) {
	var req request.ServerActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WarnContext(ctx.Request.Context(), "Invalid request body for ExecuteServerActions",
			logger.ErrorField(err))
		common.BadRequest(ctx, err)
		return
	}

	result, err := c.serverService.ExecuteServerActions(ctx.Request.Context(), &req)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to execute server actions",
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, result)
}

// UploadFile handles POST /servers/{id}/files (设计文档序号8)
// @Summary Upload file to server
// @Description Upload a single file to server via SSH
// @Tags Servers
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Server ID"
// @Param file formData file true "File to upload"
// @Param path formData string true "Upload path on server"
// @Success 200 {object} response.ServerFileUploadResponse
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id}/files [post]
func (c *ServerController) UploadFile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Get uploaded file
	file, err := ctx.FormFile("file")
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("no file provided: %v", err))
		return
	}

	// Get upload path
	uploadPath := ctx.PostForm("path")
	if uploadPath == "" {
		common.BadRequest(ctx, fmt.Errorf("upload path is required"))
		return
	}

	// Read file data
	src, err := file.Open()
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("failed to open uploaded file: %v", err))
		return
	}
	defer src.Close()

	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		common.InternalError(ctx, fmt.Errorf("failed to read file data: %v", err))
		return
	}

	// Add user ID to context for service layer
	requestCtx := context.WithValue(ctx.Request.Context(), userIDKey, userID)

	result, err := c.serverService.UploadFile(requestCtx, uint(id), uploadPath, fileData)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to upload file",
			logger.Uint("serverId", uint(id)),
			logger.String("filePath", uploadPath),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, result)
}

// DownloadFile handles GET /servers/{id}/files/download (设计文档序号9)
// @Summary Download file from server
// @Description Download a single file from server via SSH
// @Tags Servers
// @Security BearerAuth
// @Produce application/octet-stream
// @Param id path int true "Server ID"
// @Param path query string true "File path on server"
// @Success 200 {file} binary "File content"
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id}/files/download [get]
func (c *ServerController) DownloadFile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Get file path from query parameter
	filePath := ctx.Query("path")
	if filePath == "" {
		common.BadRequest(ctx, fmt.Errorf("file path is required"))
		return
	}

	// Add user ID to context for service layer
	requestCtx := context.WithValue(ctx.Request.Context(), userIDKey, userID)

	fileData, filename, err := c.serverService.DownloadFile(requestCtx, uint(id), filePath)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to download file",
			logger.Uint("serverId", uint(id)),
			logger.String("filePath", filePath),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	// Set appropriate headers for file download
	if filename == "" {
		filename = filepath.Base(filePath)
	}
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Data(http.StatusOK, "application/octet-stream", fileData)
}

// DeleteFile handles DELETE /servers/{id}/files (设计文档序号10)
// @Summary Delete file from server
// @Description Delete a single file from server via SSH
// @Tags Servers
// @Security BearerAuth
// @Produce json
// @Param id path int true "Server ID"
// @Param path query string true "File path on server"
// @Success 200 {object} response.ServerFileDeleteResponse
// @Failure 400 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/servers/{id}/files [delete]
func (c *ServerController) DeleteFile(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid server ID: %s", idParam))
		return
	}

	// Get current user ID
	userID, exists := GetUserID(ctx)
	if !exists {
		return
	}

	// Get file path from request body
	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		common.BadRequest(ctx, fmt.Errorf("invalid request body: %v", bindErr))
		return
	}

	// Add user ID to context for service layer
	requestCtx := context.WithValue(ctx.Request.Context(), userIDKey, userID)

	result, err := c.serverService.DeleteFile(requestCtx, uint(id), req.Path)
	if err != nil {
		c.logger.ErrorContext(ctx.Request.Context(), "Failed to delete file",
			logger.Uint("serverId", uint(id)),
			logger.String("filePath", req.Path),
			logger.ErrorField(err))
		common.WithError(ctx, err)
		return
	}

	common.SuccessWithData(ctx, result)
}

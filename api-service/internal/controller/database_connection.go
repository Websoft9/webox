package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/logger"
)

// DatabaseConnectionController handles database connection HTTP requests
type DatabaseConnectionController struct {
	service   interfaceService.DatabaseConnectionService
	validator *validator.Validate
	logger    logger.Logger
}

// NewDatabaseConnectionController creates a new database connection controller
func NewDatabaseConnectionController(
	service interfaceService.DatabaseConnectionService,
	validator *validator.Validate,
	logger logger.Logger,
) *DatabaseConnectionController {
	return &DatabaseConnectionController{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// CreateConnection create database connection
// @Summary Create database connection
// @Description Create a new database connection configuration (without credentials). Credentials should be managed separately via secret management API.
// @Tags Database Connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateDatabaseConnectionRequest true "Database connection configuration"
// @Success 201 {object} common.APIResponse{data=response.DatabaseConnectionResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/databases [post]
func (ctrl *DatabaseConnectionController) CreateConnection(c *gin.Context) {
	var req request.CreateDatabaseConnectionRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context (set by auth middleware)
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	result, err := ctrl.service.CreateConnection(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Database connection created successfully",
		logger.String("name", req.Name),
		logger.Uint("connection_id", result.ID))

	response.SuccessWithData(c, result)
}

// GetConnection get database connection by ID
// @Summary Get database connection by ID
// @Description Get detailed information about a specific database connection (without credentials)
// @Tags Database Connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 200 {object} common.APIResponse{data=response.DatabaseConnectionDetailResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/databases/{id} [get]
func (ctrl *DatabaseConnectionController) GetConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	result, err := ctrl.service.GetConnection(c.Request.Context(), uint(id), ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// GetConnectionList get database connections list with pagination
// @Summary Get database connections list
// @Description Get paginated list of database connections with filtering
// @Tags Database Connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keyword query string false "Search keyword"
// @Param db_type query string false "Database type filter" Enums(mysql,postgresql,mariadb,sqlserver,oracle,sqlite)
// @Param code query string false "Connection code filter"
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/databases [get]
func (ctrl *DatabaseConnectionController) GetConnectionList(c *gin.Context) {
	var req request.GetDatabaseConnectionListRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	result, err := ctrl.service.GetConnectionList(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// UpdateConnection update database connection
// @Summary Update database connection
// @Description Update an existing database connection configuration (without credentials)
// @Tags Database Connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Param request body request.UpdateDatabaseConnectionRequest true "Updated database connection configuration"
// @Success 200 {object} common.APIResponse{data=response.DatabaseConnectionResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/databases/{id} [put]
func (ctrl *DatabaseConnectionController) UpdateConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	var req request.UpdateDatabaseConnectionRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	result, err := ctrl.service.UpdateConnection(c.Request.Context(), uint(id), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Database connection updated successfully",
		logger.Uint("connection_id", uint(id)))

	response.SuccessWithData(c, result)
}

// DeleteConnection delete database connection
// @Summary Delete database connection
// @Description Delete an existing database connection by ID
// @Tags Database Connections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Connection ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/databases/{id} [delete]
func (ctrl *DatabaseConnectionController) DeleteConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, err)
		return
	}

	// Get owner ID from context
	ownerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c)
		return
	}

	// Call service
	if err := ctrl.service.DeleteConnection(c.Request.Context(), uint(id), ownerID.(uint)); err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Database connection deleted successfully",
		logger.Uint("connection_id", uint(id)))

	response.Success(c)
}

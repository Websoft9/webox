package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// EnvironmentVariableController handles environment variable HTTP requests
type EnvironmentVariableController struct {
	service   interfaceService.EnvironmentVariableService
	validator *validator.Validate
	logger    logger.Logger
}

// NewEnvironmentVariableController creates a new environment variable controller
func NewEnvironmentVariableController(
	service interfaceService.EnvironmentVariableService,
	validator *validator.Validate,
	logger logger.Logger,
) *EnvironmentVariableController {
	return &EnvironmentVariableController{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// CreatePlatformEnvVar create a platform-level environment variable
// @Summary Create platform environment variable
// @Description Create a new platform-level environment variable accessible across all projects
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreatePlatformEnvVarRequest true "Platform environment variable configuration"
// @Success 201 {object} common.APIResponse{data=response.EnvironmentVariableResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/platform [post]
func (ctrl *EnvironmentVariableController) CreatePlatformEnvVar(c *gin.Context) {
	var req request.CreatePlatformEnvVarRequest

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
	result, err := ctrl.service.CreatePlatformEnvVar(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Platform environment variable created successfully",
		logger.String("name", req.Name),
		logger.Uint("id", result.ID))

	response.SuccessWithData(c, result)
}

// CreateProjectEnvVar create a project-level environment variable
// @Summary Create project environment variable
// @Description Create a new project-level environment variable specific to a project
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateProjectEnvVarRequest true "Project environment variable configuration"
// @Success 201 {object} common.APIResponse{data=response.EnvironmentVariableResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/project [post]
func (ctrl *EnvironmentVariableController) CreateProjectEnvVar(c *gin.Context) {
	var req request.CreateProjectEnvVarRequest

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
	result, err := ctrl.service.CreateProjectEnvVar(c.Request.Context(), &req, ownerID.(uint))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Project environment variable created successfully",
		logger.String("name", req.Name),
		logger.Uint("project_id", req.ProjectID),
		logger.Uint("id", result.ID))

	response.SuccessWithData(c, result)
}

// GetEnvVar get environment variable by ID
// @Summary Get environment variable by ID
// @Description Get detailed information about a specific environment variable
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Environment Variable ID"
// @Success 200 {object} common.APIResponse{data=response.EnvironmentVariableDetailResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/{id} [get]
func (ctrl *EnvironmentVariableController) GetEnvVar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	// Call service
	result, err := ctrl.service.GetEnvVar(c.Request.Context(), uint(id))
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// GetPlatformEnvVarList get platform environment variable list
// @Summary Get platform environment variable list
// @Description Get a paginated list of platform-level environment variables with optional filters
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keywords query string false "Search keywords (name or description)"
// @Param is_sensitive query bool false "Filter by sensitive flag"
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/platform [get]
func (ctrl *EnvironmentVariableController) GetPlatformEnvVarList(c *gin.Context) {
	var req request.GetEnvVarListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	// Call service
	result, err := ctrl.service.GetPlatformEnvVarList(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// GetProjectEnvVarList get project environment variable list
// @Summary Get project environment variable list
// @Description Get a paginated list of project-level environment variables with optional filters
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id query int true "Project ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keywords query string false "Search keywords (name or description)"
// @Param is_sensitive query bool false "Filter by sensitive flag"
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/project [get]
func (ctrl *EnvironmentVariableController) GetProjectEnvVarList(c *gin.Context) {
	var req request.GetProjectEnvVarListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	// Validate project_id is required
	if err := ctrl.validator.Struct(&req); err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	// Call service
	result, err := ctrl.service.GetProjectEnvVarList(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// UpdateEnvVar update environment variable
// @Summary Update environment variable
// @Description Update an existing environment variable's value, description, or sensitivity
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Environment Variable ID"
// @Param request body request.UpdateEnvVarRequest true "Update configuration"
// @Success 200 {object} common.APIResponse{data=response.EnvironmentVariableResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/{id} [put]
func (ctrl *EnvironmentVariableController) UpdateEnvVar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	var req request.UpdateEnvVarRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.UpdateEnvVar(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Environment variable updated successfully",
		logger.Uint("id", uint(id)))

	response.SuccessWithData(c, result)
}

// DeleteEnvVar delete environment variable
// @Summary Delete environment variable
// @Description Delete an environment variable by ID
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Environment Variable ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/{id} [delete]
func (ctrl *EnvironmentVariableController) DeleteEnvVar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithError(c, errors.NewAppError(errors.CodeInvalidParameterFormat))
		return
	}

	// Call service
	err = ctrl.service.DeleteEnvVar(c.Request.Context(), uint(id))
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Environment variable deleted successfully",
		logger.Uint("id", uint(id)))

	response.Success(c)
}

// ResolveEnvVar resolve environment variable interpolation
// @Summary Resolve environment variable interpolation
// @Description Resolve ${VAR_NAME} and ${VAR_NAME:default} syntax in template string
// @Tags Environment Variables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.ResolveEnvVarRequest true "Resolution configuration"
// @Success 200 {object} common.APIResponse{data=response.ResolveEnvVarResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/environment-variables/resolve [post]
func (ctrl *EnvironmentVariableController) ResolveEnvVar(c *gin.Context) {
	var req request.ResolveEnvVarRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.ResolveEnvVar(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

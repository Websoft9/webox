package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// SecretController handles secret HTTP requests
type SecretController struct {
	service   interfaceService.SecretService
	validator *validator.Validate
	logger    logger.Logger
}

// NewSecretController creates a new secret controller
func NewSecretController(
	service interfaceService.SecretService,
	validator *validator.Validate,
	logger logger.Logger,
) *SecretController {
	return &SecretController{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// ListSecrets lists secrets with pagination and filters
// @Summary List secrets
// @Description Get a paginated list of secrets with optional filters
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keyword query string false "Keyword for name and description fuzzy search"
// @Param resource_code query string false "Filter by resource code"
// @Param start_time query string false "Start time for time range filter (RFC3339 format)"
// @Param end_time query string false "End time for time range filter (RFC3339 format)"
// @Param sort_field query string false "Sort field" default("created_at")
// @Param sort_order query string false "Sort order (ASC/DESC)" default("DESC")
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets [get]
func (ctrl *SecretController) ListSecrets(c *gin.Context) {
	var req request.ListSecretsRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.ListSecrets(c.Request.Context(), &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// CreateTextSecret creates a text type secret
// @Summary Create text secret
// @Description Create a new text type secret
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateTextSecretRequest true "Text secret configuration"
// @Success 201 {object} common.APIResponse{data=response.SecretResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/text [post]
func (ctrl *SecretController) CreateTextSecret(c *gin.Context) {
	var req request.CreateTextSecretRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.CreateTextSecret(c.Request.Context(), &req, ownerID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Text secret created successfully",
		logger.String("name", req.Name),
		logger.Uint("secret_id", result.ID))

	response.SuccessWithData(c, result)
}

// CreateAccountSecret creates an account type secret
// @Summary Create account secret
// @Description Create a new account type secret (username + password)
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateAccountSecretRequest true "Account secret configuration"
// @Success 201 {object} common.APIResponse{data=response.SecretResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/account [post]
func (ctrl *SecretController) CreateAccountSecret(c *gin.Context) {
	var req request.CreateAccountSecretRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.CreateAccountSecret(c.Request.Context(), &req, ownerID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Account secret created successfully",
		logger.String("name", req.Name),
		logger.Uint("secret_id", result.ID))

	response.SuccessWithData(c, result)
}

// CreateFileSecret creates a file type secret
// @Summary Create file secret
// @Description Create a new file type secret (certificate, key file, etc.)
// @Tags Secrets
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param resource_group_id formData int true "Resource group ID"
// @Param resource_code formData string false "Resource code (optional)"
// @Param name formData string true "Secret name"
// @Param description formData string false "Secret description"
// @Param secret_file formData file true "Secret file"
// @Param secret_password formData string false "Private key password (optional)"
// @Param authorized_users formData []int false "Authorized user IDs"
// @Param expires_at formData string false "Expiration time (RFC3339 format)"
// @Success 201 {object} common.APIResponse{data=response.SecretResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 413 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/file [post]
func (ctrl *SecretController) CreateFileSecret(c *gin.Context) {
	var req request.CreateFileSecretRequest

	// Bind multipart form data
	if err := c.ShouldBind(&req); err != nil {
		ctrl.logger.ErrorContext(c.Request.Context(), "Invalid request format", logger.ErrorField(err))
		response.WithErrorAndCode(c, errors.CodeInvalidParameterFormat, err)
		return
	}

	// Get owner ID from context
	ownerID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.CreateFileSecret(c.Request.Context(), &req, ownerID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "File secret created successfully",
		logger.String("name", req.Name),
		logger.Uint("secret_id", result.ID))

	response.SuccessWithData(c, result)
}

// GetSecret gets secret details by ID
// @Summary Get secret details
// @Description Get detailed information about a specific secret
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Secret ID"
// @Success 200 {object} common.APIResponse{data=response.SecretDetailResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [get]
func (ctrl *SecretController) GetSecret(c *gin.Context) {
	// Parse ID parameter
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.GetSecret(c.Request.Context(), id, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// UpdateSecret updates a secret
// @Summary Update secret
// @Description Update an existing secret's information
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Secret ID"
// @Param request body request.UpdateSecretRequest true "Secret update data"
// @Success 200 {object} common.APIResponse{data=response.SecretResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [put]
func (ctrl *SecretController) UpdateSecret(c *gin.Context) {
	// Parse ID parameter
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req request.UpdateSecretRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.UpdateSecret(c.Request.Context(), id, &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Secret updated successfully",
		logger.Uint("secret_id", id))

	response.SuccessWithData(c, result)
}

// DeleteSecret deletes a secret
// @Summary Delete secret
// @Description Delete a secret by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Secret ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [delete]
func (ctrl *SecretController) DeleteSecret(c *gin.Context) {
	// Parse ID parameter
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	if err := ctrl.service.DeleteSecret(c.Request.Context(), id, userID); err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Secret deleted successfully",
		logger.Uint("secret_id", id))

	response.Success(c)
}

// CreateReference creates a secret reference
// @Summary Create secret reference
// @Description Create a reference relationship between a secret and a resource
// @Tags Secrets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateReferenceRequest true "Reference configuration"
// @Success 201 {object} common.APIResponse{data=response.ReferenceResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/references [post]
func (ctrl *SecretController) CreateReference(c *gin.Context) {
	var req request.CreateReferenceRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.CreateReference(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Secret reference created successfully",
		logger.Uint("reference_id", result.ID),
		logger.Uint("secret_id", req.SecretID),
		logger.String("resource_code", req.ResourceCode))

	response.SuccessWithData(c, result)
}

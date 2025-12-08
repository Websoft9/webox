package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	interfaceService "api-service/internal/interface/service"
	"api-service/pkg/logger"
)

// CredentialController handles credential HTTP requests
type CredentialController struct {
	service   interfaceService.CredentialService
	validator *validator.Validate
	logger    logger.Logger
}

// NewCredentialController creates a new credential controller
func NewCredentialController(
	service interfaceService.CredentialService,
	validator *validator.Validate,
	logger logger.Logger,
) *CredentialController {
	return &CredentialController{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

// CreateCredential creates a new credential
// @Summary Create credential
// @Description Create a new credential with encrypted sensitive fields
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateCredentialRequest true "Credential configuration"
// @Success 201 {object} common.APIResponse{data=response.CredentialResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 409 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential [post]
func (ctrl *CredentialController) CreateCredential(c *gin.Context) {
	var req request.CreateCredentialRequest

	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get owner ID from context
	ownerID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.CreateCredential(c.Request.Context(), &req, ownerID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Credential created successfully",
		logger.String("name", req.Name),
		logger.Uint("credential_id", result.ID))

	response.SuccessWithData(c, result)
}

// UpdateCredential updates an existing credential
// @Summary Update credential
// @Description Update an existing credential's information
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Credential ID"
// @Param request body request.UpdateCredentialRequest true "Credential update data"
// @Success 200 {object} common.APIResponse{data=response.CredentialResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential/{id} [put]
func (ctrl *CredentialController) UpdateCredential(c *gin.Context) {
	// Parse ID parameter
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req request.UpdateCredentialRequest
	if !BindAndValidateRequest(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Get user ID from context
	userID, ok := GetUserID(c)
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.UpdateCredential(c.Request.Context(), id, &req, userID)
	if err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Credential updated successfully",
		logger.Uint("credential_id", id))

	response.SuccessWithData(c, result)
}

// DeleteCredential deletes a credential
// @Summary Delete credential
// @Description Delete a credential by ID
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Credential ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential/{id} [delete]
func (ctrl *CredentialController) DeleteCredential(c *gin.Context) {
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
	if err := ctrl.service.DeleteCredential(c.Request.Context(), id, userID); err != nil {
		response.WithError(c, err)
		return
	}

	ctrl.logger.InfoContext(c.Request.Context(), "Credential deleted successfully",
		logger.Uint("credential_id", id))

	response.Success(c)
}

// GetCredential gets credential details by ID
// @Summary Get credential details
// @Description Get detailed information about a specific credential
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Credential ID"
// @Success 200 {object} common.APIResponse{data=response.CredentialDetailResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential/{id} [get]
func (ctrl *CredentialController) GetCredential(c *gin.Context) {
	// Parse ID parameter
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	// Call service
	result, err := ctrl.service.GetCredential(c.Request.Context(), id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// ListCredentials lists credentials with pagination and filters
// @Summary List credentials
// @Description Get a paginated list of credentials with optional filters
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param keyword query string false "Keyword for name and description fuzzy search"
// @Param category_id query int false "Filter by category ID"
// @Param template_id query int false "Filter by template ID"
// @Param sort_by query string false "Sort field (created_at, updated_at, name)" default("created_at")
// @Param sort_order query string false "Sort order (asc, desc)" default("desc")
// @Success 200 {object} common.APIResponse{data=common.PaginationResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential [get]
func (ctrl *CredentialController) ListCredentials(c *gin.Context) {
	var req request.ListCredentialsRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.ListCredentials(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// ListCredentialTemplates lists credential templates
// @Summary List credential templates
// @Description Get a list of available credential templates
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id query int false "Filter by category ID"
// @Success 200 {object} common.APIResponse{data=[]response.CredentialTemplateResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential/templates [get]
func (ctrl *CredentialController) ListCredentialTemplates(c *gin.Context) {
	var req request.ListCredentialTemplatesRequest

	if !BindAndValidateQuery(c, &req, ctrl.validator, ctrl.logger) {
		return
	}

	// Call service
	result, err := ctrl.service.ListCredentialTemplates(c.Request.Context(), &req)
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

// ListCredentialCategories lists credential categories
// @Summary List credential categories
// @Description Get a list of all credential categories
// @Tags Credentials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} common.APIResponse{data=[]response.CredentialCategoryResponse}
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/credential/categories [get]
func (ctrl *CredentialController) ListCredentialCategories(c *gin.Context) {
	// Call service
	result, err := ctrl.service.ListCredentialCategories(c.Request.Context())
	if err != nil {
		response.WithError(c, err)
		return
	}

	response.SuccessWithData(c, result)
}

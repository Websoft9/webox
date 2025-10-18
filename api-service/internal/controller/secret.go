package controller

import (
	response "api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/interface/service"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SecretKeyController handles secret key related HTTP requests
type SecretKeyController struct {
	secretKeyService service.SecretKeyService
	logger           logger.Logger
	i18n             *i18n.I18n
	validator        *validator.Validate
}

// NewSecretKeyController creates a new secret key controller
func NewSecretKeyController(
	secretKeyService service.SecretKeyService,
	logger logger.Logger,
	i18n *i18n.I18n,
	validator *validator.Validate,
) *SecretKeyController {
	return &SecretKeyController{
		secretKeyService: secretKeyService,
		logger:           logger,
		i18n:             i18n,
		validator:        validator,
	}
}

// CreateSecretKeyText creates a new text-based secret key
// @Summary Create text secret key
// @Description Create a new text-based secret key (TEXT or ACCOUNT type)
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body request.SecretKeyCreateTextRequest true "Create text secret key request"
// @Success 201 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/text [post]
func (c *SecretKeyController) CreateSecretKeyText(ctx *gin.Context) {
	// Parse request parameters
	var req request.SecretKeyCreateTextRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Validate secret_fields based on auto-detected key_type
	if err := req.Validate(); err != nil {
		c.logger.ErrorContext(ctx, "Custom fields validation failed",
			logger.String("key_type", string(req.GetKeyType())),
			logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling create text secret key request", logger.Uint("userID", currentUserID))

	// Call service layer to create text secret key
	secretKey, err := c.secretKeyService.CreateSecretKeyText(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create text secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Text secret key created successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// CreateSecretKeyFile creates a new file-based secret key
// @Summary Create file secret key
// @Description Create a new file-based secret key (FILE type)
// @Tags Secrets
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Secret key name"
// @Param description formData string false "Secret key description"
// @Param file formData file true "Secret key file (.key, .pem, .rsa, .crt, .p12, .pfx)"
// @Param secret_fields formData string false "Secret fields as JSON string, e.g., {\"filename\": \"ssl-certificate.pem\", \"password\": \"***\"}"
// @Param resource_group_id formData integer false "Resource group ID"
// @Param expires_at formData string false "Expiration time (RFC3339 format)"
// @Param authorized_users formData string false "Comma-separated list of authorized user IDs"
// @Param resource_code formData string false "Resource code for reference"
// @Success 201 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/file [post]
func (c *SecretKeyController) CreateSecretKeyFile(ctx *gin.Context) {
	// Get file from request
	file, err := ctx.FormFile("file")
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get file from request", logger.ErrorField(err))
		response.BadRequest(ctx, errors.NewAppError(errors.CodeValidationFailed))
		return
	}

	// Create request object
	var req request.SecretKeyCreateFileRequest
	req.Name = ctx.PostForm("name")
	req.File = file

	// Handle optional description
	if description := ctx.PostForm("description"); description != "" {
		req.Description = &description
	}

	// Handle optional secret_fields (e.g., password)
	if secretFieldsStr := ctx.PostForm("secret_fields"); secretFieldsStr != "" {
		var secretFields map[string]interface{}
		err = json.Unmarshal([]byte(secretFieldsStr), &secretFields)
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to parse secret_fields", logger.ErrorField(err))
			response.WithError(ctx, errors.NewAppError(errors.CodeValidationFailed))
			return
		}
		req.SecretFields = secretFields
	}

	// Handle optional resource_code
	if resourceCode := ctx.PostForm("resource_code"); resourceCode != "" {
		req.ResourceCode = &resourceCode
	}

	// Validate request
	if err = req.Validate(); err != nil {
		c.logger.ErrorContext(ctx, "File secret key validation failed", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling create file secret key request",
		logger.Uint("userID", currentUserID),
		logger.String("filename", file.Filename))

	// Call service layer to create file secret key
	secretKey, err := c.secretKeyService.CreateSecretKeyFile(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create file secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "File secret key created successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// GetSecretKey retrieves a secret key by ID
// @Summary Get secret key
// @Description Get secret key information by ID
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [get]
func (c *SecretKeyController) GetSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling get secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to get secret key
	secretKey, err := c.secretKeyService.GetSecretKey(ctx.Request.Context(), id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to get secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key retrieved successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// UpdateSecretKey updates an existing secret key
// @Summary Update secret key
// @Description Update an existing secret key
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Param request body request.SecretKeyUpdateRequest true "Update secret key request"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 422 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [put]
func (c *SecretKeyController) UpdateSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Parse request parameters
	var req request.SecretKeyUpdateRequest
	// Bind and validate request
	if !BindAndValidateRequest(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}
	c.logger.InfoContext(ctx, "Handling update secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to update secret key
	secretKey, err := c.secretKeyService.UpdateSecretKey(ctx.Request.Context(), id, currentUserID, &req)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to update secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key updated successfully", logger.Uint("userID", currentUserID))
	response.SuccessWithData(ctx, secretKey)
}

// DeleteSecretKey deletes a secret key
// @Summary Delete secret key
// @Description Delete a secret key by ID
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Secret key ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 403 {object} common.APIResponse
// @Failure 404 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/{id} [delete]
func (c *SecretKeyController) DeleteSecretKey(ctx *gin.Context) {
	// Parse path parameter
	id, Success := ParseIDParam(ctx, "id")
	if !Success {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling delete secret key request", logger.Uint("userID", currentUserID), logger.Uint("secretKeyID", id))

	// Call service layer to delete secret key
	err := c.secretKeyService.DeleteSecretKey(ctx.Request.Context(), id, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to delete secret key", logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret key deleted successfully", logger.Uint("userID", currentUserID))
	response.Success(ctx)
}

// ListSecretKeys retrieves secret keys with pagination and filtering
// @Summary List secret keys
// @Description Get secret keys with pagination and filtering
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param key_type query string false "Key type filter" Enums(SECRET_KEY,ACCOUNT,FILE)
// @Param keyword query string false "Keyword for fuzzy search on name and description"
// @Param resource_code query string false "Resource code to filter secrets by reference"
// @Success 200 {object} common.APIResponse{data=response.SecretKeyListResponse}
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets [get]
func (c *SecretKeyController) ListSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyQueryRequest
	// Bind and validate request
	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling list secret keys request", logger.Uint("userID", currentUserID))

	// Call service layer to get secret keys
	secretKeys, err := c.secretKeyService.ListSecretKeys(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to list secret keys", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	c.logger.InfoContext(ctx, "Secret keys retrieved successfully", logger.Uint("user_id", currentUserID))
	response.SuccessWithData(ctx, secretKeys)
}

// ExportSecretKeys exports secret keys in specified format
// @Summary Export secret keys
// @Description Export secret keys in CSV, JSON, or Excel format
// @Tags Secrets
// @Security BearerAuth
// @Accept json
// @Produce application/octet-stream
// @Param format query string false "Export format" Enums(csv,json,excel) default(csv)
// @Success 200 {file} file "Exported file"
// @Failure 400 {object} common.APIResponse
// @Failure 401 {object} common.APIResponse
// @Failure 500 {object} common.APIResponse
// @Router /api/v1/secrets/export [get]
func (c *SecretKeyController) ExportSecretKeys(ctx *gin.Context) {
	// Bind request parameters
	var req request.SecretKeyExportRequest
	// Bind and validate request
	if !BindAndValidateQuery(ctx, &req, c.validator, c.logger) {
		return
	}

	// Get current user ID
	currentUserID, Success := GetUserID(ctx)
	if !Success {
		return
	}

	c.logger.InfoContext(ctx, "Handling export secret keys request", logger.Uint("userID", currentUserID))

	// Call service layer to export secret keys
	data, filename, err := c.secretKeyService.ExportSecretKeys(ctx.Request.Context(), &req, currentUserID)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to export secret keys", logger.Uint("user_id", currentUserID), logger.ErrorField(err))
		response.WithError(ctx, err)
		return
	}

	// Set response headers for file download
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Length", strconv.Itoa(len(data)))

	c.logger.InfoContext(ctx, "Secret keys exported successfully", logger.Uint("user_id", currentUserID))
	ctx.Writer.WriteHeader(http.StatusOK)
	_, _ = ctx.Writer.Write(data)
}

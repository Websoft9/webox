package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	iface "api-service/internal/interface/repository"
	svcIface "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

// Constants for validation
const (
	maxEnvVarNameLength = 64
	defaultValueIndex   = 2
)

// Compile regex once for better performance
// Variable name supports: letters (a-z, A-Z), digits (0-9), underscores (_)
// Must start with a letter or underscore
var envVarNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// environmentVariableService implements the EnvironmentVariableService interface
type environmentVariableService struct {
	repo iface.EnvironmentVariableRepository
	log  logger.Logger
}

// NewEnvironmentVariableService creates a new environment variable service
func NewEnvironmentVariableService(repo iface.EnvironmentVariableRepository, log logger.Logger) svcIface.EnvironmentVariableService {
	return &environmentVariableService{
		repo: repo,
		log:  log,
	}
}

// validateEnvVarName validates environment variable name format
// Business rule: supports letters (a-z, A-Z), digits (0-9), underscores (_)
// Must start with a letter or underscore, length 1-64 characters
func (s *environmentVariableService) validateEnvVarName(name string) error {
	if name == "" {
		return fmt.Errorf("environment variable name cannot be empty")
	}

	if len(name) > maxEnvVarNameLength {
		return fmt.Errorf("environment variable name cannot exceed %d characters", maxEnvVarNameLength)
	}

	if !envVarNamePattern.MatchString(name) {
		return errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	return nil
}

// CreatePlatformEnvVar creates a new platform-level environment variable
func (s *environmentVariableService) CreatePlatformEnvVar(
	ctx context.Context,
	req *request.CreatePlatformEnvVarRequest,
	ownerID uint,
) (*response.EnvironmentVariableResponse, error) {
	// Validate name format (business rule validation in service layer)
	if err := s.validateEnvVarName(req.Name); err != nil {
		s.log.Warn("Invalid environment variable name format",
			logger.String("name", req.Name),
			logger.String("error", err.Error()))
		return nil, err
	}

	// Check if platform variable with same name already exists
	exists, err := s.repo.ExistsByName(ctx, req.Name, model.EnvVarScopePlatform, nil, nil)
	if err != nil {
		s.log.Error("Failed to check platform variable existence",
			logger.String("name", req.Name),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Encrypt value if sensitive
	value := req.Value
	if req.IsSensitive {
		cryptoInstance := crypto.GetDefaultCrypto()
		encryptedValue, err := cryptoInstance.Encrypt(req.Value)
		if err != nil {
			s.log.Error("Failed to encrypt sensitive value",
				logger.String("name", req.Name),
				logger.String("error", err.Error()))
			return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
		}
		value = encryptedValue
	}

	// Create environment variable model
	envVar := &model.EnvironmentVariable{
		Name:        req.Name,
		Value:       value,
		Scope:       model.EnvVarScopePlatform,
		ProjectID:   nil,
		Description: req.Description,
		IsSensitive: req.IsSensitive,
		OwnerID:     ownerID,
	}

	// Create in repository
	if err := s.repo.Create(ctx, envVar); err != nil {
		s.log.Error("Failed to create platform variable",
			logger.String("name", req.Name),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.log.Info("Platform variable created successfully",
		logger.String("name", req.Name),
		logger.Uint("owner_id", ownerID))

	return s.buildEnvVarResponse(envVar), nil
}

// CreateProjectEnvVar creates a new project-level environment variable
func (s *environmentVariableService) CreateProjectEnvVar(
	ctx context.Context,
	req *request.CreateProjectEnvVarRequest,
	ownerID uint,
) (*response.EnvironmentVariableResponse, error) {
	// Validate name format (business rule validation in service layer)
	if err := s.validateEnvVarName(req.Name); err != nil {
		s.log.Warn("Invalid environment variable name format",
			logger.String("name", req.Name),
			logger.String("error", err.Error()))
		return nil, err
	}

	// Check if project variable with same name already exists
	exists, err := s.repo.ExistsByName(ctx, req.Name, model.EnvVarScopeProject, &req.ProjectID, nil)
	if err != nil {
		s.log.Error("Failed to check project variable existence",
			logger.String("name", req.Name),
			logger.Uint("project_id", req.ProjectID),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Encrypt value if sensitive
	value := req.Value
	if req.IsSensitive {
		cryptoInstance := crypto.GetDefaultCrypto()
		encryptedValue, err := cryptoInstance.Encrypt(req.Value)
		if err != nil {
			s.log.Error("Failed to encrypt sensitive value",
				logger.String("name", req.Name),
				logger.Uint("project_id", req.ProjectID),
				logger.String("error", err.Error()))
			return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
		}
		value = encryptedValue
	}

	// Create environment variable model
	projectIDVal := req.ProjectID
	envVar := &model.EnvironmentVariable{
		Name:        req.Name,
		Value:       value,
		Scope:       model.EnvVarScopeProject,
		ProjectID:   &projectIDVal,
		Description: req.Description,
		IsSensitive: req.IsSensitive,
		OwnerID:     ownerID,
	}

	// Create in repository
	if err := s.repo.Create(ctx, envVar); err != nil {
		s.log.Error("Failed to create project variable",
			logger.String("name", req.Name),
			logger.Uint("project_id", req.ProjectID),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.log.Info("Project variable created successfully",
		logger.String("name", req.Name),
		logger.Uint("project_id", req.ProjectID),
		logger.Uint("owner_id", ownerID))

	return s.buildEnvVarResponse(envVar), nil
}

// GetEnvVar retrieves an environment variable by ID
func (s *environmentVariableService) GetEnvVar(ctx context.Context, id uint) (*response.EnvironmentVariableDetailResponse, error) {
	envVar, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get variable",
			logger.Uint("id", id),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeResourceNotFound)
	}

	// Return detail response (value is masked for sensitive variables)
	return &response.EnvironmentVariableDetailResponse{
		EnvironmentVariableResponse: *s.buildEnvVarResponse(envVar),
	}, nil
}

// GetPlatformEnvVarList retrieves a paginated list of platform-level environment variables
func (s *environmentVariableService) GetPlatformEnvVarList(ctx context.Context, req *request.GetEnvVarListRequest) (*common.PaginationResponse, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = common.DEFAULT_PAGE_SIZE
	}

	// Get list from repository
	envVars, total, err := s.repo.GetPlatformList(ctx, req)
	if err != nil {
		s.log.Error("Failed to list platform variables",
			logger.Int("page", req.Page),
			logger.Int("page_size", req.PageSize),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	// Build response items
	items := make([]interface{}, 0, len(envVars))
	for _, envVar := range envVars {
		items = append(items, s.buildEnvVarResponse(envVar))
	}

	return common.NewPaginationResponse(req.Page, req.PageSize, total, items), nil
}

// GetProjectEnvVarList retrieves a paginated list of project-level environment variables
func (s *environmentVariableService) GetProjectEnvVarList(ctx context.Context, req *request.GetProjectEnvVarListRequest) (*common.PaginationResponse, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = common.DEFAULT_PAGE_SIZE
	}

	// Get list from repository
	envVars, total, err := s.repo.GetProjectList(ctx, req)
	if err != nil {
		s.log.Error("Failed to list project variables",
			logger.Uint("project_id", req.ProjectID),
			logger.Int("page", req.Page),
			logger.Int("page_size", req.PageSize),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	// Build response items
	items := make([]interface{}, 0, len(envVars))
	for _, envVar := range envVars {
		items = append(items, s.buildEnvVarResponse(envVar))
	}

	return common.NewPaginationResponse(req.Page, req.PageSize, total, items), nil
}

// UpdateEnvVar updates an existing environment variable
func (s *environmentVariableService) UpdateEnvVar(ctx context.Context, id uint, req *request.UpdateEnvVarRequest) (*response.EnvironmentVariableResponse, error) {
	// Get existing variable
	envVar, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get variable for update",
			logger.Uint("id", id),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeResourceNotFound)
	}

	// Update value if provided
	if req.Value != nil {
		value := *req.Value

		// Update is_sensitive if provided
		if req.IsSensitive != nil {
			envVar.IsSensitive = *req.IsSensitive
		}

		// Encrypt if sensitive
		if envVar.IsSensitive {
			cryptoInstance := crypto.GetDefaultCrypto()
			encryptedValue, err := cryptoInstance.Encrypt(value)
			if err != nil {
				s.log.Error("Failed to encrypt sensitive value",
					logger.Uint("id", id),
					logger.String("error", err.Error()))
				return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
			}
			value = encryptedValue
		}
		envVar.Value = value
	}

	// Update description if provided
	if req.Description != nil {
		envVar.Description = req.Description
	}

	// Update in repository
	if err := s.repo.Update(ctx, envVar); err != nil {
		s.log.Error("Failed to update variable",
			logger.Uint("id", id),
			logger.String("error", err.Error()))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.log.Info("Variable updated successfully",
		logger.Uint("id", id),
		logger.String("name", envVar.Name))

	return s.buildEnvVarResponse(envVar), nil
}

// DeleteEnvVar deletes an environment variable
func (s *environmentVariableService) DeleteEnvVar(ctx context.Context, id uint) error {
	// Check if variable exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get variable for deletion",
			logger.Uint("id", id),
			logger.String("error", err.Error()))
		return errors.NewAppErrorWrapError(err, errors.CodeResourceNotFound)
	}

	// Delete from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("Failed to delete variable",
			logger.Uint("id", id),
			logger.String("error", err.Error()))
		return errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.log.Info("Variable deleted successfully",
		logger.Uint("id", id))

	return nil
}

// ResolveEnvVar resolves environment variable interpolation in template string
// Supports ${VAR_NAME} and ${VAR_NAME:default_value} syntax
// Priority: project variables > platform variables
// Default value behavior: only used when variable exists but value is empty
// If variable not found: keeps original placeholder ${VAR_NAME}
func (s *environmentVariableService) ResolveEnvVar(ctx context.Context, req *request.ResolveEnvVarRequest) (*response.ResolveEnvVarResponse, error) {
	// Regular expression to match ${VAR_NAME} or ${VAR_NAME:default}
	// Variable name supports: letters (a-z, A-Z), digits (0-9), underscores (_)
	// Must start with a letter or underscore
	re := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*?)(?::([^}]*))?\}`)

	result := req.Template
	matches := re.FindAllStringSubmatch(req.Template, -1)

	for _, match := range matches {
		fullMatch := match[0] // ${VAR_NAME} or ${VAR_NAME:default}
		varName := match[1]   // VAR_NAME
		defaultValue := ""
		if len(match) > defaultValueIndex {
			defaultValue = match[defaultValueIndex] // default value if provided
		}

		// Try to resolve variable value
		value, err := s.resolveVariableValue(ctx, varName, req.Scope, req.ProjectID)
		if err != nil {
			// Variable not found - keep original placeholder
			// We don't use default value when variable doesn't exist, only when it's empty
			if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.CodeResourceNotFound {
				s.log.Warn("Variable not found, keeping placeholder",
					logger.String("variable", varName),
					logger.String("placeholder", fullMatch))
				continue
			} else {
				// Other errors - keep original placeholder
				s.log.Warn("Failed to resolve variable",
					logger.String("variable", varName),
					logger.String("error", err.Error()))
				continue
			}
		}

		// Variable exists - check if value is empty and default is provided
		if value == "" && defaultValue != "" {
			value = defaultValue
			s.log.Debug("Using default value for empty variable",
				logger.String("variable", varName),
				logger.String("default", defaultValue))
		}

		// Replace placeholder with value (or empty string if no default provided)
		result = strings.ReplaceAll(result, fullMatch, value)
	}

	return &response.ResolveEnvVarResponse{
		Result: result,
	}, nil
}

// resolveVariableValue resolves a single variable value with priority:
// 1. Project-level variable (if projectID provided and scope is "project")
// 2. Platform-level variable
func (s *environmentVariableService) resolveVariableValue(
	ctx context.Context,
	varName, scope string,
	projectID *uint,
) (string, error) {
	// Try project-level variable first (if scope is project and projectID provided)
	if scope == "project" && projectID != nil && *projectID > 0 {
		envVar, err := s.repo.GetByNameAndScope(ctx, varName, model.EnvVarScopeProject, projectID)
		if err == nil {
			return s.getDecryptedValue(envVar)
		}
		// If not found or error, continue to platform-level
		// Only log warning if it's not a "not found" error
		if appErr, ok := err.(*errors.AppError); !ok || appErr.Code != errors.CodeResourceNotFound {
			s.log.Warn("Error getting project variable",
				logger.String("variable", varName),
				logger.Uint("project_id", *projectID),
				logger.String("error", err.Error()))
		}
	}

	// Try platform-level variable
	envVar, err := s.repo.GetByNameAndScope(ctx, varName, model.EnvVarScopePlatform, nil)
	if err != nil {
		return "", err
	}

	return s.getDecryptedValue(envVar)
}

// getDecryptedValue decrypts the value if it's sensitive
func (s *environmentVariableService) getDecryptedValue(envVar *model.EnvironmentVariable) (string, error) {
	if !envVar.IsSensitive {
		return envVar.Value, nil
	}

	cryptoInstance := crypto.GetDefaultCrypto()
	decrypted, err := cryptoInstance.Decrypt(envVar.Value)
	if err != nil {
		s.log.Error("Failed to decrypt variable value",
			logger.Uint("id", envVar.ID),
			logger.String("name", envVar.Name),
			logger.String("error", err.Error()))
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	return decrypted, nil
}

// buildEnvVarResponse builds an EnvironmentVariableResponse from model
// Sensitive values are always masked
func (s *environmentVariableService) buildEnvVarResponse(envVar *model.EnvironmentVariable) *response.EnvironmentVariableResponse {
	value := envVar.Value

	// Mask sensitive value
	if envVar.IsSensitive {
		value = "******"
	}

	return &response.EnvironmentVariableResponse{
		ID:          envVar.ID,
		Name:        envVar.Name,
		Value:       value,
		Scope:       string(envVar.Scope),
		ProjectID:   envVar.ProjectID,
		Description: envVar.Description,
		IsSensitive: envVar.IsSensitive,
		OwnerID:     envVar.OwnerID,
		CreatedAt:   envVar.CreatedAt,
		UpdatedAt:   envVar.UpdatedAt,
	}
}

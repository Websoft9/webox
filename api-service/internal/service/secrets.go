package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/utils"

	"gorm.io/gorm"
)

const (
	// MaxFileSize is the maximum file size for secret files (5MB)
	MaxFileSize = 5 * 1024 * 1024
)

// allowedFileExtensions defines the allowed file extensions for secret files
var allowedFileExtensions = map[string]bool{
	".pem": true,
	".key": true,
	".crt": true,
	".cer": true,
	".p12": true,
	".pfx": true,
	".jks": true,
	".txt": true,
}

// secretService implements SecretService
type secretService struct {
	secretRepo        interfaceRepo.SecretRepository
	referenceRepo     interfaceRepo.SecretReferenceRepository
	authorizeRepo     interfaceRepo.SecretAuthorizeRepository
	resourceGroupRepo interfaceRepo.ResourceGroupRepository
	resourceTypeRepo  interfaceRepo.ResourceTypeRepository
	userRepo          interfaceRepo.UserRepository
	crypto            *crypto.AESCrypto
	fileStoragePath   string
	logger            logger.Logger
	db                *gorm.DB
}

// NewSecretService creates a new secret service
func NewSecretService(
	db *gorm.DB,
	secretRepo interfaceRepo.SecretRepository,
	referenceRepo interfaceRepo.SecretReferenceRepository,
	authorizeRepo interfaceRepo.SecretAuthorizeRepository,
	resourceGroupRepo interfaceRepo.ResourceGroupRepository,
	resourceTypeRepo interfaceRepo.ResourceTypeRepository,
	userRepo interfaceRepo.UserRepository,
	crypto *crypto.AESCrypto,
	fileStoragePath string,
	logger logger.Logger,
) (interfaceService.SecretService, error) {
	return &secretService{
		db:                db,
		secretRepo:        secretRepo,
		referenceRepo:     referenceRepo,
		authorizeRepo:     authorizeRepo,
		resourceGroupRepo: resourceGroupRepo,
		resourceTypeRepo:  resourceTypeRepo,
		userRepo:          userRepo,
		crypto:            crypto,
		fileStoragePath:   fileStoragePath,
		logger:            logger,
	}, nil
}

// encryptSecretField encrypts a single secret field value
func (s *secretService) encryptSecretField(value string) (string, error) {
	if value == "" {
		return "", nil
	}

	encrypted, err := s.crypto.Encrypt(value)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt secret field: %w", err)
	}

	return encrypted, nil
}

// decryptSecretField decrypts a single secret field value
// nolint:unused,godox // Reserved for future decryption needs
func (s *secretService) decryptSecretField(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	decrypted, err := s.crypto.Decrypt(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret field: %w", err)
	}

	return decrypted, nil
}

// CreateTextSecret creates a text type secret
func (s *secretService) CreateTextSecret(ctx context.Context, req *request.CreateTextSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	// Verify resource group and name uniqueness
	if err := s.verifyResourceGroupAndName(ctx, req.ResourceGroupID, req.Name, nil); err != nil {
		return nil, err
	}

	// Validate resource code before creating secret
	if req.ResourceCode != nil && *req.ResourceCode != "" {
		if validateErr := s.validateResourceCode(ctx, *req.ResourceCode); validateErr != nil {
			s.logger.ErrorContext(ctx, "Resource code validation failed",
				logger.String("resource_code", *req.ResourceCode),
				logger.ErrorField(validateErr))
			return nil, validateErr
		}
	}

	// Generate unique secret code
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	// Encrypt secret text
	encryptedText, err := s.encryptSecretField(req.SecretText)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt secret text", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	// Build secret fields JSON
	secretFields := model.JSON{
		"secret_text": encryptedText,
	}

	// Parse expiration time
	expiresAt := s.parseExpirationTime(req.ExpiresAt)

	// Create secret model
	secret := &model.Secret{
		Code:            code,
		Name:            req.Name,
		Type:            model.SecretTypeText,
		Description:     req.Description,
		SecretFields:    secretFields,
		ExpiresAt:       expiresAt,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         ownerID,
	}

	// Save to database
	if err := s.secretRepo.Create(ctx, secret); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create secret",
			logger.String("name", req.Name),
			logger.ErrorField(err))
		return nil, err
	}

	// Handle post-creation tasks
	s.handlePostCreation(ctx, secret.ID, req.AuthorizedUsers, req.ResourceCode)

	s.logger.InfoContext(ctx, "Text secret created successfully",
		logger.Uint("secret_id", secret.ID),
		logger.String("code", secret.Code),
		logger.String("name", secret.Name))

	return s.toResponse(ctx, secret), nil
}

// CreateAccountSecret creates an account type secret
func (s *secretService) CreateAccountSecret(ctx context.Context, req *request.CreateAccountSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	// Verify resource group and name uniqueness
	if err := s.verifyResourceGroupAndName(ctx, req.ResourceGroupID, req.Name, nil); err != nil {
		return nil, err
	}

	// Validate resource code before creating secret
	if req.ResourceCode != nil && *req.ResourceCode != "" {
		if validateErr := s.validateResourceCode(ctx, *req.ResourceCode); validateErr != nil {
			s.logger.ErrorContext(ctx, "Resource code validation failed",
				logger.String("resource_code", *req.ResourceCode),
				logger.ErrorField(validateErr))
			return nil, validateErr
		}
	}

	// Generate unique secret code
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	// Encrypt username and password
	encryptedUsername, err := s.encryptSecretField(req.SecretUsername)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt username", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	encryptedPassword, err := s.encryptSecretField(req.SecretPassword)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt password", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	// Build secret fields JSON
	secretFields := model.JSON{
		"secret_username": encryptedUsername,
		"secret_password": encryptedPassword,
	}

	// Parse expiration time
	expiresAt := s.parseExpirationTime(req.ExpiresAt)

	// Create secret model
	secret := &model.Secret{
		Code:            code,
		Name:            req.Name,
		Type:            model.SecretTypeAccount,
		Description:     req.Description,
		SecretFields:    secretFields,
		ExpiresAt:       expiresAt,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         ownerID,
	}

	// Save to database
	if err := s.secretRepo.Create(ctx, secret); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create secret",
			logger.String("name", req.Name),
			logger.ErrorField(err))
		return nil, err
	}

	// Handle post-creation tasks
	s.handlePostCreation(ctx, secret.ID, req.AuthorizedUsers, req.ResourceCode)

	s.logger.InfoContext(ctx, "Account secret created successfully",
		logger.Uint("secret_id", secret.ID),
		logger.String("code", secret.Code),
		logger.String("name", secret.Name))

	return s.toResponse(ctx, secret), nil
}

// validateFileUpload validates file size and type
func (s *secretService) validateFileUpload(req *request.CreateFileSecretRequest) (string, error) {
	// Validate file size
	if req.SecretFile.Size > MaxFileSize {
		return "", errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "secret.file_size_exceeded")
	}

	// Validate file type using map lookup for better performance
	ext := strings.ToLower(filepath.Ext(req.SecretFile.Filename))
	if !allowedFileExtensions[ext] {
		return "", errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "secret.file_type_not_allowed")
	}

	return ext, nil
}

// saveSecretFile saves uploaded file to storage and returns filename
func (s *secretService) saveSecretFile(ctx context.Context, req *request.CreateFileSecretRequest, ext string) (string, error) {
	// Generate UUID filename
	filename, err := utils.GenerateSecretCode()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate filename", logger.ErrorField(err))
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	filename = strings.ReplaceAll(filename, "secrets_", "") + ext

	// Ensure storage directory exists
	const dirPerm = 0755
	if mkdirErr := os.MkdirAll(s.fileStoragePath, dirPerm); mkdirErr != nil {
		s.logger.ErrorContext(ctx, "Failed to create storage directory",
			logger.String("path", s.fileStoragePath),
			logger.ErrorField(mkdirErr))
		return "", errors.NewAppErrorWrapError(mkdirErr, errors.CodeInternalError)
	}

	// Save file to storage
	filePath := filepath.Join(s.fileStoragePath, filename)
	src, err := req.SecretFile.Open()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to open uploaded file", logger.ErrorField(err))
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create file",
			logger.String("path", filePath),
			logger.ErrorField(err))
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save file", logger.ErrorField(err))
		os.Remove(filePath)
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	return filename, nil
}

// buildFileSecretFields builds secret fields JSON for file type secret
func (s *secretService) buildFileSecretFields(ctx context.Context, filename string, password *string, filePath string) (model.JSON, error) {
	secretFields := model.JSON{
		"secret_filename": filename,
	}

	if password != nil && *password != "" {
		encryptedPassword, err := s.encryptSecretField(*password)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to encrypt password", logger.ErrorField(err))
			os.Remove(filePath)
			return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
		}
		secretFields["secret_password"] = encryptedPassword
	}

	return secretFields, nil
}

// CreateFileSecret creates a file type secret
func (s *secretService) CreateFileSecret(ctx context.Context, req *request.CreateFileSecretRequest, ownerID uint) (*response.SecretResponse, error) {
	// Verify resource group and name uniqueness
	if err := s.verifyResourceGroupAndName(ctx, req.ResourceGroupID, req.Name, nil); err != nil {
		return nil, err
	}

	// Validate resource code before creating secret
	if req.ResourceCode != nil && *req.ResourceCode != "" {
		if err := s.validateResourceCode(ctx, *req.ResourceCode); err != nil {
			s.logger.ErrorContext(ctx, "Resource code validation failed",
				logger.String("resource_code", *req.ResourceCode),
				logger.ErrorField(err))
			return nil, err
		}
	}

	// Validate file upload
	ext, err := s.validateFileUpload(req)
	if err != nil {
		return nil, err
	}

	// Generate unique secret code
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	// Save file to storage
	filename, err := s.saveSecretFile(ctx, req, ext)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(s.fileStoragePath, filename)

	// Build secret fields
	secretFields, err := s.buildFileSecretFields(ctx, filename, req.SecretPassword, filePath)
	if err != nil {
		return nil, err
	}

	// Parse expiration time
	expiresAt := s.parseExpirationTime(req.ExpiresAt)

	// Create secret model
	secret := &model.Secret{
		Code:            code,
		Name:            req.Name,
		Type:            model.SecretTypeFile,
		Description:     req.Description,
		SecretFields:    secretFields,
		ExpiresAt:       expiresAt,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         ownerID,
	}

	// Save to database
	if err := s.secretRepo.Create(ctx, secret); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create secret",
			logger.String("name", req.Name),
			logger.ErrorField(err))
		os.Remove(filePath)
		return nil, err
	}

	// Handle post-creation tasks
	s.handlePostCreation(ctx, secret.ID, req.AuthorizedUsers, req.ResourceCode)

	s.logger.InfoContext(ctx, "File secret created successfully",
		logger.Uint("secret_id", secret.ID),
		logger.String("code", secret.Code),
		logger.String("name", secret.Name),
		logger.String("filename", filename))

	return s.toResponse(ctx, secret), nil
}

// verifyResourceGroupAndName verifies resource group exists and name is unique
func (s *secretService) verifyResourceGroupAndName(ctx context.Context, resourceGroupID uint, name string, excludeID *uint) error {
	// Verify resource group exists
	_, err := s.resourceGroupRepo.GetByID(ctx, resourceGroupID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Resource group not found",
			logger.Uint("resource_group_id", resourceGroupID),
			logger.ErrorField(err))
		return err
	}

	// Check name uniqueness
	exists, err := s.secretRepo.ExistsByName(ctx, resourceGroupID, name, excludeID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check secret name existence",
			logger.String("name", name),
			logger.ErrorField(err))
		return err
	}
	if exists {
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	return nil
}

// generateUniqueCode generates a unique secret code
func (s *secretService) generateUniqueCode(ctx context.Context) (string, error) {
	code, err := utils.GenerateSecretCodeWithRetry(func(candidateCode string) (bool, error) {
		codeExists, codeErr := s.secretRepo.ExistsByCode(ctx, candidateCode)
		return !codeExists, codeErr
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate unique secret code", logger.ErrorField(err))
		return "", errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	return code, nil
}

// parseExpirationTime parses expiration time from string
func (s *secretService) parseExpirationTime(expiresAtStr *string) *time.Time {
	if expiresAtStr != nil {
		if t, err := time.Parse(time.RFC3339, *expiresAtStr); err == nil {
			return &t
		}
	}
	return nil
}

// handlePostCreation handles post-creation tasks (authorizations and references)
// Resource code validation should be done before calling this method
func (s *secretService) handlePostCreation(ctx context.Context, secretID uint, authorizedUsers []uint, resourceCode *string) {
	// Handle authorized users
	if len(authorizedUsers) > 0 {
		s.createAuthorizations(ctx, secretID, authorizedUsers)
	}

	// Handle resource reference (validation already done before secret creation)
	if resourceCode != nil && *resourceCode != "" {
		if err := s.createReference(ctx, secretID, *resourceCode); err != nil {
			s.logger.ErrorContext(ctx, "Failed to create reference after secret creation",
				logger.Uint("secret_id", secretID),
				logger.String("resource_code", *resourceCode),
				logger.ErrorField(err))
			// Reference creation failed, but secret is already created
			// This should not happen as validation was done before
		}
	}
}

// createAuthorizations creates authorization records for a secret
func (s *secretService) createAuthorizations(ctx context.Context, secretID uint, userIDs []uint) {
	// Verify all users exist and are active
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil || user == nil {
			s.logger.WarnContext(ctx, "User not found",
				logger.Uint("user_id", userID),
				logger.ErrorField(err))
			continue
		}

		// Check if user is active
		if user.Status != UserStatusActive {
			s.logger.WarnContext(ctx, "User is not active, skipping authorization",
				logger.Uint("user_id", userID),
				logger.Int("status", user.Status))
			continue
		}

		// Create authorization
		authorize := &model.SecretAuthorize{
			SecretID:         secretID,
			AuthorizedUserID: userID,
		}

		if err := s.authorizeRepo.Create(ctx, authorize); err != nil {
			s.logger.WarnContext(ctx, "Failed to create authorization",
				logger.Uint("secret_id", secretID),
				logger.Uint("user_id", userID),
				logger.ErrorField(err))
		}
	}
}

// validateResourceCode validates the resource code and checks if the resource exists
// Resource code format: {resource_type}_{identifier}
// Examples: server_abc123, database_xyz789
//
// Validation steps:
// 1. Extract resource type code from the prefix
// 2. Query resource_types table to get the table name
// 3. Query the resource table to verify the resource exists
func (s *secretService) validateResourceCode(ctx context.Context, resourceCode string) error {
	if resourceCode == "" {
		return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "secret.invalid_resource_code")
	}

	// Step 1: Extract resource type code from prefix
	// Format: {resource_type}_{identifier}
	underscoreIndex := strings.Index(resourceCode, "_")
	if underscoreIndex == -1 || underscoreIndex == 0 {
		s.logger.WarnContext(ctx, "Invalid resource code format: missing underscore",
			logger.String("resource_code", resourceCode))
		return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "secret.invalid_resource_code")
	}

	resourceTypeCode := resourceCode[:underscoreIndex]
	s.logger.DebugContext(ctx, "Extracted resource type from code",
		logger.String("resource_code", resourceCode),
		logger.String("resource_type", resourceTypeCode))

	// Step 2: Query resource_types table to get the table name
	resourceType, err := s.resourceTypeRepo.GetByCode(ctx, resourceTypeCode)
	if err != nil {
		s.logger.WarnContext(ctx, "Resource type not found",
			logger.String("resource_type", resourceTypeCode),
			logger.ErrorField(err))
		return errors.NewAppErrorWithI18n(errors.CodeRecordNotFound, "secret.resource_type_not_found")
	}

	s.logger.DebugContext(ctx, "Found resource type",
		logger.String("resource_type", resourceTypeCode),
		logger.String("table_name", resourceType.Table))

	// Step 3: Query the resource table to verify the resource exists
	// Use raw SQL to query dynamic table name
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE code = ?", resourceType.Table)
	if err := s.db.WithContext(ctx).Raw(query, resourceCode).Scan(&count).Error; err != nil {
		s.logger.ErrorContext(ctx, "Failed to query resource table",
			logger.String("table_name", resourceType.Table),
			logger.String("resource_code", resourceCode),
			logger.ErrorField(err))
		return errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	if count == 0 {
		s.logger.WarnContext(ctx, "Resource not found in table",
			logger.String("table_name", resourceType.Table),
			logger.String("resource_code", resourceCode))
		return errors.NewAppErrorWithI18n(errors.CodeRecordNotFound, "secret.resource_not_found")
	}

	s.logger.InfoContext(ctx, "Resource code validated successfully",
		logger.String("resource_code", resourceCode),
		logger.String("resource_type", resourceTypeCode),
		logger.String("table_name", resourceType.Table))

	return nil
}

// createReference creates a reference record for a secret
// Resource code should be validated before calling this method
func (s *secretService) createReference(ctx context.Context, secretID uint, resourceCode string) error {
	// Check if reference already exists
	exists, err := s.referenceRepo.Exists(ctx, secretID, resourceCode)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already exists, skip
	}

	// Create reference
	reference := &model.SecretReference{
		SecretID:     secretID,
		ResourceCode: resourceCode,
	}

	return s.referenceRepo.Create(ctx, reference)
}

// toResponse converts model to response DTO
func (s *secretService) toResponse(ctx context.Context, secret *model.Secret) *response.SecretResponse {
	// Get reference count
	refCount, _ := s.secretRepo.GetReferenceCount(ctx, secret.ID)

	// Get owner name
	ownerName := ""
	if owner, err := s.userRepo.GetByID(ctx, secret.OwnerID); err == nil && owner != nil {
		ownerName = owner.Username
	}

	return &response.SecretResponse{
		ID:              secret.ID,
		Code:            secret.Code,
		Name:            secret.Name,
		Type:            string(secret.Type),
		Description:     secret.Description,
		ResourceGroupID: secret.ResourceGroupID,
		OwnerID:         secret.OwnerID,
		OwnerName:       ownerName,
		ReferenceCount:  refCount,
		CreatedAt:       secret.CreatedAt,
		UpdatedAt:       secret.UpdatedAt,
	}
}

// ListSecrets retrieves a paginated list of secrets
func (s *secretService) ListSecrets(ctx context.Context, req *request.ListSecretsRequest, userID uint) (*common.PaginationResponse, error) {
	secrets, total, err := s.secretRepo.List(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	items := make([]any, len(secrets))
	for i, secret := range secrets {
		items[i] = s.toResponse(ctx, secret)
	}

	return common.NewPaginationResponse(req.Page, req.PageSize, total, items), nil
}

// GetSecret retrieves a secret by ID with full details
func (s *secretService) GetSecret(ctx context.Context, id, userID uint) (*response.SecretDetailResponse, error) {
	// Get secret
	secret, err := s.secretRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check permission: user is owner or authorized user
	if secret.OwnerID != userID {
		isAuthorized, authErr := s.authorizeRepo.IsAuthorized(ctx, secret.ID, userID)
		if authErr != nil {
			s.logger.ErrorContext(ctx, "Failed to check authorization",
				logger.Uint("secret_id", secret.ID),
				logger.Uint("user_id", userID),
				logger.ErrorField(authErr))
			return nil, authErr
		}
		if !isAuthorized {
			return nil, errors.NewAppError(errors.CodeAccessDenied)
		}
	}

	// Get resource group name
	resourceGroupName := ""
	if rg, rgErr := s.resourceGroupRepo.GetByID(ctx, secret.ResourceGroupID); rgErr == nil && rg != nil {
		resourceGroupName = rg.Name
	}

	// Get owner name
	ownerName := ""
	if owner, ownerErr := s.userRepo.GetByID(ctx, secret.OwnerID); ownerErr == nil && owner != nil {
		ownerName = owner.Username
	}

	// Get references
	references, err := s.referenceRepo.ListBySecretID(ctx, secret.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to get references",
			logger.Uint("secret_id", secret.ID),
			logger.ErrorField(err))
		references = []*model.SecretReference{}
	}

	refResponses := make([]response.SecretReferenceResponse, len(references))
	for i, ref := range references {
		refResponses[i] = response.SecretReferenceResponse{
			ID:           ref.ID,
			ResourceCode: ref.ResourceCode,
			CreatedAt:    ref.CreatedAt,
		}
	}

	// Get authorizations
	authorizes, err := s.authorizeRepo.ListBySecretID(ctx, secret.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to get authorizations",
			logger.Uint("secret_id", secret.ID),
			logger.ErrorField(err))
		authorizes = []*model.SecretAuthorize{}
	}

	authResponses := make([]response.SecretAuthorizeResponse, len(authorizes))
	for i, auth := range authorizes {
		userName := ""
		if user, err := s.userRepo.GetByID(ctx, auth.AuthorizedUserID); err == nil && user != nil {
			userName = user.Username
		}

		authResponses[i] = response.SecretAuthorizeResponse{
			ID:        auth.ID,
			UserID:    auth.AuthorizedUserID,
			UserName:  userName,
			CreatedAt: auth.CreatedAt,
		}
	}

	return &response.SecretDetailResponse{
		ID:                secret.ID,
		Code:              secret.Code,
		Name:              secret.Name,
		Type:              string(secret.Type),
		Description:       secret.Description,
		ResourceGroupID:   secret.ResourceGroupID,
		ResourceGroupName: resourceGroupName,
		ExpiresAt:         secret.ExpiresAt,
		OwnerID:           secret.OwnerID,
		OwnerName:         ownerName,
		References:        refResponses,
		AuthorizedUsers:   authResponses,
		CreatedAt:         secret.CreatedAt,
		UpdatedAt:         secret.UpdatedAt,
	}, nil
}

// updateResourceGroup updates secret's resource group if provided
func (s *secretService) updateResourceGroup(ctx context.Context, secret *model.Secret, newGroupID *uint, secretID uint) error {
	if newGroupID == nil {
		return nil
	}

	// Verify target resource group exists
	_, err := s.resourceGroupRepo.GetByID(ctx, *newGroupID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Target resource group not found",
			logger.Uint("resource_group_id", *newGroupID),
			logger.ErrorField(err))
		return err
	}

	// Check if name exists in target resource group
	exists, err := s.secretRepo.ExistsByName(ctx, *newGroupID, secret.Name, &secretID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check secret name existence",
			logger.String("name", secret.Name),
			logger.ErrorField(err))
		return err
	}
	if exists {
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	secret.ResourceGroupID = *newGroupID
	return nil
}

// updateSecretName updates secret's name if provided
func (s *secretService) updateSecretName(ctx context.Context, secret *model.Secret, newName *string, secretID uint) error {
	if newName == nil {
		return nil
	}

	// Check if new name exists in resource group
	exists, err := s.secretRepo.ExistsByName(ctx, secret.ResourceGroupID, *newName, &secretID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check secret name existence",
			logger.String("name", *newName),
			logger.ErrorField(err))
		return err
	}
	if exists {
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	secret.Name = *newName
	return nil
}

// updateAuthorizedUsers updates secret's authorized users if provided
func (s *secretService) updateAuthorizedUsers(ctx context.Context, secretID uint, newUsers *[]uint) error {
	if newUsers == nil {
		return nil
	}

	// Delete existing authorizations
	if err := s.authorizeRepo.DeleteBySecretID(ctx, secretID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete existing authorizations",
			logger.Uint("secret_id", secretID),
			logger.ErrorField(err))
		return err
	}

	// Create new authorizations
	if len(*newUsers) > 0 {
		s.createAuthorizations(ctx, secretID, *newUsers)
	}

	return nil
}

// UpdateSecret updates an existing secret
func (s *secretService) UpdateSecret(ctx context.Context, id uint, req *request.UpdateSecretRequest, userID uint) (*response.SecretResponse, error) {
	// Get secret
	secret, err := s.secretRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check permission: only owner can update
	if secret.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Update resource group if provided
	if err := s.updateResourceGroup(ctx, secret, req.ResourceGroupID, id); err != nil {
		return nil, err
	}

	// Update name if provided
	if err := s.updateSecretName(ctx, secret, req.Name, id); err != nil {
		return nil, err
	}

	// Update description if provided
	if req.Description != nil {
		secret.Description = req.Description
	}

	// Update expiration time if provided
	if req.ExpiresAt != nil {
		secret.ExpiresAt = s.parseExpirationTime(req.ExpiresAt)
	}

	// Update authorized users if provided (nil means no update, empty slice means clear all)
	if err := s.updateAuthorizedUsers(ctx, secret.ID, req.AuthorizedUsers); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.secretRepo.Update(ctx, secret); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update secret",
			logger.Uint("secret_id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Secret updated successfully",
		logger.Uint("secret_id", id),
		logger.String("name", secret.Name))

	return s.toResponse(ctx, secret), nil
}

// DeleteSecret deletes a secret
func (s *secretService) DeleteSecret(ctx context.Context, id, userID uint) error {
	// Get secret
	secret, err := s.secretRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check permission: only owner can delete
	if secret.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	// Get file path before deletion (for file type secrets)
	var filePath string
	if secret.Type == model.SecretTypeFile {
		if filename, ok := secret.SecretFields["secret_filename"].(string); ok {
			filePath = filepath.Join(s.fileStoragePath, filename)
		}
	}

	// Check if secret has active references before transaction
	hasRefs, checkErr := s.referenceRepo.HasActiveReferences(ctx, secret.ID)
	if checkErr != nil {
		s.logger.ErrorContext(ctx, "Failed to check active references",
			logger.Uint("secret_id", secret.ID),
			logger.ErrorField(checkErr))
		return checkErr
	}

	if hasRefs {
		// Get all references for error response
		references, _ := s.referenceRepo.ListBySecretID(ctx, secret.ID)
		refCodes := make([]string, len(references))
		for i, ref := range references {
			refCodes[i] = ref.ResourceCode
		}

		s.logger.WarnContext(ctx, "Cannot delete secret with active references",
			logger.Uint("secret_id", secret.ID),
			logger.Any("references", refCodes))

		// Return 409 Conflict with reference details
		return errors.NewAppErrorWithI18n(errors.CodeResourceInUse, "secret.has_active_references")
	}

	// Execute deletion in transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Delete authorizations using transaction
		if deleteErr := tx.Where("secret_id = ?", secret.ID).Delete(&model.SecretAuthorize{}).Error; deleteErr != nil {
			s.logger.ErrorContext(ctx, "Failed to delete authorizations",
				logger.Uint("secret_id", secret.ID),
				logger.ErrorField(deleteErr))
			return errors.NewAppErrorWrapError(deleteErr, errors.CodeRecordDeleteFailed)
		}

		// Step 2: Delete secret record using transaction
		if deleteErr := tx.Delete(&model.Secret{}, secret.ID).Error; deleteErr != nil {
			s.logger.ErrorContext(ctx, "Failed to delete secret",
				logger.Uint("secret_id", secret.ID),
				logger.ErrorField(deleteErr))
			return errors.NewAppErrorWrapError(deleteErr, errors.CodeRecordDeleteFailed)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Step 4: Asynchronously delete physical file (outside transaction)
	if filePath != "" {
		go func() {
			if err := os.Remove(filePath); err != nil {
				s.logger.Warn("Failed to delete secret file",
					logger.String("path", filePath),
					logger.ErrorField(err))
			} else {
				s.logger.Info("Secret file deleted successfully",
					logger.String("path", filePath))
			}
		}()
	}

	s.logger.InfoContext(ctx, "Secret deleted successfully",
		logger.Uint("secret_id", secret.ID),
		logger.String("name", secret.Name))

	return nil
}

// CreateReference creates a secret reference
func (s *secretService) CreateReference(ctx context.Context, req *request.CreateReferenceRequest) (*response.ReferenceResponse, error) {
	// Verify secret exists
	secret, err := s.secretRepo.GetByID(ctx, req.SecretID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Secret not found",
			logger.Uint("secret_id", req.SecretID),
			logger.ErrorField(err))
		return nil, err
	}

	// Validate resource code and verify resource exists
	if validateErr := s.validateResourceCode(ctx, req.ResourceCode); validateErr != nil {
		s.logger.ErrorContext(ctx, "Resource code validation failed",
			logger.String("resource_code", req.ResourceCode),
			logger.ErrorField(validateErr))
		return nil, validateErr
	}

	// Check if reference already exists
	exists, err := s.referenceRepo.Exists(ctx, req.SecretID, req.ResourceCode)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check reference existence",
			logger.Uint("secret_id", req.SecretID),
			logger.String("resource_code", req.ResourceCode),
			logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Create reference
	reference := &model.SecretReference{
		SecretID:     req.SecretID,
		ResourceCode: req.ResourceCode,
	}

	if err := s.referenceRepo.Create(ctx, reference); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create reference",
			logger.Uint("secret_id", req.SecretID),
			logger.String("resource_code", req.ResourceCode),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Secret reference created successfully",
		logger.Uint("reference_id", reference.ID),
		logger.Uint("secret_id", req.SecretID),
		logger.String("secret_name", secret.Name),
		logger.String("resource_code", req.ResourceCode))

	return &response.ReferenceResponse{
		ID:           reference.ID,
		SecretID:     reference.SecretID,
		ResourceCode: reference.ResourceCode,
		CreatedAt:    reference.CreatedAt,
	}, nil
}

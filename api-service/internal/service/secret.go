package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// secretKeyService implements SecretKeyService interface
type secretKeyService struct {
	secretKeyRepo repository.SecretKeyRepository
	logger        logger.Logger
	i18n          *i18n.I18n
	rsaCrypto     *crypto.RSACrypto
	appConfig     *config.Config
}

// NewSecretKeyService creates a new secret key service
func NewSecretKeyService(
	secretKeyRepo repository.SecretKeyRepository,
	logger logger.Logger,
	i18n *i18n.I18n,
	cfg *config.Config,
) service.SecretKeyService {
	rsaCrypto, err := initRSACryptoFromConfig(&cfg.Security)

	if err != nil {
		rsaCrypto, err = crypto.NewRSACrypto(crypto.MinRSAKeySize)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize RSA crypto: %v", err))
		}
	}
	return &secretKeyService{
		secretKeyRepo: secretKeyRepo,
		logger:        logger,
		i18n:          i18n,
		rsaCrypto:     rsaCrypto,
		appConfig:     cfg,
	}
}

func initRSACryptoFromConfig(securityConfig *config.SecurityConfig) (*crypto.RSACrypto, error) {
	var privateKeyPEM, publicKeyPEM string

	if securityConfig.RSAPrivateKey != "" && securityConfig.RSAPublicKey != "" {
		privateKeyPEM = securityConfig.RSAPrivateKey
		publicKeyPEM = securityConfig.RSAPublicKey
	} else if securityConfig.RSAPrivateKeyFile != "" && securityConfig.RSAPublicKeyFile != "" {
		privateKeyBytes, err := os.ReadFile(securityConfig.RSAPrivateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		privateKeyPEM = string(privateKeyBytes)

		publicKeyBytes, err := os.ReadFile(securityConfig.RSAPublicKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key file: %w", err)
		}
		publicKeyPEM = string(publicKeyBytes)
	} else {
		return nil, fmt.Errorf("no RSA keys configured")
	}

	return crypto.NewRSACryptoFromKeys(privateKeyPEM, publicKeyPEM)
}

// encryptCustomFields encrypts sensitive fields in custom_fields based on key_type
func (s *secretKeyService) encryptCustomFields(ctx context.Context, keyType model.SecretKeyType, customFields map[string]interface{}) (model.CustomFields, error) {
	if customFields == nil {
		return make(model.CustomFields), nil
	}

	encryptedFields := make(model.CustomFields)

	// Copy all fields first
	for k, v := range customFields {
		encryptedFields[k] = v
	}

	// Encrypt specific fields based on key_type
	switch keyType {
	case model.SecretKeyTypeSecretKey:
		return s.encryptSecretKeyFields(ctx, encryptedFields)
	case model.SecretKeyTypeAccount:
		return s.encryptAccountFields(ctx, encryptedFields)
	case model.SecretKeyTypeFile:
		return s.encryptFileFields(ctx, encryptedFields)
	}

	return encryptedFields, nil
}

// encryptSecretKeyFields encrypts secret_key field for SECRET_KEY type
func (s *secretKeyService) encryptSecretKeyFields(ctx context.Context, fields model.CustomFields) (model.CustomFields, error) {
	if secretKey, ok := fields["secret_key"].(string); ok && secretKey != "" {
		encryptedSecretKey, err := s.rsaCrypto.EncryptString(secretKey)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to encrypt secret_key field", logger.ErrorField(err))
			return nil, errors.NewAppError(errors.CodeEncryptFailed)
		}
		fields["secret_key"] = encryptedSecretKey
	}
	return fields, nil
}

// encryptAccountFields encrypts password field for ACCOUNT type
func (s *secretKeyService) encryptAccountFields(ctx context.Context, fields model.CustomFields) (model.CustomFields, error) {
	if password, ok := fields["password"].(string); ok && password != "" {
		encryptedPassword, err := s.rsaCrypto.EncryptString(password)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to encrypt password field", logger.ErrorField(err))
			return nil, errors.NewAppError(errors.CodeEncryptFailed)
		}
		fields["password"] = encryptedPassword
	}
	return fields, nil
}

// encryptFileFields encrypts password field for FILE type (if provided)
func (s *secretKeyService) encryptFileFields(ctx context.Context, fields model.CustomFields) (model.CustomFields, error) {
	if password, ok := fields["password"].(string); ok && password != "" {
		encryptedPassword, err := s.rsaCrypto.EncryptString(password)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to encrypt password field", logger.ErrorField(err))
			return nil, errors.NewAppError(errors.CodeEncryptFailed)
		}
		fields["password"] = encryptedPassword
	}
	return fields, nil
}

// createUserSecretRelationships creates user-secret relationships for authorized users
func (s *secretKeyService) createUserSecretRelationships(ctx context.Context, secretKeyID uint, authorizedUsers []uint, userID uint, expiresAt *time.Time) {
	if len(authorizedUsers) == 0 {
		return
	}

	for _, authorizedUserID := range authorizedUsers {
		userSecret := &model.UserSecret{
			UserID:      authorizedUserID,
			SecretKeyID: secretKeyID,
			GrantedBy:   &userID,
			GrantedAt:   time.Now(),
			ExpiresAt:   expiresAt,
		}

		if err := s.secretKeyRepo.CreateUserSecret(ctx, userSecret); err != nil {
			s.logger.ErrorContext(ctx, "Failed to create user secret relationship",
				logger.Uint("secret_key_id", secretKeyID),
				logger.Uint("authorized_user_id", authorizedUserID),
				logger.ErrorField(err))
			continue
		}

		s.logger.InfoContext(ctx, "User secret relationship created",
			logger.Uint("secret_key_id", secretKeyID),
			logger.Uint("authorized_user_id", authorizedUserID))
	}

	s.logger.InfoContext(ctx, "User secret relationships creation completed",
		logger.Uint("secret_key_id", secretKeyID),
		logger.Int("total_users", len(authorizedUsers)))
}

// decryptCustomFields decrypts sensitive fields in custom_fields based on key_type
func (s *secretKeyService) decryptCustomFields(ctx context.Context, keyType model.SecretKeyType, customFields model.CustomFields) (map[string]interface{}, error) {
	decryptedFields := make(map[string]interface{})

	// Copy all fields first
	for k, v := range customFields {
		decryptedFields[k] = v
	}

	// Decrypt specific fields based on key_type
	switch keyType {
	case model.SecretKeyTypeSecretKey:
		return s.decryptSecretKeyFields(ctx, decryptedFields)
	case model.SecretKeyTypeAccount:
		return s.decryptAccountFields(ctx, decryptedFields)
	case model.SecretKeyTypeFile:
		return s.decryptFileFields(ctx, decryptedFields)
	}

	return decryptedFields, nil
}

// decryptSecretKeyFields decrypts secret_key field for SECRET_KEY type
func (s *secretKeyService) decryptSecretKeyFields(ctx context.Context, fields map[string]interface{}) (map[string]interface{}, error) {
	if encryptedSecretKey, ok := fields["secret_key"].(string); ok && encryptedSecretKey != "" {
		decryptedValue, err := s.rsaCrypto.DecryptString(encryptedSecretKey)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to decrypt secret_key field", logger.ErrorField(err))
			return nil, err
		}
		fields["secret_key"] = decryptedValue
	}
	return fields, nil
}

// decryptAccountFields decrypts password field for ACCOUNT type
func (s *secretKeyService) decryptAccountFields(ctx context.Context, fields map[string]interface{}) (map[string]interface{}, error) {
	if encryptedPassword, ok := fields["password"].(string); ok && encryptedPassword != "" {
		decryptedValue, err := s.rsaCrypto.DecryptString(encryptedPassword)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to decrypt password field", logger.ErrorField(err))
			return nil, err
		}
		fields["password"] = decryptedValue
	}
	return fields, nil
}

// decryptFileFields decrypts password field for FILE type (if exists)
func (s *secretKeyService) decryptFileFields(ctx context.Context, fields map[string]interface{}) (map[string]interface{}, error) {
	if encryptedPassword, ok := fields["password"].(string); ok && encryptedPassword != "" {
		decryptedValue, err := s.rsaCrypto.DecryptString(encryptedPassword)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to decrypt password field", logger.ErrorField(err))
			return nil, err
		}
		fields["password"] = decryptedValue
	}
	return fields, nil
}

// CreateSecretKey creates a new secret key
func (s *secretKeyService) CreateSecretKey(ctx context.Context, req *request.SecretKeyCreateRequest, userID uint) (*response.SecretKeyResponse, error) {
	s.logger.InfoContext(ctx, "Creating secret key",
		logger.String("name", req.Name),
		logger.Uint("user_id", userID))

	// Check if name already exists for this user
	exists, err := s.secretKeyRepo.ExistsByName(ctx, req.Name, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check secret key name existence",
			logger.String("name", req.Name),
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, err
	}

	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Encrypt sensitive fields in custom_fields based on key_type
	encryptedCustomFields, err := s.encryptCustomFields(ctx, req.KeyType, req.CustomFields)
	if err != nil {
		return nil, err
	}

	// Create the secret key model
	secretKey := &model.SecretKey{
		Name:            req.Name,
		KeyType:         req.KeyType,
		Description:     req.Description,
		CustomFields:    encryptedCustomFields,
		ExpiresAt:       req.ExpiresAt,
		ResourceGroupID: req.ResourceGroupID,
		OwnerID:         userID,
	}

	// Save to database
	if err := s.secretKeyRepo.Create(ctx, secretKey); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create secret key",
			logger.String("name", req.Name),
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, err
	}

	// Create user-secret relationships for authorized users
	s.createUserSecretRelationships(ctx, secretKey.ID, req.AuthorizedUsers, userID, req.ExpiresAt)

	s.logger.InfoContext(ctx, "Secret key created successfully",
		logger.Uint("secret_key_id", secretKey.ID),
		logger.Uint("user_id", userID))

	return response.ToSecretKeyResponse(secretKey), nil
}

// GetSecretKey retrieves a secret key by ID
func (s *secretKeyService) GetSecretKey(ctx context.Context, id, userID uint) (*response.SecretKeyResponse, error) {
	s.logger.InfoContext(ctx, "Getting secret key",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	secretKey, err := s.secretKeyRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	// Check ownership or user access permission
	if secretKey.OwnerID != userID {
		// If not owner, check if user has access permission in user_secret table
		hasAccess, err := s.secretKeyRepo.CheckUserSecretAccess(ctx, userID, id)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to check user secret access",
				logger.Uint("secret_key_id", id),
				logger.Uint("user_id", userID),
				logger.ErrorField(err))
			return nil, err
		}

		if !hasAccess {
			return nil, errors.NewAppError(errors.CodeAccessDenied)
		}
	}

	return response.ToSecretKeyResponse(secretKey), nil
}

// GetSecretKeyValue retrieves the decrypted value of a secret key
func (s *secretKeyService) GetSecretKeyValue(ctx context.Context, id, userID uint) (*response.SecretKeyValueResponse, error) {
	s.logger.InfoContext(ctx, "Getting secret key value",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	secretKey, err := s.secretKeyRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	// Check ownership or user access permission
	if secretKey.OwnerID != userID {
		hasAccess, accessErr := s.secretKeyRepo.CheckUserSecretAccess(ctx, userID, id)
		if accessErr != nil {
			s.logger.ErrorContext(ctx, "Failed to check user secret access",
				logger.Uint("secret_key_id", id),
				logger.Uint("user_id", userID),
				logger.ErrorField(accessErr))
			return nil, accessErr
		}

		if !hasAccess {
			return nil, errors.NewAppError(errors.CodeAccessDenied)
		}
	}

	// Check if expired
	if secretKey.IsExpired() {
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Decrypt sensitive fields in custom_fields based on key_type
	decryptedCustomFields, err := s.decryptCustomFields(ctx, secretKey.KeyType, secretKey.CustomFields)
	if err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "Secret key value accessed",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

	return response.ToSecretKeyValueResponse(secretKey.KeyType, decryptedCustomFields, secretKey.ExpiresAt), nil
}

// UpdateSecretKey updates an existing secret key
func (s *secretKeyService) UpdateSecretKey(ctx context.Context, id, userID uint, req *request.SecretKeyUpdateRequest) (*response.SecretKeyResponse, error) {
	s.logger.InfoContext(ctx, "Updating secret key",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	secretKey, err := s.secretKeyRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Encrypt sensitive fields in custom_fields based on key_type
	encryptedCustomFields := make(model.CustomFields)
	if req.CustomFields != nil {
		var err error
		encryptedCustomFields, err = s.encryptCustomFields(ctx, req.KeyType, req.CustomFields)
		if err != nil {
			return nil, err
		}
	}

	// Update fields
	secretKey.KeyType = req.KeyType
	secretKey.CustomFields = encryptedCustomFields

	// Save changes
	if err := s.secretKeyRepo.Update(ctx, secretKey); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Secret key updated successfully")
	return response.ToSecretKeyResponse(secretKey), nil
}

// DeleteSecretKey deletes a secret key
func (s *secretKeyService) DeleteSecretKey(ctx context.Context, id, userID uint) error {
	s.logger.InfoContext(ctx, "Deleting secret key",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	secretKey, err := s.secretKeyRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	// First delete all user_secret relationships for this secret key
	if err := s.secretKeyRepo.DeleteUserSecretsBySecretKeyID(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete user secret relationships",
			logger.Uint("secret_key_id", id),
			logger.ErrorField(err))
		return errors.NewAppError(errors.CodeRecordDeleteFailed)
	}

	s.logger.InfoContext(ctx, "User secret relationships deleted successfully",
		logger.Uint("secret_key_id", id))

	// Then delete the secret key itself
	if err := s.secretKeyRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Secret key deleted successfully",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

	return nil
}

// ListSecretKeys retrieves secret keys with pagination and filtering
func (s *secretKeyService) ListSecretKeys(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Listing secret keys",
		logger.Uint("user_id", userID),
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize))

	secretKeys, total, err := s.secretKeyRepo.List(ctx, req, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list secret keys",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, err
	}

	return response.ToSecretKeyListResponse(secretKeys, total, req.GetOffset(), req.GetPageSize()), nil
}

// ExportSecretKeys exports secret keys in specified format
func (s *secretKeyService) ExportSecretKeys(ctx context.Context, req *request.SecretKeyExportRequest, userID uint) (data []byte, filename string, err error) {
	s.logger.InfoContext(ctx, "Exporting secret keys",
		logger.Uint("user_id", userID),
		logger.String("format", req.Format))

	secretKeys, err := s.secretKeyRepo.ListAll(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret keys for export",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, "", err
	}

	switch req.Format {
	case constants.FormatCSV:
		return s.exportToCSV(secretKeys)
	case constants.FormatJSON:
		return s.exportToJSON(secretKeys)
	case constants.FormatExcel:
		return s.exportToExcel(secretKeys)
	default:
		return nil, "", errors.NewAppError(errors.CodeValidationFailed)
	}
}

// exportToExcel exports secret keys to Excel format
func (s *secretKeyService) exportToExcel(secretKeys []*model.SecretKey) (data []byte, filename string, err error) {
	s.logger.Info("Starting Excel export", logger.Int("count", len(secretKeys)))

	f := excelize.NewFile()
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			s.logger.Error("Failed to close Excel file", logger.ErrorField(closeErr))
		}
	}()

	sheetName := "SecretKeys"
	index, err := s.createExcelSheet(f, sheetName)
	if err != nil {
		return nil, "", err
	}

	// Set headers
	if err := s.setExcelHeaders(f, sheetName); err != nil {
		return nil, "", err
	}

	// Set data
	if err := s.setExcelData(f, sheetName, secretKeys); err != nil {
		return nil, "", err
	}

	// Finalize workbook
	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		s.logger.Warn("Failed to delete default sheet", logger.ErrorField(err))
		// 这个错误不是致命的，继续执行
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		s.logger.Error("Failed to write Excel file to buffer", logger.ErrorField(err))
		return nil, "", err
	}

	s.logger.Info("Excel file generated successfully", logger.Int("size", buf.Len()))

	filename = fmt.Sprintf("secret-keys-%s.xlsx", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

// createExcelSheet creates a new Excel sheet
func (s *secretKeyService) createExcelSheet(f *excelize.File, sheetName string) (int, error) {
	index, err := f.NewSheet(sheetName)
	if err != nil {
		s.logger.Error("Failed to create Excel sheet", logger.ErrorField(err))
		return 0, err
	}

	s.logger.Info("Excel sheet created successfully", logger.Int("index", index))
	return index, nil
}

// setExcelHeaders sets the headers in Excel sheet
func (s *secretKeyService) setExcelHeaders(f *excelize.File, sheetName string) error {
	headers := []string{"ID", "Name", "Type", "Description", "Created At", "Updated At", "Expires At"}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			s.logger.Error("Failed to set header cell",
				logger.ErrorField(err),
				logger.String("cell", cell),
				logger.String("header", header))
			return err
		}
	}

	return nil
}

// setExcelData sets the data rows in Excel sheet
func (s *secretKeyService) setExcelData(f *excelize.File, sheetName string, secretKeys []*model.SecretKey) error {
	const excelDataStartRow = 2 // Data starts from row 2 (row 1 is header)

	for i, sk := range secretKeys {
		row := i + excelDataStartRow
		if err := s.setExcelRowData(f, sheetName, row, sk); err != nil {
			return err
		}
	}
	return nil
}

// setExcelRowData sets a single row of data in Excel sheet
func (s *secretKeyService) setExcelRowData(f *excelize.File, sheetName string, row int, sk *model.SecretKey) error {
	description := ""
	if sk.Description != nil {
		description = *sk.Description
	}

	expiresAt := ""
	if sk.ExpiresAt != nil {
		expiresAt = sk.ExpiresAt.Format(time.RFC3339)
	}

	// Define the data to be set
	data := []interface{}{
		sk.ID,
		sk.Name,
		string(sk.KeyType),
		description,
		sk.CreatedAt.Format(time.RFC3339),
		sk.UpdatedAt.Format(time.RFC3339),
		expiresAt,
	}

	// Set each cell
	for col, value := range data {
		cell := fmt.Sprintf("%s%d", string(rune('A'+col)), row)
		if err := f.SetCellValue(sheetName, cell, value); err != nil {
			s.logger.Error("Failed to set data cell",
				logger.ErrorField(err),
				logger.String("cell", cell),
				logger.Int("row", row))
			return err
		}
	}

	return nil
}

// ValidateSecretKeyOwnership checks if user owns the secret key
func (s *secretKeyService) ValidateSecretKeyOwnership(ctx context.Context, secretKeyID, userID uint) error {
	secretKey, err := s.secretKeyRepo.GetByID(ctx, secretKeyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound)
		}
		return err
	}

	if secretKey.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	return nil
}

// exportToCSV exports secret keys to CSV format
func (s *secretKeyService) exportToCSV(secretKeys []*model.SecretKey) (data []byte, filename string, err error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Name", "Type", "Description", "Created At", "Updated At", "Expires At"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data
	for _, sk := range secretKeys {
		description := ""
		if sk.Description != nil {
			description = *sk.Description
		}

		expiresAt := ""
		if sk.ExpiresAt != nil {
			expiresAt = sk.ExpiresAt.Format(time.RFC3339)
		}

		record := []string{
			fmt.Sprintf("%d", sk.ID),
			sk.Name,
			string(sk.KeyType),
			description,
			sk.CreatedAt.Format(time.RFC3339),
			sk.UpdatedAt.Format(time.RFC3339),
			expiresAt,
		}

		if err := writer.Write(record); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	filename = fmt.Sprintf("secret-keys-%s.csv", time.Now().Format("20060102"))
	return buf.Bytes(), filename, nil
}

// exportToJSON exports secret keys to JSON format
func (s *secretKeyService) exportToJSON(secretKeys []*model.SecretKey) (data []byte, filename string, err error) {
	responses := make([]*response.SecretKeyResponse, 0, len(secretKeys))
	for _, sk := range secretKeys {
		responses = append(responses, response.ToSecretKeyResponse(sk))
	}

	data, err = json.MarshalIndent(responses, "", "  ")
	if err != nil {
		return nil, "", err
	}

	filename = fmt.Sprintf("secret-keys-%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}

// UploadSecretFile uploads a secret key file
func (s *secretKeyService) UploadSecretFile(ctx context.Context, file *multipart.FileHeader, fileType string, userID uint) (*response.SecretFileUploadResponse, error) {
	s.logger.InfoContext(ctx, "Uploading secret file",
		logger.String("filename", file.Filename),
		logger.String("type", fileType),
		logger.Uint("user_id", userID))

	// TODO: 调用 common 文件服务上传文件
	// 1. 验证文件类型
	// 2. 验证文件大小
	// 3. 生成唯一文件名
	// 4. 保存文件到存储路径
	// 5. 返回文件信息

	// 占位符实现
	return &response.SecretFileUploadResponse{
		Filename:     "placeholder-" + file.Filename,
		OriginalName: file.Filename,
		FilePath:     "/home/appuser/data/placeholder-" + file.Filename,
	}, nil
}

// DownloadSecretFile downloads a secret key file
func (s *secretKeyService) DownloadSecretFile(ctx context.Context, filename string, userID uint) (filePath, originalName string, err error) {
	s.logger.InfoContext(ctx, "Downloading secret file",
		logger.String("filename", filename),
		logger.Uint("user_id", userID))

	// TODO: 调用 common 文件服务下载文件
	// 1. 验证文件存在
	// 2. 验证用户权限
	// 3. 获取文件路径和原始文件名
	// 4. 返回文件信息

	// 占位符实现
	return "/home/appuser/data/" + filename, filename, nil
}

// DeleteSecretFile deletes a secret key file
func (s *secretKeyService) DeleteSecretFile(ctx context.Context, filename string, userID uint) error {
	s.logger.InfoContext(ctx, "Deleting secret file",
		logger.String("filename", filename),
		logger.Uint("user_id", userID))

	// TODO: 调用 common 文件服务删除文件
	// 1. 验证文件存在
	// 2. 验证用户权限（检查是否为文件所有者）
	// 3. 删除文件
	// 4. 记录审计日志

	// 占位符实现
	return nil
}

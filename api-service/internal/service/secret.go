package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
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

	encryptedValue, err := s.rsaCrypto.EncryptString(req.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt secret key value",
			logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeEncryptFailed)
	}

	// Create the secret key model
	secretKey := &model.SecretKey{
		Name:            req.Name,
		KeyType:         req.KeyType,
		EncryptedValue:  encryptedValue, // 使用加密后的值
		Description:     req.Description,
		CustomFields:    req.CustomFields,
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

	if len(req.AuthorizedUsers) > 0 {
		for _, authorizedUserID := range req.AuthorizedUsers {
			userSecret := &model.UserSecret{
				UserID:      authorizedUserID,
				SecretKeyID: secretKey.ID,
				GrantedBy:   &userID,
				GrantedAt:   time.Now(),
				ExpiresAt:   req.ExpiresAt,
			}

			if err := s.secretKeyRepo.CreateUserSecret(ctx, userSecret); err != nil {
				s.logger.ErrorContext(ctx, "Failed to create user secret relationship",
					logger.Uint("secret_key_id", secretKey.ID),
					logger.Uint("authorized_user_id", authorizedUserID),
					logger.ErrorField(err))
				continue
			}

			s.logger.InfoContext(ctx, "User secret relationship created",
				logger.Uint("secret_key_id", secretKey.ID),
				logger.Uint("authorized_user_id", authorizedUserID))
		}

		s.logger.InfoContext(ctx, "User secret relationships creation completed",
			logger.Uint("secret_key_id", secretKey.ID),
			logger.Int("total_users", len(req.AuthorizedUsers)))
	}

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
		// If not owner, check if user has access permission in user_secret table
		hasAccess, accessErr := s.secretKeyRepo.CheckUserSecretAccess(ctx, userID, id)
		if accessErr != nil {
			s.logger.ErrorContext(ctx, "Failed to check user secret access",
				logger.Uint("secret_key_id", id),
				logger.Uint("user_id", userID),
				logger.ErrorField(accessErr))
			return nil, err
		}

		if !hasAccess {
			return nil, errors.NewAppError(errors.CodeAccessDenied)
		}
	}

	// Check if expired
	if secretKey.IsExpired() {
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	s.logger.InfoContext(ctx, "Secret key value accessed",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

	decryptedValue, err := s.rsaCrypto.DecryptString(secretKey.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to decrypt secret key value",
			logger.ErrorField(err))
		return nil, err
	}

	return response.ToSecretKeyValueResponse(decryptedValue, secretKey.ExpiresAt), nil
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

	// Encrypt the secret value using RSA
	encryptedValue, err := s.rsaCrypto.EncryptString(req.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt secret key value",
			logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeEncryptFailed)
	}

	// Update fields based on the simplified SecretKeyUpdateRequest
	secretKey.EncryptedValue = encryptedValue
	secretKey.KeyType = req.KeyType

	// Save changes
	if err := s.secretKeyRepo.Update(ctx, secretKey); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Secret key updated successfully",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

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

	return response.ToSecretKeyListResponse(secretKeys, total, req.Page, req.PageSize), nil
}

// ExportSecretKeys exports secret keys in specified format
func (s *secretKeyService) ExportSecretKeys(ctx context.Context, req *request.SecretKeyExportRequest, userID uint) (data []byte, filename string, err error) {
	s.logger.InfoContext(ctx, "Exporting secret keys",
		logger.Uint("user_id", userID),
		logger.String("format", req.Format))

	// Create query request for filtering
	queryReq := &request.SecretKeyQueryRequest{
		Page:     1,
		PageSize: constants.DefaultExportPageSize, // Large number to get all keys
	}

	secretKeys, _, err := s.secretKeyRepo.List(ctx, queryReq, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get secret keys for export",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, "", err
	}

	switch req.Format {
	case constants.ExportFormatCsv:
		return s.exportToCSV(secretKeys)
	case constants.ExportFormatJson:
		return s.exportToJSON(secretKeys)
	case constants.ExportFormatExcel:
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

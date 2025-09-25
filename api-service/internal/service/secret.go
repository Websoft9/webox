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
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"

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
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.create_failed"))
	}

	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists, s.i18n.T(ctx, "secret.name_already_exists"))
	}

	encryptedValue, err := s.rsaCrypto.EncryptString(req.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt secret key value",
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeEncryptFailed, s.i18n.T(ctx, "business.encrypt_failed"))
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
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, s.i18n.T(ctx, "secret.create_failed"))
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
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "secret.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.get_failed"))
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "secret.access_denied"))
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
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "secret.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.get_failed"))
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "secret.access_denied"))
	}

	// Check if expired
	if secretKey.IsExpired() {
		return nil, errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "secret.expired"))
	}

	s.logger.InfoContext(ctx, "Secret key value accessed",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

	decryptedValue, err := s.rsaCrypto.DecryptString(secretKey.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to decrypt secret key value",
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeDecryptFailed, s.i18n.T(ctx, "business.decrypt_failed"))
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
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "secret.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.update_failed"))
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "secret.access_denied"))
	}

	// Encrypt the secret value using RSA
	encryptedValue, err := s.rsaCrypto.EncryptString(req.EncryptedValue)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to encrypt secret key value",
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeEncryptFailed, s.i18n.T(ctx, "business.encrypt_failed"))
	}

	// Update fields based on the simplified SecretKeyUpdateRequest
	secretKey.EncryptedValue = encryptedValue
	secretKey.KeyType = req.KeyType

	// Save changes
	if err := s.secretKeyRepo.Update(ctx, secretKey); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "secret.update_failed"))
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
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "secret.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.delete_failed"))
	}

	// Check ownership
	if secretKey.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "secret.access_denied"))
	}

	if err := s.secretKeyRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete secret key",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordDeleteFailed, s.i18n.T(ctx, "secret.delete_failed"))
	}

	s.logger.InfoContext(ctx, "Secret key deleted successfully",
		logger.Uint("secret_key_id", id),
		logger.Uint("user_id", userID))

	return nil
}

// ListSecretKeys retrieves secret keys with pagination and filtering
func (s *secretKeyService) ListSecretKeys(ctx context.Context, req *request.SecretKeyQueryRequest, userID uint) (*response.SecretKeyListResponse, error) {
	s.logger.InfoContext(ctx, "Listing secret keys",
		logger.Uint("user_id", userID),
		logger.Int("page", req.Page),
		logger.Int("page_size", req.PageSize))

	secretKeys, total, err := s.secretKeyRepo.List(ctx, req, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list secret keys",
			logger.Uint("user_id", userID),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.list_failed"))
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
		return nil, "", errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.export_failed"))
	}

	switch req.Format {
	case "csv":
		return s.exportToCSV(secretKeys)
	case "json":
		return s.exportToJSON(secretKeys)
	default:
		return nil, "", errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "secret.invalid_export_format"))
	}
}

// ValidateSecretKeyOwnership checks if user owns the secret key
func (s *secretKeyService) ValidateSecretKeyOwnership(ctx context.Context, secretKeyID, userID uint) error {
	secretKey, err := s.secretKeyRepo.GetByID(ctx, secretKeyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "secret.not_found"))
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "secret.validation_failed"))
	}

	if secretKey.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "secret.access_denied"))
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
		return nil, "", errors.NewAppError(errors.CodeInternalError, "Failed to export CSV")
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
			return nil, "", errors.NewAppError(errors.CodeInternalError, "Failed to export CSV")
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", errors.NewAppError(errors.CodeInternalError, "Failed to export CSV")
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
		return nil, "", errors.NewAppError(errors.CodeInternalError, "Failed to export JSON")
	}

	filename = fmt.Sprintf("secret-keys-%s.json", time.Now().Format("20060102"))
	return data, filename, nil
}

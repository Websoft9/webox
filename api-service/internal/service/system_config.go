package service

import (
	"api-service/internal/config"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/email"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type SystemConfigService struct {
	systemConfigRepo repository.SystemConfigRepository
	config           *config.Config
	emailService     email.EmailService
	db               *gorm.DB
	logger           logger.Logger
	i18n             *i18n.I18n
}

func NewSystemConfigService(
	systemConfigRepo repository.SystemConfigRepository,
	config *config.Config,
	db *gorm.DB,
	logger logger.Logger,
	i18n *i18n.I18n,
) *SystemConfigService {
	var emailService email.EmailService
	if config.Email.SMTP.Host != "" && config.Email.SMTP.Username != "" {
		emailService = email.NewEmailService(config, logger, i18n)
	}
	return &SystemConfigService{
		systemConfigRepo: systemConfigRepo,
		emailService:     emailService,
		config:           config,
		db:               db,
		logger:           logger,
		i18n:             i18n,
	}
}

// ListSystemConfigs lists system configurations with pagination and filtering
func (s *SystemConfigService) ListSystemConfigs(ctx context.Context, req *request.ListSystemConfigsRequest) (*response.ListSystemConfigsResponse, error) {
	filter := &request.ListSystemConfigsFilter{
		Category: req.Category,
		Keyword:  req.Keyword,
	}

	systemConfigs, err := s.systemConfigRepo.List(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list system configs", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to list system configs")
	}

	// Decrypt encrypted values
	if decryptErr := s.decryptConfigList(systemConfigs); decryptErr != nil {
		s.logger.WarnContext(ctx, "some configs could not be decrypted", logger.ErrorField(decryptErr))
	}

	respItems := make([]response.SystemConfigsResponse, 0, len(systemConfigs))
	for _, config := range systemConfigs {
		respItems = append(respItems, response.SystemConfigsResponse{
			ID:           config.ID,
			ConfigKey:    config.ConfigKey,
			ConfigValue:  config.ConfigValue,
			ConfigType:   string(config.ConfigType),
			Category:     config.Category,
			Description:  config.Description,
			IsReadonly:   config.IsReadonly,
			IsEncrypted:  config.IsEncrypted,
			DefaultValue: config.DefaultValue,
			SortOrder:    config.SortOrder,
		})
	}

	return &response.ListSystemConfigsResponse{
		Items: respItems,
	}, nil
}

// TestSMTP tests the SMTP configuration by sending a test email
func (s *SystemConfigService) TestSMTP(ctx context.Context, req *request.TestSMTPRequest) error {
	if s.emailService == nil {
		return errors.NewAppError(errors.CodeRecordCreateFailed, "SMTP service not configured")
	}

	err := s.emailService.SendEmail(ctx, req.TestEmail, "Test Email", "This is a test email from Websoft9.")
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to send test email", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to send test email")
	}

	return nil
}

// BatchUpdateSystemConfigs batch updates system configurations
func (s *SystemConfigService) BatchUpdateSystemConfigs(ctx context.Context, req *request.BatchUpdateSystemConfigsRequest) error {
	for _, config := range req.Configs {
		// Check if config exists and get readonly status
		origin, err := s.systemConfigRepo.GetByKey(ctx, config.ConfigKey)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to get system config by key",
				logger.String("config_key", config.ConfigKey), logger.ErrorField(err))
			continue // Skip non-existent configs
		}

		// First decrypt the original value to check if it's readonly
		if decryptErr := s.decryptConfigValue(origin); decryptErr != nil {
			s.logger.ErrorContext(ctx, "failed to decrypt original config value", logger.ErrorField(decryptErr))
			continue // Skip configs that can't be decrypted
		}

		if origin.IsReadonly {
			s.logger.WarnContext(ctx, "skipping readonly config",
				logger.String("config_key", config.ConfigKey))
			continue // Skip readonly configs
		}

		// Update the values
		origin.ConfigValue = config.ConfigValue
		origin.ConfigType = model.ConfigType(config.ConfigType)
		origin.Category = config.Category
		// Only update optional fields if they are provided
		if config.Description != nil {
			origin.Description = *config.Description
		}
		if config.SortOrder != nil {
			origin.SortOrder = *config.SortOrder
		}

		// Encrypt the value before saving if needed
		if encryptErr := s.encryptConfigValue(origin); encryptErr != nil {
			s.logger.ErrorContext(ctx, "failed to encrypt config value",
				logger.String("config_key", config.ConfigKey), logger.ErrorField(encryptErr))
			continue // Skip configs that can't be encrypted
		}

		err = s.systemConfigRepo.Update(ctx, origin)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to update system config",
				logger.String("config_key", config.ConfigKey), logger.ErrorField(err))
			continue // Continue with other configs even if one fails
		}
	}

	return nil
}

// encryptConfigValue encrypts the config value if needed
func (s *SystemConfigService) encryptConfigValue(config *model.SystemConfig) error {
	if !config.IsEncrypted || config.ConfigValue == "" {
		return nil
	}

	cryptoInstance := crypto.GetDefaultCrypto()
	if cryptoInstance == nil {
		return fmt.Errorf("crypto instance not initialized")
	}

	encryptedValue, err := cryptoInstance.Encrypt(config.ConfigValue)
	if err != nil {
		return fmt.Errorf("failed to encrypt config value: %w", err)
	}

	config.ConfigValue = encryptedValue
	return nil
}

// decryptConfigValue decrypts the config value if needed
func (s *SystemConfigService) decryptConfigValue(config *model.SystemConfig) error {
	if !config.IsEncrypted || config.ConfigValue == "" {
		return nil
	}

	cryptoInstance := crypto.GetDefaultCrypto()
	if cryptoInstance == nil {
		// If crypto is not initialized, just return without error (could be test environment)
		return nil
	}

	decryptedValue, err := cryptoInstance.Decrypt(config.ConfigValue)
	if err != nil {
		return fmt.Errorf("failed to decrypt config value: %w", err)
	}

	config.ConfigValue = decryptedValue
	return nil
}

// decryptConfigList decrypts a list of configs
func (s *SystemConfigService) decryptConfigList(configs []*model.SystemConfig) error {
	for _, config := range configs {
		if err := s.decryptConfigValue(config); err != nil {
			// Log the error but continue with other configs
			s.logger.WarnContext(context.Background(), "failed to decrypt config",
				logger.ErrorField(err))
			return err
		}
	}
	return nil
}

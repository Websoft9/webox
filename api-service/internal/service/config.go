package service

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	redisClient "api-service/pkg/redis"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	// Redis cache expiration time for user profile and permissions (24 hours)
	cacheExpirationHours = 24
	// String constants for boolean values in Redis
	stringTrue = "true"
	stringOne  = "1"
	// Multiplier for converting map to slice (key-value pairs)
	keyValuePairMultiplier = 2
)

// configService implements the ConfigService interface
type configService struct {
	systemConfigRepo  repository.SystemConfigRepository
	serviceConfigRepo repository.ServiceConfigRepository
	configRepo        repository.ConfigRepository
	permissionRepo    repository.PermissionRepository
	roleRepo          repository.RoleRepository
	db                *gorm.DB
	logger            logger.Logger
}

// NewConfigService creates a new ConfigService instance
func NewConfigService(
	systemConfigRepo repository.SystemConfigRepository,
	serviceConfigRepo repository.ServiceConfigRepository,
	configRepo repository.ConfigRepository,
	permissionRepo repository.PermissionRepository,
	roleRepo repository.RoleRepository,
	db *gorm.DB,
	logger logger.Logger,
) *configService {
	return &configService{
		systemConfigRepo:  systemConfigRepo,
		serviceConfigRepo: serviceConfigRepo,
		configRepo:        configRepo,
		permissionRepo:    permissionRepo,
		roleRepo:          roleRepo,
		db:                db,
		logger:            logger,
	}
}

// ==================== System Configuration Operations ====================

// GetSystemConfig retrieves a system configuration by key
// Priority: Redis -> Database -> Config File -> Constants
func (s *configService) GetSystemConfig(ctx context.Context, configKey string) (*model.SystemConfig, error) {
	// Try to get from Redis first
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SYSTEM, configKey)
	configData, err := redisClient.HGetAll(ctx, redisKey)
	if err == nil && len(configData) > 0 {
		// Found in Redis, convert to model
		config := s.hashToSystemConfig(configData)
		if decryptErr := s.decryptSystemConfigValue(config); decryptErr != nil {
			s.logger.WarnContext(ctx, "failed to decrypt config from Redis",
				logger.String("config_key", configKey), logger.ErrorField(decryptErr))
		}
		return config, nil
	}

	// Not in Redis, try database
	config, err := s.systemConfigRepo.GetByKey(ctx, configKey)
	if err != nil {
		if errors.Is(err, errors.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeConfigNotFound)
		}
		return nil, err
	}

	// Decrypt if needed
	if decryptErr := s.decryptSystemConfigValue(config); decryptErr != nil {
		s.logger.WarnContext(ctx, "failed to decrypt config from database",
			logger.String("config_key", configKey), logger.ErrorField(decryptErr))
	}

	// Write to Redis cache
	if cacheErr := s.cacheSystemConfig(ctx, config); cacheErr != nil {
		s.logger.WarnContext(ctx, "failed to cache system config",
			logger.String("config_key", configKey), logger.ErrorField(cacheErr))
	}

	return config, nil
}

// CreateSystemConfig creates a new system configuration
func (s *configService) CreateSystemConfig(ctx context.Context, config *model.SystemConfig) error {
	// Check if config already exists
	existing, err := s.systemConfigRepo.GetByKey(ctx, config.ConfigKey)
	if err == nil && existing != nil {
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Encrypt if needed
	if encryptErr := s.encryptSystemConfigValue(config); encryptErr != nil {
		return errors.NewAppErrorWrapError(encryptErr, errors.CodeEncryptFailed)
	}

	// Create in database
	if err := s.systemConfigRepo.Create(ctx, config); err != nil {
		return err
	}

	// Sync to Redis
	if cacheErr := s.cacheSystemConfig(ctx, config); cacheErr != nil {
		s.logger.WarnContext(ctx, "failed to cache new system config",
			logger.String("config_key", config.ConfigKey), logger.ErrorField(cacheErr))
	}

	s.logger.InfoContext(ctx, "system config created successfully",
		logger.String("config_key", config.ConfigKey))
	return nil
}

// UpdateSystemConfig updates an existing system configuration
func (s *configService) UpdateSystemConfig(ctx context.Context, configKey string, config *model.SystemConfig) error {
	// Check if config exists
	existing, err := s.systemConfigRepo.GetByKey(ctx, configKey)
	if err != nil {
		return err
	}

	// Check if readonly
	if existing.IsReadonly {
		return errors.NewAppError(errors.CodeConfigReadonly)
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Encrypt if needed
	if encryptErr := s.encryptSystemConfigValue(config); encryptErr != nil {
		tx.Rollback()
		return errors.NewAppErrorWrapError(encryptErr, errors.CodeEncryptFailed)
	}

	// Update in database
	config.ID = existing.ID
	if err := s.systemConfigRepo.Update(ctx, config); err != nil {
		tx.Rollback()
		return err
	}

	// Sync to Redis
	if cacheErr := s.cacheSystemConfig(ctx, config); cacheErr != nil {
		tx.Rollback()
		return errors.NewAppError(errors.CodeConfigSyncFailed)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.logger.InfoContext(ctx, "system config updated successfully",
		logger.String("config_key", configKey))
	return nil
}

// DeleteSystemConfig deletes a system configuration by key
func (s *configService) DeleteSystemConfig(ctx context.Context, configKey string) error {
	// Check if config exists
	existing, err := s.systemConfigRepo.GetByKey(ctx, configKey)
	if err != nil {
		return err
	}

	// Check if readonly
	if existing.IsReadonly {
		return errors.NewAppError(errors.CodeConfigReadonly)
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete from database
	if err := s.systemConfigRepo.Delete(ctx, existing.ID); err != nil {
		tx.Rollback()
		return err
	}

	// Delete from Redis
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SYSTEM, configKey)
	if _, err := redisClient.Del(ctx, redisKey); err != nil {
		tx.Rollback()
		return errors.NewAppError(errors.CodeConfigSyncFailed)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.logger.InfoContext(ctx, "system config deleted successfully",
		logger.String("config_key", configKey))
	return nil
}

// ==================== Service Configuration Operations ====================

// GetServiceConfig retrieves a service configuration by code
// Priority: Redis -> Database -> Config File -> Constants
func (s *configService) GetServiceConfig(ctx context.Context, code string) (*model.ServiceConfig, error) {
	// Try to get from Redis first
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SERVICE, code)
	configData, err := redisClient.HGetAll(ctx, redisKey)
	if err == nil && len(configData) > 0 {
		// Found in Redis, convert to model
		config := s.hashToServiceConfig(configData)
		if decryptErr := s.decryptServiceConfigValue(config); decryptErr != nil {
			s.logger.WarnContext(ctx, "failed to decrypt service config from Redis",
				logger.String("code", code), logger.ErrorField(decryptErr))
		}
		return config, nil
	}

	// Not in Redis, try database
	config, err := s.serviceConfigRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, errors.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeConfigNotFound)
		}
		return nil, err
	}

	// Decrypt if needed
	if decryptErr := s.decryptServiceConfigValue(config); decryptErr != nil {
		s.logger.WarnContext(ctx, "failed to decrypt service config from database",
			logger.String("code", code), logger.ErrorField(decryptErr))
	}

	// Write to Redis cache
	if cacheErr := s.cacheServiceConfig(ctx, config); cacheErr != nil {
		s.logger.WarnContext(ctx, "failed to cache service config",
			logger.String("code", code), logger.ErrorField(cacheErr))
	}

	return config, nil
}

// CreateServiceConfig creates a new service configuration
func (s *configService) CreateServiceConfig(ctx context.Context, config *model.ServiceConfig) error {
	// Check if config already exists
	existing, err := s.serviceConfigRepo.GetByCode(ctx, config.Code)
	if err == nil && existing != nil {
		return errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Encrypt if needed
	if encryptErr := s.encryptServiceConfigValue(config); encryptErr != nil {
		return errors.NewAppErrorWrapError(encryptErr, errors.CodeEncryptFailed)
	}

	// Create in database
	if err := s.serviceConfigRepo.Create(ctx, config); err != nil {
		return err
	}

	// Sync to Redis
	if cacheErr := s.cacheServiceConfig(ctx, config); cacheErr != nil {
		s.logger.WarnContext(ctx, "failed to cache new service config",
			logger.String("code", config.Code), logger.ErrorField(cacheErr))
	}

	s.logger.InfoContext(ctx, "service config created successfully",
		logger.String("code", config.Code))
	return nil
}

// UpdateServiceConfig updates an existing service configuration
func (s *configService) UpdateServiceConfig(ctx context.Context, code string, config *model.ServiceConfig) error {
	// Check if config exists
	existing, err := s.serviceConfigRepo.GetByCode(ctx, code)
	if err != nil {
		return err
	}

	// Check if readonly
	if existing.IsReadonly {
		return errors.NewAppError(errors.CodeConfigReadonly)
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Encrypt if needed
	if encryptErr := s.encryptServiceConfigValue(config); encryptErr != nil {
		tx.Rollback()
		return errors.NewAppErrorWrapError(encryptErr, errors.CodeEncryptFailed)
	}

	// Update in database
	config.ID = existing.ID
	if err := s.serviceConfigRepo.Update(ctx, config); err != nil {
		tx.Rollback()
		return err
	}

	// Sync to Redis
	if cacheErr := s.cacheServiceConfig(ctx, config); cacheErr != nil {
		tx.Rollback()
		return errors.NewAppError(errors.CodeConfigSyncFailed)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.logger.InfoContext(ctx, "service config updated successfully",
		logger.String("code", code))
	return nil
}

// DeleteServiceConfig deletes a service configuration by code
func (s *configService) DeleteServiceConfig(ctx context.Context, code string) error {
	// Check if config exists
	existing, err := s.serviceConfigRepo.GetByCode(ctx, code)
	if err != nil {
		return err
	}

	// Check if readonly
	if existing.IsReadonly {
		return errors.NewAppError(errors.CodeConfigReadonly)
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete from database
	if err := s.serviceConfigRepo.Delete(ctx, existing.ID); err != nil {
		tx.Rollback()
		return err
	}

	// Delete from Redis
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SERVICE, code)
	if _, err := redisClient.Del(ctx, redisKey); err != nil {
		tx.Rollback()
		return errors.NewAppError(errors.CodeConfigSyncFailed)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	s.logger.InfoContext(ctx, "service config deleted successfully",
		logger.String("code", code))
	return nil
}

// ==================== User Profile Configuration Operations ====================

// GetUserProfile retrieves user profile configuration
// Priority: Redis -> Database
func (s *configService) GetUserProfile(ctx context.Context, userID uint) (map[string]interface{}, error) {
	// Try to get from Redis first
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PROFILE, userID)
	profileData, err := redisClient.HGetAll(ctx, redisKey)
	if err == nil && len(profileData) > 0 {
		// Found in Redis, convert to map
		result := make(map[string]interface{})
		for k, v := range profileData {
			result[k] = v
		}
		return result, nil
	}

	// Not in Redis, try database
	profiles, err := s.configRepo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to map
	result := make(map[string]interface{})
	for _, profile := range profiles {
		result[profile.ConfigKey] = profile.GetEffectiveValue()
	}

	// Write to Redis cache
	if len(result) > 0 {
		if cacheErr := s.cacheUserProfile(ctx, userID, result); cacheErr != nil {
			s.logger.WarnContext(ctx, "failed to cache user profile",
				logger.Uint("user_id", userID), logger.ErrorField(cacheErr))
		}
	}

	return result, nil
}

// UpdateUserProfile updates user profile configuration
func (s *configService) UpdateUserProfile(ctx context.Context, userID uint, configKey string, configValue interface{}) error {
	// Check if profile exists
	existing, err := s.configRepo.GetUserProfileByKey(ctx, userID, configKey)
	if err != nil && !errors.Is(err, errors.ErrRecordNotFound) {
		return err
	}

	// Convert value to string
	valueStr := fmt.Sprintf("%v", configValue)

	if existing == nil {
		// Create new profile
		profile := &model.UserProfile{
			UserID:      userID,
			ConfigKey:   configKey,
			ConfigValue: valueStr,
			Category:    "general",
		}
		if err := s.configRepo.CreateUserProfile(ctx, profile); err != nil {
			return err
		}
	} else {
		// Update existing profile
		existing.ConfigValue = valueStr
		if err := s.configRepo.UpdateUserProfile(ctx, existing); err != nil {
			return err
		}
	}

	// Update Redis cache
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PROFILE, userID)
	if _, err := redisClient.HSet(ctx, redisKey, configKey, valueStr); err != nil {
		s.logger.WarnContext(ctx, "failed to update user profile in Redis",
			logger.Uint("user_id", userID), logger.String("config_key", configKey), logger.ErrorField(err))
	}

	s.logger.InfoContext(ctx, "user profile updated successfully",
		logger.Uint("user_id", userID), logger.String("config_key", configKey))
	return nil
}

// SyncUserProfileOnLogin synchronizes user profile to Redis when user logs in
func (s *configService) SyncUserProfileOnLogin(ctx context.Context, userID uint) error {
	// Get user profile from database
	profiles, err := s.configRepo.GetUserProfile(ctx, userID)
	if err != nil {
		return err
	}

	// Convert to map
	profileMap := make(map[string]interface{})
	for _, profile := range profiles {
		profileMap[profile.ConfigKey] = profile.GetEffectiveValue()
	}

	// Write to Redis with expiration (24 hours)
	if err := s.cacheUserProfile(ctx, userID, profileMap); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "user profile synced on login",
		logger.Uint("user_id", userID))
	return nil
}

// ==================== User Permission Configuration Operations ====================

// GetUserPermissions retrieves user permissions from Redis or database
// Priority: Redis -> Database
func (s *configService) GetUserPermissions(ctx context.Context, userID uint) (map[string]string, error) {
	// Try to get from Redis first
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PERMISSION, userID)
	permData, err := redisClient.HGetAll(ctx, redisKey)
	if err == nil && len(permData) > 0 {
		return permData, nil
	}

	// Not in Redis, query from database
	permissions, err := s.permissionRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert to map (resource -> action)
	permMap := make(map[string]string)
	for _, perm := range permissions {
		permMap[perm.Resource] = perm.Action
	}

	// Write to Redis cache
	if len(permMap) > 0 {
		if cacheErr := s.cacheUserPermissions(ctx, userID, permMap); cacheErr != nil {
			s.logger.WarnContext(ctx, "failed to cache user permissions",
				logger.Uint("user_id", userID), logger.ErrorField(cacheErr))
		}
	}

	return permMap, nil
}

// SyncUserPermissionsOnLogin synchronizes user permissions to Redis when user logs in
func (s *configService) SyncUserPermissionsOnLogin(ctx context.Context, userID uint) error {
	// Get user permissions from database
	permissions, err := s.permissionRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return err
	}

	// Convert to map
	permMap := make(map[string]string)
	for _, perm := range permissions {
		permMap[perm.Resource] = perm.Action
	}

	// Write to Redis with expiration (24 hours)
	if err := s.cacheUserPermissions(ctx, userID, permMap); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "user permissions synced on login",
		logger.Uint("user_id", userID))
	return nil
}

// SyncUserPermissionsOnChange synchronizes user permissions to Redis when permissions change
func (s *configService) SyncUserPermissionsOnChange(ctx context.Context, userID uint) error {
	// Delete old permissions from Redis
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PERMISSION, userID)
	if _, err := redisClient.Del(ctx, redisKey); err != nil {
		s.logger.WarnContext(ctx, "failed to delete old permissions from Redis",
			logger.Uint("user_id", userID), logger.ErrorField(err))
	}

	// Sync new permissions
	return s.SyncUserPermissionsOnLogin(ctx, userID)
}

// ==================== Configuration Preload Operations ====================

// PreloadConfigs preloads all configurations to Redis on system startup
func (s *configService) PreloadConfigs(ctx context.Context) error {
	startTime := time.Now()
	successCount := 0
	failCount := 0

	s.logger.InfoContext(ctx, "starting configuration preload")

	// Preload system configurations
	systemConfigs, err := s.systemConfigRepo.List(ctx, nil)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to load system configs for preload", logger.ErrorField(err))
		failCount++
	} else {
		for _, config := range systemConfigs {
			if cacheErr := s.cacheSystemConfig(ctx, config); cacheErr != nil {
				s.logger.WarnContext(ctx, "failed to preload system config",
					logger.String("config_key", config.ConfigKey), logger.ErrorField(cacheErr))
				failCount++
			} else {
				successCount++
			}
		}
	}

	// Preload service configurations
	serviceConfigs, err := s.serviceConfigRepo.List(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to load service configs for preload", logger.ErrorField(err))
		failCount++
	} else {
		for _, config := range serviceConfigs {
			if cacheErr := s.cacheServiceConfig(ctx, config); cacheErr != nil {
				s.logger.WarnContext(ctx, "failed to preload service config",
					logger.String("code", config.Code), logger.ErrorField(cacheErr))
				failCount++
			} else {
				successCount++
			}
		}
	}

	duration := time.Since(startTime)
	s.logger.InfoContext(ctx, "configuration preload completed",
		logger.Int("success_count", successCount),
		logger.Int("fail_count", failCount),
		logger.String("duration", duration.String()))

	if failCount > 0 {
		return errors.NewAppError(errors.CodeConfigPreloadFailed)
	}
	return nil
}

// ==================== Helper Methods ====================

// cacheSystemConfig caches system configuration to Redis
func (s *configService) cacheSystemConfig(ctx context.Context, config *model.SystemConfig) error {
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SYSTEM, config.ConfigKey)
	configMap := s.systemConfigToHash(config)
	_, err := redisClient.HSet(ctx, redisKey, configMap...)
	return err
}

// cacheServiceConfig caches service configuration to Redis
func (s *configService) cacheServiceConfig(ctx context.Context, config *model.ServiceConfig) error {
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_SERVICE, config.Code)
	configMap := s.serviceConfigToHash(config)
	_, err := redisClient.HSet(ctx, redisKey, configMap...)
	return err
}

// cacheUserProfile caches user profile to Redis with 24-hour expiration
func (s *configService) cacheUserProfile(ctx context.Context, userID uint, profileMap map[string]interface{}) error {
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PROFILE, userID)

	// Convert map to slice for HSet
	args := make([]interface{}, 0, len(profileMap)*keyValuePairMultiplier)
	for k, v := range profileMap {
		args = append(args, k, v)
	}

	if _, err := redisClient.HSet(ctx, redisKey, args...); err != nil {
		return err
	}

	// Set expiration to 24 hours
	return redisClient.Expire(ctx, redisKey, cacheExpirationHours*time.Hour)
}

// cacheUserPermissions caches user permissions to Redis with 24-hour expiration
func (s *configService) cacheUserPermissions(ctx context.Context, userID uint, permMap map[string]string) error {
	redisKey := fmt.Sprintf(redisClient.RK_CONFIG_PERMISSION, userID)

	// Convert map to slice for HSet
	args := make([]interface{}, 0, len(permMap)*keyValuePairMultiplier)
	for k, v := range permMap {
		args = append(args, k, v)
	}

	if _, err := redisClient.HSet(ctx, redisKey, args...); err != nil {
		return err
	}

	// Set expiration to 24 hours
	return redisClient.Expire(ctx, redisKey, cacheExpirationHours*time.Hour)
}

// systemConfigToHash converts SystemConfig to Redis hash format
func (s *configService) systemConfigToHash(config *model.SystemConfig) []interface{} {
	return []interface{}{
		"id", config.ID,
		"config_key", config.ConfigKey,
		"config_value", config.ConfigValue,
		"config_type", config.ConfigType,
		"category", config.Category,
		"description", config.Description,
		"is_readonly", config.IsReadonly,
		"is_encrypted", config.IsEncrypted,
	}
}

// hashToSystemConfig converts Redis hash to SystemConfig
func (s *configService) hashToSystemConfig(data map[string]string) *model.SystemConfig {
	config := &model.SystemConfig{
		ConfigKey:   data["config_key"],
		ConfigValue: data["config_value"],
		ConfigType:  model.ConfigType(data["config_type"]),
		Category:    data["category"],
		Description: data["description"],
	}

	if data["is_readonly"] == stringTrue || data["is_readonly"] == stringOne {
		config.IsReadonly = true
	}
	if data["is_encrypted"] == stringTrue || data["is_encrypted"] == stringOne {
		config.IsEncrypted = true
	}

	return config
}

// serviceConfigToHash converts ServiceConfig to Redis hash format
func (s *configService) serviceConfigToHash(config *model.ServiceConfig) []interface{} {
	return []interface{}{
		"id", config.ID,
		"code", config.Code,
		"config_key", config.ConfigKey,
		"config_value", config.ConfigValue,
		"config_type", config.ConfigType,
		"category", config.Category,
		"description", config.Description,
		"is_readonly", config.IsReadonly,
		"is_encrypted", config.IsEncrypted,
		"owner_id", config.OwnerID,
	}
}

// hashToServiceConfig converts Redis hash to ServiceConfig
func (s *configService) hashToServiceConfig(data map[string]string) *model.ServiceConfig {
	config := &model.ServiceConfig{
		Code:        data["code"],
		ConfigKey:   data["config_key"],
		ConfigValue: data["config_value"],
		ConfigType:  model.ConfigType(data["config_type"]),
		Category:    data["category"],
		Description: data["description"],
	}

	if data["is_readonly"] == stringTrue || data["is_readonly"] == stringOne {
		config.IsReadonly = true
	}
	if data["is_encrypted"] == stringTrue || data["is_encrypted"] == stringOne {
		config.IsEncrypted = true
	}

	return config
}

// encryptSystemConfigValue encrypts the config value if needed
func (s *configService) encryptSystemConfigValue(config *model.SystemConfig) error {
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

// decryptSystemConfigValue decrypts the config value if needed
func (s *configService) decryptSystemConfigValue(config *model.SystemConfig) error {
	if !config.IsEncrypted || config.ConfigValue == "" {
		return nil
	}

	cryptoInstance := crypto.GetDefaultCrypto()
	if cryptoInstance == nil {
		return nil // Skip decryption if crypto not initialized
	}

	decryptedValue, err := cryptoInstance.Decrypt(config.ConfigValue)
	if err != nil {
		return fmt.Errorf("failed to decrypt config value: %w", err)
	}

	config.ConfigValue = decryptedValue
	return nil
}

// encryptServiceConfigValue encrypts the service config value if needed
func (s *configService) encryptServiceConfigValue(config *model.ServiceConfig) error {
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

// decryptServiceConfigValue decrypts the service config value if needed
func (s *configService) decryptServiceConfigValue(config *model.ServiceConfig) error {
	if !config.IsEncrypted || config.ConfigValue == "" {
		return nil
	}

	cryptoInstance := crypto.GetDefaultCrypto()
	if cryptoInstance == nil {
		return nil // Skip decryption if crypto not initialized
	}

	decryptedValue, err := cryptoInstance.Decrypt(config.ConfigValue)
	if err != nil {
		return fmt.Errorf("failed to decrypt config value: %w", err)
	}

	config.ConfigValue = decryptedValue
	return nil
}

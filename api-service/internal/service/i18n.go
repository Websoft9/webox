package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"context"
)

// i18nService implements the I18nService interface
type i18nService struct {
	userProfileRepo repository.UserProfileRepository
	logger          logger.Logger
}

// NewI18nService creates a new i18n service
func NewI18nService(userProfileRepo repository.UserProfileRepository, logger logger.Logger) *i18nService {
	return &i18nService{
		userProfileRepo: userProfileRepo,
		logger:          logger,
	}
}

// GetSupportedLanguages returns all supported languages with details
func (s *i18nService) GetSupportedLanguages(ctx context.Context) (*response.SupportedLanguagesResponse, error) {
	s.logger.InfoContext(ctx, "Getting supported languages")

	supportedLangs := i18n.GetSupportedLanguages()
	languages := make([]response.LanguageInfo, 0, len(supportedLangs))

	for _, lang := range supportedLangs {
		info := i18n.GetLanguageInfo(lang)

		languageInfo := response.LanguageInfo{
			Code:       info["code"].(string),
			Name:       info["name"].(string),
			NativeName: info["native_name"].(string),
			ISOCode:    info["iso_code"].(string),
			Supported:  info["supported"].(bool),
			IsDefault:  info["default"].(bool),
		}

		// Add region if available
		if region, ok := info["region"]; ok {
			languageInfo.Region = region.(string)
		}

		languages = append(languages, languageInfo)
	}

	result := &response.SupportedLanguagesResponse{
		Languages:       languages,
		DefaultLanguage: constants.DefaultLanguage,
	}

	s.logger.InfoContext(ctx, "Successfully retrieved supported languages",
		logger.Int("languages_count", len(languages)),
		logger.String("default_language", constants.DefaultLanguage))

	return result, nil
}

// SwitchUserLanguage switches a user's language preference
func (s *i18nService) SwitchUserLanguage(ctx context.Context, userID uint, req *request.SwitchLanguageRequest) error {
	s.logger.InfoContext(ctx, "Switching user language",
		logger.Uint("user_id", userID),
		logger.String("new_language", req.Language))

	// Validate language is supported
	if !s.isLanguageSupported(req.Language) {
		s.logger.WarnContext(ctx, "Unsupported language requested",
			logger.String("language", req.Language),
			logger.Any("supported_languages", i18n.GetSupportedLanguages()))
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// Normalize language code
	normalizedLang := i18n.NormalizeLanguage(req.Language)

	// Create user profile configuration
	userProfile := &model.UserProfile{
		UserID:      userID,
		Category:    constants.UserCategory,
		ConfigKey:   constants.UserLanguage,
		ConfigValue: normalizedLang,
		Description: "User preferred language setting",
	}

	// Save to database
	if err := s.userProfileRepo.SaveUserConfig(ctx, userProfile); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save user language preference",
			logger.ErrorField(err),
			logger.Uint("user_id", userID),
			logger.String("language", normalizedLang))
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	// Update Redis cache
	redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID)
	hashValues := map[string]string{
		constants.UserLanguage: normalizedLang,
	}
	if _, err := redis.HSet(ctx, redisKey, hashValues); err != nil {
		// Log warning but don't fail the request if cache update fails
		s.logger.WarnContext(ctx, "Failed to update language cache",
			logger.ErrorField(err),
			logger.Uint("user_id", userID),
			logger.String("language", normalizedLang))
	} else {
		s.logger.InfoContext(ctx, "Updated language cache",
			logger.Uint("user_id", userID),
			logger.String("language", normalizedLang))
	}

	s.logger.InfoContext(ctx, "Successfully switched user language",
		logger.Uint("user_id", userID),
		logger.String("language", normalizedLang))

	return nil
}

// isLanguageSupported checks if a language is supported
func (s *i18nService) isLanguageSupported(lang string) bool {
	normalizedLang := i18n.NormalizeLanguage(lang)
	supportedLanguages := i18n.GetSupportedLanguages()

	for _, supported := range supportedLanguages {
		if supported == normalizedLang {
			return true
		}
	}
	return false
}

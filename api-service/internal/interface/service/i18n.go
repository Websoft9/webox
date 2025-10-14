package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// I18nService provides internationalization services
type I18nService interface {
	// GetSupportedLanguages returns all supported languages with details
	GetSupportedLanguages(ctx context.Context) (*response.SupportedLanguagesResponse, error)

	// SwitchUserLanguage switches a user's language preference
	SwitchUserLanguage(ctx context.Context, userID uint, req *request.SwitchLanguageRequest) error
}

package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"api-service/pkg/security"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	// Token settings
	tokenExpiryWarningHours = 24
	tokenRandomBytesSize    = 32
)

type apiTokenService struct {
	tokenRepo repository.APITokenRepository
	db        *gorm.DB
	logger    logger.Logger
	i18n      *i18n.I18n
}

// NewAPITokenService creates a new API token service instance
func NewAPITokenService(
	tokenRepo repository.APITokenRepository,
	db *gorm.DB,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.APITokenService {
	return &apiTokenService{
		tokenRepo: tokenRepo,
		db:        db,
		logger:    logger,
		i18n:      i18n,
	}
}

// CreateAPIToken creates a new API token
func (s *apiTokenService) CreateAPIToken(ctx context.Context, req *request.CreateAPITokenRequest, userID uint) (*response.APITokenResponse, error) {
	s.logger.InfoContext(ctx, "Creating API token",
		logger.String("service", "api-token"),
		logger.String("operation", "CreateAPIToken"),
		logger.Uint("user_id", userID))

	// Generate token
	token, tokenHash, err := s.generateToken()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to generate token")
	}

	// Convert scopes to JSON format
	scopesJSON := model.JSON{
		"scopes": req.Scopes,
	}

	// Create API token model
	apiToken := &model.APIToken{
		Name:        req.Name,
		Token:       token,
		TokenHash:   tokenHash,
		UserID:      userID,
		Scopes:      scopesJSON,
		Description: req.Description,
		ExpiresAt:   req.ExpiresAt,
	}

	// Save to database
	if err := s.tokenRepo.Create(ctx, apiToken); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to create API token")
	}

	s.logger.InfoContext(ctx, "API token created successfully",
		logger.Uint("token_id", apiToken.ID))

	// Convert to response (include full token only on creation)
	resp := response.ConvertToAPITokenResponse(apiToken)
	resp.Token = token // Include full token for first time

	return resp, nil
}

// GetAPIToken gets API token details
func (s *apiTokenService) GetAPIToken(ctx context.Context, id, userID uint) (*response.APITokenResponse, error) {
	s.logger.InfoContext(ctx, "Getting API token",
		logger.String("service", "api-token"),
		logger.String("operation", "GetAPIToken"),
		logger.Uint("token_id", id),
		logger.Uint("user_id", userID))

	// Get token from database
	token, err := s.tokenRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get API token")
	}

	// Check if user owns the token
	if token.UserID != userID {
		return nil, errors.New("token not found")
	}

	return response.ConvertToAPITokenResponse(token), nil
}

// UpdateAPIToken updates API token
func (s *apiTokenService) UpdateAPIToken(ctx context.Context, id uint, req *request.UpdateAPITokenRequest, userID uint) (*response.APITokenResponse, error) {
	s.logger.InfoContext(ctx, "Updating API token",
		logger.String("service", "api-token"),
		logger.String("operation", "UpdateAPIToken"),
		logger.Uint("token_id", id),
		logger.Uint("user_id", userID))

	// Get existing token
	token, err := s.tokenRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get API token")
	}

	// Check if user owns the token
	if token.UserID != userID {
		return nil, errors.New("token not found")
	}

	// Update fields
	if req.Name != "" {
		token.Name = req.Name
	}
	if req.Description != "" {
		token.Description = req.Description
	}
	if len(req.Scopes) > 0 {
		token.Scopes = model.JSON{
			"scopes": req.Scopes,
		}
	}
	if req.ExpiresAt != nil {
		token.ExpiresAt = req.ExpiresAt
	}

	// Save changes
	if err := s.tokenRepo.Update(ctx, token); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to update API token")
	}

	s.logger.InfoContext(ctx, "API token updated successfully",
		logger.Uint("token_id", token.ID))

	return response.ConvertToAPITokenResponse(token), nil
}

// RevokeAPIToken revokes API token
func (s *apiTokenService) RevokeAPIToken(ctx context.Context, id, userID uint) error {
	s.logger.InfoContext(ctx, "Revoking API token",
		logger.String("service", "api-token"),
		logger.String("operation", "RevokeAPIToken"),
		logger.Uint("token_id", id),
		logger.Uint("user_id", userID))

	// Get existing token
	token, err := s.tokenRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("token not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return errors.Wrap(err, "failed to get API token")
	}

	// Check if user owns the token
	if token.UserID != userID {
		return errors.New("token not found")
	}

	// Delete the token (revoke)
	if err := s.tokenRepo.Delete(ctx, token.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to revoke API token", logger.ErrorField(err))
		return errors.Wrap(err, "failed to revoke API token")
	}

	s.logger.InfoContext(ctx, "API token revoked successfully",
		logger.Uint("token_id", token.ID))

	return nil
}

// RefreshAPIToken refreshes API token
func (s *apiTokenService) RefreshAPIToken(ctx context.Context, id, userID uint) (*response.APITokenResponse, error) {
	s.logger.InfoContext(ctx, "Refreshing API token",
		logger.String("service", "api-token"),
		logger.String("operation", "RefreshAPIToken"),
		logger.Uint("token_id", id),
		logger.Uint("user_id", userID))

	// Get existing token
	token, err := s.tokenRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to get API token")
	}

	// Check if user owns the token
	if token.UserID != userID {
		return nil, errors.New("token not found")
	}

	// Generate new token
	newToken, newTokenHash, err := s.generateToken()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate new token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to generate new token")
	}

	// Update token
	token.Token = newToken
	token.TokenHash = newTokenHash
	token.LastUsedAt = nil
	token.LastUsedIP = ""

	// Extend expiration if needed
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now().Add(tokenExpiryWarningHours*time.Hour)) {
		newExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days
		token.ExpiresAt = &newExpiry
	}

	// Save changes
	if err := s.tokenRepo.Update(ctx, token); err != nil {
		s.logger.ErrorContext(ctx, "Failed to refresh API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to refresh API token")
	}

	s.logger.InfoContext(ctx, "API token refreshed successfully",
		logger.Uint("token_id", token.ID))

	// Convert to response (include new token)
	resp := response.ConvertToAPITokenResponse(token)
	resp.Token = newToken

	return resp, nil
}

// ListAPITokens gets API token list
func (s *apiTokenService) ListAPITokens(ctx context.Context, req *request.ListAPITokensRequest, userID uint) (*response.APITokenListResponse, error) {
	s.logger.InfoContext(ctx, "Listing API tokens",
		logger.String("service", "api-token"),
		logger.String("operation", "ListAPITokens"),
		logger.Uint("user_id", userID))

	// Set user ID filter
	req.UserID = &userID

	// Get tokens from database
	tokens, total, err := s.tokenRepo.List(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list API tokens", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to list API tokens")
	}

	// Convert to response
	items := make([]response.APITokenResponse, len(tokens))
	for i, token := range tokens {
		items[i] = *response.ConvertToAPITokenResponse(token)
		// Token is already masked by ConvertToAPITokenResponse
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(req.GetPageSize())))

	return &response.APITokenListResponse{
		Items:      items,
		Total:      total,
		Page:       req.GetPage(),
		PageSize:   req.GetPageSize(),
		TotalPages: totalPages,
	}, nil
}

// ValidateAPIToken validates API token
func (s *apiTokenService) ValidateAPIToken(ctx context.Context, token string) (*response.APITokenValidationResponse, error) {
	s.logger.InfoContext(ctx, "Validating API token",
		logger.String("service", "api-token"),
		logger.String("operation", "ValidateAPIToken"))

	// Hash the token for lookup
	tokenHash := security.HashToken(token)

	// Get token from database
	apiToken, err := s.tokenRepo.GetByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &response.APITokenValidationResponse{Valid: false}, nil
		}
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return nil, errors.Wrap(err, "failed to validate API token")
	}

	// Check if token is not expired
	if apiToken.IsExpired() {
		return &response.APITokenValidationResponse{Valid: false}, nil
	}

	// Extract scopes
	scopes := s.extractScopesFromToken(apiToken)

	// Update last used
	go func() {
		if err := s.tokenRepo.UpdateLastUsed(context.Background(), apiToken.ID, ""); err != nil {
			s.logger.ErrorContext(context.Background(), "Failed to update token last used", logger.ErrorField(err))
		}
	}()

	return &response.APITokenValidationResponse{
		Valid:     true,
		UserID:    apiToken.UserID,
		Username:  apiToken.Username,
		Scopes:    scopes,
		ExpiresAt: *apiToken.ExpiresAt,
	}, nil
}

// BatchRevokeAPITokens revokes multiple API tokens
func (s *apiTokenService) BatchRevokeAPITokens(ctx context.Context, ids []uint, userID uint) error {
	s.logger.InfoContext(ctx, "Batch revoking API tokens",
		logger.String("service", "api-token"),
		logger.String("operation", "BatchRevokeAPITokens"),
		logger.Uint("user_id", userID),
		logger.Int("count", len(ids)))

	if len(ids) == 0 {
		return nil
	}

	// Revoke each token individually to ensure ownership check
	var errors []error
	for _, id := range ids {
		if err := s.RevokeAPIToken(ctx, id, userID); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		s.logger.ErrorContext(ctx, "Some tokens failed to revoke",
			logger.Int("failed_count", len(errors)))
		// Return the first error
		return errors[0]
	}

	s.logger.InfoContext(ctx, "API tokens batch revoked successfully",
		logger.Int("count", len(ids)))
	return nil
}

// CleanExpiredTokens cleans expired tokens
func (s *apiTokenService) CleanExpiredTokens(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Cleaning expired API tokens",
		logger.String("service", "api-token"),
		logger.String("operation", "CleanExpiredTokens"))

	if err := s.tokenRepo.CleanExpiredTokens(ctx); err != nil {
		s.logger.ErrorContext(ctx, "Failed to clean expired tokens", logger.ErrorField(err))
		return errors.Wrap(err, "failed to clean expired tokens")
	}

	s.logger.InfoContext(ctx, "Expired API tokens cleaned successfully")
	return nil
}

func (s *apiTokenService) CheckTokenIsExists(ctx context.Context, token string) bool {
	s.logger.InfoContext(ctx, "Check API token is exists",
		logger.String("service", "api-token"),
		logger.String("operation", "CheckTokenIsExists"))

	// Hash the token for lookup
	tokenHash := security.HashToken(token)

	// Generate Redis key
	redisKey := constants.TokenRedisKeyPrefix + tokenHash

	result, _ := redis.Exists(ctx, redisKey)
	var exists = true

	if result == 0 {
		_, err := s.tokenRepo.GetByToken(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.ErrorContext(ctx, "Failed to get API token in database & redis", logger.ErrorField(err))
				exists = false
			}
		}
	}
	return exists
}

// generateToken generates a new API token and its hash
func (s *apiTokenService) generateToken() (token, tokenHash string, err error) {
	// Generate random bytes
	bytes := make([]byte, tokenRandomBytesSize)
	if _, err = rand.Read(bytes); err != nil {
		return "", "", err
	}

	// Create token with prefix
	token = fmt.Sprintf("ws9_%s", hex.EncodeToString(bytes))

	// Hash the token for storage
	tokenHash = security.HashToken(token)

	return token, tokenHash, nil
}

// extractScopesFromToken extracts scopes from API token safely
func (s *apiTokenService) extractScopesFromToken(apiToken *model.APIToken) []string {
	if apiToken.Scopes == nil {
		return []string{}
	}

	scopesData, exists := apiToken.Scopes["scopes"]
	if !exists {
		return []string{}
	}

	scopesSlice, ok := scopesData.([]interface{})
	if !ok {
		return []string{}
	}

	scopes := make([]string, 0, len(scopesSlice))
	for _, scope := range scopesSlice {
		if s, ok := scope.(string); ok {
			scopes = append(scopes, s)
		}
	}

	return scopes
}

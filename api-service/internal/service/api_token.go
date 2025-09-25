package service

import (
	"api-service/internal/config"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type apiTokenService struct {
	tokenRepo         repository.APITokenRepository
	authConfigManager *config.AuthConfigManager
	db                *gorm.DB
	logger            logger.Logger
}

// NewAPITokenService creates a new API token service instance
func NewAPITokenService(
	tokenRepo repository.APITokenRepository,
	authConfigManager *config.AuthConfigManager,
	db *gorm.DB,
	logger logger.Logger,
) service.APITokenService {
	return &apiTokenService{
		tokenRepo:         tokenRepo,
		authConfigManager: authConfigManager,
		db:                db,
		logger:            logger,
	}
}

// CleanExpiredTokens cleans expired tokens
func (s *apiTokenService) CleanExpiredTokens(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Cleaning expired API tokens",
		logger.String("service", "api-token"),
		logger.String("operation", "CleanExpiredTokens"))

	if err := s.tokenRepo.CleanExpiredTokens(ctx); err != nil {
		s.logger.ErrorContext(ctx, "Failed to clean expired tokens", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeInternalError, "failed to clean expired tokens")
	}

	s.logger.InfoContext(ctx, "Expired API tokens cleaned successfully")
	return nil
}

// RevokeAPITokenByToken revokes API token by token string
func (s *apiTokenService) RevokeAPITokenByToken(ctx context.Context, token string, userID uint) error {
	s.logger.InfoContext(ctx, "Revoking API token by token",
		logger.String("service", "api-token"),
		logger.String("operation", "RevokeAPITokenByToken"),
		logger.Uint("user_id", userID))

	// Hash the token for lookup
	tokenHash := auth.HashToken(token)

	// Get token from database
	apiToken, err := s.tokenRepo.GetByToken(ctx, tokenHash)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get API token", logger.ErrorField(err))
		return err
	}

	// Check if user owns the token
	if apiToken.UserID != userID {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "token not found")
	}

	// Delete the token (revoke)
	if err := s.tokenRepo.Delete(ctx, apiToken.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to revoke API token", logger.ErrorField(err))
		return err
	}

	// Delete from Redis
	redisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)
	_, _ = redis.Del(ctx, redisKey) // Ignore error, Redis failure shouldn't fail the operation

	s.logger.InfoContext(ctx, "API token revoked successfully",
		logger.Uint("token_id", apiToken.ID))

	return nil
}

// RefreshUserAPIToken refreshes API token for user
func (s *apiTokenService) RefreshUserAPIToken(ctx context.Context, userID uint) (*response.APITokenResponse, error) {
	s.logger.InfoContext(ctx, "Refreshing user API token",
		logger.String("service", "api-token"),
		logger.String("operation", "RefreshUserAPIToken"),
		logger.Uint("user_id", userID))

	// Get user's most recent active token
	token, err := s.tokenRepo.GetActiveTokenByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user's active API token", logger.ErrorField(err))
		return nil, err
	}

	oldTokenRedisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, token.TokenHash)

	// Generate new token
	newToken, expiresAt, err := auth.GetGlobalJWT().RefreshToken(token.Token)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate new token", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to generate new token")
	}

	// Update token
	token.Token = newToken
	token.TokenHash = auth.HashToken(newToken)
	token.LastUsedAt = nil
	token.LastUsedIP = ""
	token.ExpiresAt = &expiresAt

	// Save changes
	if updateErr := s.tokenRepo.Update(ctx, token); updateErr != nil {
		s.logger.ErrorContext(ctx, "Failed to refresh user API token", logger.ErrorField(updateErr))
		return nil, updateErr
	}

	// Store JWT token in Redis cache with expiration
	redisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, token.TokenHash)

	// Create JWT token info for Redis storage
	tokenInfo := map[string]any{
		"token":      token,
		"user_id":    userID,
		"expires_at": token.ExpiresAt,
		"type":       "jwt_login",
		"created_at": time.Now(),
		"token_id":   token.ID, // Reference to database record
	}

	// Marshal token info to JSON
	tokenJSON, err := json.Marshal(tokenInfo)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal JWT token info to JSON",
			logger.ErrorField(err), logger.String("token", token.TokenHash[:8]+"..."))
		// Don't fail if Redis storage fails, database storage is more important
	} else {
		// Delete old token from Redis and update token
		_, _ = redis.Del(ctx, oldTokenRedisKey) // Ignore error, Redis failure shouldn't fail the operation
		// Store in Redis
		err = redis.Set(ctx, redisKey, string(tokenJSON), time.Duration(auth.GetGlobalJWT().GetTokenExpiresInSeconds())*time.Second)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to store JWT token in Redis",
				logger.ErrorField(err), logger.String("token", token.TokenHash[:8]+"..."), logger.String("key", redisKey))
			// Don't fail if Redis storage fails, database storage is more important
		}
	}

	s.logger.InfoContext(ctx, "User API token refreshed successfully",
		logger.Uint("token_id", token.ID))

	// Convert to response (include new token)
	resp := response.ConvertToAPITokenResponse(token)
	resp.Token = newToken

	return resp, nil
}

func (s *apiTokenService) CheckTokenIsExists(ctx context.Context, token string) bool {
	s.logger.InfoContext(ctx, "Check API token is exists",
		logger.String("service", "api-token"),
		logger.String("operation", "CheckTokenIsExists"))

	// Hash the token for lookup
	tokenHash := auth.HashToken(token)

	// Generate Redis key
	redisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)

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

package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"
	"time"

	"gorm.io/gorm"
)

// apiTokenRepository API token data access implementation
type apiTokenRepository struct {
	db *gorm.DB
}

// NewAPITokenRepository create new API Token Repository instance
func NewAPITokenRepository(db *gorm.DB) repository.APITokenRepository {
	return &apiTokenRepository{db: db}
}

// Create create API token
func (r *apiTokenRepository) Create(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByID get API token by ID
func (r *apiTokenRepository) GetByID(ctx context.Context, id uint) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&token, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	// Set username
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}
	return &token, nil
}

// GetByToken get API token by token hash
func (r *apiTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).Preload("User").Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}
	// Set username
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}
	return &token, nil
}

// Update update API token
func (r *apiTokenRepository) Update(ctx context.Context, token *model.APIToken) error {
	result := r.db.WithContext(ctx).Updates(token)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordUpdateFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}

// Delete delete API token
func (r *apiTokenRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.APIToken{}, id)
	if result.Error != nil {
		return errors.NewAppErrorWrapError(result.Error, errors.CodeRecordDeleteFailed)
	}
	if result.RowsAffected == 0 {
		return errors.NewAppError(errors.CodeRecordNoAffected)
	}
	return nil
}

// GetByUserID get API token list by user ID
func (r *apiTokenRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error) {
	var tokens []*model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&tokens).Error

	// Set username
	for _, token := range tokens {
		if token.User.ID != 0 {
			token.Username = token.User.Username
		}
	}

	return tokens, err
}

// GetActiveTokenByUserID get most recent active API token by user ID
func (r *apiTokenRepository) GetActiveTokenByUserID(ctx context.Context, userID uint) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Where("expires_at IS NULL OR expires_at > ?", time.Now().Format(time.DateTime)).
		Order("created_at desc").
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	// Set username
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}

	return &token, nil
}

// UpdateLastUsed update token last used time
func (r *apiTokenRepository) UpdateLastUsed(ctx context.Context, id uint, ip string) error {
	now := time.Now().Format(time.DateTime)
	return r.db.WithContext(ctx).
		Model(&model.APIToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_used_at": now,
			"last_used_ip": ip,
		}).Error
}

// CleanExpiredTokens clean expired tokens
func (r *apiTokenRepository) CleanExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now().Format(time.DateTime)).
		Delete(&model.APIToken{}).Error
}

// CreateWithTx create API token with transaction
func (r *apiTokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

// UpdateWithTx update API token with transaction
func (r *apiTokenRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Updates(token).Error
}

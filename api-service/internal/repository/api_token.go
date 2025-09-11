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
			return nil, errors.NewAppErrorWithMessage(errors.CodeRecordNotFound, "token not found")
		}
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "failed to get token by ID")
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
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token_hash = ?", tokenHash).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	// Set username
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}
	return &token, nil
}

// Update update API token
func (r *apiTokenRepository) Update(ctx context.Context, token *model.APIToken) error {
	result := r.db.WithContext(ctx).Save(token)
	if result.Error != nil {
		return errors.WrapError(result.Error, errors.CodeRecordUpdateFailed, "failed to update token")
	}
	if result.RowsAffected == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNoAffected, "no rows affected")
	}
	return nil
}

// Delete delete API token
func (r *apiTokenRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.APIToken{}, id)
	if result.Error != nil {
		return errors.WrapError(result.Error, errors.CodeRecordDeleteFailed, "failed to delete token")
	}
	if result.RowsAffected == 0 {
		return errors.NewAppErrorWithMessage(errors.CodeRecordNoAffected, "no rows affected")
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

// BatchDelete batch delete API tokens
func (r *apiTokenRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.APIToken{}).Error
}

// UpdateLastUsed update token last used time
func (r *apiTokenRepository) UpdateLastUsed(ctx context.Context, id uint, ip string) error {
	now := time.Now()
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
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&model.APIToken{}).Error
}

// CreateWithTx create API token with transaction
func (r *apiTokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

// UpdateWithTx update API token with transaction
func (r *apiTokenRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Save(token).Error
}

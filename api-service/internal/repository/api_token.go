package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
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
		return nil, err
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
	return r.db.WithContext(ctx).Save(token).Error
}

// Delete delete API token
func (r *apiTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.APIToken{}, id).Error
}

// List get API token list
func (r *apiTokenRepository) List(ctx context.Context, req *request.ListAPITokensRequest) ([]*model.APIToken, int64, error) {
	var tokens []*model.APIToken
	var total int64

	query := r.db.WithContext(ctx).Model(&model.APIToken{})

	// Apply filters
	query = r.applyFilters(query, req)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	offset := (req.GetPage() - 1) * req.GetPageSize()
	err := query.
		Preload("User").
		Offset(offset).
		Limit(req.GetPageSize()).
		Order("created_at desc").
		Find(&tokens).Error

	// Set username
	for _, token := range tokens {
		if token.User.ID != 0 {
			token.Username = token.User.Username
		}
	}

	return tokens, total, err
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

// applyFilters apply query filters
func (r *apiTokenRepository) applyFilters(query *gorm.DB, req *request.ListAPITokensRequest) *gorm.DB {
	// User ID filter
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}

	// Name search
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	// Expiration status filter
	if req.Expired != nil {
		if *req.Expired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now())
		} else {
			query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
		}
	}

	return query
}

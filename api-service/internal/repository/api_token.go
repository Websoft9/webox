package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// apiTokenRepository API令牌数据访问实现
type apiTokenRepository struct {
	db *gorm.DB
}

// NewAPITokenRepository 创建新的API Token Repository实例
func NewAPITokenRepository(db *gorm.DB) repository.APITokenRepository {
	return &apiTokenRepository{db: db}
}

// Create 创建API令牌
func (r *apiTokenRepository) Create(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetByID 根据ID获取API令牌
func (r *apiTokenRepository) GetByID(ctx context.Context, id uint) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&token, id).Error
	if err != nil {
		return nil, err
	}
	// 设置用户名
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}
	return &token, nil
}

// GetByToken 根据令牌哈希获取API令牌
func (r *apiTokenRepository) GetByToken(ctx context.Context, tokenHash string) (*model.APIToken, error) {
	var token model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token_hash = ?", tokenHash).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	// 设置用户名
	if token.User.ID != 0 {
		token.Username = token.User.Username
	}
	return &token, nil
}

// Update 更新API令牌
func (r *apiTokenRepository) Update(ctx context.Context, token *model.APIToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

// Delete 删除API令牌
func (r *apiTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.APIToken{}, id).Error
}

// List 获取API令牌列表
func (r *apiTokenRepository) List(ctx context.Context, req *request.ListAPITokensRequest) ([]*model.APIToken, int64, error) {
	var tokens []*model.APIToken
	var total int64

	query := r.db.WithContext(ctx).Model(&model.APIToken{})

	// 应用过滤器
	query = r.applyFilters(query, req)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	offset := (req.GetPage() - 1) * req.GetPageSize()
	err := query.
		Preload("User").
		Offset(offset).
		Limit(req.GetPageSize()).
		Order("created_at desc").
		Find(&tokens).Error

	// 设置用户名
	for _, token := range tokens {
		if token.User.ID != 0 {
			token.Username = token.User.Username
		}
	}

	return tokens, total, err
}

// GetByUserID 根据用户ID获取API令牌列表
func (r *apiTokenRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.APIToken, error) {
	var tokens []*model.APIToken
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&tokens).Error

	// 设置用户名
	for _, token := range tokens {
		if token.User.ID != 0 {
			token.Username = token.User.Username
		}
	}

	return tokens, err
}

// BatchDelete 批量删除API令牌
func (r *apiTokenRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.APIToken{}).Error
}

// UpdateLastUsed 更新令牌最后使用时间
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

// CleanExpiredTokens 清理过期令牌
func (r *apiTokenRepository) CleanExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&model.APIToken{}).Error
}

// CreateWithTx 使用事务创建API令牌
func (r *apiTokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

// UpdateWithTx 使用事务更新API令牌
func (r *apiTokenRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, token *model.APIToken) error {
	return tx.WithContext(ctx).Save(token).Error
}

// applyFilters 应用查询过滤器
func (r *apiTokenRepository) applyFilters(query *gorm.DB, req *request.ListAPITokensRequest) *gorm.DB {
	// 用户ID过滤
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}

	// 名称搜索
	if req.Search != "" {
		query = query.Where("name LIKE ?", "%"+req.Search+"%")
	}

	// 过期状态过滤
	if req.Expired != nil {
		if *req.Expired {
			query = query.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now())
		} else {
			query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
		}
	}

	return query
}

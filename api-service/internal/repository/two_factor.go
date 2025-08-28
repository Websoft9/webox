package repository

import (
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// userTwoFactorRepository 用户双因子认证数据访问实现
type userTwoFactorRepository struct {
	db *gorm.DB
}

// NewTwoFactorRepository 创建新的Two Factor Repository实例
func NewTwoFactorRepository(db *gorm.DB) repository.UserTwoFactorRepository {
	return &userTwoFactorRepository{db: db}
}

// Create 创建双因子认证记录
func (r *userTwoFactorRepository) Create(ctx context.Context, twoFactor *model.UserTwoFactor) error {
	return r.db.WithContext(ctx).Create(twoFactor).Error
}

// GetByUserIDAndMethod 根据用户ID和方法获取双因子认证记录
func (r *userTwoFactorRepository) GetByUserIDAndMethod(ctx context.Context, userID uint, method string) (*model.UserTwoFactor, error) {
	var twoFactor model.UserTwoFactor
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND method = ?", userID, method).
		First(&twoFactor).Error
	if err != nil {
		return nil, err
	}
	return &twoFactor, nil
}

// Update 更新双因子认证记录
func (r *userTwoFactorRepository) Update(ctx context.Context, twoFactor *model.UserTwoFactor) error {
	return r.db.WithContext(ctx).Save(twoFactor).Error
}

// Delete 删除双因子认证记录
func (r *userTwoFactorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.UserTwoFactor{}, id).Error
}

// GetByUserID 根据用户ID获取所有双因子认证记录
func (r *userTwoFactorRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.UserTwoFactor, error) {
	var twoFactors []*model.UserTwoFactor
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Find(&twoFactors).Error
	return twoFactors, err
}

// EnableMethod 启用指定的双因子认证方法
func (r *userTwoFactorRepository) EnableMethod(ctx context.Context, userID uint, method, secret string) error {
	now := time.Now()
	twoFactor := &model.UserTwoFactor{
		UserID:     userID,
		Method:     method,
		Secret:     secret,
		Enabled:    true,
		VerifiedAt: &now,
	}

	return r.db.WithContext(ctx).
		Where("user_id = ? AND method = ?", userID, method).
		Assign(twoFactor).
		FirstOrCreate(twoFactor).Error
}

// DisableMethod 禁用指定的双因子认证方法
func (r *userTwoFactorRepository) DisableMethod(ctx context.Context, userID uint, method string) error {
	return r.db.WithContext(ctx).
		Model(&model.UserTwoFactor{}).
		Where("user_id = ? AND method = ?", userID, method).
		Updates(map[string]interface{}{
			"enabled":    false,
			"secret":     "",
			"updated_at": time.Now(),
		}).Error
}

// IsMethodEnabled 检查指定的双因子认证方法是否已启用
func (r *userTwoFactorRepository) IsMethodEnabled(ctx context.Context, userID uint, method string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserTwoFactor{}).
		Where("user_id = ? AND method = ? AND enabled = ?", userID, method, true).
		Count(&count).Error
	return count > 0, err
}

// CreateWithTx 使用事务创建双因子认证记录
func (r *userTwoFactorRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error {
	return tx.WithContext(ctx).Create(twoFactor).Error
}

// UpdateWithTx 使用事务更新双因子认证记录
func (r *userTwoFactorRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, twoFactor *model.UserTwoFactor) error {
	return tx.WithContext(ctx).Save(twoFactor).Error
}

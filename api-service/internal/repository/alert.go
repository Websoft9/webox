package repository

import (
	"api-service/internal/dto/request"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"context"

	"gorm.io/gorm"
)

// alertRepository implements the alert repository.
type alertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new instance of alert repository.
func NewAlertRepository(db *gorm.DB) *alertRepository {
	return &alertRepository{
		db: db,
	}
}

// CreateAlertRule creates an alert rule.
func (r *alertRepository) CreateAlertRule(ctx context.Context, rule *model.AlertRule) error {
	if err := r.db.WithContext(ctx).Create(rule).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}

	return nil
}

// GetAlertRuleByID retrieves an alert rule by its ID.
func (r *alertRepository) GetAlertRuleByID(ctx context.Context, id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &rule, nil
}

// ListAlertRules queries the list of alert rules.
func (r *alertRepository) ListAlertRules(ctx context.Context, req *request.AlertRuleQueryRequest) ([]*model.AlertRule, int64, error) {
	var rules []*model.AlertRule
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AlertRule{})

	// Apply query conditions
	if req.RuleType != "" {
		query = query.Where("rule_type = ?", req.RuleType)
	}

	if req.TargetType != "" {
		query = query.Where("target_type = ?", req.TargetType)
	}

	if req.IsEnabled != nil {
		query = query.Where("is_enabled = ?", *req.IsEnabled)
	}

	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order(req.GetSortOrder()).
		Offset(offset).Limit(limit).
		Find(&rules).Error

	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// UpdateAlertRule updates an alert rule.
func (r *alertRepository) UpdateAlertRule(ctx context.Context, id uint, updates map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&model.AlertRule{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	return nil
}

// DeleteAlertRule deletes an alert rule.
func (r *alertRepository) DeleteAlertRule(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Delete(&model.AlertRule{}, id).Error; err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}

	return nil
}

func (r *alertRepository) ListAlertRecords(ctx context.Context, req *request.AlertRecordQueryRequest) ([]*model.AlertRecord, int64, error) {
	var records []*model.AlertRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AlertRecord{})

	// Build filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Severity != "" {
		query = query.Where("severity = ?", req.Severity)
	}
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.AlertRuleID != nil {
		query = query.Where("alert_rule_id = ?", *req.AlertRuleID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := req.GetOffset()
	limit := req.GetPageSize()

	err := query.Order(req.GetSortOrder()).
		Offset(offset).Limit(limit).
		Find(&records).Error

	if err != nil {
		return nil, 0, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return records, total, nil
}

// GetAlertRecordByID retrieves an alert record by its ID
func (r *alertRepository) GetAlertRecordByID(ctx context.Context, id uint) (*model.AlertRecord, error) {
	var record model.AlertRecord

	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordQueryFailed)
	}

	return &record, nil
}

// CreateAlertRecord creates a new alert record
func (r *alertRepository) CreateAlertRecord(ctx context.Context, record *model.AlertRecord) error {
	err := r.db.Create(record).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordCreateFailed)
	}
	return nil
}

// UpdateAlertRecord updates an existing alert record
func (r *alertRepository) UpdateAlertRecord(ctx context.Context, id uint, updateData map[string]interface{}) error {
	err := r.db.Model(&model.AlertRecord{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}
	return nil
}

// DeleteAlertRecord deletes an alert record
func (r *alertRepository) DeleteAlertRecord(ctx context.Context, id uint) error {
	err := r.db.Delete(&model.AlertRecord{}, id).Error
	if err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeRecordDeleteFailed)
	}
	return nil
}

// ExistsAlertRecord checks if an alert record exists by ID
func (r *alertRepository) ExistsAlertRecord(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.AlertRecord{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

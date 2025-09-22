package repository

import (
	"api-service/internal/model"
	"context"
	"time"

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
		return err
	}

	return nil
}

// GetAlertRuleByID retrieves an alert rule by its ID.
func (r *alertRepository) GetAlertRuleByID(ctx context.Context, id uint) (*model.AlertRule, error) {
	var rule model.AlertRule
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rule).Error
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// ListAlertRules queries the list of alert rules.
func (r *alertRepository) ListAlertRules(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*model.AlertRule, int64, error) {
	var rules []*model.AlertRule
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AlertRule{})

	// Apply query conditions
	if ruleType, exists := params["rule_type"]; exists && ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}

	if targetType, exists := params["target_type"]; exists && targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	if isEnabled, exists := params["is_enabled"]; exists {
		query = query.Where("is_enabled = ?", isEnabled)
	}

	if keyword, exists := params["keyword"].(string); exists && keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated query
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rules).Error; err != nil {
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
		return err
	}

	return nil
}

// DeleteAlertRule deletes an alert rule.
func (r *alertRepository) DeleteAlertRule(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Delete(&model.AlertRule{}, id).Error; err != nil {
		return err
	}

	return nil
}

// ListAlertRecords retrieves a paginated list of alert records
func (r *alertRepository) ListAlertRecords(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.AlertRecord, int64, error) {
	var records []*model.AlertRecord
	var total int64

	query := r.db.Model(&model.AlertRecord{})

	// Apply filters
	if status, exists := filters["status"]; exists && status != "" {
		query = query.Where("status = ?", status)
	}

	if severity, exists := filters["severity"]; exists && severity != "" {
		query = query.Where("severity = ?", severity)
	}

	if startTime, exists := filters["start_time"]; exists {
		if st, ok := startTime.(time.Time); ok && !st.IsZero() {
			query = query.Where("fired_at >= ?", st)
		}
	}

	if endTime, exists := filters["end_time"]; exists {
		if et, ok := endTime.(time.Time); ok && !et.IsZero() {
			query = query.Where("fired_at <= ?", et)
		}
	}

	if alertRuleID, exists := filters["alert_rule_id"]; exists {
		query = query.Where("alert_rule_id = ?", alertRuleID)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Paginate and sort
	err = query.Order("fired_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error

	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetAlertRecordByID retrieves an alert record by its ID
func (r *alertRepository) GetAlertRecordByID(ctx context.Context, id uint) (*model.AlertRecord, error) {
	var record model.AlertRecord

	err := r.db.First(&record, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return &record, nil
}

// CreateAlertRecord creates a new alert record
func (r *alertRepository) CreateAlertRecord(ctx context.Context, record *model.AlertRecord) error {
	err := r.db.Create(record).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateAlertRecord updates an existing alert record
func (r *alertRepository) UpdateAlertRecord(ctx context.Context, id uint, updateData map[string]interface{}) error {
	err := r.db.Model(&model.AlertRecord{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteAlertRecord deletes an alert record
func (r *alertRepository) DeleteAlertRecord(ctx context.Context, id uint) error {
	err := r.db.Delete(&model.AlertRecord{}, id).Error
	if err != nil {
		return err
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

package repository

import (
	"api-service/internal/model"
	"context"
)

// AlertRepository defines the interface for the alert repository.
type AlertRepository interface {
	// CreateAlertRule creates an alert rule.
	CreateAlertRule(ctx context.Context, rule *model.AlertRule) error

	// GetAlertRuleByID retrieves an alert rule by its ID.
	GetAlertRuleByID(ctx context.Context, id uint) (*model.AlertRule, error)

	// ListAlertRules queries the list of alert rules.
	ListAlertRules(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*model.AlertRule, int64, error)

	// UpdateAlertRule updates an alert rule.
	UpdateAlertRule(ctx context.Context, id uint, updates map[string]interface{}) error

	// DeleteAlertRule deletes an alert rule.
	DeleteAlertRule(ctx context.Context, id uint) error

	// ExistsAlertRule checks if an alert rule exists by its ID.
	ListAlertRecords(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.AlertRecord, int64, error)

	// GetAlertRecordByID retrieves an alert record by its ID.
	GetAlertRecordByID(ctx context.Context, id uint) (*model.AlertRecord, error)

	// CreateAlertRecord creates an alert record.
	CreateAlertRecord(ctx context.Context, record *model.AlertRecord) error

	// UpdateAlertRecord updates an alert record.
	UpdateAlertRecord(ctx context.Context, id uint, updateData map[string]interface{}) error

	// DeleteAlertRecord deletes an alert record.
	DeleteAlertRecord(ctx context.Context, id uint) error

	// ExistsAlertRecord checks if an alert record exists by its ID.
	ExistsAlertRecord(ctx context.Context, id uint) (bool, error)
}

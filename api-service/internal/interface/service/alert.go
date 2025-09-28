package service

import (
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"context"
)

// AlertService defines the interface for alert service.
type AlertService interface {
	// CreateAlertRule creates an alert rule.
	CreateAlertRule(ctx context.Context, userID uint, req *request.AlertRuleCreateRequest) (*response.AlertRuleResponse, error)

	// GetAlertRuleByID retrieves a single alert rule by its ID.
	GetAlertRuleByID(ctx context.Context, id uint) (*response.AlertRuleResponse, error)

	// ListAlertRules retrieves a list of alert rules.
	ListAlertRules(ctx context.Context, req *request.AlertRuleQueryRequest) (*common.PaginationResponse, error)

	// UpdateAlertRule updates an alert rule.
	UpdateAlertRule(ctx context.Context, id uint, req *request.AlertRuleUpdateRequest) (*response.AlertRuleResponse, error)

	// DeleteAlertRule deletes an alert rule.
	DeleteAlertRule(ctx context.Context, id uint) error

	// ListAlertRecords retrieves a list of alert records.
	ListAlertRecords(ctx context.Context, req *request.AlertRecordQueryRequest) (*common.PaginationResponse, error)

	// AcknowledgeAlertRecord acknowledges an alert record.
	AcknowledgeAlertRecord(ctx context.Context, id, userID uint, req *request.AlertAcknowledgeRequest) error

	// ResolveAlertRecord resolves an alert record.
	ResolveAlertRecord(ctx context.Context, id, userID uint, req *request.AlertResolveRequest) error
}

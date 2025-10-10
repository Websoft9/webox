package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"time"
)

// alertService implements the alert service.
type alertService struct {
	alertRepo repository.AlertRepository
	logger    logger.Logger
	i18n      *i18n.I18n
}

// NewAlertService creates a new instance of alert service.
func NewAlertService(
	alertRepo repository.AlertRepository,
	logger logger.Logger,
	i18n *i18n.I18n,
) *alertService {
	return &alertService{
		alertRepo: alertRepo,
		logger:    logger,
		i18n:      i18n,
	}
}

// CreateAlertRule creates an alert rule.
func (s *alertService) CreateAlertRule(ctx context.Context, currentUserID uint, req *request.AlertRuleCreateRequest) (*response.AlertRuleResponse, error) {
	s.logger.InfoContext(ctx, "Creating alert rule",
		logger.String("name", req.Name),
		logger.String("ruleType", string(req.RuleType)))

	// Build alert rule model
	rule := &model.AlertRule{
		Name:                 req.Name,
		RuleType:             req.RuleType,
		TargetType:           req.TargetType,
		TargetID:             req.TargetID,
		MetricName:           req.MetricName,
		ConditionExpression:  req.ConditionExpression,
		NotificationChannels: req.NotificationChannels,
		OwnerID:              currentUserID,
	}

	if req.IsEnabled != nil {
		rule.IsEnabled = *req.IsEnabled
	} else {
		rule.IsEnabled = true
	}

	// Save to database
	if err := s.alertRepo.CreateAlertRule(ctx, rule); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create alert rule",
			logger.String("name", req.Name),
			logger.ErrorField(err))
		return nil, err
	}

	// Build response
	return s.buildAlertRuleResponse(rule), nil
}

// GetAlertRuleByID retrieves a single alert rule.
func (s *alertService) GetAlertRuleByID(ctx context.Context, id uint) (*response.AlertRuleResponse, error) {
	s.logger.InfoContext(ctx, "Getting alert rule", logger.Uint("id", id))

	rule, err := s.alertRepo.GetAlertRuleByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get alert rule",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	return s.buildAlertRuleResponse(rule), nil
}

// ListAlertRecords retrieves a paginated list of alert records
func (s *alertService) ListAlertRecords(ctx context.Context, req *request.AlertRecordQueryRequest) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Listing alert records",
		logger.String("service", "alert"),
		logger.String("operation", "ListAlertRecords"))

	// Call repository to fetch data
	records, total, err := s.alertRepo.ListAlertRecords(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list alert records", logger.ErrorField(err))
		return nil, err
	}

	// Convert to DTOs
	items := make([]response.AlertRecordResponse, len(records))
	for i, record := range records {
		items[i] = s.mapAlertRecordToDTO(record)
	}

	return common.NewPaginationResponse(
		req.GetOffset(),
		req.GetPageSize(),
		total,
		items,
	), nil
}

// UpdateAlertRule updates an alert rule.
func (s *alertService) UpdateAlertRule(ctx context.Context, id uint, req *request.AlertRuleUpdateRequest) (*response.AlertRuleResponse, error) {
	s.logger.InfoContext(ctx, "Updating alert rule", logger.Uint("id", id))

	// Check if the rule exists
	_, err := s.alertRepo.GetAlertRuleByID(ctx, id)

	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get alert rule for update",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	// Build update data
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.ConditionExpression != nil {
		updates["condition_expression"] = *req.ConditionExpression
	}

	if req.NotificationChannels != nil {
		updates["notification_channels"] = *req.NotificationChannels
	}

	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}

	// Update rule
	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		// Use a different variable name to avoid shadowing
		if updateErr := s.alertRepo.UpdateAlertRule(ctx, id, updates); updateErr != nil {
			s.logger.ErrorContext(ctx, "Failed to update alert rule",
				logger.Uint("id", id),
				logger.ErrorField(updateErr))
			return nil, updateErr
		}
	}

	// Get latest data
	updatedRule, err := s.alertRepo.GetAlertRuleByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get updated alert rule",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return nil, err
	}

	return s.buildAlertRuleResponse(updatedRule), nil
}

// DeleteAlertRule deletes an alert rule.
func (s *alertService) DeleteAlertRule(ctx context.Context, id uint) error {
	s.logger.InfoContext(ctx, "Deleting alert rule", logger.Uint("id", id))

	// Check if the rule exists
	_, err := s.alertRepo.GetAlertRuleByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get alert rule for deletion",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	// Delete rule
	if err := s.alertRepo.DeleteAlertRule(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete alert rule",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	return nil
}

// buildAlertRuleResponse builds the alert rule response.
func (s *alertService) buildAlertRuleResponse(rule *model.AlertRule) *response.AlertRuleResponse {
	return &response.AlertRuleResponse{
		ID:                   rule.ID,
		Name:                 rule.Name,
		RuleType:             rule.RuleType,
		TargetType:           rule.TargetType,
		TargetID:             rule.TargetID,
		MetricName:           rule.MetricName,
		ConditionExpression:  rule.ConditionExpression,
		NotificationChannels: rule.NotificationChannels,
		IsEnabled:            rule.IsEnabled,
		OwnerID:              rule.OwnerID,
		CreatedAt:            rule.CreatedAt,
		UpdatedAt:            rule.UpdatedAt,
	}
}

// ListAlertRules retrieves a list of alert rules.
func (s *alertService) ListAlertRules(ctx context.Context, req *request.AlertRuleQueryRequest) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Listing alert rules",
		logger.String("service", "alert"),
		logger.String("operation", "ListAlertRules"))

	// Query data
	rules, total, err := s.alertRepo.ListAlertRules(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list alert rules", logger.ErrorField(err))
		return nil, err
	}

	// Build response
	items := make([]response.AlertRuleResponse, len(rules))
	for i, rule := range rules {
		items[i] = *s.buildAlertRuleResponse(rule)
	}

	return common.NewPaginationResponse(
		req.GetOffset(),
		req.GetPageSize(),
		total,
		items,
	), nil
}

// AcknowledgeAlertRecord acknowledges an alert record
func (s *alertService) AcknowledgeAlertRecord(ctx context.Context, id, userID uint, req *request.AlertAcknowledgeRequest) error {
	s.logger.InfoContext(ctx, "Acknowledging alert record",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	// Get the record
	record, err := s.alertRepo.GetAlertRecordByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get alert record",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	// Validate status
	if record.AcknowledgedAt != nil {
		return errors.NewAppError(errors.CodeResourceStateNotAllowed)
	}

	// Check if already resolved
	if record.Status == constants.AlertStatusResolved || record.Status == constants.AlertStatusConfirmed {
		return errors.NewAppError(errors.CodeResourceStateNotAllowed)
	}

	// Prepare update data
	now := time.Now()
	updateData := map[string]interface{}{
		"acknowledged_at": now,
		"acknowledged_by": userID,
		"status":          constants.AlertStatusConfirmed,
	}

	// Add acknowledgement note if provided
	if req.Note != "" {
		updateData["resolution_note"] = req.Note
	}

	// Update the record
	err = s.alertRepo.UpdateAlertRecord(ctx, id, updateData)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update alert record for acknowledgement",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Alert record acknowledged successfully",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))
	return nil
}

// ResolveAlertRecord resolves an alert record
func (s *alertService) ResolveAlertRecord(ctx context.Context, id, userID uint, req *request.AlertResolveRequest) error {
	s.logger.InfoContext(ctx, "Resolving alert record",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))

	// Get the record
	record, err := s.alertRepo.GetAlertRecordByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get alert record",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	// Validate status
	if record.Status == constants.AlertStatusResolved || record.Status == constants.AlertStatusConfirmed {
		return errors.NewAppError(errors.CodeResourceStateNotAllowed)
	}

	// Prepare update data
	now := time.Now()
	updateData := map[string]interface{}{
		"acknowledged_by": userID,
		"status":          constants.AlertStatusResolved,
		"resolved_at":     now,
	}

	// Add resolution note if provided
	if req.ResolutionNote != "" {
		updateData["resolution_note"] = req.ResolutionNote
	}

	// Update the record
	err = s.alertRepo.UpdateAlertRecord(ctx, id, updateData)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update alert record for resolution",
			logger.Uint("id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Alert record resolved successfully",
		logger.Uint("id", id),
		logger.Uint("user_id", userID))
	return nil
}

// mapAlertRecordToDTO converts an alert record model to DTO
func (s *alertService) mapAlertRecordToDTO(record *model.AlertRecord) response.AlertRecordResponse {
	return response.AlertRecordResponse{
		ID:               record.ID,
		AlertRuleID:      record.AlertRuleID,
		AlertID:          record.AlertID,
		Title:            record.Title,
		Description:      record.Description,
		Status:           record.Status,
		Severity:         s.determineSeverity(record),
		FiredAt:          record.FiredAt,
		ResolvedAt:       record.ResolvedAt,
		AcknowledgedAt:   record.AcknowledgedAt,
		AcknowledgedBy:   record.AcknowledgedBy,
		AcknowledgeNote:  record.AcknowledgeNote,
		ResolutionNote:   record.ResolutionNote,
		NotificationSent: record.NotificationSent,
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
	}
}

// determineSeverity determines the severity level based on record status
func (s *alertService) determineSeverity(record *model.AlertRecord) string {
	// Determine severity based on business logic
	// This is a simplified implementation, actual logic may vary
	if record.Status == "FIRING" {
		return "CRITICAL"
	}
	return "INFO"
}

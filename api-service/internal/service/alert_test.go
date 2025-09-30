package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/model"
	"api-service/pkg/i18n"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock AlertRepository
type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) CreateAlertRule(ctx context.Context, rule *model.AlertRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockAlertRepository) GetAlertRuleByID(ctx context.Context, id uint) (*model.AlertRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AlertRule), args.Error(1)
}

func (m *MockAlertRepository) ListAlertRules(ctx context.Context, params map[string]interface{}, page, pageSize int) ([]*model.AlertRule, int64, error) {
	args := m.Called(ctx, params, page, pageSize)
	return args.Get(0).([]*model.AlertRule), args.Get(1).(int64), args.Error(2)
}

func (m *MockAlertRepository) UpdateAlertRule(ctx context.Context, id uint, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAlertRepository) DeleteAlertRule(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAlertRepository) ListAlertRecords(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*model.AlertRecord, int64, error) {
	args := m.Called(ctx, offset, limit, filters)
	return args.Get(0).([]*model.AlertRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockAlertRepository) GetAlertRecordByID(ctx context.Context, id uint) (*model.AlertRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AlertRecord), args.Error(1)
}

func (m *MockAlertRepository) CreateAlertRecord(ctx context.Context, record *model.AlertRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockAlertRepository) UpdateAlertRecord(ctx context.Context, id uint, updateData map[string]interface{}) error {
	args := m.Called(ctx, id, updateData)
	return args.Error(0)
}

func (m *MockAlertRepository) DeleteAlertRecord(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAlertRepository) ExistsAlertRecord(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// func (m *MockLogger) Debug(msg string, fields ...logger.Field) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) Info(msg string, fields ...logger.Field) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) Warn(msg string, fields ...logger.Field) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) Error(msg string, fields ...logger.Field) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) Fatal(msg string, fields ...logger.Field) {
// 	m.Called(msg, fields)
// }

// func (m *MockLogger) WithFields(fields ...logger.Field) logger.Logger {
// 	m.Called(fields)
// 	return m
// }

// func (m *MockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {
// 	m.Called(ctx, msg, fields)
// }

// func (m *MockLogger) InfoContext(ctx context.Context, msg string, fields ...logger.Field) {
// 	m.Called(ctx, msg, fields)
// }

// func (m *MockLogger) WarnContext(ctx context.Context, msg string, fields ...logger.Field) {
// 	m.Called(ctx, msg, fields)
// }

// func (m *MockLogger) ErrorContext(ctx context.Context, msg string, fields ...logger.Field) {
// 	m.Called(ctx, msg, fields)
// }

// func (m *MockLogger) FatalContext(ctx context.Context, msg string, fields ...logger.Field) {
// 	m.Called(ctx, msg, fields)
// }

// Setup test function
func setupAlertService() (*alertService, *MockAlertRepository, *MockLogger) {
	mockRepo := new(MockAlertRepository)
	mockLogger := new(MockLogger)
	mockI18n := i18n.NewI18n()

	// Setup default behavior for logger to avoid having to mock every call
	mockLogger.On("InfoContext", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("ErrorContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("WarnContext", mock.Anything, mock.Anything, mock.Anything).Return()

	service := NewAlertService(mockRepo, mockLogger, mockI18n)

	return service, mockRepo, mockLogger
}

func TestCreateAlertRule(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful creation
	t.Run("Success", func(t *testing.T) {
		targetID := uint(12)
		req := &request.AlertRuleCreateRequest{
			Name:                 "Test Rule",
			RuleType:             "THRESHOLD",
			TargetType:           "SERVER",
			TargetID:             &targetID,
			MetricName:           "cpu.usage",
			ConditionExpression:  "cpu.usage > 90",
			NotificationChannels: "email",
		}

		expectedRule := &model.AlertRule{
			Name:                 "Test Rule",
			RuleType:             "THRESHOLD",
			TargetType:           "SERVER",
			TargetID:             &targetID,
			MetricName:           "cpu.usage",
			ConditionExpression:  "cpu.usage > 90",
			NotificationChannels: "email",
			OwnerID:              uint(1),
			IsEnabled:            true,
		}

		mockRepo.On("CreateAlertRule", ctx, mock.MatchedBy(func(rule *model.AlertRule) bool {
			return rule.Name == expectedRule.Name &&
				rule.RuleType == expectedRule.RuleType &&
				rule.ConditionExpression == expectedRule.ConditionExpression
		})).Return(nil).Once()

		// Execute
		result, err := service.CreateAlertRule(ctx, uint(1), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.RuleType, result.RuleType)
		assert.Equal(t, req.ConditionExpression, result.ConditionExpression)
		mockRepo.AssertExpectations(t)
	})

	// Test case 3: IsEnabled set to false
	t.Run("Disabled Rule", func(t *testing.T) {
		isEnabled := false
		targetID := uint(1)
		req := &request.AlertRuleCreateRequest{
			Name:                 "Disabled Rule",
			RuleType:             "THRESHOLD",
			TargetType:           "SERVER",
			TargetID:             &targetID,
			MetricName:           "cpu.usage",
			ConditionExpression:  "cpu.usage > 90",
			NotificationChannels: "email",
			IsEnabled:            &isEnabled,
		}

		mockRepo.On("CreateAlertRule", ctx, mock.MatchedBy(func(rule *model.AlertRule) bool {
			return rule.IsEnabled == false
		})).Return(nil).Once()

		// Execute
		result, err := service.CreateAlertRule(ctx, uint(1), req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsEnabled)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAlertRuleByID(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Rule found
	t.Run("Rule Found", func(t *testing.T) {
		ruleID := uint(1)
		targetID := uint(1)
		expectedRule := &model.AlertRule{
			ID:                   ruleID,
			Name:                 "Test Rule",
			RuleType:             "THRESHOLD",
			TargetType:           "SERVER",
			TargetID:             &targetID,
			MetricName:           "cpu.usage",
			ConditionExpression:  "cpu.usage > 90",
			NotificationChannels: "email",
			IsEnabled:            true,
			OwnerID:              uint(1),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(expectedRule, nil).Once()

		// Execute
		result, err := service.GetAlertRuleByID(ctx, ruleID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedRule.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

}

func TestListAlertRules(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Empty list
	t.Run("Empty List", func(t *testing.T) {
		req := &request.AlertRuleQueryRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.On("ListAlertRules", ctx, mock.Anything, req.Page, req.PageSize).
			Return([]*model.AlertRule{}, int64(0), nil).Once()

		// Execute
		result, err := service.ListAlertRules(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(0), result.Total)

		items, ok := result.Items.([]response.AlertRuleResponse)
		assert.True(t, ok, "Items should be of type []response.AlertRuleResponse")
		assert.Empty(t, items)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Rules found with filters
	t.Run("Rules Found With Filters", func(t *testing.T) {
		req := &request.AlertRuleQueryRequest{
			Page:       1,
			PageSize:   10,
			RuleType:   "THRESHOLD",
			TargetType: "SERVER",
			Keyword:    "cpu",
		}

		isEnabled := true
		req.IsEnabled = &isEnabled
		targetID1 := uint(1)
		targetID2 := uint(2)
		rules := []*model.AlertRule{
			{
				ID:                   uint(1),
				Name:                 "CPU High Alert",
				RuleType:             "THRESHOLD",
				TargetType:           "SERVER",
				TargetID:             &targetID1,
				MetricName:           "cpu.usage",
				ConditionExpression:  "cpu.usage > 90",
				NotificationChannels: "email",
				IsEnabled:            true,
			},
			{
				ID:                   uint(2),
				Name:                 "CPU Load Alert",
				RuleType:             "THRESHOLD",
				TargetType:           "SERVER",
				TargetID:             &targetID2,
				MetricName:           "cpu.load",
				ConditionExpression:  "cpu.load > 5",
				NotificationChannels: "sms",
				IsEnabled:            true,
			},
		}

		mockRepo.On("ListAlertRules", ctx, mock.MatchedBy(func(params map[string]interface{}) bool {
			return params["rule_type"] == req.RuleType &&
				params["target_type"] == req.TargetType &&
				params["is_enabled"] == *req.IsEnabled &&
				params["keyword"] == req.Keyword
		}), req.Page, req.PageSize).Return(rules, int64(2), nil).Once()

		// Execute
		result, err := service.ListAlertRules(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.Total)

		items, ok := result.Items.([]response.AlertRuleResponse)
		assert.True(t, ok, "Items should be of type []response.AlertRuleResponse")
		assert.Len(t, items, 2)
		assert.Equal(t, "CPU High Alert", items[0].Name)
		assert.Equal(t, "CPU Load Alert", items[1].Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateAlertRule(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful update
	t.Run("Successful Update", func(t *testing.T) {
		ruleID := uint(1)
		name := "Updated Rule Name"
		condExpr := "cpu.usage > 95"
		channels := "email"
		isEnabled := false

		req := &request.AlertRuleUpdateRequest{
			Name:                 &name,
			ConditionExpression:  &condExpr,
			NotificationChannels: &channels,
			IsEnabled:            &isEnabled,
		}

		// Mock get rule for existence check
		existingRule := &model.AlertRule{
			ID:                   ruleID,
			Name:                 "Original Rule Name",
			RuleType:             "THRESHOLD",
			ConditionExpression:  "cpu.usage > 90",
			NotificationChannels: "email",
			IsEnabled:            true,
		}
		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(existingRule, nil).Once()

		// Mock update
		mockRepo.On("UpdateAlertRule", ctx, ruleID, mock.MatchedBy(func(updates map[string]interface{}) bool {
			return updates["name"] == name &&
				updates["condition_expression"] == condExpr &&
				updates["is_enabled"] == isEnabled
		})).Return(nil).Once()

		// Mock get updated rule
		updatedRule := &model.AlertRule{
			ID:                   ruleID,
			Name:                 name,
			RuleType:             "THRESHOLD",
			ConditionExpression:  condExpr,
			NotificationChannels: channels,
			IsEnabled:            isEnabled,
		}
		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(updatedRule, nil).Once()

		// Execute
		result, err := service.UpdateAlertRule(ctx, ruleID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, condExpr, result.ConditionExpression)
		assert.Equal(t, channels, result.NotificationChannels)
		assert.Equal(t, isEnabled, result.IsEnabled)
		mockRepo.AssertExpectations(t)
	})

	// Test case 2: Empty update (no changes)
	t.Run("No Changes", func(t *testing.T) {
		ruleID := uint(1)
		req := &request.AlertRuleUpdateRequest{}

		// Mock get rule for existence check
		existingRule := &model.AlertRule{
			ID:   ruleID,
			Name: "Original Rule Name",
		}
		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(existingRule, nil).Once()

		// Get updated rule (no changes)
		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(existingRule, nil).Once()

		// Execute
		result, err := service.UpdateAlertRule(ctx, ruleID, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, existingRule.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteAlertRule(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful deletion
	t.Run("Successful Deletion", func(t *testing.T) {
		ruleID := uint(1)

		// Mock get rule for existence check
		existingRule := &model.AlertRule{
			ID:   ruleID,
			Name: "Rule to Delete",
		}
		mockRepo.On("GetAlertRuleByID", ctx, ruleID).Return(existingRule, nil).Once()

		// Mock physical delete
		mockRepo.On("DeleteAlertRule", ctx, ruleID).Return(nil).Once()

		// Execute
		err := service.DeleteAlertRule(ctx, ruleID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

func TestListAlertRecords(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful list with filters
	// t.Run("Successful List", func(t *testing.T) {
	// 	ruleID := uint(1)
	// 	req := &request.AlertRecordQueryRequest{
	// 		Page:     1,
	// 		PageSize: 10,
	// 	}

	// 	// Calculate offset
	// 	offset := (req.Page - 1) * req.PageSize

	// 	// Expected filters
	// 	expectedFilters := map[string]interface{}{
	// 		"status":        req.Status,
	// 		"severity":      req.Severity,
	// 		"alert_rule_id": req.AlertRuleID,
	// 		"start_time":    req.StartTime,
	// 		"end_time":      req.EndTime,
	// 	}

	// 	// Mock records
	// 	records := []*model.AlertRecord{
	// 		{
	// 			ID:          uint(1),
	// 			AlertRuleID: ruleID,
	// 			Title:       "CPU Usage High",
	// 			Description: "CPU usage exceeded 90%",
	// 			Status:      "FIRING",
	// 			FiredAt:     time.Now().Add(-1 * time.Hour),
	// 		},
	// 		{
	// 			ID:          uint(2),
	// 			AlertRuleID: ruleID,
	// 			Title:       "Memory Usage High",
	// 			Description: "Memory usage exceeded 80%",
	// 			Status:      "FIRING",
	// 			FiredAt:     time.Now().Add(-30 * time.Minute),
	// 		},
	// 	}

	// 	mockRepo.On("ListAlertRecords", ctx, offset, req.PageSize, mock.MatchedBy(func(filters map[string]interface{}) bool {
	// 		return filters["status"] == expectedFilters["status"] &&
	// 			filters["severity"] == expectedFilters["severity"] &&
	// 			filters["alert_rule_id"] == expectedFilters["alert_rule_id"]
	// 	})).Return(records, int64(2), nil).Once()

	// 	// Execute
	// 	result, err := service.ListAlertRecords(ctx, req)

	// 	// Assert
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, result)
	// 	assert.Equal(t, int64(2), int64(2))
	// 	mockRepo.AssertExpectations(t)
	// })

	// Test case 2: Default pagination
	t.Run("Default Pagination", func(t *testing.T) {
		req := &request.AlertRecordQueryRequest{
			Page:     0, // Should be set to default 1
			PageSize: 0, // Should be set to default 20
		}

		mockRepo.On("ListAlertRecords", ctx, 0, 20, mock.Anything).
			Return([]*model.AlertRecord{}, int64(0), nil).Once()

		// Execute
		result, err := service.ListAlertRecords(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		mockRepo.AssertExpectations(t)
	})

}

func TestAcknowledgeAlertRecord(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful acknowledgment
	t.Run("Successful Acknowledgment", func(t *testing.T) {
		recordID := uint(1)
		userID := uint(10)
		req := &request.AlertAcknowledgeRequest{
			Note: "Investigating the issue",
		}

		// Mock get record
		record := &model.AlertRecord{
			ID:     recordID,
			Title:  "CPU Usage High",
			Status: "FIRING",
		}
		mockRepo.On("GetAlertRecordByID", ctx, recordID).Return(record, nil).Once()

		// Mock update
		mockRepo.On("UpdateAlertRecord", ctx, recordID, mock.MatchedBy(func(updates map[string]interface{}) bool {
			_, hasAckTime := updates["acknowledged_at"]
			return updates["acknowledged_by"] == userID &&
				updates["status"] == "CONFIRMED" &&
				updates["resolution_note"] == req.Note &&
				hasAckTime
		})).Return(nil).Once()

		// Execute
		err := service.AcknowledgeAlertRecord(ctx, recordID, userID, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

func TestResolveAlertRecord(t *testing.T) {
	service, mockRepo, _ := setupAlertService()
	ctx := context.Background()

	// Test case 1: Successful resolution
	t.Run("Successful Resolution", func(t *testing.T) {
		recordID := uint(1)
		userID := uint(10)
		req := &request.AlertResolveRequest{
			ResolutionNote: "Issue fixed by restarting the service",
		}

		// Mock get record
		record := &model.AlertRecord{
			ID:     recordID,
			Title:  "CPU Usage High",
			Status: "FIRING",
		}
		mockRepo.On("GetAlertRecordByID", ctx, recordID).Return(record, nil).Once()

		// Mock update
		mockRepo.On("UpdateAlertRecord", ctx, recordID, mock.MatchedBy(func(updates map[string]interface{}) bool {
			_, hasResolvedTime := updates["resolved_at"]
			return updates["acknowledged_by"] == userID &&
				updates["status"] == "RESOLVED" &&
				updates["resolution_note"] == req.ResolutionNote &&
				hasResolvedTime
		})).Return(nil).Once()

		// Execute
		err := service.ResolveAlertRecord(ctx, recordID, userID, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

}

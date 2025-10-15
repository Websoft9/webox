package service

import (
	"context"
	"encoding/json"
	"fmt"

	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/email"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

type notificationChannelService struct {
	channelRepo  repository.NotificationChannelRepository
	logger       logger.Logger
	emailService email.EmailService
}

// NewNotificationChannelService creates a new notification channel service instance
func NewNotificationChannelService(
	channelRepo repository.NotificationChannelRepository,
	logger logger.Logger,
	emailService email.EmailService,
) service.NotificationChannelService {
	return &notificationChannelService{
		channelRepo:  channelRepo,
		logger:       logger,
		emailService: emailService,
	}
}

// toStringPtr returns a pointer to the string value
func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetChannelList retrieves notification channels list with pagination and filtering
func (s *notificationChannelService) GetChannelList(
	ctx context.Context,
	req *request.GetNotificationChannelListRequest,
) (*common.PaginationResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification channel list",
		logger.String("service", "notification_channel"),
		logger.String("operation", "GetChannelList"),
		logger.Int("page", req.PaginationRequest.Page),
		logger.Int("page_size", req.PaginationRequest.PageSize))

	// Get channels from repository
	channels, total, err := s.channelRepo.GetList(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channels from repository", logger.ErrorField(err))
		return nil, err
	}

	// Convert to response format
	items := make([]response.NotificationChannelResponse, 0, len(channels))
	for _, channel := range channels {
		resp := response.NotificationChannelResponse(*channel)
		items = append(items, resp)
	}

	// Calculate total pages
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	result := &common.PaginationResponse{
		Page:       req.GetPage(),
		PageSize:   req.GetPageSize(),
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
	}
	return result, nil
}

// GetChannelByCode retrieves a specific notification channel by code
func (s *notificationChannelService) GetChannelByCode(
	ctx context.Context,
	code string,
) (*response.NotificationChannelResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification channel by code",
		logger.String("service", "notification_channel"),
		logger.String("operation", "GetChannelByCode"),
		logger.String("code", code))

	// Get channel from repository
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel from repository", logger.ErrorField(err))
		return nil, err
	}
	resp := response.NotificationChannelResponse(*channel)
	return &resp, nil
}

// CreateEmailChannel creates a new email notification channel
func (s *notificationChannelService) CreateEmailChannel(
	ctx context.Context,
	req *request.CreateEmailChannelRequest,
	userID uint,
) (*response.NotificationChannelResponse, error) {
	s.logger.InfoContext(ctx, "Creating email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "CreateEmailChannel"),
		logger.String("code", req.Code),
		logger.Uint("user_id", userID))

	// Check if code already exists
	exists, err := s.channelRepo.CheckCodeExists(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check code existence", logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Convert to map using JSON marshal/unmarshal (simplified)
	configBytes, err := json.Marshal(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal email config", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeInternalError)
	}

	var channelConfig model.JSON
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeInternalError)
	}

	// Create channel model
	channel := &model.NotificationChannelConfig{
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		ChannelType:   constants.NotificationChannelEmail,
		ChannelConfig: channelConfig,
		OwnerID:       userID,
		Status:        1, // Enabled by default
	}

	// Save to repository
	if err := s.channelRepo.Create(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create email channel in repository", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Email notification channel created successfully",
		logger.String("code", req.Code),
		logger.Uint("channel_id", channel.ID))

	resp := response.NotificationChannelResponse(*channel)
	return &resp, nil
}

// CreateWebhookChannel creates a new webhook notification channel
func (s *notificationChannelService) CreateWebhookChannel(
	ctx context.Context,
	req *request.CreateWebhookChannelRequest,
	userID uint,
) (*response.NotificationChannelResponse, error) {
	s.logger.InfoContext(ctx, "Creating webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "CreateWebhookChannel"),
		logger.String("code", req.Code),
		logger.Uint("user_id", userID))

	// Check if code already exists
	exists, err := s.channelRepo.CheckCodeExists(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check code existence", logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
	}

	// Convert to map using JSON marshal/unmarshal (simplified)
	configBytes, err := json.Marshal(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal webhook config", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	var channelConfig model.JSON
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	// Create channel model
	channel := &model.NotificationChannelConfig{
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		ChannelType:   constants.NotificationChannelWebhook,
		ChannelConfig: channelConfig,
		OwnerID:       userID,
		Status:        1, // Enabled by default
	}

	// Save to repository
	if err := s.channelRepo.Create(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create webhook channel in repository", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Webhook notification channel created successfully",
		logger.String("code", req.Code),
		logger.Uint("channel_id", channel.ID))

	resp := response.NotificationChannelResponse(*channel)
	return &resp, nil
}

// updateEmailConfigFields updates email config fields from request
func (s *notificationChannelService) updateEmailConfigFields(
	existingConfig *model.EmailConfig,
	req *request.UpdateEmailChannelRequest,
) model.EmailConfig {
	emailConfig := *existingConfig

	if req.SMTPHost != nil {
		emailConfig.SMTPHost = *req.SMTPHost
	}
	if req.SMTPPort != nil {
		emailConfig.SMTPPort = *req.SMTPPort
	}
	if req.SMTPSecurity != nil {
		emailConfig.SMTPSecurity = *req.SMTPSecurity
	}
	if req.SMTPTimeout != nil {
		emailConfig.SMTPTimeout = *req.SMTPTimeout
	}
	if req.SMTPUsername != nil {
		emailConfig.SMTPUsername = *req.SMTPUsername
	}
	if req.SMTPPassword != nil {
		emailConfig.SMTPPassword = *req.SMTPPassword
	}
	if req.SenderEmail != nil {
		emailConfig.SenderEmail = *req.SenderEmail
	}
	if req.SenderName != nil {
		emailConfig.SenderName = *req.SenderName
	}
	if req.Encoding != nil {
		emailConfig.Encoding = *req.Encoding
	}
	if req.RateLimit != nil {
		emailConfig.RateLimit = *req.RateLimit
	}
	if req.RetryCount != nil {
		emailConfig.RetryCount = *req.RetryCount
	}
	if req.QuietHours != nil {
		emailConfig.QuietHours = *req.QuietHours
	}

	return emailConfig
}

// UpdateEmailChannel updates an existing email notification channel
func (s *notificationChannelService) UpdateEmailChannel(
	ctx context.Context,
	code string,
	req *request.UpdateEmailChannelRequest,
	userID uint,
) (*response.NotificationChannelResponse, error) {
	s.logger.InfoContext(ctx, "Updating email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "UpdateEmailChannel"),
		logger.String("code", code),
		logger.Uint("user_id", userID))

	// Get existing channel
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, err
	}

	// Check if user owns this channel
	if channel.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Check if channel is email type
	if channel.ChannelType != constants.NotificationChannelEmail {
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Create updated email config (merge with existing config)
	var existingConfig model.EmailConfig
	if channel.ChannelConfig != nil {
		configBytes, _ := json.Marshal(channel.ChannelConfig)
		if unmarshalErr := json.Unmarshal(configBytes, &existingConfig); unmarshalErr != nil {
			s.logger.WarnContext(ctx, "Failed to unmarshal existing email config", logger.ErrorField(unmarshalErr))
			// Continue with default config
		}
	}

	// Apply updates only for non-nil fields
	if req.Name != nil {
		channel.Name = *req.Name
	}
	if req.Description != nil {
		channel.Description = toStringPtr(*req.Description)
	}

	// Update email config using helper function
	emailConfig := s.updateEmailConfigFields(&existingConfig, req)

	// Convert to map using JSON marshal/unmarshal
	configBytes, err := json.Marshal(emailConfig)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal email config", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	var channelConfig model.JSON
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	channel.ChannelConfig = channelConfig

	// Update in repository
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update email channel in repository", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeRecordUpdateFailed)
	}

	s.logger.InfoContext(ctx, "Email notification channel updated successfully",
		logger.String("code", code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	return &response.NotificationChannelResponse{
		ID:            channel.ID,
		Code:          channel.Code,
		Name:          channel.Name,
		Description:   channel.Description,
		ChannelType:   channel.ChannelType,
		ChannelConfig: channel.ChannelConfig,
		OwnerID:       channel.OwnerID,
		Status:        channel.Status,
		CreatedAt:     channel.CreatedAt,
		UpdatedAt:     channel.UpdatedAt,
	}, nil
}

// updateWebhookConfigFields updates webhook config fields from request
func (s *notificationChannelService) updateWebhookConfigFields(
	existingConfig *model.WebhookConfig,
	req *request.UpdateWebhookChannelRequest,
) model.WebhookConfig {
	webhookConfig := *existingConfig

	if req.URL != nil {
		webhookConfig.URL = *req.URL
	}
	if req.Method != nil {
		webhookConfig.Method = *req.Method
	}
	if req.Headers != nil {
		webhookConfig.Headers = *req.Headers
	}
	if req.Secret != nil {
		webhookConfig.Secret = *req.Secret
	}
	if req.Timeout != nil {
		webhookConfig.Timeout = *req.Timeout
	}
	if req.Encoding != nil {
		webhookConfig.Encoding = *req.Encoding
	}
	if req.RateLimit != nil {
		webhookConfig.RateLimit = *req.RateLimit
	}
	if req.RetryCount != nil {
		webhookConfig.RetryCount = *req.RetryCount
	}
	if req.QuietHours != nil {
		webhookConfig.QuietHours = *req.QuietHours
	}

	return webhookConfig
}

// UpdateWebhookChannel updates an existing webhook notification channel
func (s *notificationChannelService) UpdateWebhookChannel(
	ctx context.Context,
	code string,
	req *request.UpdateWebhookChannelRequest,
	userID uint,
) (*response.NotificationChannelResponse, error) {
	s.logger.InfoContext(ctx, "Updating webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "UpdateWebhookChannel"),
		logger.String("code", code),
		logger.Uint("user_id", userID))

	// Get existing channel
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, err
	}

	// Check if channel is webhook type
	if channel.ChannelType != constants.NotificationChannelWebhook {
		return nil, errors.NewAppError(errors.CodeValidationFailed)
	}

	// Create updated webhook config (merge with existing config)
	var existingConfig model.WebhookConfig
	if channel.ChannelConfig != nil {
		configBytes, _ := json.Marshal(channel.ChannelConfig)
		if unmarshalErr := json.Unmarshal(configBytes, &existingConfig); unmarshalErr != nil {
			s.logger.WarnContext(ctx, "Failed to unmarshal existing webhook config", logger.ErrorField(unmarshalErr))
			// Continue with default config
		}
	}

	// Apply updates only for non-nil fields
	if req.Name != nil {
		channel.Name = *req.Name
	}
	if req.Description != nil {
		channel.Description = toStringPtr(*req.Description)
	}

	// Update webhook config using helper function
	webhookConfig := s.updateWebhookConfigFields(&existingConfig, req)

	// Convert to map using JSON marshal/unmarshal
	configBytes, err := json.Marshal(webhookConfig)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal webhook config", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	var channelConfig model.JSON
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	channel.ChannelConfig = channelConfig

	// Update in repository
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update webhook channel in repository", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Webhook notification channel updated successfully",
		logger.String("code", code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	resp := response.NotificationChannelResponse(*channel)
	return &resp, nil
}

// DeleteChannel hard deletes a notification channel
func (s *notificationChannelService) DeleteChannel(
	ctx context.Context,
	code string,
	userID uint,
) error {
	s.logger.InfoContext(ctx, "Deleting notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "DeleteChannel"),
		logger.String("code", code),
		logger.Uint("user_id", userID))

	// Get channel to verify ownership
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return err
	}

	// Delete the channel
	if err := s.channelRepo.Delete(ctx, channel.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete notification channel", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Notification channel deleted successfully",
		logger.String("code", code))

	return nil
}

// TestEmailChannel tests email channel configuration
func (s *notificationChannelService) TestEmailChannel(
	ctx context.Context,
	req *request.TestEmailChannelRequest,
) error {
	s.logger.InfoContext(ctx, "Testing email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "TestEmailChannel"),
		logger.String("code", req.Code))

	// Get the channel configuration
	channel, err := s.channelRepo.GetByCode(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return err
	}

	// Check if channel is email type
	if channel.ChannelType != constants.NotificationChannelEmail {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// Test email configuration by sending a test email
	testSubject := req.Subject
	if testSubject == "" {
		testSubject = "Websoft9 Email Channel Test"
	}

	testContent := req.Content
	if testContent == "" {
		testContent = fmt.Sprintf("This is a test email from Websoft9 notification channel: %s", channel.Name)
	}

	// Use EmailService to send test email
	if err := s.emailService.SendEmail(ctx, req.Recipient, testSubject, testContent); err != nil {
		s.logger.ErrorContext(ctx, "Email channel test failed", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Email channel test completed successfully",
		logger.String("code", req.Code),
		logger.String("recipient", req.Recipient))

	return nil
}

// TestWebhookChannel tests webhook channel configuration
func (s *notificationChannelService) TestWebhookChannel(
	ctx context.Context,
	req *request.TestWebhookChannelRequest,
) error {
	s.logger.InfoContext(ctx, "Testing webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "TestWebhookChannel"),
		logger.String("code", req.Code))

	// Get the channel configuration
	channel, err := s.channelRepo.GetByCode(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return err
	}

	// Check if channel is webhook type
	if channel.ChannelType != constants.NotificationChannelWebhook {
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	// Parse webhook configuration
	var webhookConfig model.WebhookConfig
	if channel.ChannelConfig != nil {
		configBytes, marshalErr := json.Marshal(channel.ChannelConfig)
		if marshalErr != nil {
			s.logger.ErrorContext(ctx, "Failed to marshal channel config", logger.ErrorField(marshalErr))
			return errors.NewAppErrorWrapError(marshalErr, errors.CodeInternalError)
		}

		if unmarshalErr := json.Unmarshal(configBytes, &webhookConfig); unmarshalErr != nil {
			s.logger.ErrorContext(ctx, "Failed to unmarshal webhook config", logger.ErrorField(unmarshalErr))
			return errors.NewAppErrorWrapError(unmarshalErr, errors.CodeInternalError)
		}
	}

	// Test webhook configuration by sending a test request
	// var testPayload map[string]interface{}
	// if req.Payload == nil {
	// 	testPayload = map[string]interface{}{
	// 		"type":      "test",
	// 		"channel":   channel.Name,
	// 		"message":   "This is a test webhook from Websoft9 notification channel",
	// 		"timestamp": time.Now().Format(time.RFC3339),
	// 	}
	// } else {
	// 	// Convert interface{} to map[string]interface{}
	// 	if payload, ok := req.Payload.(map[string]interface{}); ok {
	// 		testPayload = payload
	// 	} else {
	// 		// If it's not a map, create a wrapper
	// 		testPayload = map[string]interface{}{
	// 			"data": req.Payload,
	// 		}
	// 	}
	// }

	// Send test webhook request

	return nil
}

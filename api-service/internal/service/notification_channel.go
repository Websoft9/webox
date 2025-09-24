package service

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"gorm.io/gorm"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

const (
	// Default timeout for network operations
	DefaultNetworkTimeout = 30 * time.Second
	// Maximum response body size for webhook testing
	MaxWebhookResponseSize = 1024
)

type notificationChannelService struct {
	channelRepo repository.NotificationChannelRepository
	logger      logger.Logger
	i18n        *i18n.I18n
}

// NewNotificationChannelService creates a new notification channel service instance
func NewNotificationChannelService(
	channelRepo repository.NotificationChannelRepository,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.NotificationChannelService {
	return &notificationChannelService{
		channelRepo: channelRepo,
		logger:      logger,
		i18n:        i18n,
	}
}

// toStringPtr returns a pointer to the string value
func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// modelToResponse converts a NotificationChannelConfig model to NotificationChannelResponse
func modelToResponse(channel *model.NotificationChannelConfig) response.NotificationChannelResponse {
	return response.NotificationChannelResponse{
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
	}
}

// GetChannelList retrieves notification channels list with pagination and filtering
func (s *notificationChannelService) GetChannelList(
	ctx context.Context,
	req *request.GetNotificationChannelListRequest,
) (*response.NotificationChannelListResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification channel list",
		logger.String("service", "notification_channel"),
		logger.String("operation", "GetChannelList"),
		logger.Int("page", req.PaginationRequest.Page),
		logger.Int("page_size", req.PaginationRequest.PageSize))

	// Get channels from repository
	channels, total, err := s.channelRepo.GetList(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channels from repository", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.list_failed"))
	}

	// Convert to response format
	items := make([]response.NotificationChannelResponse, 0, len(channels))
	for _, channel := range channels {
		items = append(items, modelToResponse(channel))
	}

	// Calculate total pages
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &response.NotificationChannelListResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
	}, nil
}

// GetChannelByCode retrieves a specific notification channel by code
func (s *notificationChannelService) GetChannelByCode(
	ctx context.Context,
	code string,
) (*response.NotificationChannelDetailResponse, error) {
	s.logger.InfoContext(ctx, "Getting notification channel by code",
		logger.String("service", "notification_channel"),
		logger.String("operation", "GetChannelByCode"),
		logger.String("code", code))

	// Get channel from repository
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel from repository", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	return &response.NotificationChannelDetailResponse{
		NotificationChannelResponse: modelToResponse(channel),
		// TODO: Add owner name lookup if needed
		OwnerName: nil,
	}, nil
}

// CreateEmailChannel creates a new email notification channel
func (s *notificationChannelService) CreateEmailChannel(
	ctx context.Context,
	req *request.CreateEmailChannelRequest,
	userID uint,
) (*response.NotificationChannelDetailResponse, error) {
	s.logger.InfoContext(ctx, "Creating email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "CreateEmailChannel"),
		logger.String("code", req.Code),
		logger.Uint("user_id", userID))

	// Check if code already exists
	exists, err := s.channelRepo.CheckCodeExists(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check code existence", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.check_code_failed"))
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists, s.i18n.T(ctx, "notification.channel.code_exists"))
	}

	// Create email config - convert to map for storage
	emailConfig := model.EmailConfig{
		SMTPHost:     req.SMTPHost,
		SMTPPort:     req.SMTPPort,
		SMTPSecurity: req.SMTPSecurity,
		SMTPTimeout:  req.SMTPTimeout,
		SMTPUsername: req.SMTPUsername,
		SMTPPassword: req.SMTPPassword,
		SenderEmail:  req.SenderEmail,
		SenderName:   req.SenderName,
		Encoding:     req.Encoding,
		RateLimit:    req.RateLimit,
		RetryCount:   req.RetryCount,
		QuietHours:   req.QuietHours,
	}

	// Convert to map using JSON marshal/unmarshal (simplified)
	configBytes, err := json.Marshal(emailConfig)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal email config", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_marshal_failed"))
	}

	var channelConfig model.JSONChannelConfig
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_convert_failed"))
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
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, s.i18n.T(ctx, "notification.channel.create_failed"))
	}

	s.logger.InfoContext(ctx, "Email notification channel created successfully",
		logger.String("code", req.Code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	return &response.NotificationChannelDetailResponse{
		NotificationChannelResponse: modelToResponse(channel),
		OwnerName:                   nil, // TODO: Add owner name lookup if needed
	}, nil
}

// CreateWebhookChannel creates a new webhook notification channel
func (s *notificationChannelService) CreateWebhookChannel(
	ctx context.Context,
	req *request.CreateWebhookChannelRequest,
	userID uint,
) (*response.NotificationChannelDetailResponse, error) {
	s.logger.InfoContext(ctx, "Creating webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "CreateWebhookChannel"),
		logger.String("code", req.Code),
		logger.Uint("user_id", userID))

	// Check if code already exists
	exists, err := s.channelRepo.CheckCodeExists(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check code existence", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.check_code_failed"))
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists, s.i18n.T(ctx, "notification.channel.code_exists"))
	}

	// Create webhook config - convert to map for storage
	webhookConfig := model.WebhookConfig{
		URL:        req.URL,
		Method:     req.Method,
		Headers:    req.Headers,
		Secret:     req.Secret,
		Timeout:    req.Timeout,
		Encoding:   req.Encoding,
		RateLimit:  req.RateLimit,
		RetryCount: req.RetryCount,
		QuietHours: req.QuietHours,
	}

	// Convert to map using JSON marshal/unmarshal (simplified)
	configBytes, err := json.Marshal(webhookConfig)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal webhook config", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_marshal_failed"))
	}

	var channelConfig model.JSONChannelConfig
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_convert_failed"))
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
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, s.i18n.T(ctx, "notification.channel.create_failed"))
	}

	s.logger.InfoContext(ctx, "Webhook notification channel created successfully",
		logger.String("code", req.Code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	return &response.NotificationChannelDetailResponse{
		NotificationChannelResponse: modelToResponse(channel),
		OwnerName:                   nil, // TODO: Add owner name lookup if needed
	}, nil
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
) (*response.NotificationChannelDetailResponse, error) {
	s.logger.InfoContext(ctx, "Updating email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "UpdateEmailChannel"),
		logger.String("code", code),
		logger.Uint("user_id", userID))

	// Get existing channel
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	// Check if user owns this channel
	if channel.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "common.access_denied"))
	}

	// Check if channel is email type
	if channel.ChannelType != constants.NotificationChannelEmail {
		return nil, errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "notification.channel.invalid_type"))
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
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_marshal_failed"))
	}

	var channelConfig model.JSONChannelConfig
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_convert_failed"))
	}

	channel.ChannelConfig = channelConfig

	// Update in repository
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update email channel in repository", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "notification.channel.update_failed"))
	}

	s.logger.InfoContext(ctx, "Email notification channel updated successfully",
		logger.String("code", code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	return &response.NotificationChannelDetailResponse{
		NotificationChannelResponse: modelToResponse(channel),
		OwnerName:                   nil, // TODO: Add owner name lookup if needed
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
) (*response.NotificationChannelDetailResponse, error) {
	s.logger.InfoContext(ctx, "Updating webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "UpdateWebhookChannel"),
		logger.String("code", code),
		logger.Uint("user_id", userID))

	// Get existing channel
	channel, err := s.channelRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	// Check if user owns this channel
	if channel.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "common.access_denied"))
	}

	// Check if channel is webhook type
	if channel.ChannelType != constants.NotificationChannelWebhook {
		return nil, errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "notification.channel.invalid_type"))
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
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_marshal_failed"))
	}

	var channelConfig model.JSONChannelConfig
	if err := json.Unmarshal(configBytes, &channelConfig); err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal to JSONChannelConfig", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeInternalError, s.i18n.T(ctx, "notification.channel.config_convert_failed"))
	}

	channel.ChannelConfig = channelConfig

	// Update in repository
	if err := s.channelRepo.Update(ctx, channel); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update webhook channel in repository", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "notification.channel.update_failed"))
	}

	s.logger.InfoContext(ctx, "Webhook notification channel updated successfully",
		logger.String("code", code),
		logger.Uint("channel_id", channel.ID))

	// Return response in standard format (consistent with list and detail)
	return &response.NotificationChannelDetailResponse{
		NotificationChannelResponse: modelToResponse(channel),
		OwnerName:                   nil, // TODO: Add owner name lookup if needed
	}, nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	// Check if user owns this channel
	if channel.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "common.access_denied"))
	}

	// Delete the channel
	if err := s.channelRepo.Delete(ctx, channel.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete notification channel", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordDeleteFailed, s.i18n.T(ctx, "notification.channel.delete_failed"))
	}

	s.logger.InfoContext(ctx, "Notification channel deleted successfully",
		logger.String("code", code))

	return nil
}

// TestEmailChannel tests email channel configuration
func (s *notificationChannelService) TestEmailChannel(
	ctx context.Context,
	req *request.TestEmailChannelRequest,
) (*response.TestChannelResponse, error) {
	s.logger.InfoContext(ctx, "Testing email notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "TestEmailChannel"),
		logger.String("code", req.Code))

	// Get the channel configuration
	channel, err := s.channelRepo.GetByCode(ctx, req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	// Check if channel is email type
	if channel.ChannelType != constants.NotificationChannelEmail {
		return &response.TestChannelResponse{
			Success: false,
			Message: s.i18n.T(ctx, "notification.channel.invalid_type"),
			Details: "Channel is not an email type",
		}, nil
	}

	// Parse email configuration
	var emailConfig model.EmailConfig
	if channel.ChannelConfig != nil {
		configBytes, err := json.Marshal(channel.ChannelConfig)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to marshal channel config", logger.ErrorField(err))
			return &response.TestChannelResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.channel.config_parse_failed"),
				Details: err.Error(),
			}, nil
		}

		if err := json.Unmarshal(configBytes, &emailConfig); err != nil {
			s.logger.ErrorContext(ctx, "Failed to unmarshal email config", logger.ErrorField(err))
			return &response.TestChannelResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.channel.config_parse_failed"),
				Details: err.Error(),
			}, nil
		}
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

	if err := s.sendTestEmail(&emailConfig, req.Recipient, testSubject, testContent); err != nil {
		s.logger.ErrorContext(ctx, "Email channel test failed", logger.ErrorField(err))
		return &response.TestChannelResponse{
			Success: false,
			Message: s.i18n.T(ctx, "notification.channel.email_test_failed"),
			Details: err.Error(),
		}, nil
	}

	s.logger.InfoContext(ctx, "Email channel test completed successfully",
		logger.String("code", req.Code),
		logger.String("recipient", req.Recipient))

	return &response.TestChannelResponse{
		Success: true,
		Message: s.i18n.T(ctx, "notification.channel.email_test_success"),
		Details: fmt.Sprintf("Test email sent successfully to %s", req.Recipient),
	}, nil
}

// TestWebhookChannel tests webhook channel configuration
func (s *notificationChannelService) TestWebhookChannel(
	ctx context.Context,
	req *request.TestWebhookChannelRequest,
) (*response.TestChannelResponse, error) {
	s.logger.InfoContext(ctx, "Testing webhook notification channel",
		logger.String("service", "notification_channel"),
		logger.String("operation", "TestWebhookChannel"),
		logger.String("code", req.Code))

	// Get the channel configuration
	channel, err := s.channelRepo.GetByCode(ctx, req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.channel.not_found"))
		}
		s.logger.ErrorContext(ctx, "Failed to get notification channel", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.channel.get_failed"))
	}

	// Check if channel is webhook type
	if channel.ChannelType != constants.NotificationChannelWebhook {
		return &response.TestChannelResponse{
			Success: false,
			Message: s.i18n.T(ctx, "notification.channel.invalid_type"),
			Details: "Channel is not a webhook type",
		}, nil
	}

	// Parse webhook configuration
	var webhookConfig model.WebhookConfig
	if channel.ChannelConfig != nil {
		configBytes, marshalErr := json.Marshal(channel.ChannelConfig)
		if marshalErr != nil {
			s.logger.ErrorContext(ctx, "Failed to marshal channel config", logger.ErrorField(marshalErr))
			return &response.TestChannelResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.channel.config_parse_failed"),
				Details: marshalErr.Error(),
			}, nil
		}

		if unmarshalErr := json.Unmarshal(configBytes, &webhookConfig); unmarshalErr != nil {
			s.logger.ErrorContext(ctx, "Failed to unmarshal webhook config", logger.ErrorField(unmarshalErr))
			return &response.TestChannelResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.channel.config_parse_failed"),
				Details: unmarshalErr.Error(),
			}, nil
		}
	}

	// Test webhook configuration by sending a test request
	testPayload := req.Payload
	if testPayload == nil {
		testPayload = map[string]interface{}{
			"type":      "test",
			"channel":   channel.Name,
			"message":   "This is a test webhook from Websoft9 notification channel",
			"timestamp": time.Now().Format(time.RFC3339),
		}
	}

	statusCode, responseBody, err := s.sendTestWebhook(ctx, &webhookConfig, testPayload)
	if err != nil {
		s.logger.ErrorContext(ctx, "Webhook channel test failed", logger.ErrorField(err))
		return &response.TestChannelResponse{
			Success: false,
			Message: s.i18n.T(ctx, "notification.channel.webhook_test_failed"),
			Details: err.Error(),
		}, nil
	}

	s.logger.InfoContext(ctx, "Webhook channel test completed successfully",
		logger.String("code", req.Code),
		logger.String("url", webhookConfig.URL),
		logger.Int("status_code", statusCode))

	return &response.TestChannelResponse{
		Success: true,
		Message: s.i18n.T(ctx, "notification.channel.webhook_test_success"),
		Details: fmt.Sprintf("Webhook test successful. Status: %d, Response: %s", statusCode, responseBody),
	}, nil
}

// sendTestEmail sends a test email using the provided email configuration
func (s *notificationChannelService) sendTestEmail(
	config *model.EmailConfig,
	to, subject, body string,
) error {
	// Build the email message with proper headers
	msg := s.buildEmailMessage(config.SenderEmail, to, subject, body)

	// Create SMTP authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)

	// Construct the SMTP server address
	addr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	// Send email based on security configuration
	var err error
	switch strings.ToLower(config.SMTPSecurity) {
	case "tls", "ssl":
		err = s.sendEmailWithTLS(addr, auth, config.SenderEmail, []string{to}, msg, config.SMTPHost)
	case "none", "":
		err = smtp.SendMail(addr, auth, config.SenderEmail, []string{to}, msg)
	default:
		err = smtp.SendMail(addr, auth, config.SenderEmail, []string{to}, msg)
	}

	return err
}

// sendEmailWithTLS sends email using TLS encryption
func (s *notificationChannelService) sendEmailWithTLS(
	addr string,
	auth smtp.Auth,
	from string,
	to []string,
	msg []byte,
	serverName string,
) error {
	// Create TLS configuration
	tlsConfig := &tls.Config{
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}

	// Establish TLS connection
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, serverName)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer func() {
		if quitErr := client.Quit(); quitErr != nil {
			// Log error but don't fail the operation
			s.logger.WarnContext(context.Background(), "Failed to quit SMTP client", logger.ErrorField(quitErr))
		}
	}()

	// Authenticate if auth is provided
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		// Attempt to close the writer and surface both errors if closing also fails.
		if cerr := writer.Close(); cerr != nil {
			return fmt.Errorf("failed to write message: %v; additionally failed to close writer: %w", err, cerr)
		}
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close message writer: %w", err)
	}

	return nil
}

// buildEmailMessage builds the email message with proper headers
func (s *notificationChannelService) buildEmailMessage(from, to, subject, body string) []byte {
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body

	return []byte(msg)
}

// sendTestWebhook sends a test webhook request using the provided webhook configuration
func (s *notificationChannelService) sendTestWebhook(
	ctx context.Context,
	config *model.WebhookConfig,
	payload interface{},
) (statusCode int, responseBody string, err error) {
	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Set default method if not specified
	method := config.Method
	if method == "" {
		method = "POST"
	}

	// Set timeout
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = DefaultNetworkTimeout // Default timeout
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, config.URL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set default content type
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Add signature if secret is provided
	if config.Secret != "" {
		// Simple signature - in production, you might want to use HMAC
		req.Header.Set("X-Webhook-Signature", config.Secret)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBodyBytes := make([]byte, MaxWebhookResponseSize) // Limit response size for testing
	n, _ := resp.Body.Read(responseBodyBytes)

	statusCode = resp.StatusCode
	responseBody = string(responseBodyBytes[:n])
	return statusCode, responseBody, nil
}

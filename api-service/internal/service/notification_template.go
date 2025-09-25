package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
)

type notificationTemplateService struct {
	templateRepo   interfaceRepo.NotificationTemplateRepository
	channelService interfaceService.NotificationChannelService
	logger         logger.Logger
	i18n           *i18n.I18n
}

// NewNotificationTemplateService creates a new notification template service
func NewNotificationTemplateService(
	templateRepo interfaceRepo.NotificationTemplateRepository,
	channelService interfaceService.NotificationChannelService,
	logger logger.Logger,
	i18n *i18n.I18n,
) interfaceService.NotificationTemplateService {
	return &notificationTemplateService{
		templateRepo:   templateRepo,
		channelService: channelService,
		logger:         logger,
		i18n:           i18n,
	}
}

// CreateTemplate creates a new notification template
func (s *notificationTemplateService) CreateTemplate(
	ctx context.Context,
	req *request.CreateNotificationTemplateRequest,
) (*response.NotificationTemplateDetailResponse, error) {
	s.logger.InfoContext(ctx, "Creating notification template",
		logger.String("service", "notification_template"),
		logger.String("operation", "CreateTemplate"),
		logger.String("name", req.Name))

	// Validate email template must have subject
	if req.TemplateType == "EMAIL" && (req.Subject == nil || *req.Subject == "") {
		return nil, errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "notification.template.email_subject_required"))
	}

	// Check if template name already exists
	exists, err := s.templateRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check template name existence", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.template.name_check_failed"))
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists, s.i18n.T(ctx, "notification.template.name_already_exists"))
	}

	// Create template model
	template := &model.NotificationTemplate{
		Name:         req.Name,
		TemplateType: req.TemplateType,
		Subject:      req.Subject,
		Content:      req.Content,
		IsSystem:     0, // User template
		Status:       1, // Default enabled
	}

	// Save to repository
	if err := s.templateRepo.Create(ctx, template); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create notification template", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, s.i18n.T(ctx, "notification.template.create_failed"))
	}

	s.logger.InfoContext(ctx, "Notification template created successfully",
		logger.String("name", req.Name),
		logger.Uint("template_id", template.ID))

	return s.modelToDetailResponse(template), nil
}

// GetTemplate gets a notification template by ID
func (s *notificationTemplateService) GetTemplate(
	ctx context.Context,
	id uint,
) (*response.NotificationTemplateDetailResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification template",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.template.not_found"))
	}

	return s.modelToDetailResponse(template), nil
}

// GetTemplateList gets a paginated list of notification templates
func (s *notificationTemplateService) GetTemplateList(
	ctx context.Context,
	req *request.GetNotificationTemplateListRequest,
) (*response.NotificationTemplateListResponse, error) {
	templates, total, err := s.templateRepo.GetList(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list notification templates", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.template.list_failed"))
	}

	templateResponses := make([]response.NotificationTemplateResponse, len(templates))
	for i, template := range templates {
		templateResponses[i] = *s.modelToResponse(template)
	}

	// Calculate total pages
	pageSize := req.GetPageSize()
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &response.NotificationTemplateListResponse{
		Page:       req.Page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      templateResponses,
	}, nil
}

// UpdateTemplate updates a notification template
func (s *notificationTemplateService) UpdateTemplate(
	ctx context.Context,
	id uint,
	req *request.UpdateNotificationTemplateRequest,
) (*response.NotificationTemplateDetailResponse, error) {
	// Get existing template
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification template for update",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.template.not_found"))
	}

	// Check if system template (cannot be updated)
	if template.IsSystem == 1 {
		return nil, errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "notification.template.system_template_readonly"))
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if new name already exists (excluding current template)
		exists, err := s.templateRepo.ExistsByName(ctx, *req.Name, id)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to check template name existence for update", logger.ErrorField(err))
			return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, s.i18n.T(ctx, "notification.template.name_check_failed"))
		}
		if exists {
			return nil, errors.NewAppError(errors.CodeResourceAlreadyExists, s.i18n.T(ctx, "notification.template.name_already_exists"))
		}
		template.Name = *req.Name
	}

	if req.Content != nil {
		template.Content = *req.Content
	}

	if req.Subject != nil {
		template.Subject = req.Subject
	}

	// Validate email template must have subject
	if template.TemplateType == constants.NotificationChannelEmail && (template.Subject == nil || *template.Subject == "") {
		return nil, errors.NewAppError(errors.CodeValidationFailed, s.i18n.T(ctx, "notification.template.email_subject_required"))
	}

	// Update in repository
	if err := s.templateRepo.Update(ctx, template); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update notification template", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordUpdateFailed, s.i18n.T(ctx, "notification.template.update_failed"))
	}

	s.logger.InfoContext(ctx, "Notification template updated successfully",
		logger.Uint("template_id", id))

	return s.modelToDetailResponse(template), nil
}

// DeleteTemplate deletes a notification template
func (s *notificationTemplateService) DeleteTemplate(
	ctx context.Context,
	id uint,
) error {
	// Get existing template
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification template for deletion",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.template.not_found"))
	}

	// Check if system template (cannot be deleted)
	if template.IsSystem == 1 {
		return errors.NewAppError(errors.CodeAccessDenied, s.i18n.T(ctx, "notification.template.system_template_readonly"))
	}

	// Delete from repository
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete notification template",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordDeleteFailed, s.i18n.T(ctx, "notification.template.delete_failed"))
	}

	s.logger.InfoContext(ctx, "Notification template deleted successfully",
		logger.Uint("template_id", id))

	return nil
}

// TestTemplate tests a template by validating variable data and sending through channel
func (s *notificationTemplateService) TestTemplate(
	ctx context.Context,
	id uint,
	req *request.TestNotificationTemplateRequest,
) (*response.NotificationTemplateTestResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification template for test",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordNotFound, s.i18n.T(ctx, "notification.template.not_found"))
	}

	// Validate template variables
	if validationResult := s.validateTemplateVariables(ctx, template, req); validationResult != nil {
		return validationResult, nil
	}

	// Get and validate channel
	channel, err := s.channelService.GetChannelByCode(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel for test",
			logger.String("channel_code", req.Code),
			logger.ErrorField(err))
		return s.createErrorResponse(ctx, "notification.channel.not_found", fmt.Sprintf("Channel with code '%s' not found", req.Code)), nil
	}

	// Validate template type matches channel type
	if template.TemplateType != channel.ChannelType {
		return s.createErrorResponse(ctx, "notification.template.type_mismatch",
			fmt.Sprintf("Template type '%s' does not match channel type '%s'", template.TemplateType, channel.ChannelType)), nil
	}

	// Send test notification
	return s.sendTestNotification(ctx, template, channel, req)
}

// Helper methods

func (s *notificationTemplateService) modelToResponse(template *model.NotificationTemplate) *response.NotificationTemplateResponse {
	return &response.NotificationTemplateResponse{
		ID:           template.ID,
		Name:         template.Name,
		TemplateType: template.TemplateType,
		Subject:      template.Subject,
		IsSystem:     template.IsSystem,
		Status:       template.Status,
		CreatedAt:    template.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    template.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *notificationTemplateService) modelToDetailResponse(template *model.NotificationTemplate) *response.NotificationTemplateDetailResponse {
	return &response.NotificationTemplateDetailResponse{
		NotificationTemplateResponse: response.NotificationTemplateResponse{
			ID:           template.ID,
			Name:         template.Name,
			TemplateType: template.TemplateType,
			Subject:      template.Subject,
			IsSystem:     template.IsSystem,
			Status:       template.Status,
			CreatedAt:    template.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    template.UpdatedAt.Format(time.RFC3339),
		},
		Content: template.Content,
	}
}

// renderTemplate replaces template variables with provided data
// Variables are expected to be in format {{variable_name}}
func (s *notificationTemplateService) renderTemplate(template string, data map[string]interface{}) string {
	re := regexp.MustCompile(`{{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*}}`)
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// Extract variable name from {{variable_name}}
		varName := re.FindStringSubmatch(match)
		if len(varName) > 1 {
			if value, exists := data[varName[1]]; exists {
				return fmt.Sprintf("%v", value)
			}
		}
		// If variable not found, keep original placeholder
		return match
	})
	return result
}

// extractTemplateVariables extracts variable names from template content
// Variables are expected to be in format {{variable_name}}
func (s *notificationTemplateService) extractTemplateVariables(template string) []string {
	// Regular expression to match {{variable_name}} patterns
	re := regexp.MustCompile(`{{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*}}`)
	matches := re.FindAllStringSubmatch(template, -1)

	// Use map to avoid duplicates
	varMap := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			varMap[match[1]] = true
		}
	}

	// Convert map to slice
	variables := make([]string, 0, len(varMap))
	for varName := range varMap {
		variables = append(variables, varName)
	}

	return variables
}

// validateVariableData validates that provided variable data matches template requirements
func (s *notificationTemplateService) validateVariableData(providedData map[string]interface{}, requiredVars []string) error {
	// Check if all required variables are provided
	missingVars := []string{}
	for _, requiredVar := range requiredVars {
		if _, exists := providedData[requiredVar]; !exists {
			missingVars = append(missingVars, requiredVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required variables: %v", missingVars)
	}

	// Check for extra variables (optional warning)
	for providedVar := range providedData {
		found := false
		for _, requiredVar := range requiredVars {
			if providedVar == requiredVar {
				found = true
				break
			}
		}
		// Extra variables are allowed but not used
		_ = found
	}

	// Extra variables are allowed but not used

	return nil
}

// validateTemplateVariables validates template variables and returns error response if validation fails
func (s *notificationTemplateService) validateTemplateVariables(
	ctx context.Context,
	template *model.NotificationTemplate,
	req *request.TestNotificationTemplateRequest,
) *response.NotificationTemplateTestResponse {
	// Extract template variables
	requiredVars := s.extractTemplateVariables(template.Content)
	if template.Subject != nil {
		subjectVars := s.extractTemplateVariables(*template.Subject)
		// Merge subject variables with content variables
		for _, v := range subjectVars {
			found := false
			for _, existing := range requiredVars {
				if existing == v {
					found = true
					break
				}
			}
			if !found {
				requiredVars = append(requiredVars, v)
			}
		}
	}

	// Validate variable data is provided and required
	if req.VariableData == nil {
		if len(requiredVars) > 0 {
			return &response.NotificationTemplateTestResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.template.variable_data_required"),
				Details: fmt.Sprintf("Template requires variables: %v", requiredVars),
			}
		}
	} else {
		// Validate all required variables are provided
		validationErr := s.validateVariableData(req.VariableData, requiredVars)
		if validationErr != nil {
			return &response.NotificationTemplateTestResponse{
				Success: false,
				Message: s.i18n.T(ctx, "notification.template.variable_validation_failed"),
				Details: validationErr.Error(),
			}
		}
	}

	return nil
}

// createErrorResponse creates a standardized error response
func (s *notificationTemplateService) createErrorResponse(ctx context.Context, messageKey, details string) *response.NotificationTemplateTestResponse {
	return &response.NotificationTemplateTestResponse{
		Success: false,
		Message: s.i18n.T(ctx, messageKey),
		Details: details,
	}
}

// sendTestNotification sends test notification through the specified channel
func (s *notificationTemplateService) sendTestNotification(
	ctx context.Context,
	template *model.NotificationTemplate,
	channel *response.NotificationChannelDetailResponse,
	req *request.TestNotificationTemplateRequest,
) (*response.NotificationTemplateTestResponse, error) {
	// Render template content with variable data
	renderedContent := template.Content
	var renderedSubject *string
	if template.Subject != nil {
		subject := *template.Subject
		renderedSubject = &subject
	}

	if req.VariableData != nil {
		renderedContent = s.renderTemplate(template.Content, req.VariableData)
		if template.Subject != nil {
			subject := s.renderTemplate(*template.Subject, req.VariableData)
			renderedSubject = &subject
		}
	}

	// Send test notification based on channel type
	var testErr error
	switch channel.ChannelType {
	case constants.NotificationChannelEmail:
		testReq := &request.TestEmailChannelRequest{
			Code:      req.Code,
			Recipient: req.Recipient,
			Subject:   "",
			Content:   renderedContent,
		}
		if renderedSubject != nil {
			testReq.Subject = *renderedSubject
		}
		_, testErr = s.channelService.TestEmailChannel(ctx, testReq)
	default:
		return s.createErrorResponse(ctx, "notification.channel.type_not_supported",
			fmt.Sprintf("Channel type '%s' is not supported for template testing", channel.ChannelType)), nil
	}

	if testErr != nil {
		s.logger.ErrorContext(ctx, "Failed to send test notification",
			logger.String("channel_code", req.Code),
			logger.String("template_type", template.TemplateType),
			logger.String("recipient", req.Recipient),
			logger.ErrorField(testErr))

		return s.createErrorResponse(ctx, "notification.template.test_send_failed", testErr.Error()), nil
	}

	s.logger.InfoContext(ctx, "Test notification sent successfully",
		logger.String("channel_code", req.Code),
		logger.String("template_type", template.TemplateType),
		logger.String("recipient", req.Recipient))

	return &response.NotificationTemplateTestResponse{
		Success: true,
		Message: s.i18n.T(ctx, "notification.template.test_success"),
		Details: fmt.Sprintf("Test notification sent successfully to %s via channel %s", req.Recipient, req.Code),
	}, nil
}

package service

import (
	"context"
	"fmt"
	"regexp"

	"api-service/internal/constants"
	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

type notificationTemplateService struct {
	templateRepo   interfaceRepo.NotificationTemplateRepository
	channelService interfaceService.NotificationChannelService
	logger         logger.Logger
}

// NewNotificationTemplateService creates a new notification template service
func NewNotificationTemplateService(
	templateRepo interfaceRepo.NotificationTemplateRepository,
	channelService interfaceService.NotificationChannelService,
	logger logger.Logger,
) interfaceService.NotificationTemplateService {
	return &notificationTemplateService{
		templateRepo:   templateRepo,
		channelService: channelService,
		logger:         logger,
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

	// Check if template name already exists
	exists, err := s.templateRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check template name existence", logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
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
		return nil, err
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
		return nil, err
	}

	return s.modelToDetailResponse(template), nil
}

// GetTemplateList gets a paginated list of notification templates
func (s *notificationTemplateService) GetTemplateList(
	ctx context.Context,
	req *request.GetNotificationTemplateListRequest,
) (*common.PaginationResponse, error) {
	templates, total, err := s.templateRepo.GetList(ctx, req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list notification templates", logger.ErrorField(err))
		return nil, err
	}

	templateResponses := make([]response.NotificationTemplateResponse, len(templates))
	for i, template := range templates {
		templateResponses[i] = *s.modelToResponse(template)
	}

	// Calculate total pages
	pageSize := req.GetPageSize()
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	pagination := &common.PaginationResponse{
		Items:      templateResponses,
		Total:      total,
		Page:       req.GetPage(),
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return pagination, nil
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
		return nil, err
	}

	// Check if system template (cannot be updated)
	if template.IsSystem == 1 {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if new name already exists (excluding current template)
		exists, err := s.templateRepo.ExistsByName(ctx, *req.Name, id)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to check template name existence for update", logger.ErrorField(err))
			return nil, err
		}
		if exists {
			return nil, errors.NewAppError(errors.CodeResourceAlreadyExists)
		}
		template.Name = *req.Name
	}

	if req.Content != nil {
		template.Content = *req.Content
	}

	if req.Subject != nil {
		template.Subject = req.Subject
	}

	// Update in repository
	if err := s.templateRepo.Update(ctx, template); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update notification template", logger.ErrorField(err))
		return nil, err
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
		return err
	}

	// Check if system template (cannot be deleted)
	if template.IsSystem == 1 {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	// Delete from repository
	if err := s.templateRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete notification template",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return err
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
) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification template for test",
			logger.Uint("template_id", id),
			logger.ErrorField(err))
		return err
	}

	// Validate template variables
	if validationResult := s.validateTemplateVariables(template, req); validationResult != nil {
		return validationResult
	}

	// Get and validate channel
	channel, err := s.channelService.GetChannelByCode(ctx, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get notification channel for test",
			logger.String("channel_code", req.Code),
			logger.ErrorField(err))
		return err
	}

	// Validate template type matches channel type
	if template.TemplateType != channel.ChannelType {
		return errors.NewAppErrorWrapError(fmt.Errorf("template type '%s' does not match channel type '%s'", template.TemplateType, channel.ChannelType), errors.CodeValidationFailed)
	}

	// Send test notification
	if err := s.sendTestNotification(ctx, template, channel.Code, req); err != nil {
		return errors.NewAppErrorWrapError(err, errors.CodeNotificationChannelEmailTestFailed)
	}

	return nil
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
		CreatedAt:    template.CreatedAt,
		UpdatedAt:    template.UpdatedAt,
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
			CreatedAt:    template.CreatedAt,
			UpdatedAt:    template.UpdatedAt,
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

// validateTemplateVariables validates template variables and returns error if validation fails
func (s *notificationTemplateService) validateTemplateVariables(
	template *model.NotificationTemplate,
	req *request.TestNotificationTemplateRequest,
) error {
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
			return fmt.Errorf("template requires variables: %v", requiredVars)
		}
	} else {
		// Validate all required variables are provided
		validationErr := s.validateVariableData(req.VariableData, requiredVars)
		if validationErr != nil {
			return fmt.Errorf("variable validation failed: %v", validationErr)
		}
	}

	return nil
}

// sendTestNotification sends test notification through the specified channel
func (s *notificationTemplateService) sendTestNotification(
	ctx context.Context,
	template *model.NotificationTemplate,
	channelCode string,
	req *request.TestNotificationTemplateRequest,
) error {
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

	// Get channel details first
	channel, err := s.channelService.GetChannelByCode(ctx, channelCode)
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Send test notification based on channel type
	switch channel.ChannelType {
	case constants.NotificationChannelEmail:
		testReq := &request.TestEmailChannelRequest{
			Code:      channelCode,
			Recipient: req.Recipient,
			Subject:   "",
			Content:   renderedContent,
		}
		if renderedSubject != nil {
			testReq.Subject = *renderedSubject
		}
		err := s.channelService.TestEmailChannel(ctx, testReq)
		if err != nil {
			return fmt.Errorf("failed to send test email: %w", err)
		}
	default:
		return fmt.Errorf("channel type '%s' is not supported for template testing", channel.ChannelType)
	}

	s.logger.InfoContext(ctx, "Test notification sent successfully",
		logger.String("channel_code", channelCode),
		logger.String("template_type", template.TemplateType),
		logger.String("recipient", req.Recipient))

	return nil
}

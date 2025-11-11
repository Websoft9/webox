package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"api-service/internal/dto/common"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	interfaceRepo "api-service/internal/interface/repository"
	interfaceService "api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/crypto"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
)

const (
	// credentialNamePattern is the regex pattern for credential name validation
	// #nosec G101 -- This is a regex pattern, not a hardcoded credential
	credentialNamePattern = `^[a-zA-Z0-9_]{3,50}$`

	// credentialReferencePattern is the regex pattern for credential reference
	// Format: {{ credentials.credential_name.parameter_name }}
	// #nosec G101 -- This is a regex pattern for parsing references, not a hardcoded credential
	credentialReferencePattern = `\{\{\s*credentials\.([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\s*\}\}`

	// MaskedValue is the value shown for encrypted fields
	MaskedValue = "******"

	// referenceMatchGroups is the expected number of regex match groups
	referenceMatchGroups = 3
)

// credentialService implements CredentialService
type credentialService struct {
	credentialRepo        interfaceRepo.CredentialRepository
	categoryRepo          interfaceRepo.CredentialCategoryRepository
	templateRepo          interfaceRepo.CredentialTemplateRepository
	crypto                *crypto.AESCrypto
	logger                logger.Logger
	credentialNamePattern *regexp.Regexp
	credentialRefPattern  *regexp.Regexp
}

// NewCredentialService creates a new credential service
func NewCredentialService(
	credentialRepo interfaceRepo.CredentialRepository,
	categoryRepo interfaceRepo.CredentialCategoryRepository,
	templateRepo interfaceRepo.CredentialTemplateRepository,
	crypto *crypto.AESCrypto,
	logger logger.Logger,
) (interfaceService.CredentialService, error) {
	// Use MustCompile for const patterns - will panic if pattern is invalid
	namePattern := regexp.MustCompile(credentialNamePattern)
	refPattern := regexp.MustCompile(credentialReferencePattern)

	return &credentialService{
		credentialRepo:        credentialRepo,
		categoryRepo:          categoryRepo,
		templateRepo:          templateRepo,
		crypto:                crypto,
		logger:                logger,
		credentialNamePattern: namePattern,
		credentialRefPattern:  refPattern,
	}, nil
}

// validateCredentialName validates credential name format
func (s *credentialService) validateCredentialName(name string) error {
	if !s.credentialNamePattern.MatchString(name) {
		return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "credential.invalid_name_format")
	}
	return nil
}

// validateTemplate validates template and returns form schema
func (s *credentialService) validateTemplate(ctx context.Context, templateID uint) (*model.CredentialTemplate, []model.FormField, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Template not found",
			logger.Uint("template_id", templateID),
			logger.ErrorField(err))
		return nil, nil, errors.NewAppErrorWithI18n(errors.CodeRecordNotFound, "credential.template_not_found")
	}

	// Parse form schema
	var formSchema []model.FormField
	formSchemaBytes, err := json.Marshal(template.FormSchema)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal form schema",
			logger.Uint("template_id", templateID),
			logger.ErrorField(err))
		return nil, nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}
	if err := json.Unmarshal(formSchemaBytes, &formSchema); err != nil {
		s.logger.ErrorContext(ctx, "Failed to parse form schema",
			logger.Uint("template_id", templateID),
			logger.ErrorField(err))
		return nil, nil, errors.NewAppErrorWrapError(err, errors.CodeInternalError)
	}

	return template, formSchema, nil
}

// validateParameters validates parameters against template schema
func (s *credentialService) validateParameters(formSchema []model.FormField, parameters []request.CredentialParameterItem) error {
	// Build schema map for quick lookup
	schemaMap := make(map[string]model.FormField)
	for i := range formSchema {
		schemaMap[formSchema[i].InputName] = formSchema[i]
	}

	// Check all required fields are provided
	if err := s.validateRequiredFields(formSchema, parameters); err != nil {
		return err
	}

	// Validate each parameter
	return s.validateParameterValues(schemaMap, parameters)
}

// validateRequiredFields checks if all required fields are provided
func (s *credentialService) validateRequiredFields(formSchema []model.FormField, parameters []request.CredentialParameterItem) error {
	for i := range formSchema {
		if formSchema[i].InputRequired {
			found := false
			for j := range parameters {
				if parameters[j].InputName == formSchema[i].InputName {
					found = true
					break
				}
			}
			if !found {
				return errors.NewAppErrorWithI18n(errors.CodeRequiredParameterMissing, "credential.parameter_missing")
			}
		}
	}
	return nil
}

// validateParameterValues validates each parameter value against schema rules
func (s *credentialService) validateParameterValues(schemaMap map[string]model.FormField, parameters []request.CredentialParameterItem) error {
	for i := range parameters {
		field, exists := schemaMap[parameters[i].InputName]
		if !exists {
			return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "credential.parameter_not_in_schema")
		}

		// Validate length
		valueLen := len(parameters[i].InputValue)
		if field.InputMinLength > 0 && valueLen < field.InputMinLength {
			return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterLength, "credential.parameter_too_short")
		}
		if field.InputMaxLength > 0 && valueLen > field.InputMaxLength {
			return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterLength, "credential.parameter_too_long")
		}

		// Validate pattern
		if field.InputPattern != "" {
			matched, matchErr := regexp.MatchString(field.InputPattern, parameters[i].InputValue)
			if matchErr != nil || !matched {
				return errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "credential.parameter_pattern_mismatch")
			}
		}
	}
	return nil
}

// encryptParameters encrypts parameters that need encryption
func (s *credentialService) encryptParameters(
	ctx context.Context,
	formSchema []model.FormField,
	parameters []request.CredentialParameterItem,
) ([]model.CredentialParameter, error) {
	// Build schema map
	schemaMap := make(map[string]model.FormField)
	for i := range formSchema {
		schemaMap[formSchema[i].InputName] = formSchema[i]
	}

	result := make([]model.CredentialParameter, 0, len(parameters))
	for _, param := range parameters {
		field := schemaMap[param.InputName]
		value := param.InputValue

		// Encrypt if needed
		if field.IsEncrypted {
			encrypted, err := s.crypto.Encrypt(value)
			if err != nil {
				s.logger.ErrorContext(ctx, "Failed to encrypt parameter",
					logger.String("parameter", param.InputName),
					logger.ErrorField(err))
				return nil, errors.NewAppErrorWrapError(err, errors.CodeEncryptFailed)
			}
			value = encrypted
		}

		result = append(result, model.CredentialParameter{
			InputName:   param.InputName,
			InputValue:  value,
			IsEncrypted: field.IsEncrypted,
		})
	}

	return result, nil
}

// CreateCredential creates a new credential
func (s *credentialService) CreateCredential(ctx context.Context, req *request.CreateCredentialRequest, ownerID uint) (*response.CredentialResponse, error) {
	// Validate credential name format
	if err := s.validateCredentialName(req.Name); err != nil {
		return nil, err
	}

	// Check name uniqueness
	exists, err := s.credentialRepo.ExistsByName(ctx, req.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewAppErrorWithI18n(errors.CodeResourceAlreadyExists, "credential.name_already_exists")
	}

	// Validate template and get form schema
	template, formSchema, templateErr := s.validateTemplate(ctx, req.TemplateID)
	if templateErr != nil {
		return nil, templateErr
	}

	// Validate parameters against schema
	if validateErr := s.validateParameters(formSchema, req.Parameters); validateErr != nil {
		return nil, validateErr
	}

	// Encrypt parameters
	encryptedParams, err := s.encryptParameters(ctx, formSchema, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Convert parameters array to JSON map for storage
	// Store as map with numeric keys to maintain array structure
	paramsMap := make(model.JSON)
	for i, param := range encryptedParams {
		paramsMap[fmt.Sprintf("%d", i)] = map[string]interface{}{
			"input_name":   param.InputName,
			"input_value":  param.InputValue,
			"is_encrypted": param.IsEncrypted,
		}
	}

	// Create credential model
	credential := &model.Credential{
		Name:        req.Name,
		Description: req.Description,
		TemplateID:  req.TemplateID,
		Parameters:  paramsMap,
		OwnerID:     ownerID,
	}

	// Save to database
	if err := s.credentialRepo.Create(ctx, credential); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create credential",
			logger.String("name", req.Name),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Credential created successfully",
		logger.Uint("credential_id", credential.ID),
		logger.String("name", credential.Name))

	// Build response
	return s.buildCredentialResponse(ctx, credential, template), nil
}

// buildCredentialResponse builds credential response
func (s *credentialService) buildCredentialResponse(_ context.Context, credential *model.Credential, template *model.CredentialTemplate) *response.CredentialResponse {
	ownerName := ""
	if credential.Owner != nil {
		ownerName = credential.Owner.Username
	}

	templateName := ""
	categoryID := uint(0)
	categoryName := ""
	if template != nil {
		templateName = template.Name
		categoryID = template.CategoryID
		if template.Category != nil {
			categoryName = template.Category.Name
		}
	}

	return &response.CredentialResponse{
		ID:           credential.ID,
		Name:         credential.Name,
		Description:  credential.Description,
		TemplateID:   credential.TemplateID,
		TemplateName: templateName,
		CategoryID:   categoryID,
		CategoryName: categoryName,
		OwnerID:      credential.OwnerID,
		OwnerName:    ownerName,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}
}

// UpdateCredential updates an existing credential
func (s *credentialService) UpdateCredential(ctx context.Context, id uint, req *request.UpdateCredentialRequest, userID uint) (*response.CredentialResponse, error) {
	// Get credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check permission: only owner can update
	if credential.OwnerID != userID {
		return nil, errors.NewAppError(errors.CodeAccessDenied)
	}

	// Update description if provided
	if req.Description != nil {
		credential.Description = req.Description
	}

	// Update parameters if provided
	if len(req.Parameters) > 0 {
		// Get template and form schema
		_, formSchema, templateErr := s.validateTemplate(ctx, credential.TemplateID)
		if templateErr != nil {
			return nil, templateErr
		}

		// Validate parameters
		if validateErr := s.validateParameters(formSchema, req.Parameters); validateErr != nil {
			return nil, validateErr
		}

		// Encrypt parameters
		encryptedParams, err := s.encryptParameters(ctx, formSchema, req.Parameters)
		if err != nil {
			return nil, err
		}

		// Convert parameters array to JSON map for storage
		paramsMap := make(model.JSON)
		for i, param := range encryptedParams {
			paramsMap[fmt.Sprintf("%d", i)] = map[string]interface{}{
				"input_name":   param.InputName,
				"input_value":  param.InputValue,
				"is_encrypted": param.IsEncrypted,
			}
		}

		credential.Parameters = paramsMap
	}

	// Save to database
	if err := s.credentialRepo.Update(ctx, credential); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update credential",
			logger.Uint("credential_id", id),
			logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Credential updated successfully",
		logger.Uint("credential_id", id),
		logger.String("name", credential.Name))

	return s.buildCredentialResponse(ctx, credential, credential.Template), nil
}

// DeleteCredential deletes a credential
func (s *credentialService) DeleteCredential(ctx context.Context, id, userID uint) error {
	// Get credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check permission: only owner can delete
	if credential.OwnerID != userID {
		return errors.NewAppError(errors.CodeAccessDenied)
	}

	// TODO: Check if credential is being referenced by workflows or applications
	// For now, we allow deletion

	// Delete credential
	if err := s.credentialRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete credential",
			logger.Uint("credential_id", id),
			logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Credential deleted successfully",
		logger.Uint("credential_id", id),
		logger.String("name", credential.Name))

	return nil
}

// GetCredential retrieves a credential by ID
func (s *credentialService) GetCredential(ctx context.Context, id uint) (*response.CredentialDetailResponse, error) {
	// Get credential
	credential, err := s.credentialRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Parse parameters from JSON map
	var params []model.CredentialParameter
	for key := range credential.Parameters {
		if paramMap, ok := credential.Parameters[key].(map[string]interface{}); ok {
			param := model.CredentialParameter{
				InputName:   paramMap["input_name"].(string),
				InputValue:  paramMap["input_value"].(string),
				IsEncrypted: paramMap["is_encrypted"].(bool),
			}
			params = append(params, param)
		}
	}

	// Mask encrypted fields
	paramResponses := make([]response.CredentialParameterResponse, len(params))
	for i, param := range params {
		value := param.InputValue
		if param.IsEncrypted {
			value = MaskedValue
		}
		paramResponses[i] = response.CredentialParameterResponse{
			InputName:   param.InputName,
			InputValue:  value,
			IsEncrypted: param.IsEncrypted,
		}
	}

	// Build response
	ownerName := ""
	if credential.Owner != nil {
		ownerName = credential.Owner.Username
	}

	templateName := ""
	categoryID := uint(0)
	categoryName := ""
	if credential.Template != nil {
		templateName = credential.Template.Name
		categoryID = credential.Template.CategoryID
		if credential.Template.Category != nil {
			categoryName = credential.Template.Category.Name
		}
	}

	return &response.CredentialDetailResponse{
		ID:           credential.ID,
		Name:         credential.Name,
		Description:  credential.Description,
		TemplateID:   credential.TemplateID,
		TemplateName: templateName,
		CategoryID:   categoryID,
		CategoryName: categoryName,
		Parameters:   paramResponses,
		OwnerID:      credential.OwnerID,
		OwnerName:    ownerName,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}, nil
}

// ListCredentials retrieves a paginated list of credentials
func (s *credentialService) ListCredentials(ctx context.Context, req *request.ListCredentialsRequest) (*common.PaginationResponse, error) {
	credentials, total, err := s.credentialRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]any, len(credentials))
	for i, credential := range credentials {
		items[i] = s.buildCredentialResponse(ctx, credential, credential.Template)
	}

	return common.NewPaginationResponse(req.GetPage(), req.GetPageSize(), total, items), nil
}

// ListCredentialTemplates retrieves credential templates
func (s *credentialService) ListCredentialTemplates(ctx context.Context, req *request.ListCredentialTemplatesRequest) ([]response.CredentialTemplateResponse, error) {
	templates, err := s.templateRepo.List(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.CredentialTemplateResponse, len(templates))
	for i, template := range templates {
		// Parse form schema
		var formSchema []model.FormField
		formSchemaBytes, err := json.Marshal(template.FormSchema)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to marshal form schema",
				logger.Uint("template_id", template.ID),
				logger.ErrorField(err))
			continue
		}
		if err := json.Unmarshal(formSchemaBytes, &formSchema); err != nil {
			s.logger.ErrorContext(ctx, "Failed to parse form schema",
				logger.Uint("template_id", template.ID),
				logger.ErrorField(err))
			continue
		}

		// Convert to response format
		formFieldResponses := make([]response.FormFieldResponse, len(formSchema))
		for j := range formSchema {
			formFieldResponses[j] = response.FormFieldResponse{
				InputLabel:        formSchema[j].InputLabel,
				InputName:         formSchema[j].InputName,
				InputType:         formSchema[j].InputType,
				InputDefaultValue: formSchema[j].InputDefaultValue,
				InputMinLength:    formSchema[j].InputMinLength,
				InputMaxLength:    formSchema[j].InputMaxLength,
				InputPlaceholder:  formSchema[j].InputPlaceholder,
				InputRequired:     formSchema[j].InputRequired,
				InputPattern:      formSchema[j].InputPattern,
				IsEncrypted:       formSchema[j].IsEncrypted,
			}
		}

		categoryName := ""
		if template.Category != nil {
			categoryName = template.Category.Name
		}

		responses[i] = response.CredentialTemplateResponse{
			ID:           template.ID,
			Name:         template.Name,
			Description:  template.Description,
			CategoryID:   template.CategoryID,
			CategoryName: categoryName,
			FormSchema:   formFieldResponses,
			CreatedAt:    template.CreatedAt,
			UpdatedAt:    template.UpdatedAt,
		}
	}

	return responses, nil
}

// ListCredentialCategories retrieves credential categories
func (s *credentialService) ListCredentialCategories(ctx context.Context) ([]response.CredentialCategoryResponse, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]response.CredentialCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = response.CredentialCategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		}
	}

	return responses, nil
}

// ResolveCredentialReference resolves credential reference expression
// Input: {{ credentials.my_db.username }}
// Output: decrypted actual value
func (s *credentialService) ResolveCredentialReference(ctx context.Context, expression string) (string, error) {
	// Validate expression format
	matches := s.credentialRefPattern.FindStringSubmatch(expression)
	if len(matches) != referenceMatchGroups {
		return "", errors.NewAppErrorWithI18n(errors.CodeInvalidParameterFormat, "credential.invalid_reference_format")
	}

	credentialName := matches[1]
	parameterName := matches[2]

	// Get credential by name
	credential, err := s.credentialRepo.GetByName(ctx, credentialName)
	if err != nil {
		s.logger.ErrorContext(ctx, "Credential not found for reference",
			logger.String("credential_name", credentialName),
			logger.ErrorField(err))
		return "", errors.NewAppErrorWithI18n(errors.CodeRecordNotFound, "credential.not_found")
	}

	// Parse parameters from JSON map
	var params []model.CredentialParameter
	for key := range credential.Parameters {
		if paramMap, ok := credential.Parameters[key].(map[string]interface{}); ok {
			param := model.CredentialParameter{
				InputName:   paramMap["input_name"].(string),
				InputValue:  paramMap["input_value"].(string),
				IsEncrypted: paramMap["is_encrypted"].(bool),
			}
			params = append(params, param)
		}
	}

	// Find parameter
	var targetParam *model.CredentialParameter
	for i := range params {
		if params[i].InputName == parameterName {
			targetParam = &params[i]
			break
		}
	}

	if targetParam == nil {
		return "", errors.NewAppErrorWithI18n(errors.CodeRecordNotFound, "credential.parameter_not_found")
	}

	// Decrypt if encrypted
	value := targetParam.InputValue
	if targetParam.IsEncrypted {
		decrypted, err := s.crypto.Decrypt(value)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to decrypt parameter",
				logger.String("credential", credentialName),
				logger.String("parameter", parameterName),
				logger.ErrorField(err))
			return "", errors.NewAppErrorWrapError(err, errors.CodeDecryptFailed)
		}
		value = decrypted
	}

	s.logger.InfoContext(ctx, "Credential reference resolved",
		logger.String("credential", credentialName),
		logger.String("parameter", parameterName))

	return value, nil
}

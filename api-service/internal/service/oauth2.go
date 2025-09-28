package service

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth2Service handles OAuth2 authentication flows
type OAuth2Service struct {
	authConfigManager *config.AuthConfigManager
	logger            logger.Logger
	httpClient        *http.Client
}

// OAuth2UserInfo represents user information from OAuth2 provider
type OAuth2UserInfo struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	Username string                 `json:"username"`
	Avatar   string                 `json:"avatar"`
	Raw      map[string]interface{} `json:"raw"`
}

// OAuth2TokenResponse represents OAuth2 token response
type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// NewOAuth2Service creates a new OAuth2 service instance
func NewOAuth2Service(authConfigManager *config.AuthConfigManager, logger logger.Logger) *OAuth2Service {
	return &OAuth2Service{
		authConfigManager: authConfigManager,
		logger:            logger,
		httpClient: &http.Client{
			Timeout: constants.HTTPClientTimeout * time.Second,
		},
	}
}

// GetAuthorizationURL generates OAuth2 authorization URL for a provider
func (s *OAuth2Service) GetAuthorizationURL(ctx context.Context, providerName, state string) (string, error) {
	s.logger.InfoContext(ctx, "Generating OAuth2 authorization URL",
		logger.String("provider", providerName),
		logger.String("state", state))

	// Get provider configuration
	provider, exists := s.authConfigManager.GetOAuth2Provider(providerName)
	if !exists {
		return "", fmt.Errorf("OAuth2 provider '%s' not found", providerName)
	}

	if !provider.Enabled {
		return "", fmt.Errorf("OAuth2 provider '%s' is disabled", providerName)
	}

	// Build authorization URL
	params := url.Values{}
	params.Add("client_id", provider.ClientID)
	params.Add("redirect_uri", provider.RedirectURI)
	params.Add("response_type", "code")
	params.Add("state", state)

	if len(provider.Scopes) > 0 {
		params.Add("scope", strings.Join(provider.Scopes, " "))
	}

	// Add provider-specific parameters
	switch providerName {
	case constants.ProviderGoogle:
		params.Add("access_type", "offline")
		params.Add("prompt", "consent")
	case constants.ProviderGitHub:
		// GitHub specific parameters can be added here
	case constants.ProviderAzureAD:
		params.Add("response_mode", "query")
	}

	authURL := provider.AuthorizeURL + "?" + params.Encode()

	s.logger.InfoContext(ctx, "Generated OAuth2 authorization URL successfully",
		logger.String("provider", providerName))

	return authURL, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *OAuth2Service) ExchangeCodeForToken(ctx context.Context, providerName, code string) (*OAuth2TokenResponse, error) {
	s.logger.InfoContext(ctx, "Exchanging OAuth2 code for token",
		logger.String("provider", providerName))

	// Get provider configuration
	provider, exists := s.authConfigManager.GetOAuth2Provider(providerName)
	if !exists {
		return nil, fmt.Errorf("OAuth2 provider '%s' not found", providerName)
	}

	// Prepare token exchange request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", provider.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", provider.RedirectURI)

	// Make token exchange request
	req, err := http.NewRequestWithContext(ctx, "POST", provider.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create token exchange request", logger.ErrorField(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Token exchange request failed", logger.ErrorField(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to read token response", logger.ErrorField(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.ErrorContext(ctx, "Token exchange failed",
			logger.Int("status_code", resp.StatusCode),
			logger.String("response", string(body)))
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	// Parse token response
	var tokenResp OAuth2TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		s.logger.ErrorContext(ctx, "Failed to parse token response", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "OAuth2 token exchange successful",
		logger.String("provider", providerName))

	return &tokenResp, nil
}

// GetUserInfo retrieves user information from OAuth2 provider
func (s *OAuth2Service) GetUserInfo(ctx context.Context, providerName, accessToken string) (*OAuth2UserInfo, error) {
	s.logger.InfoContext(ctx, "Retrieving OAuth2 user info",
		logger.String("provider", providerName))

	// Get provider configuration
	provider, exists := s.authConfigManager.GetOAuth2Provider(providerName)
	if !exists {
		return nil, fmt.Errorf("OAuth2 provider '%s' not found", providerName)
	}

	// Make user info request
	req, err := http.NewRequestWithContext(ctx, "GET", provider.UserInfoURL, http.NoBody)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user info request", logger.ErrorField(err))
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "User info request failed", logger.ErrorField(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to read user info response", logger.ErrorField(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.ErrorContext(ctx, "User info request failed",
			logger.Int("status_code", resp.StatusCode),
			logger.String("response", string(body)))
		return nil, fmt.Errorf("user info request failed: %s", string(body))
	}

	// Parse user info response
	var rawUserInfo map[string]interface{}
	if unmarshalErr := json.Unmarshal(body, &rawUserInfo); unmarshalErr != nil {
		s.logger.ErrorContext(ctx, "Failed to parse user info response", logger.ErrorField(unmarshalErr))
		return nil, unmarshalErr
	}

	// Map user info according to provider configuration
	userInfo, err := s.mapUserInfo(providerName, provider, rawUserInfo)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to map user info", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "OAuth2 user info retrieved successfully",
		logger.String("provider", providerName),
		logger.String("user_id", userInfo.ID))

	return userInfo, nil
}

// mapUserInfo maps raw user info to standardized format according to provider mapping
func (s *OAuth2Service) mapUserInfo(providerName string, provider *config.OAuth2ProviderConfig, rawUserInfo map[string]interface{}) (*OAuth2UserInfo, error) {
	userInfo := &OAuth2UserInfo{
		Raw: rawUserInfo,
	}

	// Apply user mapping configuration
	s.mapUserFields(provider.UserMapping, rawUserInfo, userInfo)

	// Set ID based on provider-specific logic
	s.mapUserID(providerName, rawUserInfo, userInfo)

	// Validate required fields
	if userInfo.ID == "" {
		return nil, errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	// Use ID as username if username is not mapped
	if userInfo.Username == "" {
		userInfo.Username = userInfo.ID
	}

	return userInfo, nil
}

// mapUserFields applies user field mappings from provider configuration
func (s *OAuth2Service) mapUserFields(userMapping map[string]string, rawUserInfo map[string]interface{}, userInfo *OAuth2UserInfo) {
	for field, sourceField := range userMapping {
		if value, exists := rawUserInfo[sourceField]; exists {
			s.mapSingleUserField(field, value, userInfo)
		}
	}
}

// mapSingleUserField maps a single user field to the standardized format
func (s *OAuth2Service) mapSingleUserField(field string, value interface{}, userInfo *OAuth2UserInfo) {
	strVal, ok := value.(string)
	if !ok {
		return
	}

	switch field {
	case constants.UserMappingUsername:
		userInfo.Username = strVal
	case constants.UserMappingEmail:
		userInfo.Email = strVal
	case "name":
		userInfo.Name = strVal
	case "avatar":
		userInfo.Avatar = strVal
	}
}

// mapUserID extracts and maps user ID from raw user info based on provider
func (s *OAuth2Service) mapUserID(providerName string, rawUserInfo map[string]interface{}, userInfo *OAuth2UserInfo) {
	switch providerName {
	case constants.ProviderGoogle:
		s.mapProviderSpecificID(rawUserInfo, userInfo, "id", "string")
	case constants.ProviderGitHub:
		s.mapProviderSpecificID(rawUserInfo, userInfo, "id", "number")
	case constants.ProviderAzureAD:
		s.mapProviderSpecificID(rawUserInfo, userInfo, "id", "string")
	default:
		s.mapGenericUserID(rawUserInfo, userInfo)
	}
}

// mapProviderSpecificID maps user ID for known providers
func (s *OAuth2Service) mapProviderSpecificID(rawUserInfo map[string]interface{}, userInfo *OAuth2UserInfo, field, expectedType string) {
	if id, exists := rawUserInfo[field]; exists {
		switch expectedType {
		case "string":
			if strVal, ok := id.(string); ok {
				userInfo.ID = strVal
			}
		case "number":
			if numVal, ok := id.(float64); ok {
				userInfo.ID = fmt.Sprintf("%.0f", numVal)
			}
		}
	}
}

// mapGenericUserID attempts to extract user ID from common field names
func (s *OAuth2Service) mapGenericUserID(rawUserInfo map[string]interface{}, userInfo *OAuth2UserInfo) {
	// Try to get ID from various possible fields
	for _, idField := range []string{"id", "sub", "user_id", "userId"} {
		if id, exists := rawUserInfo[idField]; exists {
			switch idVal := id.(type) {
			case string:
				userInfo.ID = idVal
			case float64:
				userInfo.ID = fmt.Sprintf("%.0f", idVal)
			}
			if userInfo.ID != "" {
				break
			}
		}
	}
}

// ValidateState validates OAuth2 state parameter to prevent CSRF attacks
func (s *OAuth2Service) ValidateState(ctx context.Context, receivedState, expectedState string) error {
	if receivedState == "" {
		return errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	if receivedState != expectedState {
		s.logger.WarnContext(ctx, "OAuth2 state validation failed",
			logger.String("received", receivedState),
			logger.String("expected", expectedState))
		return errors.NewAppError(errors.CodeValidationFailed)
	}

	return nil
}

// RefreshToken refreshes OAuth2 access token using refresh token
func (s *OAuth2Service) RefreshToken(ctx context.Context, providerName, refreshToken string) (*OAuth2TokenResponse, error) {
	s.logger.InfoContext(ctx, "Refreshing OAuth2 token",
		logger.String("provider", providerName))

	// Get provider configuration
	provider, exists := s.authConfigManager.GetOAuth2Provider(providerName)
	if !exists {
		return nil, fmt.Errorf("OAuth2 provider '%s' not found", providerName)
	}

	// Prepare token refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", provider.ClientID)
	data.Set("client_secret", provider.ClientSecret)
	data.Set("refresh_token", refreshToken)

	// Make token refresh request
	req, err := http.NewRequestWithContext(ctx, "POST", provider.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create token refresh request", logger.ErrorField(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "Token refresh request failed", logger.ErrorField(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to read token refresh response", logger.ErrorField(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.ErrorContext(ctx, "Token refresh failed",
			logger.Int("status_code", resp.StatusCode),
			logger.String("response", string(body)))
		return nil, fmt.Errorf("token refresh failed: %s", string(body))
	}

	// Parse token response
	var tokenResp OAuth2TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		s.logger.ErrorContext(ctx, "Failed to parse token refresh response", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "OAuth2 token refresh successful",
		logger.String("provider", providerName))

	return &tokenResp, nil
}

// IsProviderEnabled checks if an OAuth2 provider is enabled
func (s *OAuth2Service) IsProviderEnabled(providerName string) bool {
	provider, exists := s.authConfigManager.GetOAuth2Provider(providerName)
	return exists && provider.Enabled
}

// GetEnabledProviders returns list of all enabled OAuth2 providers
func (s *OAuth2Service) GetEnabledProviders() []config.OAuth2ProviderConfig {
	return s.authConfigManager.GetEnabledOAuth2Providers()
}

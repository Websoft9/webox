package service

import (
	"api-service/internal/config"
	"api-service/internal/constants"
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/email"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"api-service/pkg/redis"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	authConfig *config.AuthConfig // auth config
)

const (
	tokenRandomBytesSize    = 32 // Size in bytes for verification token generation
	secondsToMinutesConvert = 60 // Conversion factor from seconds to minutes
	hashKeyValuePairs       = 2  // Number of elements per key-value pair for Redis hash
)

// VerificationToken represents email verification or password reset token
type VerificationToken struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	Type      string    `json:"type"` // "email_verification" or "password_reset"
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Used      bool      `json:"used"`
}

type userAuthService struct {
	userRepo          repository.UserRepository
	apiTokenRepo      repository.APITokenRepository
	userProfileRepo   repository.UserProfileRepository
	systemConfigRepo  repository.SystemConfigRepository
	emailService      email.EmailService
	oauth2Service     *OAuth2Service
	logger            logger.Logger
	appConfig         *config.Config
	authConfigManager *config.AuthConfigManager
	authPolicy        *auth.AuthPolicy
}

func NewUserAuthService(
	userRepo repository.UserRepository,
	apiTokenRepo repository.APITokenRepository,
	userProfileRepo repository.UserProfileRepository,
	systemConfigRepo repository.SystemConfigRepository,
	oauth2Service *OAuth2Service,
	zapLogger logger.Logger,
	appConfig *config.Config,
	authConfigManager *config.AuthConfigManager,
) service.UserAuthService {
	// Create email service internally
	var emailService email.EmailService
	if appConfig.Email.SMTP.Host != "" && appConfig.Email.SMTP.Username != "" {
		emailService = email.NewEmailService(appConfig, zapLogger)
	}
	authConfig = authConfigManager.GetConfig()

	service := &userAuthService{
		userRepo:          userRepo,
		apiTokenRepo:      apiTokenRepo,
		userProfileRepo:   userProfileRepo,
		systemConfigRepo:  systemConfigRepo,
		emailService:      emailService,
		oauth2Service:     oauth2Service,
		logger:            zapLogger,
		appConfig:         appConfig,
		authConfigManager: authConfigManager,
		authPolicy:        auth.NewAuthPolicy(authConfigManager, zapLogger),
	}

	return service
}

// Register handles user registration with email verification
func (s *userAuthService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Starting user registration", logger.String("email", req.Username))

	// 1. Validate email format
	if err := s.authPolicy.ValidateEmail(req.Username); err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", logger.String("email", req.Username))
		return nil, errors.ErrInvalidEmailFormat
	}

	// 2. Validate password strength using configured policy
	if err := s.authPolicy.ValidatePasswordWithPolicy(ctx, req.Password); err != nil {
		s.logger.WarnContext(ctx, "Password policy validation failed", logger.ErrorField(err))
		return nil, err
	}

	// 3. Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check email existence", logger.ErrorField(err))
		return nil, err
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// 4. Generate username from email prefix
	username := s.generateUsernameFromEmail(req.Username)

	// 5. Hash password
	hashedPassword := auth.HashToken(req.Password)

	// 6. Create user
	user := &model.User{
		Username:     username,
		Email:        req.Username,
		PasswordHash: hashedPassword,
		Status:       UserStatusActive,
		Language:     "zh-CN",
		Timezone:     "Asia/Shanghai",
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user", logger.ErrorField(err))
		return nil, err
	}

	// 7. Set email verification lock in Redis to prevent duplicate verification attempts
	if authConfig.UserAuth.EmailAuth.Enabled {
		if lockErr := s.setEmailVerificationLock(ctx, req.Username, req.Username, authConfig.UserAuth.EmailAuth.ExpiresIn); lockErr != nil {
			s.logger.WarnContext(ctx, "Failed to set email verification lock",
				logger.ErrorField(lockErr), logger.String("email", req.Username))
			// Don't return error, user registration is already successful
		}

		// 8. Generate email verification token
		verificationToken, err := s.generateVerificationToken(ctx, req.Username, "email_verification")
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to generate verification token", logger.ErrorField(err))
			// Don't return error, user is already created successfully
		} else {
			// 9. Send verification email
			expiresMinutes := authConfig.UserAuth.EmailAuth.ExpiresIn / secondsToMinutesConvert
			if err := s.emailService.SendVerificationEmail(ctx, expiresMinutes, req.Username, verificationToken.Token, s.appConfig.App.BaseURL, user.Language); err != nil {
				s.logger.ErrorContext(ctx, "Failed to send verification email", logger.ErrorField(err))
				// Don't return error, user is already created successfully
			}
		}
	}

	s.logger.InfoContext(ctx, "User registration successful",
		logger.Uint("user_id", user.ID),
		logger.String("email", req.Username))

	return response.BuildUserResponse(user), nil
}

// Login handles user authentication with username/email and password
func (s *userAuthService) Login(ctx context.Context, req *request.UserLoginRequest, clientIP string) (*response.UserLoginResponse, error) {
	s.logger.InfoContext(ctx, "Starting user login", logger.String("username_or_email", req.Username))

	// 1. Login security checks (IP whitelist, time restrictions)
	if err := s.authPolicy.CheckLoginSecurity(ctx, req.Username, clientIP); err != nil {
		s.logger.WarnContext(ctx, "Login security check failed", logger.ErrorField(err))
		return nil, err
	}

	// 2. Find user by username or email
	user, err := s.userRepo.GetByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "User not found", logger.String("username_or_email", req.Username))
			return nil, errors.ErrInvalidCredentials
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return nil, errors.ErrRecordQueryFailed
	}

	// 3. Check if user has pending email verification
	hasLock, err := s.hasEmailVerificationLock(ctx, user.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check email verification lock", logger.ErrorField(err), logger.String("email", req.Username))
		// Continue with login process even if Redis check fails
	}
	if hasLock {
		s.logger.WarnContext(ctx, "User has pending email verification", logger.String("email", req.Username))
		return nil, errors.ErrEmailNotVerified
	}

	// 4. Verify password
	passwordValid := user.PasswordHash == auth.HashToken(req.Password)

	// Use auth policy to handle password validation and attempts tracking
	if pwdValidErr := s.authPolicy.ValidatePassword(ctx, req.Username, passwordValid); pwdValidErr != nil {
		return nil, pwdValidErr
	}

	// 5. Check user status
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "User account is inactive", logger.String("username_or_email", req.Username))
		return nil, errors.ErrAccountDisabled
	}

	// Get user roles
	userRoleIDs := make([]uint, len(user.Roles))
	for i := range user.Roles {
		userRoleIDs[i] = user.Roles[i].ID
	}

	// 6. Generate JWT token using configured expiration time
	token, expiresAt, err := auth.GetGlobalJWT().GenerateTokenWithUserInfo(user.ID, user.Username, userRoleIDs)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate JWT token", logger.ErrorField(err))
		return nil, errors.ErrRecordCreateFailed
	}

	// Store JWT token in database and Redis cache (according to security design)
	if err := s.storeJWTToken(ctx, token, user.ID, expiresAt); err != nil {
		s.logger.WarnContext(ctx, "Failed to store JWT token in database/Redis", logger.ErrorField(err))
		// Don't fail the login process, as token is still valid for authentication
	}

	// 7. Update last login time and IP
	now := time.Now()
	updateUser := model.User{
		ID:          user.ID,
		LastLoginAt: &now,
		LastLoginIP: clientIP,
	}
	if err := s.userRepo.Update(ctx, &updateUser); err != nil {
		s.logger.WarnContext(ctx, "Failed to update login info", logger.ErrorField(err))
		// Don't return error as login is already successful
	}

	// 8. Store user preferences in Redis
	if err := s.storeUserPreferencesInRedis(ctx, user.ID); err != nil {
		s.logger.WarnContext(ctx, "Failed to store user preferences in Redis", logger.ErrorField(err))
		// Don't return error as login is already successful
	}

	// 9. Build response
	s.logger.InfoContext(ctx, "User login successful",
		logger.Uint("user_id", user.ID),
		logger.String("username", user.Username),
		logger.String("email", user.Email))

	return &response.UserLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *response.BuildUserResponse(user),
	}, nil
}

// ForgotPassword handles password reset request
func (s *userAuthService) ForgotPassword(ctx context.Context, req *request.ForgotPasswordRequest) error {
	s.logger.InfoContext(ctx, "Processing forgot password request", logger.String("email", req.Username))

	// 1. Validate email format
	if err := s.authPolicy.ValidateEmail(req.Username); err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", logger.String("email", req.Username))
		return errors.ErrInvalidEmailFormat
	}

	// 2. Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// For security reasons, return success even if user doesn't exist
			s.logger.WarnContext(ctx, "User not found for password reset", logger.String("email", req.Username))
			return nil
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return errors.ErrRecordQueryFailed
	}

	// 3. Check user status
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "Inactive user requested password reset", logger.String("email", req.Username))
		// For security reasons, don't reveal account status
		return nil
	}

	// 4. Generate password reset token
	resetToken, err := s.generateVerificationToken(ctx, req.Username, "password_reset")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate reset token", logger.ErrorField(err))
		return errors.ErrRecordCreateFailed
	}

	// 5. Send password reset email
	expiresMinutes := authConfig.UserAuth.EmailAuth.ExpiresIn / secondsToMinutesConvert
	if err := s.emailService.SendPasswordResetEmail(ctx, expiresMinutes, req.Username, resetToken.Token, s.appConfig.App.BaseURL, user.Language); err != nil {
		s.logger.ErrorContext(ctx, "Failed to send password reset email", logger.ErrorField(err))
		return errors.ErrRecordCreateFailed
	}

	s.logger.InfoContext(ctx, "Password reset email sent successfully", logger.String("email", req.Username))
	return nil
}

// ValidateResetToken validates a password reset token without marking it as used
func (s *userAuthService) ValidateResetToken(ctx context.Context, token string) error {
	s.logger.InfoContext(ctx, "Validating password reset token", logger.String("token", token[:8]+"..."))

	// Validate reset token without marking as used
	_, err := s.validateVerificationToken(ctx, token, "password_reset")
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid password reset token", logger.ErrorField(err))
		if strings.Contains(err.Error(), "expired") {
			return errors.ErrTokenExpired
		} else if strings.Contains(err.Error(), "used") {
			return errors.ErrTokenAlreadyUsed
		}
		return errors.ErrInvalidToken
	}

	s.logger.InfoContext(ctx, "Password reset token validation successful")
	return nil
}

// ResetPassword handles password reset with token
func (s *userAuthService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	s.logger.InfoContext(ctx, "Processing password reset", logger.String("token", req.Token[:8]+"..."))

	// 1. Validate password strength using configured policy
	if err := s.authPolicy.ValidatePasswordWithPolicy(ctx, req.NewPassword); err != nil {
		s.logger.WarnContext(ctx, "Password policy validation failed", logger.ErrorField(err))
		return err
	}

	// 2. Validate reset token
	verificationToken, err := s.validateVerificationToken(ctx, req.Token, "password_reset")
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid reset token", logger.ErrorField(err))
		if strings.Contains(err.Error(), "expired") {
			return errors.ErrTokenExpired
		} else if strings.Contains(err.Error(), "used") {
			return errors.ErrTokenAlreadyUsed
		}
		return errors.ErrInvalidToken
	}

	// 3. Find user
	user, err := s.userRepo.GetByEmail(ctx, verificationToken.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "User not found for password reset", logger.String("email", verificationToken.Email))
			return errors.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return errors.ErrRecordQueryFailed
	}

	// 4. Update password
	user.PasswordHash = auth.HashToken(req.NewPassword)

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password", logger.ErrorField(err))
		return errors.ErrRecordUpdateFailed
	}

	// 5. Mark token as used
	if err := s.markTokenAsUsed(ctx, req.Token); err != nil {
		s.logger.WarnContext(ctx, "Failed to mark token as used", logger.ErrorField(err))
		// Don't return error, password reset is already successful
	}

	s.logger.InfoContext(ctx, "Password reset successful",
		logger.Uint("user_id", user.ID),
		logger.String("email", verificationToken.Email))

	return nil
}

// VerifyEmail handles email verification with token
func (s *userAuthService) VerifyEmail(ctx context.Context, req *request.VerifyEmailRequest) error {
	s.logger.InfoContext(ctx, "Processing email verification", logger.String("token", req.Token[:8]+"..."))

	// 1. Validate email verification token
	verificationToken, err := s.validateVerificationToken(ctx, req.Token, "email_verification")
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid verification token", logger.ErrorField(err))
		if strings.Contains(err.Error(), "expired") {
			return errors.ErrTokenExpired
		} else if strings.Contains(err.Error(), "used") {
			return errors.ErrTokenAlreadyUsed
		}
		return errors.ErrInvalidToken
	}

	// 2. Find user
	user, err := s.userRepo.GetByEmail(ctx, verificationToken.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "User not found for email verification", logger.String("email", verificationToken.Email))
			return errors.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return errors.ErrRecordQueryFailed
	}

	// 3. Check if user has pending email verification lock
	hasLock, err := s.hasEmailVerificationLock(ctx, verificationToken.Email)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to check email verification lock",
			logger.ErrorField(err), logger.String("email", verificationToken.Email))
		// Continue with verification process even if Redis check fails
	} else if !hasLock {
		s.logger.InfoContext(ctx, "User does not have pending email verification, may be already verified",
			logger.String("email", verificationToken.Email))
		return errors.ErrValidationFailed
	}

	// 4. Update user status
	user.Status = UserStatusActive // Activate account

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update email verification status", logger.ErrorField(err))
		return errors.ErrRecordUpdateFailed
	}

	// 5. Mark token as used
	if err := s.markTokenAsUsed(ctx, req.Token); err != nil {
		s.logger.WarnContext(ctx, "Failed to mark token as used", logger.ErrorField(err))
		// Don't return error, email verification is already successful
	}

	// 6. Delete email verification lock from Redis
	if err := s.deleteEmailVerificationLock(ctx, verificationToken.Email); err != nil {
		s.logger.WarnContext(ctx, "Failed to delete email verification lock",
			logger.ErrorField(err), logger.String("email", verificationToken.Email))
		// Don't return error, email verification is already successful
	}

	s.logger.InfoContext(ctx, "Email verification successful",
		logger.Uint("user_id", user.ID),
		logger.String("email", verificationToken.Email))

	return nil
}

// ResendVerificationEmail handles resending email verification
func (s *userAuthService) ResendVerificationEmail(ctx context.Context, req *request.ResendVerificationRequest) error {
	s.logger.InfoContext(ctx, "Resending verification email", logger.String("email", req.Username))

	// 1. Validate email format
	if err := s.authPolicy.ValidateEmail(req.Username); err != nil {
		s.logger.WarnContext(ctx, "Invalid email format", logger.String("email", req.Username))
		return errors.ErrInvalidEmailFormat
	}

	// 2. Find user
	user, err := s.userRepo.GetByEmail(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.WarnContext(ctx, "User not found for resend verification", logger.String("email", req.Username))
			return errors.ErrRecordNotFound
		}
		s.logger.ErrorContext(ctx, "Failed to find user", logger.ErrorField(err))
		return errors.ErrRecordQueryFailed
	}

	// 3. Check if user has pending email verification lock
	hasLock, err := s.hasEmailVerificationLock(ctx, req.Username)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to check email verification lock",
			logger.ErrorField(err), logger.String("email", req.Username))
		// Continue with resend process even if Redis check fails
	} else if !hasLock {
		s.logger.InfoContext(ctx, "User does not have pending email verification, may be already verified",
			logger.String("email", req.Username))
		return errors.ErrValidationFailed
	}

	// 4. Generate new verification token
	verificationToken, err := s.generateVerificationToken(ctx, req.Username, "email_verification")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate verification token", logger.ErrorField(err))
		return errors.ErrRecordCreateFailed
	}

	// 5. Send verification email
	expiresMinutes := authConfig.UserAuth.EmailAuth.ExpiresIn / secondsToMinutesConvert
	if err := s.emailService.SendVerificationEmail(ctx, expiresMinutes, req.Username, verificationToken.Token, s.appConfig.App.BaseURL, user.Language); err != nil {
		s.logger.ErrorContext(ctx, "Failed to send verification email", logger.ErrorField(err))
		return errors.ErrRecordCreateFailed
	}

	s.logger.InfoContext(ctx, "Verification email resent successfully", logger.String("email", req.Username))
	return nil
}

// OAuth2Login handles OAuth2 authentication
func (s *userAuthService) OAuth2Login(ctx context.Context, req *request.OAuth2LoginRequest, clientIP string) (*response.UserLoginResponse, error) {
	s.logger.InfoContext(ctx, "Processing OAuth2 login",
		logger.String("provider", req.Provider))

	// Check if OAuth2 service is available
	if s.oauth2Service == nil {
		s.logger.ErrorContext(ctx, "OAuth2 service not available")
		return nil, errors.ErrThirdPartyServiceUnavailable
	}

	// 1. Validate state parameter (prevent CSRF attacks)
	if err := s.oauth2Service.ValidateState(ctx, req.State, req.State); err != nil {
		s.logger.WarnContext(ctx, "OAuth2 state validation failed", logger.ErrorField(err))
		return nil, errors.ErrValidationFailed
	}

	// 2. Exchange authorization code for access token
	tokenResp, err := s.oauth2Service.ExchangeCodeForToken(ctx, req.Provider, req.Code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to exchange OAuth2 code for token", logger.ErrorField(err))
		return nil, errors.ErrValidationFailed
	}

	// 3. Get user information
	userInfo, err := s.oauth2Service.GetUserInfo(ctx, req.Provider, tokenResp.AccessToken)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get OAuth2 user info", logger.ErrorField(err))
		return nil, errors.ErrRecordQueryFailed
	}

	// 4. Find or create user
	user, err := s.findOrCreateOAuth2User(ctx, userInfo, req.Provider)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to find or create OAuth2 user", logger.ErrorField(err))
		return nil, err
	}

	// 5. Check user status
	if user.Status != UserStatusActive {
		s.logger.WarnContext(ctx, "OAuth2 user account is inactive", logger.String("email", user.Email))
		return nil, errors.ErrAccountDisabled
	}

	// Get user roles
	userRoleIDs := make([]uint, len(user.Roles))
	for i := range user.Roles {
		userRoleIDs[i] = user.Roles[i].ID
	}

	// 6. Generate JWT token using configured expiration time
	token, expiresAt, err := auth.GetGlobalJWT().GenerateTokenWithUserInfo(user.ID, user.Username, userRoleIDs)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate JWT token", logger.ErrorField(err))
		return nil, errors.ErrRecordCreateFailed
	}

	// 7. Store JWT token in database and Redis cache (according to security design)
	if err := s.storeJWTToken(ctx, token, user.ID, expiresAt); err != nil {
		s.logger.WarnContext(ctx, "Failed to store JWT token in database/Redis", logger.ErrorField(err))
		// Don't fail the login process, as token is still valid for authentication
	}

	// 8. Update last login time
	now := time.Now()
	updateUser := model.User{
		ID:          user.ID,
		LastLoginAt: &now,
		LastLoginIP: clientIP,
	}
	if err := s.userRepo.Update(ctx, &updateUser); err != nil {
		s.logger.WarnContext(ctx, "Failed to update login info", logger.ErrorField(err))
		// Don't return error as login is already successful
	}

	s.logger.InfoContext(ctx, "OAuth2 login successful",
		logger.Uint("user_id", user.ID),
		logger.String("provider", req.Provider))

	return &response.UserLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *response.BuildUserResponse(user),
	}, nil
}

// generateUsernameFromEmail 从邮箱生成用户名
func (s *userAuthService) generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		username := parts[0]
		// 确保用户名唯一
		counter := 1
		originalUsername := username
		for {
			exists, err := s.userRepo.ExistsByUsername(context.Background(), username)
			if err != nil || !exists {
				break
			}
			username = fmt.Sprintf("%s%d", originalUsername, counter)
			counter++
		}
		return username
	}
	return "user" + fmt.Sprintf("%d", time.Now().Unix())
}

// findOrCreateOAuth2User 查找或创建OAuth2用户
func (s *userAuthService) findOrCreateOAuth2User(ctx context.Context, userInfo *OAuth2UserInfo, provider string) (*model.User, error) {
	// 1. 先尝试通过邮箱查找用户
	if userInfo.Email != "" {
		user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
		if err == nil {
			// 用户已存在，更新OAuth2信息
			s.logger.InfoContext(ctx, "Found existing user for OAuth2 login",
				logger.String("email", userInfo.Email),
				logger.String("provider", provider))
			return user, nil
		} else if err != gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound)
		}
	}

	if !authConfig.UserAuth.OAuth2.AutoRegister {
		return nil, errors.NewAppError(errors.CodeRecordNotFound)
	}

	// 3. 创建新用户
	username := s.generateUsernameFromEmail(userInfo.Email)
	if username == "" && userInfo.Username != "" {
		username = userInfo.Username
	}
	if username == "" {
		username = fmt.Sprintf("oauth2_user_%d", time.Now().Unix())
	}

	user := &model.User{
		Username: username,
		Email:    userInfo.Email,
		Nickname: userInfo.Name,
		Avatar:   userInfo.Avatar,
		Status:   UserStatusActive,
		Language: "zh-CN",
		Timezone: "Asia/Shanghai",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to create OAuth2 user", logger.ErrorField(err))
		return nil, errors.NewAppError(errors.CodeRecordCreateFailed)
	}

	s.logger.InfoContext(ctx, "Created new OAuth2 user",
		logger.Uint("user_id", user.ID),
		logger.String("email", userInfo.Email),
		logger.String("provider", provider))

	return user, nil
}

// setEmailVerificationLock sets email verification lock in Redis
func (s *userAuthService) setEmailVerificationLock(ctx context.Context, email, value string, ttl int) error {
	key := redis.FormatRedisKey(redis.RK_EMAIL_VERILOCK, email)
	err := redis.Set(ctx, key, value, time.Duration(ttl)*time.Second)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to set email verification lock in Redis",
			logger.ErrorField(err), logger.String("email", email))
		return err
	}

	s.logger.InfoContext(ctx, "Email verification lock set successfully",
		logger.String("email", email), logger.String("key", key))
	return nil
}

// hasEmailVerificationLock checks if email verification lock exists in Redis
func (s *userAuthService) hasEmailVerificationLock(ctx context.Context, email string) (bool, error) {
	key := redis.FormatRedisKey(redis.RK_EMAIL_VERILOCK, email)
	result, err := redis.Get(ctx, key)
	if err != nil {
		if err.Error() == redis.RedisNilError {
			// Key does not exist, no lock
			return false, nil
		}
		s.logger.ErrorContext(ctx, "Failed to check email verification lock in Redis",
			logger.ErrorField(err), logger.String("email", email))
		return false, err
	}
	hasLock := result != ""
	s.logger.InfoContext(ctx, "Email verification lock check completed",
		logger.String("email", email))
	return hasLock, nil
}

// deleteEmailVerificationLock removes email verification lock from Redis
func (s *userAuthService) deleteEmailVerificationLock(ctx context.Context, email string) error {
	key := redis.FormatRedisKey(redis.RK_EMAIL_VERILOCK, email)
	_, err := redis.Del(ctx, key)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete email verification lock from Redis",
			logger.ErrorField(err), logger.String("email", email))
		return err
	}

	s.logger.InfoContext(ctx, "Email verification lock deleted successfully",
		logger.String("email", email), logger.String("key", key))
	return nil
}

// storeTokenInRedis stores verification token in Redis with expiration
func (s *userAuthService) storeTokenInRedis(ctx context.Context, token string, verificationToken *VerificationToken) error {
	// Generate hash of the token for Redis key
	tokenHash := auth.HashToken(token)
	key := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)

	// Marshal token to JSON
	tokenJSON, err := json.Marshal(verificationToken)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal token to JSON",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."))
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Store in Redis with TTL
	tokenExpiresSeconds := authConfig.UserAuth.EmailAuth.ExpiresIn
	ttl := time.Duration(tokenExpiresSeconds) * time.Second
	err = redis.Set(ctx, key, string(tokenJSON), ttl)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to store token in Redis",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."), logger.String("key", key))
		return fmt.Errorf("failed to store token in redis: %w", err)
	}

	s.logger.InfoContext(ctx, "Token stored in Redis successfully",
		logger.String("token", token[:8]+"..."), logger.String("key", key), logger.Duration("ttl", ttl))
	return nil
}

// getTokenFromRedis retrieves verification token from Redis
func (s *userAuthService) getTokenFromRedis(ctx context.Context, token string) (*VerificationToken, error) {
	// Generate hash of the token for Redis key
	tokenHash := auth.HashToken(token)
	key := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)

	// Get from Redis
	result, err := redis.Get(ctx, key)
	if err != nil {
		if err.Error() == redis.RedisNilError {
			s.logger.WarnContext(ctx, "Token not found in Redis",
				logger.String("token", token[:8]+"..."), logger.String("key", key))
			return nil, fmt.Errorf("token not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get token from Redis",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."), logger.String("key", key))
		return nil, fmt.Errorf("failed to get token from redis: %w", err)
	}

	// Unmarshal JSON to token
	var verificationToken VerificationToken
	err = json.Unmarshal([]byte(result), &verificationToken)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to unmarshal token from JSON",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."))
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	s.logger.InfoContext(ctx, "Token retrieved from Redis successfully",
		logger.String("token", token[:8]+"..."), logger.String("key", key))
	return &verificationToken, nil
}

// storeJWTToken stores JWT token in database (api_tokens table) and Redis cache
func (s *userAuthService) storeJWTToken(ctx context.Context, token string, userID uint, expiresAt time.Time) error {
	// Generate hash of the token for database storage
	tokenHash := auth.HashToken(token)

	// 1. Store JWT token in api_tokens database table
	apiToken := &model.APIToken{
		Name:        "JWT Login Token",
		Token:       token,
		TokenHash:   tokenHash,
		UserID:      userID,
		Scopes:      model.JSON{"scopes": []string{"api:access", "user:profile"}},
		Description: "JWT token generated during user login",
		ExpiresAt:   &expiresAt,
	}

	// Store in database using APITokenRepository
	if err := s.apiTokenRepo.Create(ctx, apiToken); err != nil {
		s.logger.ErrorContext(ctx, "Failed to store JWT token in database",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."))
		return fmt.Errorf("failed to store JWT token in database: %w", err)
	}

	s.logger.InfoContext(ctx, "JWT token stored in database successfully",
		logger.String("token", token[:8]+"..."), logger.Uint("token_id", apiToken.ID))

	// 2. Store JWT token in Redis cache with expiration
	// Generate Redis key for JWT token
	redisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)

	// Create JWT token info for Redis storage
	tokenInfo := map[string]interface{}{
		"token":      token,
		"user_id":    userID,
		"expires_at": expiresAt,
		"type":       "jwt_login",
		"created_at": time.Now(),
		"token_id":   apiToken.ID, // Reference to database record
	}

	// Marshal token info to JSON
	tokenJSON, err := json.Marshal(tokenInfo)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to marshal JWT token info to JSON",
			logger.ErrorField(err), logger.String("token", token[:8]+"..."))
		// Don't fail if Redis storage fails, database storage is more important
	} else {
		// Calculate TTL for Redis from configured expiration time
		ttl := time.Duration(auth.GetGlobalJWT().GetTokenExpiresInSeconds()) * time.Second
		// Store in Redis
		err = redis.Set(ctx, redisKey, string(tokenJSON), ttl)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to store JWT token in Redis",
				logger.ErrorField(err), logger.String("token", token[:8]+"..."), logger.String("key", redisKey))
			// Don't fail if Redis storage fails, database storage is more important
		} else {
			s.logger.InfoContext(ctx, "JWT token stored in Redis successfully",
				logger.String("token", token[:8]+"..."), logger.String("key", redisKey), logger.Duration("ttl", ttl))
		}
	}

	return nil
}

// Logout handles user logout by deleting JWT token from database and Redis
func (s *userAuthService) Logout(ctx context.Context, token string) error {
	// Generate token hash for lookup
	tokenHash := auth.HashToken(token)

	// 1. Delete JWT token from database (api_tokens table)
	apiToken, err := s.apiTokenRepo.GetByToken(ctx, tokenHash)
	if err == nil && apiToken != nil {
		// Delete from database using repository method
		deleteErr := s.apiTokenRepo.Delete(ctx, apiToken.ID)
		if deleteErr != nil {
			s.logger.WarnContext(ctx, "Failed to delete JWT token from database", logger.ErrorField(deleteErr))
		} else {
			s.logger.InfoContext(ctx, "JWT token deleted from database successfully",
				logger.Uint("token_id", apiToken.ID))
		}
	} else {
		s.logger.WarnContext(ctx, "JWT token not found in database", logger.ErrorField(err))
	}

	// 2. Delete JWT token from Redis cache
	redisKey := redis.FormatRedisKey(redis.RK_AUTH_TOKEN, tokenHash)

	_, err = redis.Del(ctx, redisKey)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to delete JWT token from Redis",
			logger.ErrorField(err), logger.String("key", redisKey))
	} else {
		s.logger.InfoContext(ctx, "JWT token deleted from Redis successfully",
			logger.String("token", token[:8]+"..."), logger.String("key", redisKey))
	}

	return nil
}

// generateVerificationToken generates a verification token for email verification or password reset
func (s *userAuthService) generateVerificationToken(ctx context.Context, email, tokenType string) (*VerificationToken, error) {
	// Generate random token
	tokenBytes := make([]byte, tokenRandomBytesSize)
	if _, err := rand.Read(tokenBytes); err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate random token", logger.ErrorField(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenExpiresSeconds := authConfig.UserAuth.EmailAuth.ExpiresIn

	token := hex.EncodeToString(tokenBytes)

	// Create token object
	verificationToken := &VerificationToken{
		Token:     token,
		Email:     email,
		Type:      tokenType,
		ExpiresAt: time.Now().Add(time.Duration(tokenExpiresSeconds) * time.Second),
		CreatedAt: time.Now(),
		Used:      false,
	}

	// Store token in Redis with expiration
	if err := s.storeTokenInRedis(ctx, token, verificationToken); err != nil {
		s.logger.ErrorContext(ctx, "Failed to store token in Redis", logger.ErrorField(err))
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	s.logger.InfoContext(ctx, "Verification token generated successfully",
		logger.String("email", email),
		logger.String("type", tokenType),
		logger.String("token", token[:8]+"..."))

	return verificationToken, nil
}

// validateVerificationToken validates a verification token
func (s *userAuthService) validateVerificationToken(ctx context.Context, token, tokenType string) (*VerificationToken, error) {
	s.logger.InfoContext(ctx, "Validating verification token",
		logger.String("token", token[:8]+"..."),
		logger.String("type", tokenType))

	// Get token from Redis
	verificationToken, err := s.getTokenFromRedis(ctx, token)
	if err != nil {
		s.logger.WarnContext(ctx, "Token not found in Redis",
			logger.String("token", token[:8]+"..."), logger.ErrorField(err))
		return nil, fmt.Errorf("invalid token")
	}

	// Check token type
	if verificationToken.Type != tokenType {
		s.logger.WarnContext(ctx, "Token type mismatch",
			logger.String("expected", tokenType),
			logger.String("actual", verificationToken.Type))
		return nil, fmt.Errorf("invalid token type")
	}

	// Check if already used
	if verificationToken.Used {
		s.logger.WarnContext(ctx, "Token already used", logger.String("token", token[:8]+"..."))
		return nil, fmt.Errorf("token already used")
	}

	// Check if expired
	if time.Now().After(verificationToken.ExpiresAt) {
		s.logger.WarnContext(ctx, "Token expired",
			logger.String("token", token[:8]+"..."),
			logger.Time("expires_at", verificationToken.ExpiresAt))
		return nil, fmt.Errorf("token expired")
	}

	s.logger.InfoContext(ctx, "Token validation successful",
		logger.String("email", verificationToken.Email))

	return verificationToken, nil
}

// markTokenAsUsed marks a token as used
func (s *userAuthService) markTokenAsUsed(ctx context.Context, token string) error {
	s.logger.InfoContext(ctx, "Marking token as used", logger.String("token", token[:8]+"..."))

	// Get token from Redis
	verificationToken, err := s.getTokenFromRedis(ctx, token)
	if err != nil {
		s.logger.WarnContext(ctx, "Token not found when marking as used",
			logger.String("token", token[:8]+"..."), logger.ErrorField(err))
		return fmt.Errorf("token not found")
	}

	// Mark as used
	verificationToken.Used = true

	// Update token in Redis
	if err := s.storeTokenInRedis(ctx, token, verificationToken); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update token as used in Redis",
			logger.String("token", token[:8]+"..."), logger.ErrorField(err))
		return fmt.Errorf("failed to update token: %w", err)
	}

	s.logger.InfoContext(ctx, "Token marked as used", logger.String("token", token[:8]+"..."))
	return nil
}

// storeUserPreferencesInRedis stores user preferences in Redis using HASH structure
func (s *userAuthService) storeUserPreferencesInRedis(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "Storing user preferences in Redis", logger.Uint("user_id", userID))

	// 1. Query user preferences from database using UserProfileRepository
	// Get language and timezone preferences (general category)
	preferences, err := s.userProfileRepo.GetUserConfigsByCategory(ctx, userID, constants.UserCategory)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user preferences from database",
			logger.ErrorField(err), logger.Uint("user_id", userID))
		return err
	}

	// 2. Check for missing preferences and fill with system defaults
	preferences = s.fillMissingPreferences(ctx, preferences)

	// 3. Generate Redis key for user preferences
	redisKey := redis.FormatRedisKeyWithID(redis.RK_USER_PREFERENCES, userID)

	// 4. Convert preferences to Redis hash field-value pairs
	preferencesMap := make(map[string]interface{})
	for _, pref := range preferences {
		if pref.ConfigKey != "" {
			preferencesMap[pref.ConfigKey] = pref.ConfigValue
		}
	}

	// 5. Store preferences in Redis using HASH structure
	if len(preferencesMap) > 0 {
		// Convert map to alternating key-value slice for HSet
		hashValues := make([]interface{}, 0, len(preferencesMap)*hashKeyValuePairs)
		for key, value := range preferencesMap {
			hashValues = append(hashValues, key, value)
		}

		_, err = redis.HSet(ctx, redisKey, hashValues...)
		if err != nil {
			s.logger.ErrorContext(ctx, "Failed to store preferences in Redis HASH",
				logger.ErrorField(err), logger.String("key", redisKey))
			return err
		}

		s.logger.InfoContext(ctx, "User preferences stored in Redis successfully",
			logger.Uint("user_id", userID),
			logger.String("key", redisKey),
			logger.Int("preferences_count", len(preferencesMap)))
	} else {
		s.logger.InfoContext(ctx, "No user preferences found to store in Redis",
			logger.Uint("user_id", userID))
	}

	return nil
}

// fillMissingPreferences checks for missing language and timezone preferences
// and fills them with system defaults if needed
func (s *userAuthService) fillMissingPreferences(ctx context.Context, preferences []*model.UserProfile) []*model.UserProfile {
	// Create a map to track existing preferences
	existingPrefs := make(map[string]string)
	for _, pref := range preferences {
		if pref.ConfigKey != "" {
			existingPrefs[pref.ConfigKey] = pref.ConfigValue
		}
	}

	preferencesMapping := map[string]string{
		constants.UserLanguage: constants.SystemLanguage,
		constants.UserTimezone: constants.SystemTimezone,
	}

	for userPref, sysConfig := range preferencesMapping {
		if _, hasKey := existingPrefs[userPref]; !hasKey {
			configValue := s.getSystemConfigByKey(ctx, sysConfig)

			// If system config is empty, use hardcoded defaults
			if configValue == "" {
				switch userPref {
				case constants.UserLanguage:
					configValue = constants.DefaultLanguage
				case constants.UserTimezone:
					configValue = constants.DefaultTimeZone
				}
			}

			preferences = append(preferences, &model.UserProfile{
				ConfigKey:   userPref,
				ConfigValue: configValue,
			})
			s.logger.InfoContext(ctx, "Added default preference value",
				logger.String(userPref, configValue))
		}
	}
	return preferences
}

// getSystemConfigByKey retrieves system configuration by key
func (s *userAuthService) getSystemConfigByKey(ctx context.Context, configKey string) string {
	config, err := s.systemConfigRepo.GetByKey(ctx, configKey)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.WarnContext(ctx, "Failed to get system config", logger.ErrorField(err))
	}

	if config != nil && config.ConfigValue != "" {
		return config.ConfigValue
	}
	return ""
}

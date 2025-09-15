package service

import (
	"api-service/internal/constants"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/i18n"
	"api-service/pkg/logger"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	// Email verification code settings
	emailCodeExpiry    = 5 * time.Minute
	backupCodesCount   = 10
	backupCodeByteSize = 4
	emailCodeByteSize  = 3
	emailCodeModulo    = 1000000
)

type twoFactorService struct {
	twoFactorRepo repository.UserTwoFactorRepository
	db            *gorm.DB
	logger        logger.Logger
	i18n          *i18n.I18n
}

// NewTwoFactorService creates a new two-factor authentication service instance
func NewTwoFactorService(
	twoFactorRepo repository.UserTwoFactorRepository,
	db *gorm.DB,
	logger logger.Logger,
	i18n *i18n.I18n,
) service.TwoFactorService {
	return &twoFactorService{
		twoFactorRepo: twoFactorRepo,
		db:            db,
		logger:        logger,
		i18n:          i18n,
	}
}

// GetTwoFactorStatus gets user's two-factor authentication status
func (s *twoFactorService) GetTwoFactorStatus(ctx context.Context, userID uint) (*response.TwoFactorStatusResponse, error) {
	s.logger.InfoContext(ctx, "Getting two-factor status",
		logger.String("service", "two-factor"),
		logger.String("operation", "GetTwoFactorStatus"),
		logger.Uint("user_id", userID))

	methods, err := s.twoFactorRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get two-factor methods", logger.ErrorField(err))
		return nil, err
	}

	return response.ConvertToTwoFactorStatusResponse(methods), nil
}

// EnableTOTP enables TOTP two-factor authentication
func (s *twoFactorService) EnableTOTP(ctx context.Context, userID uint) (*response.TOTPSetupResponse, error) {
	s.logger.InfoContext(ctx, "Enabling TOTP",
		logger.String("service", "two-factor"),
		logger.String("operation", "EnableTOTP"),
		logger.Uint("user_id", userID))

	// Check if TOTP is already enabled
	existing, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil && !errors.Is(err, errors.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing TOTP", logger.ErrorField(err))
		return nil, err
	}

	if existing != nil && existing.Enabled {
		return nil, errors.NewAppErrorWithMessage(errors.CodeResourceStateNotAllowed, "TOTP is already enabled")
	}

	// Generate TOTP secret
	secret, err := auth.GenerateSimpleTOTPSecret()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate TOTP secret", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to generate TOTP secret")
	}

	// Generate QR code URL
	qrCodeURL := auth.GenerateTOTPQRCode(secret, fmt.Sprintf("user_%d", userID), "Websoft9")

	// Generate backup codes
	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate backup codes", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to generate backup codes")
	}

	// Create or update TOTP record (not enabled yet)
	twoFactor := &model.UserTwoFactor{
		UserID:      userID,
		Method:      "TOTP",
		Secret:      secret,
		BackupCodes: model.JSON{"codes": backupCodes},
		Enabled:     false, // Will be enabled after confirmation
	}

	if existing != nil {
		twoFactor.ID = existing.ID
		if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
			s.logger.ErrorContext(ctx, "Failed to update TOTP record", logger.ErrorField(err))
			return nil, err
		}
	} else {
		if err := s.twoFactorRepo.Create(ctx, twoFactor); err != nil {
			s.logger.ErrorContext(ctx, "Failed to create TOTP record", logger.ErrorField(err))
			return nil, err
		}
	}

	s.logger.InfoContext(ctx, "TOTP setup initiated successfully",
		logger.Uint("user_id", userID))

	return &response.TOTPSetupResponse{
		Secret:      secret,
		QRCodeURL:   qrCodeURL,
		BackupCodes: backupCodes,
	}, nil
}

// ConfirmTOTP confirms TOTP setup with verification code
func (s *twoFactorService) ConfirmTOTP(ctx context.Context, userID uint, code string) (*response.TOTPConfirmResponse, error) {
	s.logger.InfoContext(ctx, "Confirming TOTP",
		logger.String("service", "two-factor"),
		logger.String("operation", "ConfirmTOTP"),
		logger.Uint("user_id", userID))

	// Get TOTP record
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get TOTP record", logger.ErrorField(err))
		return nil, err
	}

	// Verify TOTP code
	valid, err := auth.ValidateSimpleTOTP(twoFactor.Secret, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to validate TOTP code", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeValidationFailed, "failed to validate TOTP code")
	}

	if !valid {
		return nil, errors.NewAppErrorWithMessage(errors.CodeValidationFailed, "invalid code")
	}

	// Enable TOTP
	twoFactor.Enabled = true
	now := time.Now()
	twoFactor.VerifiedAt = &now

	if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "TOTP confirmed and enabled successfully",
		logger.Uint("user_id", userID))

	// Extract backup codes
	var backupCodes []string
	if twoFactor.BackupCodes != nil {
		if codes, exists := twoFactor.BackupCodes["codes"]; exists {
			if codesSlice, ok := codes.([]interface{}); ok {
				backupCodes = make([]string, len(codesSlice))
				for i, code := range codesSlice {
					if s, ok := code.(string); ok {
						backupCodes[i] = s
					}
				}
			}
		}
	}

	return &response.TOTPConfirmResponse{
		Enabled:     true,
		BackupCodes: backupCodes,
	}, nil
}

// DisableTOTP disables TOTP two-factor authentication
func (s *twoFactorService) DisableTOTP(ctx context.Context, userID uint, code string) error {
	s.logger.InfoContext(ctx, "Disabling TOTP",
		logger.String("service", "two-factor"),
		logger.String("operation", "DisableTOTP"),
		logger.Uint("user_id", userID))

	// Get TOTP record
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get TOTP record", logger.ErrorField(err))
		return err
	}

	// Verify TOTP code
	valid, err := auth.ValidateSimpleTOTP(twoFactor.Secret, code)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to validate TOTP code", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeValidationFailed, "failed to validate TOTP code")
	}

	if !valid {
		return errors.NewAppErrorWithMessage(errors.CodeValidationFailed, "invalid code")
	}

	// Delete TOTP record
	if err := s.twoFactorRepo.Delete(ctx, twoFactor.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete TOTP record", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "TOTP disabled successfully",
		logger.Uint("user_id", userID))

	return nil
}

// EnableEmailTwoFactor enables email two-factor authentication
func (s *twoFactorService) EnableEmailTwoFactor(ctx context.Context, userID uint, email string) error {
	s.logger.InfoContext(ctx, "Enabling email two-factor",
		logger.String("service", "two-factor"),
		logger.String("operation", "EnableEmailTwoFactor"),
		logger.Uint("user_id", userID))

	// Check if email 2FA is already enabled
	existing, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "EMAIL")
	if err != nil && !errors.Is(err, errors.ErrRecordNotFound) {
		s.logger.ErrorContext(ctx, "Failed to check existing email 2FA", logger.ErrorField(err))
		return err
	}

	// Create or update email 2FA record
	twoFactor := &model.UserTwoFactor{
		UserID:  userID,
		Method:  "EMAIL",
		Email:   email,
		Enabled: true,
	}

	if existing != nil {
		twoFactor.ID = existing.ID
		if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
			s.logger.ErrorContext(ctx, "Failed to update email 2FA record", logger.ErrorField(err))
			return err
		}
	} else {
		if err := s.twoFactorRepo.Create(ctx, twoFactor); err != nil {
			s.logger.ErrorContext(ctx, "Failed to create email 2FA record", logger.ErrorField(err))
			return err
		}
	}

	s.logger.InfoContext(ctx, "Email two-factor enabled successfully",
		logger.Uint("user_id", userID))

	return nil
}

// DisableEmailTwoFactor disables email two-factor authentication
func (s *twoFactorService) DisableEmailTwoFactor(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "Disabling email two-factor",
		logger.String("service", "two-factor"),
		logger.String("operation", "DisableEmailTwoFactor"),
		logger.Uint("user_id", userID))

	// Get email 2FA record
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "EMAIL")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get email 2FA record", logger.ErrorField(err))
		return err
	}

	// Delete email 2FA record
	if err := s.twoFactorRepo.Delete(ctx, twoFactor.ID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete email 2FA record", logger.ErrorField(err))
		return err
	}

	s.logger.InfoContext(ctx, "Email two-factor disabled successfully",
		logger.Uint("user_id", userID))

	return nil
}

// SendEmailCode sends email verification code
func (s *twoFactorService) SendEmailCode(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "Sending email code",
		logger.String("service", "two-factor"),
		logger.String("operation", "SendEmailCode"),
		logger.Uint("user_id", userID))

	// Get email 2FA record
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "EMAIL")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get email 2FA record", logger.ErrorField(err))
		return err
	}

	// Generate verification code
	code := s.generateEmailCode()

	// Store code with expiration (5 minutes)
	expiry := time.Now().Add(emailCodeExpiry)
	twoFactor.Secret = code
	twoFactor.VerifiedAt = &expiry // Reuse this field for code expiry

	if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
		s.logger.ErrorContext(ctx, "Failed to store email code", logger.ErrorField(err))
		return err
	}

	// TODO: Send email with code
	// This would integrate with an email service
	s.logger.InfoContext(ctx, "Email code generated (email sending not implemented)",
		logger.Uint("user_id", userID),
		logger.String("code", code))

	return nil
}

// VerifyTwoFactor verifies two-factor authentication code
func (s *twoFactorService) VerifyTwoFactor(ctx context.Context, userID uint, code, method string) (*response.TwoFactorVerificationResponse, error) {
	s.logger.InfoContext(ctx, "Verifying two-factor code",
		logger.String("service", "two-factor"),
		logger.String("operation", "VerifyTwoFactor"),
		logger.Uint("user_id", userID),
		logger.String("method", method))

	switch method {
	case constants.TwoFactorMethodTOTP:
		return s.verifyTOTP(ctx, userID, code)
	case constants.TwoFactorMethodEmail:
		return s.verifyEmail(ctx, userID, code)
	case constants.TwoFactorMethodBackup:
		return s.verifyBackupCode(ctx, userID, code)
	default:
		return nil, errors.NewAppErrorWithMessage(errors.CodeInvalidParameterFormat, "unsupported 2FA method")
	}
}

// GenerateBackupCodes generates backup codes for two-factor authentication
func (s *twoFactorService) GenerateBackupCodes(ctx context.Context, userID uint) (*response.BackupCodesResponse, error) {
	s.logger.InfoContext(ctx, "Generating backup codes",
		logger.String("service", "two-factor"),
		logger.String("operation", "GenerateBackupCodes"),
		logger.Uint("user_id", userID))

	// Get TOTP record (backup codes are associated with TOTP)
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get TOTP record", logger.ErrorField(err))
		return nil, err
	}

	// Generate new backup codes
	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to generate backup codes", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "failed to generate backup codes")
	}

	// Update backup codes
	twoFactor.BackupCodes = model.JSON{"codes": backupCodes}

	if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update backup codes", logger.ErrorField(err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "Backup codes generated successfully",
		logger.Uint("user_id", userID))

	return &response.BackupCodesResponse{
		Codes: backupCodes,
	}, nil
}

// Helper methods

func (s *twoFactorService) verifyTOTP(ctx context.Context, userID uint, code string) (*response.TwoFactorVerificationResponse, error) {
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil {
		return &response.TwoFactorVerificationResponse{Valid: false}, nil
	}

	valid, err := auth.ValidateSimpleTOTP(twoFactor.Secret, code)
	if err != nil || !valid {
		return &response.TwoFactorVerificationResponse{Valid: false}, nil
	}

	return &response.TwoFactorVerificationResponse{
		Valid:  true,
		UserID: userID,
		Method: "totp",
	}, nil
}

func (s *twoFactorService) verifyEmail(ctx context.Context, userID uint, code string) (*response.TwoFactorVerificationResponse, error) {
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "EMAIL")
	if err != nil {
		return &response.TwoFactorVerificationResponse{Valid: false}, nil
	}

	// Check if code matches and hasn't expired
	if twoFactor.Secret != code || twoFactor.VerifiedAt == nil || time.Now().After(*twoFactor.VerifiedAt) {
		return &response.TwoFactorVerificationResponse{Valid: false}, nil
	}

	return &response.TwoFactorVerificationResponse{
		Valid:  true,
		UserID: userID,
		Method: "email",
	}, nil
}

func (s *twoFactorService) verifyBackupCode(ctx context.Context, userID uint, code string) (*response.TwoFactorVerificationResponse, error) {
	twoFactor, err := s.twoFactorRepo.GetByUserIDAndMethod(ctx, userID, "TOTP")
	if err != nil {
		return &response.TwoFactorVerificationResponse{Valid: false}, nil
	}

	// Check backup codes
	if twoFactor.BackupCodes != nil {
		if codes, exists := twoFactor.BackupCodes["codes"]; exists {
			if codesSlice, ok := codes.([]interface{}); ok {
				for i, backupCode := range codesSlice {
					if codeStr, ok := backupCode.(string); ok && codeStr == code {
						// Remove used backup code
						codesSlice = append(codesSlice[:i], codesSlice[i+1:]...)
						twoFactor.BackupCodes["codes"] = codesSlice
						if err := s.twoFactorRepo.Update(ctx, twoFactor); err != nil {
							return nil, err
						}

						return &response.TwoFactorVerificationResponse{
							Valid:  true,
							UserID: userID,
							Method: "backup",
						}, nil
					}
				}
			}
		}
	}

	return &response.TwoFactorVerificationResponse{Valid: false}, nil
}

func (s *twoFactorService) generateBackupCodes() ([]string, error) {
	codes := make([]string, backupCodesCount)
	for i := 0; i < backupCodesCount; i++ {
		bytes := make([]byte, backupCodeByteSize)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = hex.EncodeToString(bytes)
	}
	return codes, nil
}

func (s *twoFactorService) generateEmailCode() string {
	bytes := make([]byte, emailCodeByteSize)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to time-based generation if crypto/rand fails
		return fmt.Sprintf("%06d", time.Now().UnixNano()%emailCodeModulo)
	}
	return fmt.Sprintf("%06d", int(bytes[0])<<16|int(bytes[1])<<8|int(bytes[2]))[:6]
}

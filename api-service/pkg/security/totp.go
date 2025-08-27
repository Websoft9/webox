package security

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	// TOTP constants
	DefaultTOTPSecretSize = 32
	DefaultTOTPPeriod     = 30
	QRCodeImageSize       = 200
	BackupCodeBytesSize   = 8
	BackupCodeMaxLength   = 8
	BackupCodeModulo      = 100
)

// TOTPConfig TOTP configuration
type TOTPConfig struct {
	Issuer      string
	AccountName string
	SecretSize  uint
	Algorithm   otp.Algorithm
	Digits      otp.Digits
	Period      uint
}

// DefaultTOTPConfig default TOTP configuration
var DefaultTOTPConfig = TOTPConfig{
	Issuer:     "Websoft9",
	SecretSize: DefaultTOTPSecretSize,
	Algorithm:  otp.AlgorithmSHA1,
	Digits:     otp.DigitsSix,
	Period:     DefaultTOTPPeriod,
}

// GenerateTOTPSecret generates TOTP secret
func GenerateTOTPSecret(accountName string, config *TOTPConfig) (*otp.Key, error) {
	if config == nil {
		config = &DefaultTOTPConfig
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.Issuer,
		AccountName: accountName,
		SecretSize:  config.SecretSize,
		Algorithm:   config.Algorithm,
		Digits:      config.Digits,
		Period:      config.Period,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	return key, nil
}

// ValidateTOTP validates TOTP code
func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}

// GenerateQRCode generates QR code image data
func GenerateQRCode(key *otp.Key) ([]byte, error) {
	// Generate QR code image
	img, err := key.Image(QRCodeImageSize, QRCodeImageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}

	// Encode image as PNG format
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code image: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateBackupCodes generates backup codes
func GenerateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)

	for i := 0; i < count; i++ {
		code, err := generateBackupCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate backup code: %w", err)
		}
		codes[i] = code
	}

	return codes, nil
}

// generateBackupCode generates single backup code
func generateBackupCode() (string, error) {
	// Generate 8 bytes of random data
	bytes := make([]byte, BackupCodeBytesSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert to 8-digit numeric string
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%02d", int(b)%BackupCodeModulo)
	}

	// Take first 8 digits
	if len(code) > BackupCodeMaxLength {
		code = code[:BackupCodeMaxLength]
	}

	return code, nil
}

// ValidateBackupCode validates backup code
func ValidateBackupCode(code string, backupCodes []string) bool {
	for _, backupCode := range backupCodes {
		if code == backupCode {
			return true
		}
	}
	return false
}

// RemoveUsedBackupCode removes used backup code
func RemoveUsedBackupCode(usedCode string, backupCodes []string) []string {
	result := make([]string, 0, len(backupCodes))
	for _, code := range backupCodes {
		if code != usedCode {
			result = append(result, code)
		}
	}
	return result
}

// EncodeSecret encodes secret to Base32 string
func EncodeSecret(secret []byte) string {
	return base32.StdEncoding.EncodeToString(secret)
}

// DecodeSecret decodes Base32 string to secret
func DecodeSecret(encodedSecret string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(encodedSecret)
}

// GenerateSimpleTOTPSecret generates a TOTP secret (simplified version)
func GenerateSimpleTOTPSecret() (string, error) {
	key, err := GenerateTOTPSecret("user", nil)
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// GenerateTOTPQRCode generates a QR code URL for TOTP setup
func GenerateTOTPQRCode(secret, accountName, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, accountName, secret, issuer)
}

// ValidateSimpleTOTP validates a TOTP code (simplified version)
func ValidateSimpleTOTP(secret, code string) (bool, error) {
	valid := totp.Validate(code, secret)
	return valid, nil
}

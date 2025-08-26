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

// TOTPConfig TOTP配置
type TOTPConfig struct {
	Issuer      string
	AccountName string
	SecretSize  uint
	Algorithm   otp.Algorithm
	Digits      otp.Digits
	Period      uint
}

// DefaultTOTPConfig 默认TOTP配置
var DefaultTOTPConfig = TOTPConfig{
	Issuer:     "Websoft9",
	SecretSize: 32,
	Algorithm:  otp.AlgorithmSHA1,
	Digits:     otp.DigitsSix,
	Period:     30,
}

// GenerateTOTPSecret 生成TOTP密钥
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

// ValidateTOTP 验证TOTP代码
func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}

// GenerateQRCode 生成二维码图片数据
func GenerateQRCode(key *otp.Key) ([]byte, error) {
	// 生成二维码图片
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code image: %w", err)
	}

	// 将图片编码为PNG格式
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code image: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateBackupCodes 生成备用码
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

// generateBackupCode 生成单个备用码
func generateBackupCode() (string, error) {
	// 生成8字节随机数据
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// 转换为8位数字字符串
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%02d", int(b)%100)
	}

	// 取前8位
	if len(code) > 8 {
		code = code[:8]
	}

	return code, nil
}

// ValidateBackupCode 验证备用码
func ValidateBackupCode(code string, backupCodes []string) bool {
	for _, backupCode := range backupCodes {
		if code == backupCode {
			return true
		}
	}
	return false
}

// RemoveUsedBackupCode 移除已使用的备用码
func RemoveUsedBackupCode(usedCode string, backupCodes []string) []string {
	result := make([]string, 0, len(backupCodes))
	for _, code := range backupCodes {
		if code != usedCode {
			result = append(result, code)
		}
	}
	return result
}

// EncodeSecret 编码密钥为Base32字符串
func EncodeSecret(secret []byte) string {
	return base32.StdEncoding.EncodeToString(secret)
}

// DecodeSecret 解码Base32字符串为密钥
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

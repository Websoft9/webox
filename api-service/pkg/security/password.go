package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	MaxLength        int  `json:"max_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireLowercase bool `json:"require_lowercase"`
	RequireNumbers   bool `json:"require_numbers"`
	RequireSymbols   bool `json:"require_symbols"`
}

// DefaultPasswordPolicy 默认密码策略
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:        8,
	MaxLength:        128,
	RequireUppercase: true,
	RequireLowercase: true,
	RequireNumbers:   true,
	RequireSymbols:   false,
}

// ValidatePassword 验证密码是否符合策略
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = &DefaultPasswordPolicy
	}

	// 检查长度
	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}
	if len(password) > policy.MaxLength {
		return fmt.Errorf("password must be no more than %d characters long", policy.MaxLength)
	}

	// 检查大写字母
	if policy.RequireUppercase {
		matched, _ := regexp.MatchString(`[A-Z]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one uppercase letter")
		}
	}

	// 检查小写字母
	if policy.RequireLowercase {
		matched, _ := regexp.MatchString(`[a-z]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one lowercase letter")
		}
	}

	// 检查数字
	if policy.RequireNumbers {
		matched, _ := regexp.MatchString(`[0-9]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one number")
		}
	}

	// 检查特殊字符
	if policy.RequireSymbols {
		matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
		if !matched {
			return fmt.Errorf("password must contain at least one special character")
		}
	}

	return nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int, includeSymbols bool) (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers   = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	charset := lowercase + uppercase + numbers
	if includeSymbols {
		charset += symbols
	}

	password := make([]byte, length)
	for i := range password {
		randomIndex, err := randomInt(len(charset))
		if err != nil {
			return "", fmt.Errorf("failed to generate random password: %w", err)
		}
		password[i] = charset[randomIndex]
	}

	return string(password), nil
}

// randomInt 生成随机整数
func randomInt(maxVal int) (int, error) {
	bytes := make([]byte, 1)
	_, err := rand.Read(bytes)
	if err != nil {
		return 0, err
	}
	return int(bytes[0]) % maxVal, nil
}

// GenerateAPIToken 生成API令牌
func GenerateAPIToken(prefix string) (token, hashedToken string, err error) {
	// 生成32字节随机数据
	randomBytes := make([]byte, 32)
	if _, err = rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 转换为十六进制字符串
	tokenSuffix := hex.EncodeToString(randomBytes)

	// 组合完整令牌
	token = fmt.Sprintf("%s_%s", prefix, tokenSuffix)

	// 生成令牌哈希用于存储
	hash := sha256.Sum256([]byte(token))
	hashedToken = hex.EncodeToString(hash[:])

	return token, hashedToken, nil
}

// HashAPIToken 对API令牌进行哈希
func HashAPIToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

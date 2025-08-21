package validator

import (
	"api-service/pkg/errors"
	"regexp"
	"strings"
	"unicode"
)

// 密码验证相关常量
const (
	MinPasswordRequirements = 2   // 最少满足的密码要求数量
	EmailPartsCount         = 2   // 邮箱地址@分割后的部分数量
	MaxApplications         = 10  // 应用最大数量
	MaxWorkflows            = 5   // 工作流最大数量
	MaxFiles                = 100 // 文件最大数量
)

// 系统保留用户名
var reservedUsernames = map[string]bool{
	"admin":     true,
	"root":      true,
	"system":    true,
	"websoft9":  true,
	"api":       true,
	"www":       true,
	"ftp":       true,
	"mail":      true,
	"test":      true,
	"guest":     true,
	"anonymous": true,
}

// 允许的邮箱域名白名单（如果为空则允许所有域名）
var allowedEmailDomains = []string{
	// "company.com",
	// "websoft9.com",
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	// 检查是否为空
	if username == "" {
		return errors.NewAppError(errors.CodeValidationError, "用户名不能为空")
	}

	// 检查长度
	if len(username) < 3 || len(username) > 20 {
		return errors.NewAppError(errors.CodeValidationError, "用户名长度必须在3-20字符之间")
	}

	// 检查字符规则：只允许字母、数字、下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.NewAppError(errors.CodeValidationError, "用户名只能包含字母、数字和下划线")
	}

	// 检查是否以字母开头
	if !unicode.IsLetter(rune(username[0])) {
		return errors.NewAppError(errors.CodeValidationError, "用户名必须以字母开头")
	}

	// 检查是否为保留用户名
	if reservedUsernames[strings.ToLower(username)] {
		return errors.ErrUsernameReserved
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	// 检查是否为空
	if password == "" {
		return errors.NewAppError(errors.CodeValidationError, "密码不能为空")
	}

	// 检查长度
	if len(password) < 6 || len(password) > 50 {
		return errors.NewAppError(errors.CodeValidationError, "密码长度必须在6-50字符之间")
	}

	// 检查密码复杂度
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 至少包含大写字母、小写字母、数字中的两种
	requirements := 0
	if hasUpper {
		requirements++
	}
	if hasLower {
		requirements++
	}
	if hasNumber {
		requirements++
	}
	if hasSpecial {
		requirements++
	}

	if requirements < MinPasswordRequirements {
		return errors.ErrPasswordTooWeak
	}

	return nil
}

// ValidateEmail 验证邮箱格式和域名
func ValidateEmail(email string) error {
	// 检查是否为空
	if email == "" {
		return errors.NewAppError(errors.CodeValidationError, "邮箱不能为空")
	}

	// 基本格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.ErrInvalidEmail
	}

	// 域名白名单验证（如果配置了白名单）
	if len(allowedEmailDomains) > 0 {
		parts := strings.Split(email, "@")
		if len(parts) != EmailPartsCount {
			return errors.ErrInvalidEmail
		}

		domain := strings.ToLower(parts[1])
		allowed := false
		for _, allowedDomain := range allowedEmailDomains {
			if strings.EqualFold(domain, allowedDomain) {
				allowed = true
				break
			}
		}

		if !allowed {
			return errors.NewAppError(errors.CodeValidationError, "不允许使用该邮箱域名")
		}
	}

	return nil
}

// ValidateUserStatus 验证用户状态转换
func ValidateUserStatus(currentStatus, newStatus string) error {
	// 定义允许的状态转换
	allowedTransitions := map[string][]string{
		"inactive": {"active", "banned"},
		"active":   {"inactive", "banned"},
		"banned":   {"inactive"},
	}

	validStatuses := []string{"active", "inactive", "banned"}

	// 检查新状态是否有效
	isValidStatus := false
	for _, status := range validStatuses {
		if newStatus == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return errors.NewAppError(errors.CodeValidationError, "无效的用户状态")
	}

	// 检查状态转换是否被允许
	if allowedNextStatuses, exists := allowedTransitions[currentStatus]; exists {
		for _, allowedStatus := range allowedNextStatuses {
			if newStatus == allowedStatus {
				return nil
			}
		}
		return errors.NewAppError(errors.CodeValidationError, "不允许的状态转换")
	}

	return errors.NewAppError(errors.CodeValidationError, "当前状态不支持转换")
}

// ValidateUserPermission 验证用户权限
func ValidateUserPermission(userRole, requiredPermission string) error {
	// 定义角色权限映射
	rolePermissions := map[string][]string{
		"admin": {
			"user:create", "user:read", "user:update", "user:delete",
			"app:create", "app:read", "app:update", "app:delete",
		},
		"user":  {"user:read", "app:create", "app:read", "app:update"},
		"guest": {"user:read", "app:read"},
	}

	permissions, exists := rolePermissions[userRole]
	if !exists {
		return errors.NewAppError(errors.CodeValidationError, "无效的用户角色")
	}

	for _, permission := range permissions {
		if permission == requiredPermission {
			return nil
		}
	}

	return errors.ErrForbidden
}

// ValidateResourceQuota 验证资源配额
func ValidateResourceQuota(userID uint, resourceType string, currentCount int) error {
	// 定义资源配额限制
	quotaLimits := map[string]int{
		"applications": MaxApplications,
		"workflows":    MaxWorkflows,
		"files":        MaxFiles,
	}

	limit, exists := quotaLimits[resourceType]
	if !exists {
		return errors.NewAppError(errors.CodeValidationError, "未知的资源类型")
	}

	if currentCount >= limit {
		return errors.ErrUserQuotaExceeded
	}

	return nil
}

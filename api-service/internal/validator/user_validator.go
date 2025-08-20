package validator

import (
	"api-service/internal/service"
	"api-service/pkg/response"
	"errors"
	"regexp"
	"unicode"
)

// UserValidator 用户验证器
type UserValidator struct{}

// NewUserValidator 创建用户验证器
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateCreateUser 验证创建用户请求
func (v *UserValidator) ValidateCreateUser(req service.CreateUserRequest) []response.ValidationErrorDetail {
	var errors []response.ValidationErrorDetail

	// 验证用户名: 3-32字符，字母数字下划线
	if err := v.validateUsername(req.Username); err != nil {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "username",
			Message: err.Error(),
			Code:    "INVALID_USERNAME",
		})
	}

	// 验证邮箱格式
	if err := v.validateEmail(req.Email); err != nil {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "email",
			Message: err.Error(),
			Code:    "INVALID_EMAIL",
		})
	}

	// 验证密码: 8-32字符，包含大小写字母和数字
	if err := v.validatePassword(req.Password); err != nil {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "password",
			Message: err.Error(),
			Code:    "INVALID_PASSWORD",
		})
	}

	// 验证手机号（如果提供）
	if req.Phone != "" {
		if err := v.validatePhone(req.Phone); err != nil {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "phone",
				Message: err.Error(),
				Code:    "INVALID_PHONE",
			})
		}
	}

	// 验证性别
	if req.Gender < 0 || req.Gender > 2 {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "gender",
			Message: "性别值必须为0(未知)、1(男)、2(女)",
			Code:    "INVALID_GENDER",
		})
	}

	// 验证状态
	if req.Status < 0 || req.Status > 1 {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "status",
			Message: "状态值必须为0(禁用)或1(启用)",
			Code:    "INVALID_STATUS",
		})
	}

	return errors
}

// ValidateUpdateUser 验证更新用户请求
func (v *UserValidator) ValidateUpdateUser(req service.UpdateUserRequest) []response.ValidationErrorDetail {
	var errors []response.ValidationErrorDetail

	// 验证用户名（如果提供）
	if req.Username != nil {
		if err := v.validateUsername(*req.Username); err != nil {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "username",
				Message: err.Error(),
				Code:    "INVALID_USERNAME",
			})
		}
	}

	// 验证邮箱（如果提供）
	if req.Email != nil {
		if err := v.validateEmail(*req.Email); err != nil {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "email",
				Message: err.Error(),
				Code:    "INVALID_EMAIL",
			})
		}
	}

	// 验证手机号（如果提供）
	if req.Phone != nil && *req.Phone != "" {
		if err := v.validatePhone(*req.Phone); err != nil {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "phone",
				Message: err.Error(),
				Code:    "INVALID_PHONE",
			})
		}
	}

	// 验证性别（如果提供）
	if req.Gender != nil {
		if *req.Gender < 0 || *req.Gender > 2 {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "gender",
				Message: "性别值必须为0(未知)、1(男)、2(女)",
				Code:    "INVALID_GENDER",
			})
		}
	}

	// 验证状态（如果提供）
	if req.Status != nil {
		if *req.Status < 0 || *req.Status > 1 {
			errors = append(errors, response.ValidationErrorDetail{
				Field:   "status",
				Message: "状态值必须为0(禁用)或1(启用)",
				Code:    "INVALID_STATUS",
			})
		}
	}

	return errors
}

// ValidateChangePassword 验证修改密码请求
func (v *UserValidator) ValidateChangePassword(req service.ChangePasswordRequest) []response.ValidationErrorDetail {
	var errors []response.ValidationErrorDetail

	// 验证新密码
	if err := v.validatePassword(req.NewPassword); err != nil {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "new_password",
			Message: err.Error(),
			Code:    "INVALID_PASSWORD",
		})
	}

	// 验证密码确认
	if req.NewPassword != req.ConfirmPassword {
		errors = append(errors, response.ValidationErrorDetail{
			Field:   "confirm_password",
			Message: "密码确认不匹配",
			Code:    "PASSWORD_MISMATCH",
		})
	}

	return errors
}

// validateUsername 验证用户名
func (v *UserValidator) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return errors.New("用户名长度必须在3-32字符之间")
	}

	// 用户名只能包含字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	// 用户名必须以字母开头
	if !unicode.IsLetter(rune(username[0])) {
		return errors.New("用户名必须以字母开头")
	}

	return nil
}

// validateEmail 验证邮箱格式
func (v *UserValidator) validateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	// 简单的邮箱格式验证
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// validatePassword 验证密码强度
func (v *UserValidator) validatePassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("密码长度必须在8-32字符之间")
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.New("密码必须包含至少一个大写字母")
	}
	if !hasLower {
		return errors.New("密码必须包含至少一个小写字母")
	}
	if !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}

	return nil
}

// validatePhone 验证手机号格式
func (v *UserValidator) validatePhone(phone string) error {
	if phone == "" {
		return nil // 手机号是可选的
	}

	// 中国手机号格式验证
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if !matched {
		return errors.New("手机号格式不正确")
	}

	return nil
}

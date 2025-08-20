package errors

import "fmt"

// 用户管理相关错误码 (按照API文档5.2.4节定义)
const (
	// 用户相关错误码
	ErrUserNotFound      = 4001 // 用户不存在
	ErrUserAlreadyExists = 4002 // 用户已存在
	ErrUserDisabled      = 4003 // 用户已禁用
	ErrInvalidPassword   = 4004 // 密码不正确
	ErrUsernameExists    = 4005 // 用户名已存在
	ErrEmailExists       = 4006 // 邮箱已存在
	ErrInvalidUserID     = 4007 // 无效的用户ID
	ErrPasswordMismatch  = 4008 // 密码确认不匹配
	ErrOldPasswordWrong  = 4009 // 旧密码错误

	// 权限相关错误码
	ErrPermissionDenied = 4031 // 权限不足
	ErrUnauthorized     = 4032 // 未认证

	// 验证相关错误码
	ErrValidationFailed = 4221 // 验证失败
	ErrInvalidParameter = 4222 // 无效参数
)

// UserError 用户管理错误
type UserError struct {
	Code    int
	Message string
	Details string
}

func (e UserError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewUserError 创建用户错误
func NewUserError(code int, message string, details ...string) *UserError {
	err := &UserError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// 预定义的错误实例
var (
	ErrUserNotFoundError      = NewUserError(ErrUserNotFound, "用户不存在")
	ErrUserAlreadyExistsError = NewUserError(ErrUserAlreadyExists, "用户已存在")
	ErrUserDisabledError      = NewUserError(ErrUserDisabled, "用户已禁用")
	ErrInvalidPasswordError   = NewUserError(ErrInvalidPassword, "密码不正确")
	ErrUsernameExistsError    = NewUserError(ErrUsernameExists, "用户名已存在")
	ErrEmailExistsError       = NewUserError(ErrEmailExists, "邮箱已存在")
	ErrInvalidUserIDError     = NewUserError(ErrInvalidUserID, "无效的用户ID")
	ErrPasswordMismatchError  = NewUserError(ErrPasswordMismatch, "密码确认不匹配")
	ErrOldPasswordWrongError  = NewUserError(ErrOldPasswordWrong, "旧密码错误")
	ErrPermissionDeniedError  = NewUserError(ErrPermissionDenied, "权限不足")
	ErrUnauthorizedError      = NewUserError(ErrUnauthorized, "未认证")
	ErrValidationFailedError  = NewUserError(ErrValidationFailed, "验证失败")
	ErrInvalidParameterError  = NewUserError(ErrInvalidParameter, "无效参数")
)

// IsUserError 判断是否为用户错误
func IsUserError(err error) bool {
	_, ok := err.(*UserError)
	return ok
}

// GetUserErrorCode 获取用户错误码
func GetUserErrorCode(err error) int {
	if userErr, ok := err.(*UserError); ok {
		return userErr.Code
	}
	return 0
}
